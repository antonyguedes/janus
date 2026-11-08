package besu

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"

	"janus/node/internal/attest"
	"janus/node/internal/besu/bindings"
)

// ToBinding converte a atestação para o tipo do binding.
func ToBinding(a attest.FabricLockAttestation) bindings.FabricLockAttestation {
	return bindings.FabricLockAttestation{
		SwapId: a.SwapID, Hashlock: a.Hashlock, AssetIdHash: a.AssetIDHash, RecipientIdHash: a.RecipientIDHash,
		Deadline: a.Deadline, ChannelHash: a.ChannelHash, NamespaceHash: a.NamespaceHash,
		FabricTxId: a.FabricTxID, BlockNumber: a.BlockNumber,
	}
}

// SigsToBinding converte as assinaturas (já ordenadas por mspIdHash).
func SigsToBinding(sigs []attest.OrgSignature) []bindings.OrgSignature {
	out := make([]bindings.OrgSignature, len(sigs))
	for i, s := range sigs {
		out[i] = bindings.OrgSignature{MspIdHash: s.MSPIDHash, KeyIndex: s.KeyIndex, R: s.R, S: s.S}
	}
	return out
}

// AttestationDomain é o domínio das atestações deste deployment.
func (c *Client) AttestationDomain() attest.Domain {
	return attest.Domain{ChainID: c.ChainID(), Verifier: attest.Address(c.d.FabricProofVerifier)}
}

// SettleAttested é Settle com os tipos do pacote attest.
func (c *Client) SettleAttested(ctx context.Context, swapID, preimage [32]byte, a attest.FabricLockAttestation, sigs []attest.OrgSignature) (*types.Receipt, error) {
	return c.Settle(ctx, swapID, preimage, ToBinding(a), SigsToBinding(sigs))
}

// VerifyAttestation chama FabricProofVerifier.isValid (diagnóstico, sem gas).
func (c *Client) VerifyAttestation(ctx context.Context, a attest.FabricLockAttestation, sigs []attest.OrgSignature) (bool, error) {
	return call(ctx, c.verifierC, c.verifier.PackIsValid(ToBinding(a), SigsToBinding(sigs)), c.verifier.UnpackIsValid)
}

// Attestor é uma org atestadora a cadastrar no verificador.
type Attestor struct {
	MSPID    string
	Qx, Qy   [32]byte
	CertHash [32]byte // sha256 do certificado X.509 (DER) validado fora da cadeia
}

// RegisterAttestors cadastra orgs e chaves e define o threshold. Exige GOVERNOR_ROLE (admin).
// Idempotente por org: uma org já cadastrada só recebe a chave.
func (c *Client) RegisterAttestors(ctx context.Context, threshold uint32, atts ...Attestor) error {
	for _, a := range atts {
		h := attest.MSPIDHash(a.MSPID)
		org, err := call(ctx, c.verifierC, c.verifier.PackGetOrg(h), c.verifier.UnpackGetOrg)
		if err != nil {
			return err
		}
		if !org.Registered {
			if _, err := c.send(ctx, c.verifierC, c.verifier.PackRegisterOrg(a.MSPID)); err != nil {
				return fmt.Errorf("registerOrg %s: %w", a.MSPID, err)
			}
		}
		if _, err := c.send(ctx, c.verifierC, c.verifier.PackAddAttestorKey(h, a.Qx, a.Qy, a.CertHash)); err != nil {
			return fmt.Errorf("addAttestorKey %s: %w", a.MSPID, err)
		}
	}
	if _, err := c.send(ctx, c.verifierC, c.verifier.PackSetThreshold(threshold)); err != nil {
		return fmt.Errorf("setThreshold: %w", err)
	}
	return nil
}
