// Package attest produz e verifica as atestações de lock no Fabric usadas pela Abordagem B2
// (ADR-0003). O digest é bit a bit idêntico a FabricProofVerifier.digest (Solidity):
//
//	sha256(abi.encode(DOMAIN_TAG, chainId, verifier, swapId, hashlock, assetIdHash,
//	                  recipientIdHash, deadline, channelHash, namespaceHash, fabricTxId, blockNumber))
//
// Como todos os campos são estáticos, abi.encode equivale a concatenar palavras de 32 bytes,
// com inteiros e endereços alinhados à direita (big-endian, zeros à esquerda).
package attest

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"slices"

	"golang.org/x/crypto/sha3"
)

// DomainTag = sha256("janus.fabric-lock-attestation.v1"), igual a FabricProofVerifier.DOMAIN_TAG.
var DomainTag = sha256.Sum256([]byte("janus.fabric-lock-attestation.v1"))

// Address é um endereço EVM de 20 bytes.
type Address [20]byte

// FabricLockAttestation espelha a struct Solidity homônima.
type FabricLockAttestation struct {
	SwapID          [32]byte `json:"swapId"`
	Hashlock        [32]byte `json:"hashlock"`
	AssetIDHash     [32]byte `json:"assetIdHash"`
	RecipientIDHash [32]byte `json:"recipientIdHash"`
	Deadline        uint64   `json:"deadline"`
	ChannelHash     [32]byte `json:"channelHash"`
	NamespaceHash   [32]byte `json:"namespaceHash"`
	FabricTxID      [32]byte `json:"fabricTxId"`
	BlockNumber     uint64   `json:"blockNumber"`
}

// Domain vincula a atestação a uma cadeia e a um deployment do verificador.
type Domain struct {
	ChainID  uint64
	Verifier Address
}

// OrgSignature espelha a struct Solidity homônima.
type OrgSignature struct {
	MSPIDHash [32]byte
	KeyIndex  uint32
	R, S      [32]byte
}

// HashString devolve sha256(utf8(s)), usado para assetId, recipientId, canal e namespace.
func HashString(s string) [32]byte { return sha256.Sum256([]byte(s)) }

// HashBytes devolve sha256(b).
func HashBytes(b []byte) [32]byte { return sha256.Sum256(b) }

// MSPIDHash devolve keccak256(bytes(mspId)), igual ao registro no FabricProofVerifier.
func MSPIDHash(mspID string) [32]byte {
	h := sha3.NewLegacyKeccak256()
	h.Write([]byte(mspID))
	var out [32]byte
	copy(out[:], h.Sum(nil))
	return out
}

func word(u uint64) []byte {
	w := make([]byte, 32)
	new(big.Int).SetUint64(u).FillBytes(w)
	return w
}

// Digest calcula o digest assinado pelos atestadores.
func Digest(d Domain, a FabricLockAttestation) [32]byte {
	var buf bytes.Buffer
	buf.Grow(12 * 32)
	buf.Write(DomainTag[:])
	buf.Write(word(d.ChainID))
	addr := make([]byte, 32)
	copy(addr[12:], d.Verifier[:])
	buf.Write(addr)
	buf.Write(a.SwapID[:])
	buf.Write(a.Hashlock[:])
	buf.Write(a.AssetIDHash[:])
	buf.Write(a.RecipientIDHash[:])
	buf.Write(word(a.Deadline))
	buf.Write(a.ChannelHash[:])
	buf.Write(a.NamespaceHash[:])
	buf.Write(a.FabricTxID[:])
	buf.Write(word(a.BlockNumber))
	return sha256.Sum256(buf.Bytes())
}

var (
	p256N     = elliptic.P256().Params().N
	p256HalfN = new(big.Int).Rsh(p256N, 1)
)

// ErrNotP256 indica uma chave que não está na curva P-256.
var ErrNotP256 = errors.New("attest: key is not P-256")

// Signer assina digests de atestação com a chave P-256 do atestador de uma org.
type Signer struct {
	MSPID    string
	KeyIndex uint32
	Key      *ecdsa.PrivateKey
}

// Sign assina o digest e normaliza para low-s, como exigem o OpenZeppelin P256 e o Fabric.
func (s Signer) Sign(d Domain, a FabricLockAttestation) (OrgSignature, error) {
	if s.Key == nil || s.Key.Curve != elliptic.P256() {
		return OrgSignature{}, ErrNotP256
	}
	h := Digest(d, a)
	r, ss, err := ecdsa.Sign(rand.Reader, s.Key, h[:])
	if err != nil {
		return OrgSignature{}, fmt.Errorf("attest: sign: %w", err)
	}
	if ss.Cmp(p256HalfN) > 0 {
		ss = new(big.Int).Sub(p256N, ss)
	}
	sig := OrgSignature{MSPIDHash: MSPIDHash(s.MSPID), KeyIndex: s.KeyIndex}
	r.FillBytes(sig.R[:])
	ss.FillBytes(sig.S[:])
	return sig, nil
}

// PublicKeyXY devolve as coordenadas (qx, qy) da chave pública, para cadastro no verificador.
func PublicKeyXY(pub *ecdsa.PublicKey) (x, y [32]byte, err error) {
	if pub == nil || pub.Curve != elliptic.P256() {
		return x, y, ErrNotP256
	}
	raw, err := pub.Bytes() // 0x04 || X || Y
	if err != nil {
		return x, y, err
	}
	copy(x[:], raw[1:33])
	copy(y[:], raw[33:65])
	return x, y, nil
}

// Verify confere uma assinatura (low-s obrigatório), espelhando o contrato.
func Verify(pub *ecdsa.PublicKey, d Domain, a FabricLockAttestation, sig OrgSignature) bool {
	r := new(big.Int).SetBytes(sig.R[:])
	ss := new(big.Int).SetBytes(sig.S[:])
	if r.Sign() <= 0 || r.Cmp(p256N) >= 0 || ss.Sign() <= 0 || ss.Cmp(p256HalfN) > 0 {
		return false
	}
	h := Digest(d, a)
	return ecdsa.Verify(pub, h[:], r, ss)
}

// SortSignatures ordena por mspIdHash crescente, a ordem exigida por FabricProofVerifier.requireValid.
func SortSignatures(sigs []OrgSignature) {
	slices.SortFunc(sigs, func(a, b OrgSignature) int { return bytes.Compare(a.MSPIDHash[:], b.MSPIDHash[:]) })
}
