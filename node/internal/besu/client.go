// Package besu é o cliente do nó Janus para a rede Besu (Drex): envio de transações para
// HTLCPayment (A) e AttestedSettlement (B), leitura de estado no formato de swap.BesuView e
// acompanhamento de eventos com checkpoint e backfill.
package besu

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"janus/node/internal/besu/bindings"
	"janus/node/internal/swap"
)

// Deployment são os endereços gravados por contracts/script/Deploy.s.sol.
type Deployment struct {
	ChainID             uint64         `json:"chainId"`
	DrexToken           common.Address `json:"DrexToken"`
	HTLCPayment         common.Address `json:"HTLCPayment"`
	FabricProofVerifier common.Address `json:"FabricProofVerifier"`
	AttestedSettlement  common.Address `json:"AttestedSettlement"`
	FabricChannel       string         `json:"fabricChannel"`
	FabricNamespace     string         `json:"fabricNamespace"`
}

// LoadDeployment lê deployments/<chainId>.local.json.
func LoadDeployment(path string) (*Deployment, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	d := &Deployment{}
	return d, json.Unmarshal(raw, d)
}

// Client é o cliente Besu de uma instituição (uma chave secp256k1).
type Client struct {
	eth     *ethclient.Client
	chainID *big.Int
	key     *ecdsa.PrivateKey
	from    common.Address
	d       *Deployment

	token      *bindings.DrexToken
	htlc       *bindings.HTLCPayment
	settlement *bindings.AttestedSettlement
	verifier   *bindings.FabricProofVerifier

	tokenC, htlcC, settlementC, verifierC *bind.BoundContract

	// ReceiptTimeout: espera máxima pela inclusão de uma transação.
	ReceiptTimeout time.Duration
}

// Dial conecta ao RPC e confere o chainId contra o deployment.
func Dial(ctx context.Context, rpcURL string, key *ecdsa.PrivateKey, d *Deployment) (*Client, error) {
	eth, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("besu: dial: %w", err)
	}
	id, err := eth.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("besu: chainId: %w", err)
	}
	if d.ChainID != 0 && id.Uint64() != d.ChainID {
		return nil, fmt.Errorf("besu: chainId %s != deployment %d", id, d.ChainID)
	}
	return newClient(eth, id, key, d), nil
}

const defaultReceiptTimeout = 60 * time.Second

func newClient(eth *ethclient.Client, chainID *big.Int, key *ecdsa.PrivateKey, d *Deployment) *Client {
	c := &Client{eth: eth, chainID: chainID, key: key, d: d, ReceiptTimeout: defaultReceiptTimeout}
	c.from = bind.NewKeyedTransactor(key, chainID).From
	c.token = bindings.NewDrexToken()
	c.htlc = bindings.NewHTLCPayment()
	c.settlement = bindings.NewAttestedSettlement()
	c.verifier = bindings.NewFabricProofVerifier()
	c.tokenC = c.token.Instance(eth, d.DrexToken)
	c.htlcC = c.htlc.Instance(eth, d.HTLCPayment)
	c.settlementC = c.settlement.Instance(eth, d.AttestedSettlement)
	c.verifierC = c.verifier.Instance(eth, d.FabricProofVerifier)
	return c
}

func (c *Client) Close()                  { c.eth.Close() }
func (c *Client) From() common.Address    { return c.from }
func (c *Client) ChainID() uint64         { return c.chainID.Uint64() }
func (c *Client) Deployment() *Deployment { return c.d }
func (c *Client) Eth() *ethclient.Client  { return c.eth }

func (c *Client) opts(ctx context.Context) *bind.TransactOpts {
	o := bind.NewKeyedTransactor(c.key, c.chainID)
	o.Context = ctx
	o.GasPrice = big.NewInt(0) // rede com gas zero (legacy tx)
	return o
}

// ErrReverted indica uma transação incluída com status 0.
var ErrReverted = errors.New("besu: transaction reverted")

// send transaciona e espera o recibo.
func (c *Client) send(ctx context.Context, bc *bind.BoundContract, data []byte) (*types.Receipt, error) {
	tx, err := bind.Transact(bc, c.opts(ctx), data)
	if err != nil {
		return nil, fmt.Errorf("besu: send: %w", err)
	}
	wctx, cancel := context.WithTimeout(ctx, c.ReceiptTimeout)
	defer cancel()
	r, err := bind.WaitMined(wctx, c.eth, tx.Hash())
	if err != nil {
		return nil, fmt.Errorf("besu: wait %s: %w", tx.Hash(), err)
	}
	if r.Status != types.ReceiptStatusSuccessful {
		return r, fmt.Errorf("%w: %s", ErrReverted, tx.Hash())
	}
	return r, nil
}

// call executa uma leitura.
func call[T any](ctx context.Context, bc *bind.BoundContract, data []byte, unpack func([]byte) (T, error)) (T, error) {
	return bind.Call(bc, &bind.CallOpts{Context: ctx}, data, unpack)
}

// --------------------------------------------------------------------- Drex

func (c *Client) Balance(ctx context.Context, who common.Address) (*big.Int, error) {
	return call(ctx, c.tokenC, c.token.PackBalanceOf(who), c.token.UnpackBalanceOf)
}

func (c *Client) Approve(ctx context.Context, spender common.Address, amount *big.Int) (*types.Receipt, error) {
	return c.send(ctx, c.tokenC, c.token.PackApprove(spender, amount))
}

// --------------------------------------------------------------------- Abordagem A

