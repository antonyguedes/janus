package besu

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// EventKind classifica os eventos dos contratos de pagamento.
type EventKind string

const (
	EvLocked   EventKind = "locked"   // HTLCPayment.PaymentLocked / AttestedSettlement.PaymentDeposited
	EvClaimed  EventKind = "claimed"  // HTLCPayment.PaymentClaimed / AttestedSettlement.PaymentSettled (traz S)
	EvRefunded EventKind = "refunded" // PaymentRefunded
)

// Event é um evento normalizado de pagamento.
type Event struct {
	Kind     EventKind   `json:"kind"`
	Contract string      `json:"contract"` // "htlc" | "settlement"
	SwapID   [32]byte    `json:"swapId"`
	Preimage *[32]byte   `json:"preimage,omitempty"`
	Block    uint64      `json:"block"`
	TxHash   common.Hash `json:"txHash"`
}

// Checkpoint persiste o último bloco processado: no boot o watcher faz backfill a partir dele.
type Checkpoint interface {
	Load() (uint64, error)
	Save(block uint64) error
}

// FileCheckpoint guarda o checkpoint em um arquivo JSON, com escrita atômica (tmp + rename).
type FileCheckpoint struct{ Path string }

func (f FileCheckpoint) Load() (uint64, error) {
	raw, err := os.ReadFile(f.Path)
	if errors.Is(err, os.ErrNotExist) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	var v struct{ Block uint64 }
	return v.Block, json.Unmarshal(raw, &v)
}

func (f FileCheckpoint) Save(block uint64) error {
	if err := os.MkdirAll(filepath.Dir(f.Path), 0o700); err != nil {
		return err
	}
	tmp := f.Path + ".tmp"
	raw, _ := json.Marshal(struct{ Block uint64 }{block})
	if err := os.WriteFile(tmp, raw, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, f.Path)
}

// WatchOptions configura o watcher.
type WatchOptions struct {
	FromBlock     uint64        // bloco inicial quando não há checkpoint (ex.: bloco do deploy)
	Confirmations uint64        // QBFT tem finalidade imediata: 0 é seguro
	Poll          time.Duration // intervalo de polling
	Chunk         uint64        // tamanho máximo da janela de eth_getLogs
}

// Watch acompanha os eventos de HTLCPayment e AttestedSettlement. Faz backfill a partir do
// checkpoint e depois polling. Só avança o checkpoint depois de entregar os eventos da janela,
// então um crash nunca perde evento (no pior caso, reentrega: os consumidores são idempotentes).
func (c *Client) Watch(ctx context.Context, cp Checkpoint, o WatchOptions, out chan<- Event) error {
	if o.Poll == 0 {
		o.Poll = time.Second
	}
	if o.Chunk == 0 {
		o.Chunk = 2000
	}
	last, err := cp.Load()
	if err != nil {
		return fmt.Errorf("besu: checkpoint: %w", err)
	}
	next := max(last+1, o.FromBlock)
	if last == 0 && o.FromBlock == 0 {
		next = 0
	}
	t := time.NewTicker(o.Poll)
	defer t.Stop()
	for {
		head, err := c.eth.BlockNumber(ctx)
		if err == nil && head >= o.Confirmations {
			safe := head - o.Confirmations
			for next <= safe {
				to := min(next+o.Chunk-1, safe)
				evs, err := c.Events(ctx, next, to)
				if err != nil {
					break // tenta de novo no próximo tick
				}
				for _, e := range evs {
					select {
					case out <- e:
					case <-ctx.Done():
						return ctx.Err()
					}
				}
				if err := cp.Save(to); err != nil {
					return fmt.Errorf("besu: save checkpoint: %w", err)
				}
				next = to + 1
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
		}
	}
}

// Events lê e decodifica os eventos de pagamento no intervalo [from, to].
func (c *Client) Events(ctx context.Context, from, to uint64) ([]Event, error) {
	logs, err := c.eth.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(from),
		ToBlock:   new(big.Int).SetUint64(to),
		Addresses: []common.Address{c.d.HTLCPayment, c.d.AttestedSettlement},
	})
	if err != nil {
		return nil, fmt.Errorf("besu: getLogs [%d,%d]: %w", from, to, err)
	}
	var out []Event
	for i := range logs {
		if e, ok := c.decode(&logs[i]); ok {
			out = append(out, e)
		}
	}
	return out, nil
}

func (c *Client) decode(l *types.Log) (Event, bool) {
	e := Event{Block: l.BlockNumber, TxHash: l.TxHash}
	switch l.Address {
	case c.d.HTLCPayment:
		e.Contract = "htlc"
		if v, err := c.htlc.UnpackPaymentLockedEvent(l); err == nil {
			e.Kind, e.SwapID = EvLocked, v.SwapId
		} else if v, err := c.htlc.UnpackPaymentClaimedEvent(l); err == nil {
			p := v.Preimage
			e.Kind, e.SwapID, e.Preimage = EvClaimed, v.SwapId, &p
		} else if v, err := c.htlc.UnpackPaymentRefundedEvent(l); err == nil {
			e.Kind, e.SwapID = EvRefunded, v.SwapId
		} else {
			return e, false
		}
	case c.d.AttestedSettlement:
		e.Contract = "settlement"
		if v, err := c.settlement.UnpackPaymentDepositedEvent(l); err == nil {
			e.Kind, e.SwapID = EvLocked, v.SwapId
		} else if v, err := c.settlement.UnpackPaymentSettledEvent(l); err == nil {
			p := v.Preimage
			e.Kind, e.SwapID, e.Preimage = EvClaimed, v.SwapId, &p
		} else if v, err := c.settlement.UnpackPaymentRefundedEvent(l); err == nil {
			e.Kind, e.SwapID = EvRefunded, v.SwapId
		} else {
			return e, false
		}
	default:
		return e, false
	}
	return e, true
}
