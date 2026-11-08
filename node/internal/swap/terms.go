// Package swap contém o núcleo puro do protocolo: termos canônicos, política de tempo, a visão
// observada das duas cadeias e a função de decisão NextAction. Nada aqui faz I/O, então tudo é
// testável sem rede. As implementações de SettlementProtocol (protocol/htlc, protocol/cacti) só
// observam as cadeias e executam as ações decididas aqui.
package swap

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

// Protocol seleciona a abordagem de liquidação.
type Protocol string

const (
	ProtocolHTLC  Protocol = "htlc"  // Abordagem A
	ProtocolCacti Protocol = "cacti" // Abordagem B2 (HTLC atestado, ADR-0003)
)

// Address é um endereço EVM.
type Address [20]byte

// Party identifica uma instituição nas duas redes.
type Party struct {
	MSPID       string  `json:"mspId"`
	FSCIdentity []byte  `json:"fscIdentity"` // identidade de owner no Panurus (bytes serializados)
	EVM         Address `json:"evm"`
}

// Amount é o preço em unidades mínimas do token (centavos de Drex).
type Amount struct {
	Token Address `json:"token"`
	Value uint64  `json:"value"`
}

// NetworkRefs fixa os deployments envolvidos, para que os termos não possam ser reaproveitados
// em outra rede.
type NetworkRefs struct {
	FabricChannel   string  `json:"fabricChannel"`
	FabricNamespace string  `json:"fabricNamespace"`
	BesuChainID     uint64  `json:"besuChainId"`
	PaymentContract Address `json:"paymentContract"` // HTLCPayment (A) ou AttestedSettlement (B)
}

// SwapTerms são os termos canônicos trocados por sessão P2P do FSC e assinados pelos dois lados.
type SwapTerms struct {
	SwapID   [32]byte    `json:"swapId"`
	Protocol Protocol    `json:"protocol"`
	AssetID  string      `json:"assetId"`
	Seller   Party       `json:"seller"`
	Buyer    Party       `json:"buyer"`
	Price    Amount      `json:"price"`
	Hashlock [32]byte    `json:"hashlock"` // H = SHA256(S), usado em A e em B2
	Nonce    [32]byte    `json:"nonce"`
	T1       int64       `json:"t1"` // unix s: deadline do lock do NFT no Fabric
	T2       int64       `json:"t2"` // unix s: timelock do pagamento no Besu
	Networks NetworkRefs `json:"networks"`
}

// ComputeSwapID = sha256(hashlock || len||seller.fscIdentity || len||buyer.fscIdentity || nonce).
func ComputeSwapID(hashlock [32]byte, seller, buyer Party, nonce [32]byte) [32]byte {
	var b bytes.Buffer
	b.Write(hashlock[:])
	writeBytes(&b, seller.FSCIdentity)
	writeBytes(&b, buyer.FSCIdentity)
	b.Write(nonce[:])
	return sha256.Sum256(b.Bytes())
}

var termsDomain = []byte("janus.swap-terms.v1")

// CanonicalBytes é a codificação determinística (campos de tamanho fixo ou com prefixo de tamanho,
// em ordem fixa) assinada pelas duas partes. Não depende de JSON nem de ordem de mapas.
func (t *SwapTerms) CanonicalBytes() []byte {
	var b bytes.Buffer
	writeBytes(&b, termsDomain)
	b.Write(t.SwapID[:])
	writeBytes(&b, []byte(t.Protocol))
	writeBytes(&b, []byte(t.AssetID))
	for _, p := range []Party{t.Seller, t.Buyer} {
		writeBytes(&b, []byte(p.MSPID))
		writeBytes(&b, p.FSCIdentity)
		b.Write(p.EVM[:])
	}
	b.Write(t.Price.Token[:])
	writeU64(&b, t.Price.Value)
	b.Write(t.Hashlock[:])
	b.Write(t.Nonce[:])
	writeU64(&b, uint64(t.T1))
	writeU64(&b, uint64(t.T2))
	writeBytes(&b, []byte(t.Networks.FabricChannel))
	writeBytes(&b, []byte(t.Networks.FabricNamespace))
	writeU64(&b, t.Networks.BesuChainID)
	b.Write(t.Networks.PaymentContract[:])
	return b.Bytes()
}

// Hash é o hash canônico dos termos (o que as partes assinam).
func (t *SwapTerms) Hash() [32]byte { return sha256.Sum256(t.CanonicalBytes()) }

// T1Time e T2Time convertem os prazos para time.Time.
func (t *SwapTerms) T1Time() time.Time { return time.Unix(t.T1, 0) }
func (t *SwapTerms) T2Time() time.Time { return time.Unix(t.T2, 0) }

// Erros de validação dos termos.
var (
	ErrBadSwapID   = errors.New("terms: swapId does not match sha256(hashlock||seller||buyer||nonce)")
	ErrBadTimes    = errors.New("terms: invalid timelocks")
	ErrBadParties  = errors.New("terms: invalid parties")
	ErrBadPrice    = errors.New("terms: invalid price")
	ErrBadProtocol = errors.New("terms: unknown protocol")
)

// Validate confere os termos contra a política de tempo no instante `now`, que é quando cada
// parte decide se aceita. Os dois lados chamam Validate antes de assinar.
func (t *SwapTerms) Validate(p TimePolicy, now time.Time) error {
	switch t.Protocol {
	case ProtocolHTLC, ProtocolCacti:
	default:
		return ErrBadProtocol
	}
	if t.AssetID == "" || t.Hashlock == ([32]byte{}) {
		return fmt.Errorf("terms: assetId and hashlock are required")
	}
	if len(t.Seller.FSCIdentity) == 0 || len(t.Buyer.FSCIdentity) == 0 ||
		bytes.Equal(t.Seller.FSCIdentity, t.Buyer.FSCIdentity) ||
		t.Seller.EVM == (Address{}) || t.Buyer.EVM == (Address{}) || t.Seller.EVM == t.Buyer.EVM {
		return ErrBadParties
	}
	if t.Price.Value == 0 || t.Price.Token == (Address{}) {
		return ErrBadPrice
	}
	if ComputeSwapID(t.Hashlock, t.Seller, t.Buyer, t.Nonce) != t.SwapID {
		return ErrBadSwapID
	}
	if t.T1 <= t.T2 {
		return fmt.Errorf("%w: T1 must be after T2", ErrBadTimes)
	}
	if gap := time.Duration(t.T1-t.T2) * time.Second; gap < p.Delta() {
		return fmt.Errorf("%w: T1-T2=%s < Δ=%s", ErrBadTimes, gap, p.Delta())
	}
	if left := t.T2Time().Sub(now); left < p.BuyerLockMargin+p.RevealMargin {
		return fmt.Errorf("%w: only %s left until T2", ErrBadTimes, left)
	}
	if t.Networks.FabricChannel == "" || t.Networks.FabricNamespace == "" || t.Networks.BesuChainID == 0 ||
		t.Networks.PaymentContract == (Address{}) {
		return fmt.Errorf("terms: incomplete network refs")
	}
	return nil
}

func writeBytes(b *bytes.Buffer, v []byte) {
	writeU64(b, uint64(len(v)))
	b.Write(v)
}

func writeU64(b *bytes.Buffer, v uint64) {
	var x [8]byte
	binary.BigEndian.PutUint64(x[:], v)
	b.Write(x[:])
}
