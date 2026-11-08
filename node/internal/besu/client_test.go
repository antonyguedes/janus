package besu

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"net"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"janus/node/internal/hashlock"
	"janus/node/internal/swap"
)

// Chaves de DESENVOLVIMENTO do anvil (públicas). Iguais às da allowlist do infra/besu.
var (
	devAdmin  = mustKey("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	devSeller = mustKey("59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d")
	devBuyer  = mustKey("5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a")
)

type ecdsaKey = ecdsa.PrivateKey

func mustKey(h string) *ecdsaKey {
	k, err := crypto.HexToECDSA(h)
	if err != nil {
		panic(err)
	}
	return k
}

func freePort(t *testing.T) int {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

// startAnvil sobe um anvil com EVM Osaka e gas zero (mesmas condições do Besu QBFT local).
func startAnvil(t *testing.T) string {
	t.Helper()
	bin, err := exec.LookPath("anvil")
	if err != nil {
		t.Skip("anvil não encontrado no PATH")
	}
	port := freePort(t)
	cmd := exec.Command(bin, "--port", fmt.Sprint(port), "--hardfork", "osaka",
		"--block-base-fee-per-gas", "0", "--gas-price", "0", "--chain-id", "1337", "--silent")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill(); _ = cmd.Wait() })
	url := fmt.Sprintf("http://127.0.0.1:%d", port)
	for range 100 {
		if c, err := ethclient.Dial(url); err == nil {
			if _, err := c.ChainID(context.Background()); err == nil {
				c.Close()
				return url
			}
			c.Close()
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("anvil não respondeu")
	return ""
}

func addr(k *ecdsaKey) common.Address { return crypto.PubkeyToAddress(k.PublicKey) }

func setup(t *testing.T) (ctx context.Context, url string, d *Deployment) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	t.Cleanup(cancel)
	url = startAnvil(t)
	eth, err := ethclient.Dial(url)
	if err != nil {
		t.Fatal(err)
	}
	defer eth.Close()
	d, err = DeployAll(ctx, eth, devAdmin, DeploySpec{
		Channel: "recebiveis", Namespace: "tcc",
		Allow: []common.Address{addr(devSeller), addr(devBuyer)},
		Mint:  map[common.Address]*big.Int{addr(devBuyer): big.NewInt(10_000_000_00)},
	})
	if err != nil {
		t.Fatal(err)
	}
	return ctx, url, d
}

func terms(t *testing.T, ctx context.Context, c *Client, d *Deployment, h hashlock.Hash, t2 time.Duration) *swap.SwapTerms {
	hdr, err := c.Eth().HeaderByNumber(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	now := int64(hdr.Time)
	seller := swap.Party{MSPID: "VendedoraMSP", FSCIdentity: []byte("vendedora"), EVM: swap.Address(addr(devSeller))}
	buyer := swap.Party{MSPID: "BancoMSP", FSCIdentity: []byte("banco"), EVM: swap.Address(addr(devBuyer))}
	tt := &swap.SwapTerms{
		Protocol: swap.ProtocolHTLC, AssetID: "recebivel-1", Seller: seller, Buyer: buyer,
		Price: swap.Amount{Token: swap.Address(d.DrexToken), Value: 1_250_000_00}, Hashlock: h,
		T2: now + int64(t2/time.Second), T1: now + 2*int64(t2/time.Second),
		Networks: swap.NetworkRefs{FabricChannel: "recebiveis", FabricNamespace: "tcc", BesuChainID: d.ChainID,
			PaymentContract: swap.Address(d.HTLCPayment)},
	}
	tt.SwapID = swap.ComputeSwapID(h, seller, buyer, tt.Nonce)
	return tt
}

func TestHTLCLockClaimAndWatch(t *testing.T) {
	ctx, url, d := setup(t)
	buyer, err := Dial(ctx, url, devBuyer, d)
	if err != nil {
		t.Fatal(err)
	}
	seller, err := Dial(ctx, url, devSeller, d)
	if err != nil {
		t.Fatal(err)
	}
	s, h, _ := hashlock.New()
	tt := terms(t, ctx, buyer, d, h, 3*time.Minute)

	// Watcher do Banco (watchtower), com checkpoint em arquivo.
	cp := FileCheckpoint{Path: filepath.Join(t.TempDir(), "besu.checkpoint")}
	events := make(chan Event, 16)
	wctx, stop := context.WithCancel(ctx)
	go func() { _ = buyer.Watch(wctx, cp, WatchOptions{Poll: 100 * time.Millisecond}, events) }()

	if _, err := buyer.LockPayment(ctx, tt); err != nil {
		t.Fatal(err)
	}
	v, err := seller.View(ctx, tt)
	if err != nil {
		t.Fatal(err)
	}
	if err := swap.CheckBesuLock(tt, swap.DevPolicy(), v); err != nil {
		t.Fatalf("a Vendedora deveria aceitar o lock: %v (%+v)", err, v)
	}
	if _, err := seller.ClaimPayment(ctx, tt.SwapID, [32]byte{1}); err == nil {
		t.Fatal("preimagem errada deveria reverter")
	}
	if _, err := seller.ClaimPayment(ctx, tt.SwapID, s); err != nil {
		t.Fatal(err)
	}

	var got []Event
	deadline := time.After(10 * time.Second)
	for len(got) < 2 {
		select {
		case e := <-events:
			got = append(got, e)
		case <-deadline:
			t.Fatalf("eventos insuficientes: %+v", got)
		}
	}
	stop()
	if got[0].Kind != EvLocked || got[1].Kind != EvClaimed || got[1].Preimage == nil || *got[1].Preimage != s {
		t.Fatalf("eventos: %+v", got)
	}
	v, _ = buyer.View(ctx, tt)
	if v.State != swap.BesuClaimed || v.Preimage == nil || *v.Preimage != s {
		t.Fatalf("view: %+v", v)
	}
	bal, _ := buyer.Balance(ctx, addr(devSeller))
	if bal.Uint64() != tt.Price.Value {
		t.Fatalf("saldo da vendedora: %s", bal)
	}

	// Reinício: o watcher retoma do checkpoint e não reentrega eventos já processados.
	last, _ := cp.Load()
	if last == 0 {
		t.Fatal("checkpoint não avançou")
	}
	evs := make(chan Event, 4)
	w2, stop2 := context.WithTimeout(ctx, 700*time.Millisecond)
	defer stop2()
	_ = buyer.Watch(w2, cp, WatchOptions{Poll: 100 * time.Millisecond}, evs)
	if len(evs) != 0 {
		t.Fatalf("reentregou %d eventos depois do checkpoint", len(evs))
	}
}

func TestHTLCRefundAfterT2(t *testing.T) {
	ctx, url, d := setup(t)
	buyer, err := Dial(ctx, url, devBuyer, d)
	if err != nil {
		t.Fatal(err)
	}
	_, h, _ := hashlock.New()
	tt := terms(t, ctx, buyer, d, h, 3*time.Minute)
	if _, err := buyer.LockPayment(ctx, tt); err != nil {
		t.Fatal(err)
	}
	if _, err := buyer.RefundPayment(ctx, tt.SwapID); err == nil {
		t.Fatal("refund antes de T2 deveria reverter")
	}
	rpc := buyer.Eth().Client()
	if err := rpc.CallContext(ctx, nil, "evm_increaseTime", 181); err != nil {
		t.Fatal(err)
	}
	if err := rpc.CallContext(ctx, nil, "evm_mine"); err != nil {
		t.Fatal(err)
	}
	if _, err := buyer.RefundPayment(ctx, tt.SwapID); err != nil {
		t.Fatal(err)
	}
	v, _ := buyer.View(ctx, tt)
	if v.State != swap.BesuRefunded {
		t.Fatalf("state: %s", v.State)
	}
}
