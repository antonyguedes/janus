package token

import (
	"context"
	"crypto"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/LFDT-Panurus/panurus/token"
	"github.com/LFDT-Panurus/panurus/token/services/identity"
	"github.com/LFDT-Panurus/panurus/token/services/interop/encoding"
	"github.com/LFDT-Panurus/panurus/token/services/interop/htlc"
	"github.com/LFDT-Panurus/panurus/token/services/network"
	"github.com/LFDT-Panurus/panurus/token/services/nfttx"
	"github.com/LFDT-Panurus/panurus/token/services/nfttx/uniqueness"
	"github.com/LFDT-Panurus/panurus/token/services/ttx"
	token2 "github.com/LFDT-Panurus/panurus/token/token"
	"github.com/hyperledger-labs/fabric-smart-client/platform/common/utils/collections/iterators"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/id"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/view"
)

// AuditorNode é o identificador FSC do nó auditor (a Registradora).
const AuditorNode = "registradora"

// Registros em memória do processo (MVP). As factories geradas pelo NWO não recebem dependências,
// por isso ficam no nível do pacote.
var (
	DefaultGravames    = NewGravames()
	DefaultAudit       = DefaultAuditPolicy()
	DefaultExpectation = NewExpectations()
)

func auditorIdentity(ctx view.Context) (view.Identity, error) {
	p, err := id.GetProvider(ctx)
	if err != nil {
		return nil, fmt.Errorf("id provider: %w", err)
	}
	return p.Identity(AuditorNode), nil
}

func htlcHashOpts(h []byte) []token.TransferOption {
	return []token.TransferOption{
		htlc.WithHash(h),
		htlc.WithHashFunc(crypto.SHA256),
		htlc.WithHashEncoding(encoding.None),
	}
}

func unmarshalInto(in []byte, v any) error {
	if len(in) == 0 {
		return nil
	}
	if err := json.Unmarshal(in, v); err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}
	return nil
}

// =============================================================================== emissão

// Issue: a Registradora emite um recebível (NFT) para um owner.
type Issue struct {
	Recipient string    `json:"recipient"` // id do nó FSC do owner (ex.: "vendedora")
	Recebivel Recebivel `json:"recebivel"`
}

type IssueView struct{ *Issue }

func (v *IssueView) Call(ctx view.Context) (any, error) {
	r := v.Recebivel
	if err := r.Validate(); err != nil {
		return nil, err
	}
	idp, err := id.GetProvider(ctx)
	if err != nil {
		return nil, err
	}
	recipient, err := nfttx.RequestRecipientIdentity(ctx, idp.Identity(v.Recipient))
	if err != nil {
		return nil, fmt.Errorf("issue: recipient identity: %w", err)
	}
	auditor, err := auditorIdentity(ctx)
	if err != nil {
		return nil, err
	}
	tx, err := nfttx.NewAnonymousTransaction(ctx, nfttx.WithAuditor(auditor))
	if err != nil {
		return nil, fmt.Errorf("issue: new tx: %w", err)
	}
	wallet := nfttx.MyIssuerWallet(ctx)
	if wallet == nil {
		return nil, fmt.Errorf("issue: issuer wallet not found")
	}
	// O ID único é derivado do número de registro: um mesmo registro não pode ser emitido duas vezes.
	uid, err := uniqueness.GetService(ctx).ComputeID(ctx.Context(), r.Registro)
	if err != nil {
		return nil, fmt.Errorf("issue: unique id: %w", err)
	}
	if err := tx.Issue(ctx.Context(), wallet, &r, recipient, nfttx.WithUniqueID(uid)); err != nil {
		return nil, fmt.Errorf("issue: %w", err)
	}
	if _, err := ctx.RunView(nfttx.NewCollectEndorsementsView(tx)); err != nil {
		return nil, fmt.Errorf("issue: endorse: %w", err)
	}
	if _, err := ctx.RunView(nfttx.NewOrderingAndFinalityView(tx)); err != nil {
		return nil, fmt.Errorf("issue: commit: %w", err)
	}
	return r.LinearID, nil
}

