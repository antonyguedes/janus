package besu

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
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

// SettledAttestation devolve a atestação usada no settle de um swap, decodificando o calldata da
// transação que emitiu PaymentSettled (auditoria/E2E: confere o que a cadeia de fato aceitou).
func (c *Client) SettledAttestation(ctx context.Context, swapID [32]byte) (attest.FabricLockAttestation, error) {
	var out attest.FabricLockAttestation
	e, err := c.settleEvent(ctx, swapID)
	if err != nil {
		return out, err
	}
	return c.decodeSettle(ctx, e)
}

// SettleSender devolve quem submeteu o settle do swap (ex.: a conta de relay do Cacti).
func (c *Client) SettleSender(ctx context.Context, swapID [32]byte) (common.Address, error) {
	e, err := c.settleEvent(ctx, swapID)
	if err != nil {
		return common.Address{}, err
	}
	tx, _, err := c.eth.TransactionByHash(ctx, e.TxHash)
	if err != nil {
		return common.Address{}, err
	}
	return c.Sender(tx)
}

// settleEvent encontra o PaymentSettled do swap.
func (c *Client) settleEvent(ctx context.Context, swapID [32]byte) (Event, error) {
	// O Besu rejeita toBlock além do topo e limita o intervalo de eth_getLogs (5000 por
	// padrão): busca em janelas de 2000 blocos, do topo para trás.
	head, err := c.eth.BlockNumber(ctx)
	if err != nil {
		return Event{}, err
	}
	const window = 2000
	var evs []Event
	for to := head; ; {
		from := uint64(0)
		if to >= window {
			from = to - window + 1
		}
		chunk, err := c.Events(ctx, from, to)
		if err != nil {
			return Event{}, err
		}
		found := false
		for _, e := range chunk {
			if e.Contract == "settlement" && e.Kind == EvClaimed && e.SwapID == swapID {
				evs, found = []Event{e}, true
			}
		}
		if found || from == 0 {
			break
		}
		to = from - 1
	}
	if len(evs) == 0 {
		return Event{}, fmt.Errorf("settle do swap %x não encontrado", swapID[:8])
	}
	return evs[0], nil
}

// decodeSettle decodifica a atestação do calldata da transação de settle.
func (c *Client) decodeSettle(ctx context.Context, e Event) (attest.FabricLockAttestation, error) {
	var out attest.FabricLockAttestation
	tx, _, err := c.eth.TransactionByHash(ctx, e.TxHash)
	if err != nil {
		return out, err
	}
	parsed, err := bindings.AttestedSettlementMetaData.ParseABI()
	if err != nil {
		return out, err
	}
	data := tx.Data()
	if len(data) < 4 {
		return out, fmt.Errorf("settle: calldata curto")
	}
	args, err := parsed.Methods["settle"].Inputs.Unpack(data[4:])
	if err != nil {
		return out, err
	}
	raw, err := json.Marshal(args[2])
	if err != nil {
		return out, err
	}
	var b bindings.FabricLockAttestation
	if err := json.Unmarshal(raw, &b); err != nil {
		return out, err
	}
	return attest.FabricLockAttestation{SwapID: b.SwapId, Hashlock: b.Hashlock, AssetIDHash: b.AssetIdHash,
		RecipientIDHash: b.RecipientIdHash, Deadline: b.Deadline, ChannelHash: b.ChannelHash,
		NamespaceHash: b.NamespaceHash, FabricTxID: b.FabricTxId, BlockNumber: b.BlockNumber}, nil
}
