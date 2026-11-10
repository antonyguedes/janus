package token

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/LFDT-Panurus/panurus/token"
	"github.com/LFDT-Panurus/panurus/token/services/interop/htlc"
	"github.com/LFDT-Panurus/panurus/token/services/network"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/comm/session"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/id"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/view"

	"janus/node/internal/attest"
	"janus/node/internal/cacti"
)

// Abordagem B2 (ADR-0003): atestação do lock do NFT no Fabric para o AttestedSettlement no Besu.
//
// - Atestadores (Registradora e Banco) assinam uma FabricLockAttestation com a chave P-256 de
//   atestador da org (o certificado X.509 é emitido pela CA da org e validado no cadastro on-chain).
//   - A Registradora usa a evidência que ela mesma auditou (LockRegistry, preenchido pelo
//     AuditView) e confere no ledger que a transação foi comitada como válida.
//   - O Banco usa o próprio vault (ele é o recipient do script HTLC).
// - O relay (a Vendedora, ou um componente de transporte como o conector Fabric do Cacti) coleta
//   as assinaturas por sessão FSC e as submete ao Besu.

// LockRecord é um lock HTLC visto pela auditoria.
type LockRecord struct {
	Hash      []byte
	LinearID  string
	Sender    []byte
	Recipient []byte
	Deadline  time.Time
	TxID      string
}

// LockRegistry guarda os locks auditados (memória do processo, MVP).
type LockRegistry struct {
	mu sync.RWMutex
	m  map[string]LockRecord
	kv KV // durável (opcional)
}

func NewLockRegistry() *LockRegistry { return &LockRegistry{m: map[string]LockRecord{}} }

// Put registra o lock. Com armazenamento durável, só vale depois de persistido.
func (r *LockRegistry) Put(x LockRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.kv != nil {
		raw, err := json.Marshal(x)
		if err != nil {
			return err
		}
		if err := r.kv.Put(context.Background(), bucketLocks, hashKey(x.Hash), raw); err != nil {
			return fmt.Errorf("lock auditado: persistir: %w", err)
		}
	}
	r.m[string(x.Hash)] = x
	return nil
}

func (r *LockRegistry) Get(hash []byte) (LockRecord, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	x, ok := r.m[string(hash)]
	return x, ok
}

// DefaultLocks é preenchido pelo AuditView da Registradora.
var DefaultLocks = NewLockRegistry()

// --------------------------------------------------------------------- hooks do atestador

// LedgerProver obtém a prova de uma transação no ledger comitado (conector Fabric do Cacti).
type LedgerProver interface {
	TxProof(ctx context.Context, txID string) (cacti.TxProof, error)
}

// AttestorHooks configuram o papel de atestador (Registradora e Banco).
type AttestorHooks struct {
	Signer    attest.Signer
	Channel   string
	Namespace string
	// Ledger (opcional, mas obrigatório quando configurado): o conector Fabric do Cacti da
	// própria org. O atestador só assina se a prova do ledger confirmar a transação de lock
	// válida, do chaincode certo e endossada pelas orgs exigidas, e o número do bloco entra na
	// atestação (ADR-0005).
	Ledger            LedgerProver
	RequiredEndorsers []string
}

var (
	attMu     sync.RWMutex
	attestorH *AttestorHooks
)

// SetAttestor habilita este nó como atestador.
func SetAttestor(h AttestorHooks) {
	attMu.Lock()
	defer attMu.Unlock()
	attestorH = &h
}

func currentAttestor() (*AttestorHooks, error) {
	attMu.RLock()
	defer attMu.RUnlock()
	if attestorH == nil || attestorH.Signer.Key == nil {
		return nil, errors.New("attest: this node is not an attestor")
	}
	return attestorH, nil
}

// --------------------------------------------------------------------- mensagens

// AttestRequest é o pedido de atestação enviado pelo relay.
type AttestRequest struct {
	SwapID   [32]byte       `json:"swapId"`
	Hash     [32]byte       `json:"hash"`
	ChainID  uint64         `json:"chainId"`
	Verifier attest.Address `json:"verifier"`
}

// AttestResponse traz a atestação construída pela org e a assinatura dela.
type AttestResponse struct {
	Attestation attest.FabricLockAttestation `json:"attestation"`
	Signature   attest.OrgSignature          `json:"signature"`
}

// Collected é o resultado da coleta, pronto para AttestedSettlement.settle.
type Collected struct {
	Attestation attest.FabricLockAttestation `json:"attestation"`
	Signatures  []attest.OrgSignature        `json:"signatures"`
}

// --------------------------------------------------------------------- relay (coleta)

// Collect: o relay pede atestações aos nós atestadores e confere que todas atestam os mesmos fatos.
type Collect struct {
	Request   AttestRequest `json:"request"`
	Attestors []string      `json:"attestors"` // ids FSC dos nós atestadores
}

type CollectView struct{ *Collect }

func (v *CollectView) Call(ctx view.Context) (any, error) {
	if len(v.Attestors) == 0 {
		return nil, errors.New("collect: no attestors")
	}
	idp, err := id.GetProvider(ctx)
	if err != nil {
		return nil, err
	}
	var out Collected
	for i, node := range v.Attestors {
		s, err := session.NewJSON(ctx, v, idp.Identity(node))
		if err != nil {
			return nil, fmt.Errorf("collect: session %s: %w", node, err)
		}
		if err := s.Send(ctx.Context(), &v.Request); err != nil {
			return nil, err
		}
		var res AttestResponse
		if err := s.ReceiveWithTimeout(&res, 30*time.Second); err != nil {
			return nil, fmt.Errorf("collect: %s: %w", node, err)
		}
		if i == 0 {
			out.Attestation = res.Attestation
		} else if res.Attestation != out.Attestation {
			return nil, fmt.Errorf("collect: attestor %s disagrees on the lock facts", node)
		}
		out.Signatures = append(out.Signatures, res.Signature)
	}
	attest.SortSignatures(out.Signatures)
	return &out, nil
}

