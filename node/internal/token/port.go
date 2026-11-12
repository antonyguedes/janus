package token

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	fscview "github.com/hyperledger-labs/fabric-smart-client/platform/view/services/view"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/view"

	"janus/node/internal/swap"
)

// ViewInitiator é a parte do view.Manager do FSC usada pela porta (permite fakes nos testes).
type ViewInitiator interface {
	InitiateView(ctx context.Context, v view.View) (any, error)
}

var _ ViewInitiator = (*fscview.Manager)(nil)

// FSCPort implementa protocol/htlc.FabricPort chamando as views deste pacote no próprio processo
// do nó FSC (sem passar pela rede).
type FSCPort struct {
	Views ViewInitiator
	// Counterparty é o identificador FSC do nó do Banco (destino do lock). Só a Vendedora usa.
	Counterparty string
	// DeadlineSkew: falha injetada (cenário 6) que desloca o deadline gravado no lock.
	DeadlineSkew time.Duration
}

func (p *FSCPort) Status(ctx context.Context, hash [32]byte) (swap.FabricView, error) {
	res, err := p.Views.InitiateView(ctx, &StatusView{Status: &Status{Hash: hash[:]}})
	if err != nil {
		return swap.FabricView{}, fmt.Errorf("fabric status: %w", err)
	}
	raw, ok := res.([]byte)
	if !ok {
		return swap.FabricView{}, fmt.Errorf("fabric status: unexpected result %T", res)
	}
	var st LockStatus
	if err := json.Unmarshal(raw, &st); err != nil {
		return swap.FabricView{}, err
	}
	return toFabricView(hash, st), nil
}

func toFabricView(hash [32]byte, st LockStatus) swap.FabricView {
	v := swap.FabricView{AssetID: st.LinearID, Sender: st.Sender, Recipient: st.Recipient,
		Hashlock: hash, HashOK: st.HashOK, Deadline: st.Deadline}
	switch st.State {
	case "locked":
		v.State = swap.FabricLocked
	case "claimed":
		v.State = swap.FabricClaimed
	case "reclaimed":
		v.State = swap.FabricReclaimed
	default:
		v.State = swap.FabricNone
	}
	return v
}

func (p *FSCPort) Expect(ctx context.Context, t *swap.SwapTerms, tolerance time.Duration) error {
	_, err := p.Views.InitiateView(ctx, &ExpectLockView{Expectation: &Expectation{
		Hash: t.Hashlock[:], LinearID: t.AssetID, Deadline: t.T1Time(), Tolerance: tolerance,
	}})
	return err
}

func (p *FSCPort) Lock(ctx context.Context, t *swap.SwapTerms) error {
	if p.Counterparty == "" {
		return fmt.Errorf("fabric lock: counterparty node not configured")
	}
	_, err := p.Views.InitiateView(ctx, &LockView{Lock: &Lock{
		LinearID: t.AssetID, Recipient: p.Counterparty, Hash: t.Hashlock[:], Deadline: t.T1Time().Add(p.DeadlineSkew),
	}})
	return err
}

func (p *FSCPort) Claim(ctx context.Context, preimage [32]byte) error {
	_, err := p.Views.InitiateView(ctx, &ClaimView{Claim: &Claim{Preimage: preimage[:]}})
	return err
}

func (p *FSCPort) Reclaim(ctx context.Context, hash [32]byte) error {
	_, err := p.Views.InitiateView(ctx, &ReclaimView{Reclaim: &Reclaim{Hash: hash[:]}})
	return err
}
