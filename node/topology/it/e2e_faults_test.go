package it

import (
	"bufio"
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	"janus/node/internal/api"
	"janus/node/internal/besu"
	"janus/node/internal/engine"
	"janus/node/internal/swap"
)

// collectSteps acompanha o SSE de um nó e guarda os passos do motor (para asserções de mecanismo).
type collectSteps struct {
	mu    sync.Mutex
	steps []api.StepDTO
}

func watchSteps(t *testing.T, url string) *collectSteps {
	t.Helper()
	c := &collectSteps{}
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url+"/api/v1/events", nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		defer res.Body.Close()
		sc := bufio.NewScanner(res.Body)
		sc.Buffer(make([]byte, 1<<20), 1<<20)
		for sc.Scan() {
			line := sc.Text()
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			var st api.StepDTO
			if json.Unmarshal([]byte(strings.TrimPrefix(line, "data: ")), &st) == nil && st.SwapID != "" {
				c.mu.Lock()
				c.steps = append(c.steps, st)
				c.mu.Unlock()
			}
		}
	}()
	return c
}

func (c *collectSteps) find(f func(api.StepDTO) bool) (api.StepDTO, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, s := range c.steps {
		if f(s) {
			return s, true
		}
	}
	return api.StepDTO{}, false
}

func waitTerminal(t *testing.T, url, id string, timeout time.Duration) api.SwapDTO {
	t.Helper()
	deadline := time.Now().Add(timeout)
	var last api.SwapDTO
	for time.Now().Before(deadline) {
		if d, err := getSwap(url, id); err == nil {
			last = d
			if d.Terminal {
				return d
			}
		}
		time.Sleep(2 * time.Second)
	}
	t.Fatalf("%s: swap %s não encerrou (status %s)", url, id[:12], last.Status)
	return last
}

// TestE2EBuyerNeverLocks (cenário 2, A): o Banco nunca trava o Drex. A Vendedora reavê o NFT
// depois de T1 + tolerância da auditoria, e nenhum Drex muda de mãos.
func TestE2EBuyerNeverLocks(t *testing.T) {
	o := e2eOpts{proto: swap.ProtocolHTLC, registro: "REG-E2E-2", policy: "e2e",
		faults: map[string]engine.Faults{"banco": {NoLock: true}}}
	env := bootE2E(t, o)
	before := drexOf(t, env, keyBuyer)
	id, proposed := issueAndPropose(t, env, o)
	waitStatus(t, env.urls["vendedora"], proposed.SwapID, swap.StatusAssetLocked, time.Minute)
	waitStatus(t, env.urls["vendedora"], proposed.SwapID, swap.StatusReverted, 5*time.Minute)
	if got := drexOf(t, env, keyBuyer); got != before {
		t.Fatalf("Drex do Banco mudou: %d != %d", got, before)
	}
	if got := drexOf(t, env, keySeller); got != 0 {
		t.Fatalf("a Vendedora não deveria ter recebido: %d", got)
	}
	nftWith(t, env, "vendedora", "banco", id)
}

// TestE2ERevealAtT2Minus (cenário 4, A): a Vendedora segura S até T2 − 25 s (dentro da margem de
// revelação de 20 s da política e2e). O Banco, como watchtower, ainda faz o claim antes de T1.
func TestE2ERevealAtT2Minus(t *testing.T) {
	o := e2eOpts{proto: swap.ProtocolHTLC, registro: "REG-E2E-4", policy: "e2e",
		faults: map[string]engine.Faults{"vendedora": {RevealAtT2MinusSec: 25}}}
	env := bootE2E(t, o)
	id, proposed := issueAndPropose(t, env, o)
	waitStatus(t, env.urls["banco"], proposed.SwapID, swap.StatusCompleted, 4*time.Minute)
	waitStatus(t, env.urls["vendedora"], proposed.SwapID, swap.StatusCompleted, time.Minute)
	nftWith(t, env, "banco", "vendedora", id)

	// mecanismo: o claim (revelação de S) entrou num bloco em T2 − 25 s ou depois, e antes de T2
	k, _ := crypto.HexToECDSA(keyBuyer)
	bc, err := besu.Dial(context.Background(), env.rpc, k, env.dep)
	if err != nil {
		t.Fatal(err)
	}
	defer bc.Close()
	head, _ := bc.Eth().BlockNumber(context.Background())
	evs, err := bc.Events(context.Background(), 0, head)
	if err != nil {
		t.Fatal(err)
	}
	sid := hex32(t, "0x"+proposed.SwapID)
	for _, e := range evs {
		if e.SwapID != sid || e.Kind != besu.EvClaimed {
			continue
		}
		hdr, err := bc.Eth().HeaderByNumber(context.Background(), new(big.Int).SetUint64(e.Block))
		if err != nil {
			t.Fatal(err)
		}
		at, t2 := int64(hdr.Time), proposed.Terms.T2
		if at < t2-26 || at >= t2 {
			t.Fatalf("revelação em %d, esperado em [T2-25, T2) = [%d, %d)", at, t2-25, t2)
		}
		t.Logf("S revelado em T2-%ds; o Banco ainda tinha %ds até T1 para o claim", t2-at, proposed.Terms.T1-at)
		return
	}
	t.Fatal("evento de claim não encontrado")
}

// TestE2EWrongDeadline (cenário 6, A): a Vendedora trava o NFT com deadline 60 s antes do T1
// combinado. O Banco recusa o lock (fora da tolerância da expectativa), nada é travado e ninguém
// paga; os dois motores abortam.
func TestE2EWrongDeadline(t *testing.T) {
	o := e2eOpts{proto: swap.ProtocolHTLC, registro: "REG-E2E-6", policy: "e2e",
		faults: map[string]engine.Faults{"vendedora": {LockDeadlineSkewSec: -60}}}
	env := bootE2E(t, o)
	seller := watchSteps(t, env.urls["vendedora"])
	before := drexOf(t, env, keyBuyer)
	id, proposed := issueAndPropose(t, env, o)
	s := waitTerminal(t, env.urls["vendedora"], proposed.SwapID, 3*time.Minute)
	b := waitTerminal(t, env.urls["banco"], proposed.SwapID, 2*time.Minute)
	if s.Status != swap.StatusNegotiated || b.Status != swap.StatusNegotiated {
		t.Fatalf("nada deveria ter sido travado: vendedora=%s banco=%s", s.Status, b.Status)
	}
	// mecanismo: o lock foi tentado e recusado pelo Banco por causa do deadline
	if _, ok := seller.find(func(st api.StepDTO) bool {
		return st.SwapID == proposed.SwapID && st.Decision.Kind == swap.ActLockAsset &&
			strings.Contains(st.Error, "deadline differs")
	}); !ok {
		t.Fatal("esperava a recusa do lock pelo Banco (deadline fora da tolerância)")
	}
	if got := drexOf(t, env, keyBuyer); got != before {
		t.Fatalf("Drex do Banco mudou: %d != %d", got, before)
	}
	nftWith(t, env, "vendedora", "banco", id)
}