type CollectViewFactory struct{}

func (f *CollectViewFactory) NewView(in []byte) (view.View, error) {
	v := &CollectView{Collect: &Collect{}}
	return v, unmarshalInto(in, v.Collect)
}

// --------------------------------------------------------------------- atestador (responder)

// AttestResponderView responde ao CollectView. Só atesta um lock comitado e consistente.
type AttestResponderView struct{}

func (a *AttestResponderView) Call(ctx view.Context) (any, error) {
	s := session.JSON(ctx)
	var req AttestRequest
	if err := s.ReceiveWithTimeout(&req, 30*time.Second); err != nil {
		return nil, err
	}
	res, err := buildAttestation(ctx, req)
	if err != nil {
		_ = s.SendError(ctx.Context(), err.Error())
		return nil, err
	}
	return nil, s.Send(ctx.Context(), res)
}

func buildAttestation(ctx view.Context, req AttestRequest) (*AttestResponse, error) {
	h, err := currentAttestor()
	if err != nil {
		return nil, err
	}
	rec, err := lockEvidence(ctx, req.Hash)
	if err != nil {
		return nil, err
	}
	if err := ledgerValid(ctx, rec.TxID); err != nil {
		return nil, err
	}
	txid, err := parseTxID(rec.TxID)
	if err != nil {
		return nil, err
	}
	var block uint64
	if h.Ledger != nil {
		proof, err := h.Ledger.TxProof(ctx.Context(), rec.TxID)
		if err != nil {
			return nil, fmt.Errorf("attest: prova do ledger via Cacti: %w", err)
		}
		if err := proof.Satisfies(h.Namespace, h.RequiredEndorsers); err != nil {
			return nil, err
		}
		block = proof.BlockNumber
	}
	a := attest.FabricLockAttestation{
		SwapID:          req.SwapID,
		Hashlock:        req.Hash,
		AssetIDHash:     attest.HashString(rec.LinearID),
		RecipientIDHash: attest.HashBytes(rec.Recipient),
		Deadline:        uint64(rec.Deadline.Unix()),
		ChannelHash:     attest.HashString(h.Channel),
		NamespaceHash:   attest.HashString(h.Namespace),
		FabricTxID:      txid,
		BlockNumber:     block,
	}
	sig, err := h.Signer.Sign(attest.Domain{ChainID: req.ChainID, Verifier: req.Verifier}, a)
	if err != nil {
		return nil, err
	}
	return &AttestResponse{Attestation: a, Signature: sig}, nil
}

// lockEvidence: a auditoria (Registradora) ou o vault (Banco, recipient do script).
func lockEvidence(ctx view.Context, hash [32]byte) (LockRecord, error) {
	if rec, ok := DefaultLocks.Get(hash[:]); ok {
		return rec, nil
	}
	if _, err := currentAttestor(); err == nil && htlcWalletMissing(ctx) {
		return LockRecord{}, errors.New("attest: lock not seen by this node's audit")
	}
	res, err := (&StatusView{Status: &Status{Hash: hash[:]}}).Call(ctx)
	if err != nil {
		return LockRecord{}, err
	}
	raw, _ := res.([]byte)
	st, err := decodeStatus(raw)
	if err != nil {
		return LockRecord{}, err
	}
	if st.State != "locked" || st.TxID == "" || !st.HashOK {
		return LockRecord{}, fmt.Errorf("attest: no committed htlc lock for hash (state=%s)", st.State)
	}
	return LockRecord{Hash: hash[:], LinearID: st.LinearID, Sender: st.Sender, Recipient: st.Recipient,
		Deadline: st.Deadline, TxID: st.TxID}, nil
}

func htlcWalletMissing(ctx view.Context) bool {
	return htlc.GetWallet(ctx, "") == nil
}

func ledgerValid(ctx view.Context, txID string) error {
	tms, err := token.GetManagementService(ctx)
	if err != nil {
		return err
	}
	nw := network.GetInstance(ctx, tms.Network(), tms.Channel())
	if nw == nil {
		return errors.New("attest: network not found")
	}
	l, err := nw.Ledger()
	if err != nil {
		return err
	}
	code, msg, err := l.Status(ctx.Context(), txID)
	if err != nil {
		return fmt.Errorf("attest: ledger status: %w", err)
	}
	if code != network.Valid {
		return fmt.Errorf("attest: lock tx %s not valid on ledger (%d %s)", txID, code, msg)
	}
	return nil
}

// parseTxID converte o ID de transação do Fabric (hex de 32 bytes) em bytes32.
func parseTxID(s string) ([32]byte, error) {
	var out [32]byte
	b, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	if err != nil || len(b) != 32 {
		return sha256.Sum256([]byte(s)), nil // IDs fora do formato padrão: usa o hash
	}
	copy(out[:], b)
	return out, nil
}

// AttestLocal atesta no próprio processo (o Banco atesta a si mesmo sem sessão de rede).
func AttestLocal(ctx context.Context, vi ViewInitiator, req AttestRequest) (*AttestResponse, error) {
	res, err := vi.InitiateView(ctx, &localAttestView{req: req})
	if err != nil {
		return nil, err
	}
	r, _ := res.(*AttestResponse)
	return r, nil
}

type localAttestView struct{ req AttestRequest }

func (v *localAttestView) Call(ctx view.Context) (any, error) { return buildAttestation(ctx, v.req) }
