package engine

import (
	"context"
	"crypto/sha256"
	"errors"
	"sync"
	"testing"
	"time"

	"janus/node/internal/protocol/htlc"
	"janus/node/internal/swap"
)

// ------------------------------------------------------------------ cadeias falsas compartilhadas

type chains struct {
	mu   sync.Mutex
	fab  swap.FabricView
	besu swap.BesuView
}

type fakeFabric struct {
	c    *chains
	me   swap.Role
	exp  bool
	lock time.Duration // latência de commit
}

func (f *fakeFabric) Status(_ context.Context, _ [32]byte) (swap.FabricView, error) {
	f.c.mu.Lock()
	defer f.c.mu.Unlock()
	return f.c.fab, nil
}
func (f *fakeFabric) Expect(context.Context, *swap.SwapTerms, time.Duration) error {
	f.exp = true
	return nil
}
func (f *fakeFabric) Lock(_ context.Context, t *swap.SwapTerms) error {
	f.c.mu.Lock()
	defer f.c.mu.Unlock()
	if f.c.fab.State != swap.FabricNone {
		return nil // idempotente
	}
	f.c.fab = swap.FabricView{State: swap.FabricLocked, AssetID: t.AssetID, Sender: t.Seller.FSCIdentity,
		Recipient: t.Buyer.FSCIdentity, Hashlock: t.Hashlock, HashOK: true, Deadline: t.T1Time()}
	return nil
}
func (f *fakeFabric) Claim(_ context.Context, s [32]byte) error {
	f.c.mu.Lock()
	defer f.c.mu.Unlock()
	if f.c.fab.State != swap.FabricLocked || !time.Now().Before(f.c.fab.Deadline) || sha256.Sum256(s[:]) != f.c.fab.Hashlock {
		return errors.New("fabric: claim rejected")
	}
	f.c.fab.State = swap.FabricClaimed
	return nil
}
func (f *fakeFabric) Reclaim(context.Context, [32]byte) error {
	f.c.mu.Lock()
	defer f.c.mu.Unlock()
	if f.c.fab.State != swap.FabricLocked || time.Now().Before(f.c.fab.Deadline) {
		return errors.New("fabric: reclaim rejected")
	}
	f.c.fab.State = swap.FabricReclaimed
	return nil
}

type fakeBesu struct{ c *chains }

func (b *fakeBesu) View(context.Context, *swap.SwapTerms) (swap.BesuView, error) {
	b.c.mu.Lock()
	defer b.c.mu.Unlock()
	return b.c.besu, nil
}
func (b *fakeBesu) LockPayment(_ context.Context, t *swap.SwapTerms) error {
	b.c.mu.Lock()
	defer b.c.mu.Unlock()
	if b.c.besu.State != swap.BesuNone {
		return nil
	}
	b.c.besu = swap.BesuView{State: swap.BesuLocked, Contract: t.Networks.PaymentContract, Token: t.Price.Token,
		Sender: t.Buyer.EVM, Recipient: t.Seller.EVM, Amount: t.Price.Value, Hashlock: t.Hashlock, Timelock: t.T2Time()}
	return nil
}
func (b *fakeBesu) ClaimPayment(_ context.Context, _ [32]byte, s [32]byte) error {
	b.c.mu.Lock()
	defer b.c.mu.Unlock()
	if b.c.besu.State != swap.BesuLocked || !time.Now().Before(b.c.besu.Timelock) || sha256.Sum256(s[:]) != b.c.besu.Hashlock {
		return errors.New("besu: claim reverted")
	}
	b.c.besu.State, b.c.besu.Preimage = swap.BesuClaimed, &s
	return nil
}
func (b *fakeBesu) RefundPayment(context.Context, [32]byte) error {
	b.c.mu.Lock()
	defer b.c.mu.Unlock()
	if b.c.besu.State != swap.BesuLocked || time.Now().Before(b.c.besu.Timelock) {
		return errors.New("besu: refund reverted")
	}
	b.c.besu.State = swap.BesuRefunded
	return nil
}

