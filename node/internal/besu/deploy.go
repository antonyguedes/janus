package besu

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"janus/node/internal/besu/bindings"
)

// DeploySpec descreve um deploy local (testes/E2E). Em ambiente compartilhado, use
// contracts/script/Deploy.s.sol.
type DeploySpec struct {
	Channel, Namespace string
	Allow              []common.Address            // contas na allowlist do Drex
	Mint               map[common.Address]*big.Int // saldos iniciais (centavos)
}

// DeployAll faz o deploy dos quatro contratos com a chave admin (emissora do Drex e governança).
func DeployAll(ctx context.Context, eth *ethclient.Client, admin *ecdsa.PrivateKey, spec DeploySpec) (*Deployment, error) {
	chainID, err := eth.ChainID(ctx)
	if err != nil {
		return nil, err
	}
	opts := bind.NewKeyedTransactor(admin, chainID)
	opts.Context = ctx
	opts.GasPrice = big.NewInt(0)
	adminAddr := opts.From

	deploy := func(meta *bind.MetaData, input []byte) (common.Address, error) {
		addr, tx, err := bind.DeployContract(opts, common.FromHex(meta.Bin), eth, input)
		if err != nil {
			return common.Address{}, err
		}
		if _, err := bind.WaitDeployed(ctx, eth, tx.Hash()); err != nil {
			return common.Address{}, err
		}
		return addr, nil
	}
	d := &Deployment{ChainID: chainID.Uint64(), FabricChannel: spec.Channel, FabricNamespace: spec.Namespace}
	if d.DrexToken, err = deploy(&bindings.DrexTokenMetaData, bindings.NewDrexToken().PackConstructor(adminAddr)); err != nil {
		return nil, fmt.Errorf("deploy DrexToken: %w", err)
	}
	if d.HTLCPayment, err = deploy(&bindings.HTLCPaymentMetaData, bindings.NewHTLCPayment().PackConstructor(d.DrexToken)); err != nil {
		return nil, fmt.Errorf("deploy HTLCPayment: %w", err)
	}
	if d.FabricProofVerifier, err = deploy(&bindings.FabricProofVerifierMetaData, bindings.NewFabricProofVerifier().PackConstructor(adminAddr)); err != nil {
		return nil, fmt.Errorf("deploy FabricProofVerifier: %w", err)
	}
	ch, ns := sha256.Sum256([]byte(spec.Channel)), sha256.Sum256([]byte(spec.Namespace))
	if d.AttestedSettlement, err = deploy(&bindings.AttestedSettlementMetaData,
		bindings.NewAttestedSettlement().PackConstructor(d.DrexToken, d.FabricProofVerifier, ch, ns)); err != nil {
		return nil, fmt.Errorf("deploy AttestedSettlement: %w", err)
	}

	c := newClient(eth, chainID, admin, d) // reaproveita send() com espera de recibo
	tok := c.token
	for _, a := range append(spec.Allow, d.HTLCPayment, d.AttestedSettlement) {
		if _, err := c.send(ctx, c.tokenC, tok.PackSetAllowed(a, true)); err != nil {
			return nil, fmt.Errorf("allowlist %s: %w", a, err)
		}
	}
	for a, v := range spec.Mint {
		if _, err := c.send(ctx, c.tokenC, tok.PackMint(a, v)); err != nil {
			return nil, fmt.Errorf("mint %s: %w", a, err)
		}
	}
	return d, nil
}
