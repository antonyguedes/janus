package swap

import (
	"bytes"
	"fmt"
	"time"
)

// Role é o papel do nó no swap.
type Role string

const (
	RoleSeller Role = "vendedora"
	RoleBuyer  Role = "banco"
)

// FabricState é o estado do NFT do recebível em relação ao swap.
type FabricState int

const (
	FabricNone      FabricState = iota // ainda não travado para este swap
	FabricLocked                       // travado no script HTLC do Panurus
	FabricClaimed                      // o Banco fez o claim com S (NFT é do Banco)
	FabricReclaimed                    // a Vendedora fez o reclaim depois de T1
)

func (s FabricState) String() string {
	return [...]string{"none", "locked", "claimed", "reclaimed"}[s]
}

// BesuState é o estado do pagamento no contrato do Besu.
type BesuState int

const (
	BesuNone     BesuState = iota
	BesuLocked             // HTLCPayment.Locked ou AttestedSettlement.Deposited
	BesuClaimed            // pago à Vendedora; S publicado
	BesuRefunded           // devolvido ao Banco depois de T2
)

func (s BesuState) String() string {
	return [...]string{"none", "locked", "claimed", "refunded"}[s]
}

// FabricView é o que o nó lê do vault/ledger do Fabric sobre o lock do NFT.
type FabricView struct {
	State     FabricState
	AssetID   string
	Sender    []byte    // script.Sender
	Recipient []byte    // script.Recipient
	Hashlock  [32]byte  // script.HashInfo.Hash
	HashOK    bool      // HashFunc == SHA256 e HashEncoding == None (ADR-0002)
	Deadline  time.Time // script.Deadline
	TxID      string
	Block     uint64
}

// BesuView é o que o nó lê do contrato de pagamento.
type BesuView struct {
	State     BesuState
	Contract  Address
	Token     Address
	Sender    Address
	Recipient Address
	Amount    uint64
	Hashlock  [32]byte
	Timelock  time.Time
	// MinFabricDeadline existe só na Abordagem B (AttestedSettlement.minFabricDeadline).
	MinFabricDeadline time.Time
	Preimage          *[32]byte // publicado em PaymentClaimed / PaymentSettled
	TxHash            string
}

// Observation é o estado combinado das duas cadeias mais o conhecimento local do nó.
type Observation struct {
	Now    time.Time
	Fabric FabricView
	Besu   BesuView
	// Preimage conhecida localmente: o segredo da Vendedora, ou a preimagem que o Banco aprendeu.
	Preimage *[32]byte
	// AttestationReady (só na B): atestações suficientes coletadas para o settle.
	AttestationReady bool
}

// CheckFabricLock confere o lock do NFT contra os termos. Quem chama é o Banco, antes de travar
// o pagamento (A) ou de atestar (B), e a Vendedora, para confirmar o próprio lock.
func CheckFabricLock(t *SwapTerms, p TimePolicy, f FabricView) error {
	if f.State != FabricLocked {
		return fmt.Errorf("fabric: asset not locked (%s)", f.State)
	}
	if f.AssetID != t.AssetID {
		return fmt.Errorf("fabric: asset %q != %q", f.AssetID, t.AssetID)
	}
	if !bytes.Equal(f.Recipient, t.Buyer.FSCIdentity) {
		return fmt.Errorf("fabric: htlc recipient is not the buyer")
	}
	if !bytes.Equal(f.Sender, t.Seller.FSCIdentity) {
		return fmt.Errorf("fabric: htlc sender is not the seller")
	}
	if !f.HashOK || f.Hashlock != t.Hashlock {
		return fmt.Errorf("fabric: hashlock/hash function mismatch")
	}
	if d := absDur(f.Deadline.Sub(t.T1Time())); d > p.DeadlineTolerance {
		return fmt.Errorf("fabric: deadline differs from T1 by %s (> %s)", d, p.DeadlineTolerance)
	}
	return nil
}

// CheckBesuLock confere o lock do pagamento contra os termos. Quem chama é a Vendedora, antes
// de revelar S.
func CheckBesuLock(t *SwapTerms, p TimePolicy, b BesuView) error {
	if b.State != BesuLocked {
		return fmt.Errorf("besu: payment not locked (%s)", b.State)
	}
	switch {
	case b.Contract != t.Networks.PaymentContract:
		return fmt.Errorf("besu: unexpected contract")
	case b.Token != t.Price.Token:
		return fmt.Errorf("besu: unexpected token")
	case b.Amount != t.Price.Value:
		return fmt.Errorf("besu: amount %d != %d", b.Amount, t.Price.Value)
	case b.Recipient != t.Seller.EVM:
		return fmt.Errorf("besu: recipient is not the seller")
	case b.Sender != t.Buyer.EVM:
		return fmt.Errorf("besu: sender is not the buyer")
	case b.Hashlock != t.Hashlock:
		return fmt.Errorf("besu: hashlock mismatch")
	}
	// O timelock do Besu não pode passar de T2: um T2 maior encurta a margem Δ da Vendedora.
	if b.Timelock.After(t.T2Time()) {
		return fmt.Errorf("besu: timelock %s after agreed T2", b.Timelock.UTC())
	}
	if gap := t.T1Time().Sub(b.Timelock); gap < p.Delta() {
		return fmt.Errorf("besu: T1 - timelock = %s < Δ", gap)
	}
	if t.Protocol == ProtocolCacti {
		if b.MinFabricDeadline.After(t.T1Time()) {
			return fmt.Errorf("besu: minFabricDeadline after T1 (settle would be impossible)")
		}
	}
	return nil
}

// Status é o resumo do swap, usado pela UI e pela API.
type Status string

const (
	StatusNegotiated    Status = "negotiated"
	StatusAssetLocked   Status = "asset_locked"
	StatusPaymentLocked Status = "payment_locked"
	StatusBothLocked    Status = "both_locked"
	StatusRevealed      Status = "revealed" // pago no Besu, NFT ainda travado
	StatusCompleted     Status = "completed"
	StatusReverted      Status = "reverted"
	// Estados anômalos: uma das partes ficou com as duas pernas.
	StatusSellerHasBoth Status = "anomaly_seller_has_both"
	StatusBuyerHasBoth  Status = "anomaly_buyer_has_both"
)

// Derive resume a observação. É uma função pura do estado das cadeias.
func Derive(o Observation) Status {
	f, b := o.Fabric.State, o.Besu.State
	switch {
	case f == FabricClaimed && b == BesuClaimed:
		return StatusCompleted
	case f == FabricReclaimed && b == BesuClaimed:
		return StatusSellerHasBoth
	case f == FabricClaimed && b != BesuClaimed:
		return StatusBuyerHasBoth
	case b == BesuClaimed:
		return StatusRevealed
	case (f == FabricReclaimed || f == FabricNone) && (b == BesuRefunded || b == BesuNone):
		if f == FabricNone && b == BesuNone {
			return StatusNegotiated
		}
		return StatusReverted
	case f == FabricLocked && b == BesuLocked:
		return StatusBothLocked
	case f == FabricLocked:
		return StatusAssetLocked
	case b == BesuLocked:
		return StatusPaymentLocked
	}
	return StatusNegotiated
}

// Terminal informa se o swap não tem mais transições possíveis para as partes.
func (s Status) Terminal() bool {
	switch s {
	case StatusCompleted, StatusReverted, StatusSellerHasBoth, StatusBuyerHasBoth:
		return true
	}
	return false
}

func absDur(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
