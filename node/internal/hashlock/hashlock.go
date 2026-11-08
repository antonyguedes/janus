// Package hashlock implementa o par (S, H) da Abordagem A, de forma interoperável:
//   - S: exatamente 32 bytes aleatórios (crypto/rand);
//   - H = SHA256(S), em bytes crus.
//
// Isso equivale a `sha256(abi.encodePacked(bytes32 S))` no Solidity e ao HTLC do Panurus configurado
// com HashFunc=crypto.SHA256 e HashEncoding=encoding.None (ADR-0002). Os vetores de teste compartilhados
// ficam em testdata/vectors/htlc.json.
package hashlock

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// Size é o tamanho fixo da preimagem e do hashlock.
const Size = 32

// Preimage é o segredo S.
type Preimage [Size]byte

// Hash é o hashlock H = SHA256(S).
type Hash [Size]byte

// ErrMismatch indica que a preimagem não corresponde ao hashlock.
var ErrMismatch = errors.New("hashlock: preimage does not match hashlock")

// New gera uma preimagem aleatória e o hashlock dela.
func New() (Preimage, Hash, error) {
	var s Preimage
	if _, err := rand.Read(s[:]); err != nil {
		return Preimage{}, Hash{}, fmt.Errorf("hashlock: random: %w", err)
	}
	return s, s.Hash(), nil
}

// Hash calcula H = SHA256(S).
func (s Preimage) Hash() Hash { return sha256.Sum256(s[:]) }

// Verify confere, em tempo constante, que SHA256(s) == h.
func Verify(s Preimage, h Hash) error {
	got := s.Hash()
	if subtle.ConstantTimeCompare(got[:], h[:]) != 1 {
		return ErrMismatch
	}
	return nil
}

func (s Preimage) Hex() string { return "0x" + hex.EncodeToString(s[:]) }
func (h Hash) Hex() string     { return "0x" + hex.EncodeToString(h[:]) }
func (h Hash) IsZero() bool    { return h == Hash{} }

// ParseHash32 decodifica um valor hexadecimal de 32 bytes, com ou sem prefixo 0x.
func ParseHash32(s string) ([Size]byte, error) {
	var out [Size]byte
	b, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	if err != nil {
		return out, fmt.Errorf("hashlock: hex: %w", err)
	}
	if len(b) != Size {
		return out, fmt.Errorf("hashlock: expected %d bytes, got %d", Size, len(b))
	}
	copy(out[:], b)
	return out, nil
}
