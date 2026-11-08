// Package htlc implementa a Abordagem A: HTLC nativo do Panurus no Fabric + HTLCPayment no Besu.
package htlc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"janus/node/internal/swap"
)

// FabricPort é o lado Fabric visto pelo protocolo. Em produção, chamadas às views do FSC
// (internal/token). Nos testes, um fake em memória.
type FabricPort interface {
	// Status devolve o estado do lock do hashlock visto por este nó.
	Status(ctx context.Context, hash [32]byte) (swap.FabricView, error)
	// Expect registra (Banco) os termos aceitos para um lock futuro.
	Expect(ctx context.Context, t *swap.SwapTerms, tolerance time.Duration) error
	Lock(ctx context.Context, t *swap.SwapTerms) error
	Claim(ctx context.Context, preimage [32]byte) error
	Reclaim(ctx context.Context, hash [32]byte) error
}

// BesuPort é o lado Besu visto pelo protocolo (implementado por besu.Client).
type BesuPort interface {
	View(ctx context.Context, t *swap.SwapTerms) (swap.BesuView, error)
	LockPayment(ctx context.Context, t *swap.SwapTerms) error
	ClaimPayment(ctx context.Context, swapID, preimage [32]byte) error
	RefundPayment(ctx context.Context, swapID [32]byte) error
}

// SecretStore guarda a preimagem S (a Vendedora gera, o Banco aprende).
type SecretStore interface {
	Get(ctx context.Context, swapID [32]byte) (*[32]byte, error)
	Put(ctx context.Context, swapID, preimage [32]byte) error
}

// Protocol é a Abordagem A.
type Protocol struct {
	Fabric  FabricPort
	Besu    BesuPort
	Secrets SecretStore
	Policy  swap.TimePolicy
}

func (p *Protocol) Name() swap.Protocol { return swap.ProtocolHTLC }

func (p *Protocol) Observe(ctx context.Context, role swap.Role, t *swap.SwapTerms) (swap.Observation, error) {
	o := swap.Observation{Now: time.Now()}
	var err error
	if o.Fabric, err = p.Fabric.Status(ctx, t.Hashlock); err != nil {
		return o, fmt.Errorf("observe fabric: %w", err)
	}
	if o.Besu, err = p.Besu.View(ctx, t); err != nil {
		return o, fmt.Errorf("observe besu: %w", err)
	}
	if o.Preimage, err = p.Secrets.Get(ctx, t.SwapID); err != nil {
		return o, fmt.Errorf("observe secrets: %w", err)
	}
	// O Banco aprende S do Besu e o persiste antes de qualquer outra coisa (watchtower).
	if role == swap.RoleBuyer && o.Preimage == nil && o.Besu.Preimage != nil {
		if err := p.Secrets.Put(ctx, t.SwapID, *o.Besu.Preimage); err != nil {
			return o, err
		}
		o.Preimage = o.Besu.Preimage
	}
	return o, nil
}

// ErrNotExecutable indica uma decisão sem efeito nas cadeias (wait/done/abort/alert).
var ErrNotExecutable = errors.New("htlc: decision has no on-chain action")

func (p *Protocol) Execute(ctx context.Context, role swap.Role, t *swap.SwapTerms, d swap.Decision) error {
	pre := func() ([32]byte, error) {
		s, err := p.Secrets.Get(ctx, t.SwapID)
		if err != nil || s == nil {
			return [32]byte{}, fmt.Errorf("preimage unavailable: %v", err)
		}
		return *s, nil
	}
	switch d.Kind {
	case swap.ActLockAsset:
		return p.Fabric.Lock(ctx, t)
	case swap.ActLockPayment:
		return p.Besu.LockPayment(ctx, t)
	case swap.ActClaimPayment:
		s, err := pre()
		if err != nil {
			return err
		}
		return p.Besu.ClaimPayment(ctx, t.SwapID, s)
	case swap.ActClaimAsset:
		s, err := pre()
		if err != nil {
			return err
		}
		return p.Fabric.Claim(ctx, s)
	case swap.ActRefundPayment:
		return p.Besu.RefundPayment(ctx, t.SwapID)
	case swap.ActReclaimAsset:
		return p.Fabric.Reclaim(ctx, t.Hashlock)
	case swap.ActRequestAttest, swap.ActSettlePayment:
		return fmt.Errorf("htlc: %s belongs to the cacti protocol", d.Kind)
	}
	return ErrNotExecutable
}

// Accept prepara o nó do Banco para receber o lock (registra a expectativa a partir dos termos).
func (p *Protocol) Accept(ctx context.Context, t *swap.SwapTerms) error {
	return p.Fabric.Expect(ctx, t, p.Policy.DeadlineTolerance)
}
