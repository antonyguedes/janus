package token

import (
	"bytes"
	"crypto"
	"fmt"
	"sync"
	"time"

	"github.com/LFDT-Panurus/panurus/token/services/interop/encoding"
	"github.com/LFDT-Panurus/panurus/token/services/interop/htlc"
)

// AuditPolicy é a política da Registradora (auditora) para transações HTLC.
// É a camada 2 da mitigação de tempo (ADR-0004): o auditor confere tudo com o PRÓPRIO relógio,
// independentemente do relógio dos endossantes.
type AuditPolicy struct {
	// ReclaimGrace: o auditor só aprova um reclaim se now >= deadline + ReclaimGrace.
	// É conservador a favor do recipient (Banco).
	ReclaimGrace time.Duration
	// MinLock/MaxLock: limites de sanidade para o deadline de um novo lock.
	MinLock, MaxLock time.Duration
}

// DefaultAuditPolicy usa a mesma tolerância de skew da TimePolicy (30 s).
func DefaultAuditPolicy() AuditPolicy {
	return AuditPolicy{ReclaimGrace: 30 * time.Second, MinLock: time.Minute, MaxLock: 7 * 24 * time.Hour}
}

// HTLCOp é a operação inferida de uma entrada HTLC.
type HTLCOp int

const (
	OpClaim HTLCOp = iota + 1
	OpReclaim
)

// CheckLockScript valida o script de um NOVO lock (saída HTLC).
func (p AuditPolicy) CheckLockScript(s *htlc.Script, now time.Time) error {
	if s.HashInfo.HashFunc != crypto.SHA256 || s.HashInfo.HashEncoding != encoding.None {
		return fmt.Errorf("audit: htlc must use SHA256/None (ADR-0002)")
	}
	if len(s.HashInfo.Hash) != 32 {
		return fmt.Errorf("audit: hashlock must be 32 bytes")
	}
	if s.Sender.IsNone() || s.Recipient.IsNone() || s.Sender.Equal(s.Recipient) {
		return fmt.Errorf("audit: invalid htlc parties")
	}
	d := s.Deadline.Sub(now)
	if d < p.MinLock || d > p.MaxLock {
		return fmt.Errorf("audit: lock deadline %s out of bounds [%s, %s]", d, p.MinLock, p.MaxLock)
	}
	return nil
}

// ClassifySpend infere se o gasto de uma entrada HTLC é claim ou reclaim pelo dono da saída.
func ClassifySpend(s *htlc.Script, outOwner []byte) (HTLCOp, error) {
	switch {
	case bytes.Equal(outOwner, s.Recipient):
		return OpClaim, nil
	case bytes.Equal(outOwner, s.Sender):
		return OpReclaim, nil
	}
	return 0, fmt.Errorf("audit: htlc spent to neither sender nor recipient")
}

// CheckSpend aplica as regras de tempo do auditor ao gasto de uma entrada HTLC.
func (p AuditPolicy) CheckSpend(s *htlc.Script, op HTLCOp, now time.Time) error {
	switch op {
	case OpClaim:
		if !now.Before(s.Deadline) {
			return fmt.Errorf("audit: claim after deadline (%s >= %s)", now.UTC(), s.Deadline.UTC())
		}
	case OpReclaim:
		if earliest := s.Deadline.Add(p.ReclaimGrace); now.Before(earliest) {
			return fmt.Errorf("audit: early reclaim refused (now %s < deadline+grace %s)", now.UTC(), earliest.UTC())
		}
	default:
		return fmt.Errorf("audit: unknown op")
	}
	return nil
}

// Gravames é o registro de ônus da Registradora: um recebível com gravame não pode ser
// transferido nem travado. Em memória no MVP. Em produção, a base da Registradora.
type Gravames struct {
	mu  sync.RWMutex
	set map[string]string // LinearID -> motivo
}

func NewGravames() *Gravames { return &Gravames{set: map[string]string{}} }

func (g *Gravames) Set(linearID, motivo string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.set[linearID] = motivo
}

func (g *Gravames) Clear(linearID string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.set, linearID)
}

// Check devolve erro se o recebível tiver gravame.
func (g *Gravames) Check(linearID string) error {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if m, ok := g.set[linearID]; ok {
		return fmt.Errorf("audit: recebível %s com gravame (%s)", linearID, m)
	}
	return nil
}
