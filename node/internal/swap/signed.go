package swap

import (
	"errors"
	"fmt"
)

// Signer assina o hash canônico dos termos com a identidade FSC/X.509 da parte.
// Em produção é implementado sobre o SigService do FSC. Nos testes, com ECDSA em memória.
type Signer interface {
	Sign(digest []byte) ([]byte, error)
}

// Verifier confere a assinatura de uma identidade FSC sobre um digest.
type Verifier interface {
	Verify(identity, digest, signature []byte) error
}

// SignedTerms são os termos com as assinaturas das duas partes. Nenhum lock acontece antes de
// SignedTerms.Verify passar nos dois nós.
type SignedTerms struct {
	Terms     SwapTerms `json:"terms"`
	SellerSig []byte    `json:"sellerSig"`
	BuyerSig  []byte    `json:"buyerSig"`
}

var ErrMissingSignature = errors.New("terms: missing signature")

// SignAs adiciona a assinatura do papel indicado.
func (s *SignedTerms) SignAs(role Role, signer Signer) error {
	h := s.Terms.Hash()
	sig, err := signer.Sign(h[:])
	if err != nil {
		return fmt.Errorf("terms: sign as %s: %w", role, err)
	}
	switch role {
	case RoleSeller:
		s.SellerSig = sig
	case RoleBuyer:
		s.BuyerSig = sig
	default:
		return fmt.Errorf("terms: role %s cannot sign", role)
	}
	return nil
}

// Verify confere as duas assinaturas contra as identidades FSC presentes nos próprios termos.
func (s *SignedTerms) Verify(v Verifier) error {
	if len(s.SellerSig) == 0 || len(s.BuyerSig) == 0 {
		return ErrMissingSignature
	}
	h := s.Terms.Hash()
	if err := v.Verify(s.Terms.Seller.FSCIdentity, h[:], s.SellerSig); err != nil {
		return fmt.Errorf("terms: seller signature: %w", err)
	}
	if err := v.Verify(s.Terms.Buyer.FSCIdentity, h[:], s.BuyerSig); err != nil {
		return fmt.Errorf("terms: buyer signature: %w", err)
	}
	return nil
}
