// Package cacti implementa a Abordagem B2, "HTLC atestado" (ADR-0003): o NFT fica travado no HTLC do
// Panurus e o AttestedSettlement no Besu só libera o Drex com a preimagem S e uma atestação do
// lock no Fabric assinada pela política k-de-n de orgs (FabricProofVerifier, P-256 nativo).
//
// O transporte das atestações (relay) é feito pelo nó Janus via sessões FSC. A interface
// Collector permite trocar o transporte por um conector do Hyperledger Cacti sem mudar o protocolo.
package cacti

import (
	"context"
	"fmt"
	"sync"
	"time"

	"janus/node/internal/attest"
	"janus/node/internal/protocol/htlc"
	"janus/node/internal/swap"
)

// SettlementPort é o AttestedSettlement visto pelo protocolo (implementado sobre besu.Client).
type SettlementPort interface {
	View(ctx context.Context, t *swap.SwapTerms) (swap.BesuView, error)
	Deposit(ctx context.Context, t *swap.SwapTerms, assetIDHash, recipientIDHash [32]byte, minFabricDeadline int64) error
	Settle(ctx context.Context, swapID, preimage [32]byte, att attest.FabricLockAttestation, sigs []attest.OrgSignature) error
	Refund(ctx context.Context, swapID [32]byte) error
	// Domain é o domínio da atestação (chainId + endereço do verificador).
	Domain() attest.Domain
}

// Collector obtém as atestações do lock (relay).
type Collector interface {
	Collect(ctx context.Context, swapID, hash [32]byte, d attest.Domain) (attest.FabricLockAttestation, []attest.OrgSignature, error)
}

// Protocol é a Abordagem B2.
type Protocol struct {
	Fabric    htlc.FabricPort
	Payment   SettlementPort
	Collector Collector
	Secrets   htlc.SecretStore
	Policy    swap.TimePolicy

	mu        sync.Mutex
	collected map[[32]byte]collected
}

type collected struct {
	att  attest.FabricLockAttestation
	sigs []attest.OrgSignature
}

func (p *Protocol) Name() swap.Protocol { return swap.ProtocolCacti }

func (p *Protocol) cached(id [32]byte) (collected, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	c, ok := p.collected[id]
	return c, ok
}

func (p *Protocol) Observe(ctx context.Context, role swap.Role, t *swap.SwapTerms) (swap.Observation, error) {
	o := swap.Observation{Now: time.Now()}
	var err error
	if o.Fabric, err = p.Fabric.Status(ctx, t.Hashlock); err != nil {
		return o, fmt.Errorf("observe fabric: %w", err)
	}
	if o.Besu, err = p.Payment.View(ctx, t); err != nil {
		return o, fmt.Errorf("observe besu: %w", err)
	}
	if o.Preimage, err = p.Secrets.Get(ctx, t.SwapID); err != nil {
		return o, err
	}
	if role == swap.RoleBuyer && o.Preimage == nil && o.Besu.Preimage != nil {
		if err := p.Secrets.Put(ctx, t.SwapID, *o.Besu.Preimage); err != nil {
			return o, err
		}
		o.Preimage = o.Besu.Preimage
	}
	_, o.AttestationReady = p.cached(t.SwapID)
	return o, nil
}

func (p *Protocol) Execute(ctx context.Context, role swap.Role, t *swap.SwapTerms, d swap.Decision) error {
	switch d.Kind {
	case swap.ActLockAsset:
		return p.Fabric.Lock(ctx, t)
	case swap.ActLockPayment:
		// Os termos do depósito amarram a atestação: ativo, recipient (identidade de owner do
		// Banco no script HTLC) e o deadline mínimo do lock (T1 − tolerância).
		return p.Payment.Deposit(ctx, t, attest.HashString(t.AssetID), attest.HashBytes(t.Buyer.FSCIdentity),
			t.T1-int64(p.Policy.DeadlineTolerance/time.Second))
	case swap.ActRequestAttest:
		att, sigs, err := p.Collector.Collect(ctx, t.SwapID, t.Hashlock, p.Payment.Domain())
		if err != nil {
			return fmt.Errorf("collect attestations: %w", err)
		}
		p.mu.Lock()
		if p.collected == nil {
			p.collected = map[[32]byte]collected{}
		}
		p.collected[t.SwapID] = collected{att, sigs}
		p.mu.Unlock()
		return nil
	case swap.ActSettlePayment:
		c, ok := p.cached(t.SwapID)
		if !ok {
			return fmt.Errorf("settle: attestations not collected")
		}
		s, err := p.Secrets.Get(ctx, t.SwapID)
		if err != nil || s == nil {
			return fmt.Errorf("settle: preimage unavailable: %v", err)
		}
		return p.Payment.Settle(ctx, t.SwapID, *s, c.att, c.sigs)
	case swap.ActClaimAsset:
		s, err := p.Secrets.Get(ctx, t.SwapID)
		if err != nil || s == nil {
			return fmt.Errorf("claim: preimage unavailable: %v", err)
		}
		return p.Fabric.Claim(ctx, *s)
	case swap.ActRefundPayment:
		return p.Payment.Refund(ctx, t.SwapID)
	case swap.ActReclaimAsset:
		return p.Fabric.Reclaim(ctx, t.Hashlock)
	case swap.ActClaimPayment:
		return fmt.Errorf("cacti: %s belongs to the htlc protocol", d.Kind)
	}
	return htlc.ErrNotExecutable
}

// Accept prepara o Banco para receber o lock.
func (p *Protocol) Accept(ctx context.Context, t *swap.SwapTerms) error {
	return p.Fabric.Expect(ctx, t, p.Policy.DeadlineTolerance)
}
