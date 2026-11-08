package token

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/LFDT-Panurus/panurus/token"
	"github.com/LFDT-Panurus/panurus/token/services/ttx"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/comm/session"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/id"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/view"

	"janus/node/internal/swap"
)

// Negociação dos SwapTerms por sessão P2P do FSC (canal autenticado). Substitui a PDC: o valor do
// acordo não vai ao ledger. A preimagem S nunca trafega. As duas partes assinam o hash canônico
// dos termos com a chave da identidade de owner, a mesma que aparece no script HTLC.

// NegotiationHooks conectam as views ao nó Janus (store, motor, política de aceite). São
// configurados pela SDK do Janus no PostStart. As views de responder são instanciadas pelo
// registry do FSC sem dependências.
type NegotiationHooks struct {
	// LocalParty devolve MSP e endereço EVM deste nó (a identidade FSC vem da wallet).
	LocalParty func() swap.Party
	// Policy é a política de tempo usada para validar termos recebidos.
	Policy swap.TimePolicy
	// Accept decide se o Banco aceita os termos (preço, ativo, limites). nil = aceita.
	Accept func(ctx context.Context, t *swap.SwapTerms) error
	// OnSigned é chamado quando os termos têm as duas assinaturas (persistir e iniciar o motor).
	OnSigned func(ctx context.Context, role swap.Role, st *swap.SignedTerms) error
}

var (
	hooksMu sync.RWMutex
	hooks   NegotiationHooks
)

// SetNegotiationHooks instala os hooks (chamado pela SDK do Janus).
func SetNegotiationHooks(h NegotiationHooks) {
	hooksMu.Lock()
	defer hooksMu.Unlock()
	hooks = h
}

func currentHooks() (NegotiationHooks, error) {
	hooksMu.RLock()
	defer hooksMu.RUnlock()
	if hooks.LocalParty == nil {
		return hooks, errors.New("negotiation: janus node not configured")
	}
	return hooks, nil
}

// --------------------------------------------------------------------- assinatura FSC

type ownerSigner struct {
	ctx context.Context
	sig token.Signer
}

func (s ownerSigner) Sign(digest []byte) ([]byte, error) { return s.sig.Sign(digest) }

type ownerVerifier struct {
	ctx context.Context
	tms *token.ManagementService
}

func (v ownerVerifier) Verify(identity, digest, sig []byte) error {
	ver, err := v.tms.SigService().OwnerVerifier(v.ctx, identity)
	if err != nil {
		return fmt.Errorf("no verifier for identity: %w", err)
	}
	return ver.Verify(digest, sig)
}

func myOwnerIdentity(ctx view.Context) (token.Identity, *token.ManagementService, error) {
	tms, err := token.GetManagementService(ctx)
	if err != nil {
		return nil, nil, err
	}
	w := ttx.MyWallet(ctx)
	if w == nil {
		return nil, nil, errors.New("negotiation: owner wallet not found")
	}
	me, err := w.GetRecipientIdentity(ctx.Context())
	if err != nil {
		return nil, nil, err
	}
	return me, tms, nil
}

// --------------------------------------------------------------------- mensagens

type partyMsg struct {
	Party swap.Party `json:"party"`
}

// --------------------------------------------------------------------- Vendedora (iniciador)

// Propose: a Vendedora propõe um swap ao Banco. Hashlock e Nonce vêm do chamador (S já foi
// gerado e guardado cifrado antes de qualquer mensagem).
type Propose struct {
	BuyerNode string           `json:"buyerNode"`
	Protocol  swap.Protocol    `json:"protocol"`
	AssetID   string           `json:"assetId"`
	Price     swap.Amount      `json:"price"`
	Hashlock  [32]byte         `json:"hashlock"`
	Nonce     [32]byte         `json:"nonce"`
	T1        int64            `json:"t1,omitempty"` // 0 = agora + política
	T2        int64            `json:"t2,omitempty"`
	Networks  swap.NetworkRefs `json:"networks"`
}

type ProposeView struct{ *Propose }