type memSecrets struct {
	mu sync.Mutex
	m  map[[32]byte][32]byte
}

func (m *memSecrets) Get(_ context.Context, id [32]byte) (*[32]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.m[id]; ok {
		return &s, nil
	}
	return nil, nil
}
func (m *memSecrets) Put(_ context.Context, id, s [32]byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m[id] = s
	return nil
}

// ------------------------------------------------------------------ cenário

// fastPolicy: prazos em segundos para o teste rodar rápido.
func fastPolicy() swap.TimePolicy {
	return swap.TimePolicy{
		T1Offset: 4 * time.Second, T2Offset: 2 * time.Second,
		Detection: 300 * time.Millisecond, FabricFinality: 200 * time.Millisecond, Relay: 100 * time.Millisecond,
		Skew: 200 * time.Millisecond, Operational: 200 * time.Millisecond,
		RevealMargin: 500 * time.Millisecond, BuyerLockMargin: 500 * time.Millisecond,
		DeadlineTolerance: 200 * time.Millisecond,
	}
}

type fixture struct {
	t               *swap.SwapTerms
	c               *chains
	sellerSecrets   *memSecrets
	buyerSecrets    *memSecrets
	seller, buyer   *Engine
	sellerP, buyerP *htlc.Protocol
}

func newFixture(sellerFaults, buyerFaults Faults) *fixture {
	p := fastPolicy()
	s := sha256.Sum256([]byte("engine-secret"))
	h := sha256.Sum256(s[:])
	now := time.Now().Truncate(time.Second).Add(time.Second)
	seller := swap.Party{MSPID: "VendedoraMSP", FSCIdentity: []byte("vendedora"), EVM: swap.Address{1}}
	buyer := swap.Party{MSPID: "BancoMSP", FSCIdentity: []byte("banco"), EVM: swap.Address{2}}
	t := &swap.SwapTerms{Protocol: swap.ProtocolHTLC, AssetID: "r1", Seller: seller, Buyer: buyer,
		Price: swap.Amount{Token: swap.Address{9}, Value: 100}, Hashlock: h,
		T1: now.Add(p.T1Offset).Unix(), T2: now.Add(p.T2Offset).Unix(),
		Networks: swap.NetworkRefs{FabricChannel: "c", FabricNamespace: "n", BesuChainID: 1, PaymentContract: swap.Address{7}}}
	t.SwapID = swap.ComputeSwapID(h, seller, buyer, t.Nonce)

	f := &fixture{t: t, c: &chains{}, sellerSecrets: &memSecrets{m: map[[32]byte][32]byte{t.SwapID: s}},
		buyerSecrets: &memSecrets{m: map[[32]byte][32]byte{}}}
	f.sellerP = &htlc.Protocol{Fabric: &fakeFabric{c: f.c, me: swap.RoleSeller}, Besu: &fakeBesu{c: f.c}, Secrets: f.sellerSecrets, Policy: p}
	f.buyerP = &htlc.Protocol{Fabric: &fakeFabric{c: f.c, me: swap.RoleBuyer}, Besu: &fakeBesu{c: f.c}, Secrets: f.buyerSecrets, Policy: p}
	mk := func(role swap.Role, pr *htlc.Protocol, fl Faults) *Engine {
		return New(Config{Role: role, Policy: p, Poll: 50 * time.Millisecond, Faults: fl,
			Protocol: map[swap.Protocol]swap.SettlementProtocol{swap.ProtocolHTLC: pr}})
	}
	f.seller, f.buyer = mk(swap.RoleSeller, f.sellerP, sellerFaults), mk(swap.RoleBuyer, f.buyerP, buyerFaults)
	return f
}

func (f *fixture) runBoth(t *testing.T) (sd, bd swap.Decision, berr error) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); sd, _ = f.seller.Run(ctx, f.t) }()
	go func() { defer wg.Done(); bd, berr = f.buyer.Run(ctx, f.t) }()
	wg.Wait()
	return
}