type IssueViewFactory struct{}

func (f *IssueViewFactory) NewView(in []byte) (view.View, error) {
	v := &IssueView{Issue: &Issue{}}
	return v, unmarshalInto(in, v.Issue)
}

// AcceptIssuedView: responder do owner para IssueView.
type AcceptIssuedView struct{}

func (a *AcceptIssuedView) Call(ctx view.Context) (any, error) {
	me, err := nfttx.RespondRequestRecipientIdentity(ctx)
	if err != nil {
		return nil, err
	}
	tx, err := nfttx.ReceiveTransaction(ctx)
	if err != nil {
		return nil, err
	}
	outs, err := tx.Outputs(ctx.Context())
	if err != nil {
		return nil, err
	}
	if outs.Count() != 1 || outs.ByRecipient(me).Count() != 1 {
		return nil, fmt.Errorf("accept issue: expected exactly one output for me")
	}
	r := &Recebivel{}
	if err := outs.StateAt(0, r); err != nil {
		return nil, fmt.Errorf("accept issue: state: %w", err)
	}
	if err := r.Validate(); err != nil || r.LinearID == "" {
		return nil, fmt.Errorf("accept issue: invalid recebivel: %v", err)
	}
	if _, err := ctx.RunView(nfttx.NewAcceptView(tx)); err != nil {
		return nil, err
	}
	_, err = ctx.RunView(nfttx.NewFinalityView(tx))
	return nil, err
}

// =============================================================================== consulta

type Get struct {
	LinearID string `json:"linearId"`
}

type GetView struct{ *Get }

// Call devolve o recebível se ele estiver na wallet do nó (fora de lock).
func (v *GetView) Call(ctx view.Context) (any, error) {
	w := nfttx.MyWallet(ctx)
	if w == nil {
		return nil, fmt.Errorf("get: owner wallet not found")
	}
	r := &Recebivel{}
	if err := w.QueryByKey(ctx.Context(), r, "LinearID", v.LinearID); err != nil {
		return nil, fmt.Errorf("get: recebivel %s not found: %w", v.LinearID, err)
	}
	return json.Marshal(r)
}

type GetViewFactory struct{}

func (f *GetViewFactory) NewView(in []byte) (view.View, error) {
	v := &GetView{Get: &Get{}}
	return v, unmarshalInto(in, v.Get)
}

// =============================================================================== lock (Vendedora)

// Lock: a Vendedora trava o NFT no HTLC do Panurus com H e deadline absoluto T1.
type Lock struct {
	LinearID  string    `json:"linearId"`
	Recipient string    `json:"recipient"` // id do nó FSC do Banco
	Hash      []byte    `json:"hash"`      // H = SHA256(S), 32 bytes
	Deadline  time.Time `json:"deadline"`  // T1
}

// LockResult volta para o chamador.
type LockResult struct {
	TxID     string    `json:"txId"`
	Hash     []byte    `json:"hash"`
	Deadline time.Time `json:"deadline"`
}

type LockView struct{ *Lock }