func (v *ProposeView) Call(ctx view.Context) (any, error) {
	h, err := currentHooks()
	if err != nil {
		return nil, err
	}
	me, tms, err := myOwnerIdentity(ctx)
	if err != nil {
		return nil, err
	}
	idp, err := id.GetProvider(ctx)
	if err != nil {
		return nil, err
	}
	s, err := session.NewJSON(ctx, v, idp.Identity(v.BuyerNode))
	if err != nil {
		return nil, fmt.Errorf("propose: open session: %w", err)
	}
	// 1) pede a identidade de owner e o endereço EVM do Banco
	if err := s.Send(ctx.Context(), &partyMsg{}); err != nil {
		return nil, err
	}
	var buyer partyMsg
	if err := s.ReceiveWithTimeout(&buyer, 30*time.Second); err != nil {
		return nil, fmt.Errorf("propose: buyer party: %w", err)
	}
	// 2) monta e assina os termos canônicos
	seller := h.LocalParty()
	seller.FSCIdentity = me
	t := swap.SwapTerms{Protocol: v.Protocol, AssetID: v.AssetID, Seller: seller, Buyer: buyer.Party,
		Price: v.Price, Hashlock: v.Hashlock, Nonce: v.Nonce, T1: v.T1, T2: v.T2, Networks: v.Networks}
	now := time.Now()
	if t.T1 == 0 || t.T2 == 0 {
		t.T2, t.T1 = now.Add(h.Policy.T2Offset).Unix(), now.Add(h.Policy.T1Offset).Unix()
	}
	t.SwapID = swap.ComputeSwapID(t.Hashlock, t.Seller, t.Buyer, t.Nonce)
	if err := t.Validate(h.Policy, now); err != nil {
		return nil, fmt.Errorf("propose: %w", err)
	}
	signer, err := tms.SigService().GetSigner(ctx.Context(), me)
	if err != nil {
		return nil, err
	}
	st := &swap.SignedTerms{Terms: t}
	if err := st.SignAs(swap.RoleSeller, ownerSigner{ctx.Context(), signer}); err != nil {
		return nil, err
	}
	if err := s.Send(ctx.Context(), st); err != nil {
		return nil, err
	}
	// 3) recebe a assinatura do Banco e confere as duas
	var back swap.SignedTerms
	if err := s.ReceiveWithTimeout(&back, 30*time.Second); err != nil {
		return nil, fmt.Errorf("propose: buyer answer: %w", err)
	}
	if back.Terms.Hash() != st.Terms.Hash() {
		return nil, errors.New("propose: buyer changed the terms")
	}
	st.BuyerSig = back.BuyerSig
	if err := st.Verify(ownerVerifier{ctx.Context(), tms}); err != nil {
		return nil, fmt.Errorf("propose: %w", err)
	}
	if h.OnSigned != nil {
		if err := h.OnSigned(ctx.Context(), swap.RoleSeller, st); err != nil {
			return nil, err
		}
	}
	return json.Marshal(st)
}

type ProposeViewFactory struct{}

func (f *ProposeViewFactory) NewView(in []byte) (view.View, error) {
	v := &ProposeView{Propose: &Propose{}}
	return v, unmarshalInto(in, v.Propose)
}

// --------------------------------------------------------------------- Banco (responder)

// AcceptProposalView: responder do Banco para ProposeView.
type AcceptProposalView struct{}

func (a *AcceptProposalView) Call(ctx view.Context) (any, error) {
	h, err := currentHooks()
	if err != nil {
		return nil, err
	}
	me, tms, err := myOwnerIdentity(ctx)
	if err != nil {
		return nil, err
	}
	s := session.JSON(ctx)
	var hello partyMsg
	if err := s.ReceiveWithTimeout(&hello, 30*time.Second); err != nil {
		return nil, err
	}
	local := h.LocalParty()
	local.FSCIdentity = me
	if err := s.Send(ctx.Context(), &partyMsg{Party: local}); err != nil {
		return nil, err
	}
	var st swap.SignedTerms
	if err := s.ReceiveWithTimeout(&st, 30*time.Second); err != nil {
		return nil, err
	}
	t := &st.Terms
	reject := func(err error) (any, error) {
		_ = s.SendError(ctx.Context(), err.Error())
		return nil, err
	}
	// O Banco confere que os termos o nomeiam (identidade e EVM) e que são válidos agora.
	if string(t.Buyer.FSCIdentity) != string(me) || t.Buyer.EVM != local.EVM || t.Buyer.MSPID != local.MSPID {
		return reject(errors.New("accept: terms do not name this node as buyer"))
	}
	if err := t.Validate(h.Policy, time.Now()); err != nil {
		return reject(err)
	}
	if err := (ownerVerifier{ctx.Context(), tms}).Verify(t.Seller.FSCIdentity, hashOf(t), st.SellerSig); err != nil {
		return reject(fmt.Errorf("accept: seller signature: %w", err))
	}
	if h.Accept != nil {
		if err := h.Accept(ctx.Context(), t); err != nil {
			return reject(err)
		}
	}
	signer, err := tms.SigService().GetSigner(ctx.Context(), me)
	if err != nil {
		return reject(err)
	}
	if err := st.SignAs(swap.RoleBuyer, ownerSigner{ctx.Context(), signer}); err != nil {
		return reject(err)
	}
	// Registra a expectativa ANTES de responder: o lock da Vendedora pode chegar logo depois.
	DefaultExpectation.Put(Expectation{Hash: t.Hashlock[:], LinearID: t.AssetID, Deadline: t.T1Time(),
		Tolerance: h.Policy.DeadlineTolerance})
	if h.OnSigned != nil {
		if err := h.OnSigned(ctx.Context(), swap.RoleBuyer, &st); err != nil {
			return reject(err)
		}
	}
	return nil, s.Send(ctx.Context(), &st)
}

func hashOf(t *swap.SwapTerms) []byte {
	h := t.Hash()
	return h[:]
}
