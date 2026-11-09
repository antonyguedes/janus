package janus

import (
	"context"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"janus/node/internal/besu/bindings"
	cacticonn "janus/node/internal/cacti"
	"os"

	fscview "github.com/hyperledger-labs/fabric-smart-client/platform/view/services/view"

	"janus/node/internal/attest"
	"janus/node/internal/besu"
	"janus/node/internal/swap"
	"janus/node/internal/token"
)

// settlementAdapter implementa cacti.SettlementPort sobre besu.Client.
// settlementAdapter implementa cacti.SettlementPort sobre besu.Client. Com relay configurado, o
// settle é submetido pelo conector Besu do Hyperledger Cacti (ADR-0006) e o nó confirma a inclusão
// e o status pela própria cadeia.
type settlementAdapter struct {
	c        *besu.Client
	relay    *cacticonn.BesuConnector
	fallback bool // submete direto se o relay falhar (vivacidade); padrão: não
	log      *slog.Logger
}

func (a settlementAdapter) View(ctx context.Context, t *swap.SwapTerms) (swap.BesuView, error) {
	return a.c.SettlementView(ctx, t.SwapID)
}
func (a settlementAdapter) Deposit(ctx context.Context, t *swap.SwapTerms, asset, recip [32]byte, minDeadline int64) error {
	_, err := a.c.Deposit(ctx, t, asset, recip, minDeadline)
	return err
}
func (a settlementAdapter) Settle(ctx context.Context, id, s [32]byte, att attest.FabricLockAttestation, sigs []attest.OrgSignature) error {
	if a.relay == nil {
		_, err := a.c.SettleAttested(ctx, id, s, att, sigs)
		return err
	}
	err := a.settleViaRelay(ctx, id, s, att, sigs)
	if err != nil && a.fallback {
		a.log.Warn("relay Cacti falhou; submetendo o settle diretamente (fallback)", "err", err)
		_, err = a.c.SettleAttested(ctx, id, s, att, sigs)
	}
	return err
}

func (a settlementAdapter) settleViaRelay(ctx context.Context, id, s [32]byte, att attest.FabricLockAttestation, sigs []attest.OrgSignature) error {
	dep := a.c.Deployment()
	res, err := a.relay.Send(ctx, bindings.AttestedSettlementMetaData.ABI, dep.AttestedSettlement.Hex(), "settle",
		cacticonn.SettleParams(id, s, att, sigs))
	if err != nil {
		return fmt.Errorf("settle via relay Cacti: %w", err)
	}
	// O nó não confia só na resposta do relay: confirma inclusão e status na cadeia.
	rc, err := a.c.WaitReceipt(ctx, common.HexToHash(res.TxHash), 60*time.Second)
	if err != nil {
		return err
	}
	if rc.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("settle via relay Cacti revertido: %s", res.TxHash)
	}
	a.log.Info("settle submetido pelo relay Cacti", "tx", res.TxHash, "relay", a.relay.Account)
	return nil
}
func (a settlementAdapter) Refund(ctx context.Context, id [32]byte) error {
	_, err := a.c.RefundDeposit(ctx, id)
	return err
}
func (a settlementAdapter) Domain() attest.Domain { return a.c.AttestationDomain() }

// fscCollector implementa cacti.Collector chamando a CollectView (relay por sessões FSC).
type fscCollector struct {
	views     *fscview.Manager
	attestors []string
}

func (c *fscCollector) Collect(ctx context.Context, swapID, hash [32]byte, d attest.Domain) (attest.FabricLockAttestation, []attest.OrgSignature, error) {
	if len(c.attestors) == 0 {
		return attest.FabricLockAttestation{}, nil, errors.New("no attestors configured (janus.json: attestors)")
	}
	res, err := c.views.InitiateView(ctx, &token.CollectView{Collect: &token.Collect{
		Request:   token.AttestRequest{SwapID: swapID, Hash: hash, ChainID: d.ChainID, Verifier: d.Verifier},
		Attestors: c.attestors,
	}})
	if err != nil {
		return attest.FabricLockAttestation{}, nil, err
	}
	out, ok := res.(*token.Collected)
	if !ok {
		return attest.FabricLockAttestation{}, nil, fmt.Errorf("collect: unexpected %T", res)
	}
	return out.Attestation, out.Signatures, nil
}

// loadOrCreateP256 lê uma chave EC P-256 (PEM PKCS#8) ou cria uma (DEV).
func loadOrCreateP256(path string) (*ecdsa.PrivateKey, error) {
	if raw, err := os.ReadFile(path); err == nil {
		b, _ := pem.Decode(raw)
		if b == nil {
			return nil, fmt.Errorf("attestor key %s: not PEM", path)
		}
		k, err := x509.ParsePKCS8PrivateKey(b.Bytes)
		if err != nil {
			return nil, err
		}
		ek, ok := k.(*ecdsa.PrivateKey)
		if !ok || ek.Curve != elliptic.P256() {
			return nil, fmt.Errorf("attestor key %s: not P-256", path)
		}
		return ek, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	der, err := x509.MarshalPKCS8PrivateKey(k)
	if err != nil {
		return nil, err
	}
	return k, os.WriteFile(path, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), 0o600)
}