func TestEngineHappyPath(t *testing.T) {
	f := newFixture(Faults{}, Faults{})
	steps, unsub := f.buyer.Subscribe()
	defer unsub()
	sd, bd, _ := f.runBoth(t)
	if sd.Kind != swap.ActDone || bd.Kind != swap.ActDone {
		t.Fatalf("seller=%+v buyer=%+v", sd, bd)
	}
	if st := swap.Derive(swap.Observation{Fabric: f.c.fab, Besu: f.c.besu}); st != swap.StatusCompleted {
		t.Fatalf("status %s", st)
	}
	if len(steps) == 0 {
		t.Fatal("o motor deveria publicar passos para o SSE")
	}
}

func TestEngineSellerNeverReveals(t *testing.T) {
	f := newFixture(Faults{NoReveal: true}, Faults{})
	sd, bd, _ := f.runBoth(t)
	if sd.Kind != swap.ActDone || bd.Kind != swap.ActDone {
		t.Fatalf("seller=%+v buyer=%+v", sd, bd)
	}
	if st := swap.Derive(swap.Observation{Fabric: f.c.fab, Besu: f.c.besu}); st != swap.StatusReverted {
		t.Fatalf("status %s (fab=%s besu=%s)", st, f.c.fab.State, f.c.besu.State)
	}
}

func TestEngineBuyerNeverLocks(t *testing.T) {
	f := newFixture(Faults{}, Faults{NoLock: true})
	sd, _, _ := f.runBoth(t)
	if sd.Kind != swap.ActDone || f.c.fab.State != swap.FabricReclaimed {
		t.Fatalf("seller=%+v fab=%s", sd, f.c.fab.State)
	}
}

// Cenário 5: o Banco cai logo depois de ver S e volta. A reconciliação no boot é só rodar o
// motor de novo: a decisão deriva do estado das cadeias e do S persistido.
func TestEngineBuyerCrashAfterRevealAndRecovers(t *testing.T) {
	f := newFixture(Faults{}, Faults{})
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	// O processo do Banco "morre" no instante em que ia fazer o claim do NFT.
	crashCtx, crash := context.WithCancel(ctx)
	crashy := New(Config{Role: swap.RoleBuyer, Policy: fastPolicy(), Poll: 50 * time.Millisecond,
		Protocol: map[swap.Protocol]swap.SettlementProtocol{swap.ProtocolHTLC: &crashBeforeClaim{f.buyerP, crash}}})
	done := make(chan swap.Decision, 1)
	go func() { d, _ := f.seller.Run(ctx, f.t); done <- d }()
	if _, err := crashy.Run(crashCtx, f.t); !errors.Is(err, context.Canceled) {
		t.Fatalf("esperava queda simulada, err=%v", err)
	}
	if f.c.fab.State != swap.FabricLocked || f.c.besu.State != swap.BesuClaimed {
		t.Fatalf("estado na queda: fab=%s besu=%s", f.c.fab.State, f.c.besu.State)
	}
	if s, _ := f.buyerSecrets.Get(ctx, f.t.SwapID); s == nil {
		t.Fatal("o Banco deveria ter persistido S antes de cair")
	}
	// reinício (reconciliação)
	bd, err := f.buyer.Run(ctx, f.t)
	if err != nil || bd.Kind != swap.ActDone || f.c.fab.State != swap.FabricClaimed {
		t.Fatalf("recuperação: %+v err=%v fab=%s", bd, err, f.c.fab.State)
	}
	if sd := <-done; sd.Kind != swap.ActDone {
		t.Fatalf("seller=%+v", sd)
	}
}

// crashBeforeClaim simula a queda do processo no momento em que o Banco ia fazer o claim do NFT.
type crashBeforeClaim struct {
	*htlc.Protocol
	crash context.CancelFunc
}

func (c *crashBeforeClaim) Execute(ctx context.Context, r swap.Role, t *swap.SwapTerms, d swap.Decision) error {
	if d.Kind == swap.ActClaimAsset {
		c.crash()
		return context.Canceled
	}
	return c.Protocol.Execute(ctx, r, t, d)
}
