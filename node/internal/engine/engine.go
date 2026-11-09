// Package engine executa o laço de um swap em um nó: observa as duas cadeias, decide com
// swap.NextAction, executa com o SettlementProtocol e publica cada passo (para a API/SSE).
// O mesmo laço faz o papel de watchtower (Banco) e de watchdog de reversões (os dois papéis).
package engine

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"janus/node/internal/swap"
)

// Step é um passo publicado pelo motor.
type Step struct {
	SwapID   [32]byte          `json:"swapId"`
	Role     swap.Role         `json:"role"`
	At       time.Time         `json:"at"`
	Status   swap.Status       `json:"status"`
	Decision swap.Decision     `json:"decision"`
	Fabric   swap.FabricState  `json:"fabric"`
	Besu     swap.BesuState    `json:"besu"`
	Error    string            `json:"error,omitempty"`
	Executed bool              `json:"executed"`
	Meta     map[string]string `json:"meta,omitempty"`
}

// Faults injeta falhas para os cenários E2E (--no-lock, --no-reveal, --crash-after-reveal).
type Faults struct {
	NoLock           bool // o nó nunca trava a sua perna
	NoReveal         bool // a Vendedora nunca revela S
	CrashAfterReveal bool // o nó para logo depois da revelação (o reinício faz a reconciliação)
}

func (f Faults) blocks(d swap.Decision) bool {
	switch d.Kind {
	case swap.ActLockAsset, swap.ActLockPayment:
		return f.NoLock
	case swap.ActClaimPayment, swap.ActSettlePayment:
		return f.NoReveal
	}
	return false
}

// Config do motor.
type Config struct {
	Role     swap.Role
	Policy   swap.TimePolicy
	Protocol map[swap.Protocol]swap.SettlementProtocol
	Poll     time.Duration // reavaliação periódica (além dos eventos)
	Faults   Faults
	Log      *slog.Logger
}

// Engine roda um laço por swap ativo.
type Engine struct {
	cfg  Config
	mu   sync.Mutex
	subs map[chan Step]struct{}
	runs map[[32]byte]context.CancelFunc
	kick map[[32]byte]chan struct{}
	// ErrCrash é devolvido por Run quando Faults.CrashAfterReveal dispara.
}

// ErrCrash simula a queda do nó (Faults.CrashAfterReveal).
var ErrCrash = errors.New("engine: simulated crash after reveal")

func New(cfg Config) *Engine {
	if cfg.Poll == 0 {
		cfg.Poll = 2 * time.Second
	}
	if cfg.Log == nil {
		cfg.Log = slog.Default()
	}
	return &Engine{cfg: cfg, subs: map[chan Step]struct{}{}, runs: map[[32]byte]context.CancelFunc{}, kick: map[[32]byte]chan struct{}{}}
}

// Subscribe devolve um canal com todos os passos (a API usa para o SSE).
func (e *Engine) Subscribe() (<-chan Step, func()) {
	ch := make(chan Step, 64)
	e.mu.Lock()
	e.subs[ch] = struct{}{}
	e.mu.Unlock()
	return ch, func() {
		e.mu.Lock()
		delete(e.subs, ch)
		e.mu.Unlock()
	}
}

func (e *Engine) publish(s Step) {
	e.mu.Lock()
	defer e.mu.Unlock()
	for ch := range e.subs {
		select {
		case ch <- s:
		default: // assinante lento não bloqueia o motor
		}
	}
}

// Kick força uma reavaliação imediata (chamado quando chega um evento de uma das cadeias).
func (e *Engine) Kick(swapID [32]byte) {
	e.mu.Lock()
	ch := e.kick[swapID]
	e.mu.Unlock()
	if ch != nil {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

// Run executa o laço até um estado terminal para este nó (done/abort) ou até o contexto acabar.
// Na reconciliação de boot, Run é chamado para cada swap não terminal do armazenamento: como a
// decisão deriva só do estado das cadeias, o nó retoma exatamente de onde parou.
func (e *Engine) Run(ctx context.Context, t *swap.SwapTerms) (swap.Decision, error) {
	proto, ok := e.cfg.Protocol[t.Protocol]
	if !ok {
		return swap.Decision{}, fmt.Errorf("engine: protocol %q not configured", t.Protocol)
	}
	ctx, cancel := context.WithCancel(ctx)
	kick := make(chan struct{}, 1)
	e.mu.Lock()
	if prev, ok := e.runs[t.SwapID]; ok {
		prev()
	}
	e.runs[t.SwapID], e.kick[t.SwapID] = cancel, kick
	e.mu.Unlock()
	defer func() {
		cancel()
		e.mu.Lock()
		delete(e.runs, t.SwapID)
		delete(e.kick, t.SwapID)
		e.mu.Unlock()
	}()

	var last swap.Decision
	for {
		obs, err := proto.Observe(ctx, e.cfg.Role, t)
		step := Step{SwapID: t.SwapID, Role: e.cfg.Role, At: time.Now()}
		if err != nil {
			step.Error = err.Error()
			e.publish(step)
		} else {
			d := swap.NextAction(e.cfg.Role, t, e.cfg.Policy, obs)
			step.Status, step.Decision, step.Fabric, step.Besu = swap.Derive(obs), d, obs.Fabric.State, obs.Besu.State
			last = d
			switch d.Kind {
			case swap.ActDone, swap.ActAbort:
				e.publish(step)
				return d, nil
			case swap.ActWait, swap.ActAlert:
			default:
				if e.cfg.Faults.blocks(d) {
					step.Meta = map[string]string{"fault": "ação suprimida por falha injetada"}
				} else if err := proto.Execute(ctx, e.cfg.Role, t, d); err != nil {
					step.Error = err.Error()
					// também no log: a operação acompanha falhas sem precisar do SSE (runbook)
					e.cfg.Log.Warn("ação falhou; nova tentativa no próximo ciclo",
						"swap", fmt.Sprintf("%x", t.SwapID[:8]), "acao", d.Kind, "err", err)
				} else {
					step.Executed = true
					if e.cfg.Faults.CrashAfterReveal && (d.Kind == swap.ActClaimPayment || d.Kind == swap.ActSettlePayment || d.Kind == swap.ActClaimAsset) {
						e.publish(step)
						return d, ErrCrash
					}
				}
			}
			e.publish(step)
			if step.Executed {
				continue // reobserva imediatamente depois de agir
			}
		}
		wait := e.cfg.Poll
		if !last.WakeAt.IsZero() {
			if until := time.Until(last.WakeAt); until > 0 && until < wait {
				wait = until
			}
		}
		select {
		case <-ctx.Done():
			return last, ctx.Err()
		case <-kick:
		case <-time.After(wait):
		}
	}
}