func (v *LockView) Call(ctx view.Context) (any, error) {
	if len(v.Hash) != 32 {
		return nil, fmt.Errorf("lock: hash must be 32 bytes")
	}
	d := time.Until(v.Deadline)
	if d <= 0 {
		return nil, fmt.Errorf("lock: deadline in the past")
	}
	idp, err := id.GetProvider(ctx)
	if err != nil {
		return nil, err
	}
	me, recipient, err := htlc.ExchangeRecipientIdentities(ctx, "", idp.Identity(v.Recipient))
	if err != nil {
		return nil, fmt.Errorf("lock: exchange identities: %w", err)
	}
	nw := nfttx.MyWallet(ctx)
	if nw == nil {
		return nil, fmt.Errorf("lock: owner wallet not found")
	}
	r := &Recebivel{}
	if err := nw.QueryByKey(ctx.Context(), r, "LinearID", v.LinearID); err != nil {
		return nil, fmt.Errorf("lock: recebivel %s not owned: %w", v.LinearID, err)
	}
	typ, err := TokenType(r)
	if err != nil {
		return nil, err
	}
	auditor, err := auditorIdentity(ctx)
	if err != nil {
		return nil, err
	}
	tx, err := htlc.NewAnonymousTransaction(ctx, ttx.WithAuditor(auditor))
	if err != nil {
		return nil, fmt.Errorf("lock: new tx: %w", err)
	}
	if _, err := tx.Lock(ctx.Context(), nw.OwnerWallet, me, typ, 1, recipient, d, htlcHashOpts(v.Hash)...); err != nil {
		return nil, fmt.Errorf("lock: %w", err)
	}
	if _, err := ctx.RunView(htlc.NewCollectEndorsementsView(tx)); err != nil {
		return nil, fmt.Errorf("lock: endorse: %w", err)
	}
	if _, err := ctx.RunView(htlc.NewOrderingAndFinalityView(tx)); err != nil {
		return nil, fmt.Errorf("lock: commit: %w", err)
	}
	outs, err := tx.Outputs(ctx.Context())
	if err != nil {
		return nil, err
	}
	s := outs.ByScript().ScriptAt(0)
	if s == nil {
		return nil, fmt.Errorf("lock: no htlc output")
	}
	return &LockResult{TxID: tx.ID(), Hash: s.HashInfo.Hash, Deadline: s.Deadline}, nil
}

type LockViewFactory struct{}

func (f *LockViewFactory) NewView(in []byte) (view.View, error) {
	v := &LockView{Lock: &Lock{}}
	return v, unmarshalInto(in, v.Lock)
}

// =============================================================================== expectativas (Banco)

// Expectation é o que o Banco espera de um lock, a partir dos SwapTerms assinados.
type Expectation struct {
	Hash      []byte        `json:"hash"`
	LinearID  string        `json:"linearId"`
	Deadline  time.Time     `json:"deadline"` // T1
	Tolerance time.Duration `json:"tolerance"`
}

// Expectations: registro em memória, por hashlock, dos locks que o Banco aceita receber.
type Expectations struct {
	mu sync.RWMutex
	m  map[string]Expectation
	kv KV // durável (opcional)
}

func NewExpectations() *Expectations { return &Expectations{m: map[string]Expectation{}} }

// Put registra a expectativa. Com armazenamento durável, só vale depois de persistida.
func (e *Expectations) Put(x Expectation) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.kv != nil {
		raw, err := json.Marshal(x)
		if err != nil {
			return err
		}
		if err := e.kv.Put(context.Background(), bucketExpectations, hashKey(x.Hash), raw); err != nil {
			return fmt.Errorf("expectativa: persistir: %w", err)
		}
	}
	e.m[string(x.Hash)] = x
	return nil
}

func (e *Expectations) Get(hash []byte) (Expectation, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	x, ok := e.m[string(hash)]
	return x, ok
}

// Verify confere um script de lock recebido contra a expectativa.
func (x Expectation) Verify(s *htlc.Script, linearID string) error {
	if linearID != x.LinearID {
		return fmt.Errorf("expect: asset %q != %q", linearID, x.LinearID)
	}
	if s.HashInfo.HashFunc != crypto.SHA256 || s.HashInfo.HashEncoding != encoding.None {
		return fmt.Errorf("expect: hash func must be SHA256/None")
	}
	diff := s.Deadline.Sub(x.Deadline)
	if diff < 0 {
		diff = -diff
	}
	if diff > x.Tolerance {
		return fmt.Errorf("expect: deadline differs from T1 by %s (> %s)", diff, x.Tolerance)
	}
	return nil
}

type ExpectLockView struct{ *Expectation }