// LockPayment aprova e trava o pagamento no HTLCPayment (Banco).
func (c *Client) LockPayment(ctx context.Context, t *swap.SwapTerms) (*types.Receipt, error) {
	amt := new(big.Int).SetUint64(t.Price.Value)
	if _, err := c.Approve(ctx, c.d.HTLCPayment, amt); err != nil {
		return nil, err
	}
	return c.send(ctx, c.htlcC, c.htlc.PackLock(t.SwapID, common.Address(t.Seller.EVM), amt, t.Hashlock, uint64(t.T2)))
}

// ClaimPayment revela S no HTLCPayment (Vendedora; permissionless).
func (c *Client) ClaimPayment(ctx context.Context, swapID, preimage [32]byte) (*types.Receipt, error) {
	return c.send(ctx, c.htlcC, c.htlc.PackClaim(swapID, preimage))
}

// RefundPayment devolve o pagamento depois de T2 (Banco; permissionless).
func (c *Client) RefundPayment(ctx context.Context, swapID [32]byte) (*types.Receipt, error) {
	return c.send(ctx, c.htlcC, c.htlc.PackRefund(swapID))
}

// HTLCView lê o lock e devolve no formato da máquina de estados.
func (c *Client) HTLCView(ctx context.Context, swapID [32]byte) (swap.BesuView, error) {
	l, err := call(ctx, c.htlcC, c.htlc.PackGetLock(swapID), c.htlc.UnpackGetLock)
	if err != nil {
		return swap.BesuView{}, err
	}
	v := swap.BesuView{
		State: swap.BesuState(l.State), Contract: swap.Address(c.d.HTLCPayment), Token: swap.Address(c.d.DrexToken),
		Sender: swap.Address(l.Sender), Recipient: swap.Address(l.Recipient), Amount: l.Amount.Uint64(),
		Hashlock: l.Hashlock, Timelock: time.Unix(int64(l.Timelock), 0),
	}
	if v.State == swap.BesuClaimed {
		pre, err := call(ctx, c.htlcC, c.htlc.PackRevealedPreimage(swapID), c.htlc.UnpackRevealedPreimage)
		if err != nil {
			return v, err
		}
		v.Preimage = &pre
	}
	return v, nil
}

// --------------------------------------------------------------------- Abordagem B2

// Deposit aprova e deposita no AttestedSettlement (Banco).
func (c *Client) Deposit(ctx context.Context, t *swap.SwapTerms, assetIDHash, recipientIDHash [32]byte, minFabricDeadline int64) (*types.Receipt, error) {
	amt := new(big.Int).SetUint64(t.Price.Value)
	if _, err := c.Approve(ctx, c.d.AttestedSettlement, amt); err != nil {
		return nil, err
	}
	return c.send(ctx, c.settlementC, c.settlement.PackDeposit(
		t.SwapID, common.Address(t.Seller.EVM), amt, t.Hashlock, assetIDHash, recipientIDHash, uint64(t.T2), uint64(minFabricDeadline)))
}

// Settle revela S e apresenta a atestação do lock (Vendedora ou relay; permissionless).
func (c *Client) Settle(ctx context.Context, swapID, preimage [32]byte, att bindings.FabricLockAttestation, sigs []bindings.OrgSignature) (*types.Receipt, error) {
	return c.send(ctx, c.settlementC, c.settlement.PackSettle(swapID, preimage, att, sigs))
}

func (c *Client) RefundDeposit(ctx context.Context, swapID [32]byte) (*types.Receipt, error) {
	return c.send(ctx, c.settlementC, c.settlement.PackRefund(swapID))
}

// SettlementView lê o depósito e devolve no formato da máquina de estados.
func (c *Client) SettlementView(ctx context.Context, swapID [32]byte) (swap.BesuView, error) {
	d, err := call(ctx, c.settlementC, c.settlement.PackGetDeposit(swapID), c.settlement.UnpackGetDeposit)
	if err != nil {
		return swap.BesuView{}, err
	}
	v := swap.BesuView{
		State: swap.BesuState(d.State), Contract: swap.Address(c.d.AttestedSettlement), Token: swap.Address(c.d.DrexToken),
		Sender: swap.Address(d.Buyer), Recipient: swap.Address(d.Seller), Amount: d.Amount.Uint64(),
		Hashlock: d.Hashlock, Timelock: time.Unix(int64(d.Timelock), 0),
		MinFabricDeadline: time.Unix(int64(d.MinFabricDeadline), 0),
	}
	if v.State == swap.BesuClaimed {
		pre, err := call(ctx, c.settlementC, c.settlement.PackRevealedPreimage(swapID), c.settlement.UnpackRevealedPreimage)
		if err != nil {
			return v, err
		}
		v.Preimage = &pre
	}
	return v, nil
}

// View escolhe o contrato conforme o protocolo dos termos.
func (c *Client) View(ctx context.Context, t *swap.SwapTerms) (swap.BesuView, error) {
	if t.Protocol == swap.ProtocolCacti {
		return c.SettlementView(ctx, t.SwapID)
	}
	return c.HTLCView(ctx, t.SwapID)
}

// WaitReceipt espera a inclusão de uma transação enviada por outro componente (ex.: o relay do
// Cacti, que responde com o ack do pool) e devolve o recibo.
func (c *Client) WaitReceipt(ctx context.Context, h common.Hash, timeout time.Duration) (*types.Receipt, error) {
	wctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	r, err := bind.WaitMined(wctx, c.eth, h)
	if err != nil {
		return nil, fmt.Errorf("besu: wait %s: %w", h, err)
	}
	return r, nil
}

// Sender recupera o remetente de uma transação (assinatura EIP-155/legacy ou tipada).
func (c *Client) Sender(tx *types.Transaction) (common.Address, error) {
	return types.Sender(types.LatestSignerForChainID(c.chainID), tx)
}