func (v *ExpectLockView) Call(view.Context) (any, error) {
	if len(v.Hash) != 32 || v.LinearID == "" || v.Deadline.IsZero() {
		return nil, fmt.Errorf("expect: incomplete expectation")
	}
	if err := DefaultExpectation.Put(*v.Expectation); err != nil {
		return nil, err
	}
	return true, nil
}

type ExpectLockViewFactory struct{}

func (f *ExpectLockViewFactory) NewView(in []byte) (view.View, error) {
	v := &ExpectLockView{Expectation: &Expectation{}}
	return v, unmarshalInto(in, v.Expectation)
}

// LockAcceptView: responder do Banco para LockView. Só aceita locks esperados (termos assinados).
type LockAcceptView struct{}

func (a *LockAcceptView) Call(ctx view.Context) (any, error) {
	me, sender, err := htlc.RespondExchangeRecipientIdentities(ctx)
	if err != nil {
		return nil, err
	}
	tx, err := htlc.ReceiveTransaction(ctx)
	if err != nil {
		return nil, err
	}
	outs, err := tx.Outputs(ctx.Context())
	if err != nil {
		return nil, err
	}
	scripts := outs.ByScript()
	if scripts.Count() != 1 {
		return nil, fmt.Errorf("lock accept: expected exactly one htlc output, got %d", scripts.Count())
	}
	s := scripts.ScriptAt(0)
	if err := s.Validate(time.Now()); err != nil {
		return nil, fmt.Errorf("lock accept: %w", err)
	}
	if !me.Equal(s.Recipient) || !sender.Equal(s.Sender) {
		return nil, fmt.Errorf("lock accept: unexpected htlc parties")
	}
	r, err := FromTokenType(scripts.At(0).Type)
	if err != nil {
		return nil, fmt.Errorf("lock accept: %w", err)
	}
	x, ok := DefaultExpectation.Get(s.HashInfo.Hash)
	if !ok {
		return nil, fmt.Errorf("lock accept: no signed terms for this hashlock")
	}
	if err := x.Verify(s, r.LinearID); err != nil {
		return nil, err
	}
	if _, err := ctx.RunView(htlc.NewAcceptView(tx)); err != nil {
		return nil, err
	}
	_, err = ctx.RunView(htlc.NewFinalityView(tx))
	return nil, err
}

// =============================================================================== claim (Banco)

type Claim struct {
	Preimage []byte `json:"preimage"`
}

type ClaimView struct{ *Claim }

func (v *ClaimView) Call(ctx view.Context) (any, error) {
	w := htlc.GetWallet(ctx, "")
	if w == nil {
		return nil, fmt.Errorf("claim: wallet not found")
	}
	matched, err := htlc.Wallet(ctx, w).ListByPreImage(ctx.Context(), v.Preimage)
	if err != nil {
		return nil, fmt.Errorf("claim: list by preimage: %w", err)
	}
	if matched.Count() != 1 {
		return nil, fmt.Errorf("claim: expected one htlc for preimage, got %d", matched.Count())
	}
	auditor, err := auditorIdentity(ctx)
	if err != nil {
		return nil, err
	}
	tx, err := htlc.NewAnonymousTransaction(ctx, ttx.WithAuditor(auditor))
	if err != nil {
		return nil, err
	}
	if err := tx.Claim(ctx.Context(), w, matched.At(0), v.Preimage); err != nil {
		return nil, fmt.Errorf("claim: %w", err)
	}
	if _, err := ctx.RunView(htlc.NewCollectEndorsementsView(tx)); err != nil {
		return nil, fmt.Errorf("claim: endorse: %w", err)
	}
	if _, err := ctx.RunView(htlc.NewOrderingAndFinalityView(tx)); err != nil {
		return nil, fmt.Errorf("claim: commit: %w", err)
	}
	return tx.ID(), nil
}

type ClaimViewFactory struct{}

func (f *ClaimViewFactory) NewView(in []byte) (view.View, error) {
	v := &ClaimView{Claim: &Claim{}}
	return v, unmarshalInto(in, v.Claim)
}

// =============================================================================== reclaim (Vendedora)

type Reclaim struct {
	Hash []byte `json:"hash"`
}

type ReclaimView struct{ *Reclaim }

func (v *ReclaimView) Call(ctx view.Context) (any, error) {
	w := htlc.GetWallet(ctx, "")
	if w == nil {
		return nil, fmt.Errorf("reclaim: wallet not found")
	}
	expired, err := htlc.Wallet(ctx, w).GetExpiredByHash(ctx.Context(), v.Hash)
	if err != nil || expired == nil {
		return nil, fmt.Errorf("reclaim: no expired htlc for hash: %v", err)
	}
	auditor, err := auditorIdentity(ctx)
	if err != nil {
		return nil, err
	}
	tx, err := htlc.NewAnonymousTransaction(ctx, ttx.WithAuditor(auditor))
	if err != nil {
		return nil, err
	}
	if err := tx.Reclaim(ctx.Context(), w, expired); err != nil {
		return nil, fmt.Errorf("reclaim: %w", err)
	}
	if _, err := ctx.RunView(htlc.NewCollectEndorsementsView(tx)); err != nil {
		return nil, fmt.Errorf("reclaim: endorse: %w", err)
	}
	if _, err := ctx.RunView(htlc.NewOrderingAndFinalityView(tx)); err != nil {
		return nil, fmt.Errorf("reclaim: commit: %w", err)
	}
	return tx.ID(), nil
}

type ReclaimViewFactory struct{}

func (f *ReclaimViewFactory) NewView(in []byte) (view.View, error) {
	v := &ReclaimView{Reclaim: &Reclaim{}}
	return v, unmarshalInto(in, v.Reclaim)
}

// =============================================================================== status

// Status: estado do lock de um hashlock, visto por este nó (alimenta swap.Observation).
type Status struct {
	Hash []byte `json:"hash"`
}

// LockStatus é a resposta de StatusView.
type LockStatus struct {
	State     string    `json:"state"` // none | locked | claimed | reclaimed
	LinearID  string    `json:"linearId,omitempty"`
	Sender    []byte    `json:"sender,omitempty"`
	Recipient []byte    `json:"recipient,omitempty"`
	Deadline  time.Time `json:"deadline,omitzero"`
	HashOK    bool      `json:"hashOk"`
	Preimage  []byte    `json:"preimage,omitempty"` // publicado no ledger pelo claim
	TxID      string    `json:"txId,omitempty"`     // transação que criou o output HTLC
}

func decodeStatus(raw []byte) (LockStatus, error) {
	var st LockStatus
	return st, json.Unmarshal(raw, &st)
}

type StatusView struct{ *Status }

func (v *StatusView) Call(ctx view.Context) (any, error) {
	st := &LockStatus{State: "none"}
	// 1) O ledger é a fonte da verdade do claim: quem não participou da transação de claim
	//    (a Vendedora) pode ter um token HTLC desatualizado no vault.
	tms, err := token.GetManagementService(ctx)
	if err != nil {
		return nil, err
	}
	nw := network.GetInstance(ctx, tms.Network(), tms.Channel())
	if nw == nil {
		return nil, fmt.Errorf("status: network not found")
	}
	ledger, err := nw.Ledger()
	if err != nil {
		return nil, err
	}
	get := func(k string) ([]byte, error) {
		mk, err := ledger.TransferMetadataKey(k)
		if err != nil {
			return nil, err
		}
		vals, err := ledger.GetStates(ctx.Context(), tms.Namespace(), mk)
		if err != nil || len(vals) == 0 {
			return nil, err
		}
		return vals[0], nil
	}
	if pre, err := get(htlc.ClaimKey(v.Hash)); err != nil {
		return nil, fmt.Errorf("status: claim key: %w", err)
	} else if len(pre) > 0 {
		st.State, st.Preimage = "claimed", pre
		return json.Marshal(st)
	}
	// 2) Ainda existe um token HTLC com este hash (como sender ou recipient)?
	w := htlc.GetWallet(ctx, "")
	if w == nil {
		return nil, fmt.Errorf("status: wallet not found")
	}
	hw := htlc.Wallet(ctx, w)
	// As listagens padrão do Panurus só trazem locks NÃO expirados. Depois de T1 o lock continua
	// no ledger até o reclaim, por isso os expirados também precisam ser consultados. Sem isso, o
	// lock expirado parecia "reclaimed" e a Vendedora encerrava sem reaver o NFT.
	for _, list := range []func() ([]*token2.UnspentToken, error){
		func() ([]*token2.UnspentToken, error) { return collect(hw.ListTokensAsSender(ctx.Context())) },
		func() ([]*token2.UnspentToken, error) { return collect(hw.ListTokensIterator(ctx.Context())) },
		func() ([]*token2.UnspentToken, error) { return unspent(hw.ListExpired(ctx.Context())) },
		func() ([]*token2.UnspentToken, error) { return unspent(hw.ListExpiredReceivedTokens(ctx.Context())) },
	} {
		toks, err := list()
		if err != nil {
			return nil, fmt.Errorf("status: list htlc: %w", err)
		}
		for _, t := range toks {
			s, err := scriptOf(t.Owner)
			if err != nil || string(s.HashInfo.Hash) != string(v.Hash) {
				continue
			}
			st.State, st.Sender, st.Recipient, st.Deadline, st.TxID = "locked", s.Sender, s.Recipient, s.Deadline, t.Id.TxId
			st.HashOK = s.HashInfo.HashFunc == crypto.SHA256 && s.HashInfo.HashEncoding == encoding.None
			if r, err := FromTokenType(t.Type); err == nil {
				st.LinearID = r.LinearID
			}
			return json.Marshal(st)
		}
	}
	if lk, err := get(htlc.LockKey(v.Hash)); err != nil {
		return nil, fmt.Errorf("status: lock key: %w", err)
	} else if len(lk) > 0 {
		st.State = "reclaimed"
	}
	return json.Marshal(st)
}

type StatusViewFactory struct{}

func (f *StatusViewFactory) NewView(in []byte) (view.View, error) {
	v := &StatusView{Status: &Status{}}
	return v, unmarshalInto(in, v.Status)
}

func collect(it iterators.Iterator[*token2.UnspentToken], err error) ([]*token2.UnspentToken, error) {
	if err != nil {
		return nil, err
	}
	return iterators.ReadAllPointers(it)
}

var errNotHTLC = errors.New("not an htlc owner")

func unspent(u *token2.UnspentTokens, err error) ([]*token2.UnspentToken, error) {
	if err != nil || u == nil {
		return nil, err
	}
	return u.Tokens, nil
}

func scriptOf(owner []byte) (*htlc.Script, error) {
	ti, err := identity.UnmarshalTypedIdentity(owner)
	if err != nil {
		return nil, errNotHTLC // identidade sem tipo: não é script
	}
	if ti.Type != htlc.ScriptType {
		return nil, errNotHTLC
	}
	s := &htlc.Script{}
	return s, s.FromBytes(ti.Identity)
}

// =============================================================================== gravame (Registradora)

type Gravame struct {
	LinearID string `json:"linearId"`
	Motivo   string `json:"motivo"`
	Clear    bool   `json:"clear"`
}

type GravameView struct{ *Gravame }

func (v *GravameView) Call(view.Context) (any, error) {
	if v.LinearID == "" {
		return nil, fmt.Errorf("gravame: linearId obrigatório")
	}
	var err error
	if v.Clear {
		err = DefaultGravames.Clear(v.LinearID)
	} else {
		err = DefaultGravames.Set(v.LinearID, v.Motivo)
	}
	if err != nil {
		return nil, err
	}
	return true, nil
}

type GravameViewFactory struct{}

func (f *GravameViewFactory) NewView(in []byte) (view.View, error) {
	v := &GravameView{Gravame: &Gravame{}}
	return v, unmarshalInto(in, v.Gravame)
}
