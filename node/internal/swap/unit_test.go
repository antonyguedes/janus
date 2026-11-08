package swap

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"testing"
	"time"
)

func TestPoliciesValid(t *testing.T) {
	for name, p := range map[string]TimePolicy{"prod": ProductionPolicy(), "dev": DevPolicy()} {
		if err := p.Validate(); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
	}
	if d := DevPolicy().Delta(); d != 90*time.Second {
		t.Fatalf("dev Δ = %s, want 90s (ADR-0004)", d)
	}
	bad := DevPolicy()
	bad.T1Offset = bad.T2Offset + 10*time.Second
	if bad.Validate() == nil {
		t.Fatal("T1-T2 < Δ must be rejected")
	}
	bad = DevPolicy()
	bad.DeadlineTolerance = time.Hour
	if bad.Validate() == nil {
		t.Fatal("tolerance > skew must be rejected")
	}
}

func validTerms(t *testing.T) (*SwapTerms, TimePolicy, time.Time) {
	w := newWorld(ProtocolHTLC)
	return w.t, w.p, time.Unix(base, 0)
}

func TestTermsValidate(t *testing.T) {
	terms, p, now := validTerms(t)
	if err := terms.Validate(p, now); err != nil {
		t.Fatal(err)
	}
	cases := map[string]func(x *SwapTerms){
		"swapId":     func(x *SwapTerms) { x.SwapID[0] ^= 1 },
		"protocol":   func(x *SwapTerms) { x.Protocol = "x" },
		"price":      func(x *SwapTerms) { x.Price.Value = 0 },
		"sameEVM":    func(x *SwapTerms) { x.Buyer.EVM = x.Seller.EVM },
		"t1<=t2":     func(x *SwapTerms) { x.T1 = x.T2 },
		"delta":      func(x *SwapTerms) { x.T1 = x.T2 + 30 },
		"noNetworks": func(x *SwapTerms) { x.Networks.FabricChannel = "" },
	}
	for name, mut := range cases {
		x := *terms
		mut(&x)
		if err := x.Validate(p, now); err == nil {
			t.Fatalf("%s: expected error", name)
		}
	}
	// tarde demais para aceitar
	if err := terms.Validate(p, terms.T2Time().Add(-time.Minute)); !errors.Is(err, ErrBadTimes) {
		t.Fatalf("late acceptance: %v", err)
	}
}

func TestCanonicalHashIsDeterministicAndBinding(t *testing.T) {
	terms, _, _ := validTerms(t)
	h := terms.Hash()
	cp := *terms
	if cp.Hash() != h {
		t.Fatal("not deterministic")
	}
	muts := []func(x *SwapTerms){
		func(x *SwapTerms) { x.AssetID += "x" },
		func(x *SwapTerms) { x.Price.Value++ },
		func(x *SwapTerms) { x.T1++ },
		func(x *SwapTerms) { x.Networks.BesuChainID++ },
		func(x *SwapTerms) { x.Networks.PaymentContract[0] ^= 1 },
		func(x *SwapTerms) { x.Seller.MSPID = "Outra" },
		// fronteira de campos com prefixo de tamanho: "ab"+"c" != "a"+"bc"
		func(x *SwapTerms) { x.Networks.FabricChannel, x.Networks.FabricNamespace = "recebivei", "stcc" },
	}
	for i, m := range muts {
		x := *terms
		m(&x)
		if x.Hash() == h {
			t.Fatalf("mutation %d did not change the hash", i)
		}
	}
}

// ecdsaSigner/Verifier de teste (em produção: SigService do FSC).
type keyring map[string]*ecdsa.PrivateKey

func (k keyring) Verify(id, digest, sig []byte) error {
	key, ok := k[string(id)]
	if !ok || !ecdsa.VerifyASN1(&key.PublicKey, digest, sig) {
		return errors.New("bad signature")
	}
	return nil
}

type signerFn func([]byte) ([]byte, error)

func (f signerFn) Sign(d []byte) ([]byte, error) { return f(d) }

func TestSignedTerms(t *testing.T) {
	terms, _, _ := validTerms(t)
	kr := keyring{}
	for _, id := range []string{"vendedora", "banco", "intruso"} {
		kr[id], _ = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	}
	signer := func(id string) Signer {
		return signerFn(func(d []byte) ([]byte, error) { return ecdsa.SignASN1(rand.Reader, kr[id], d) })
	}
	st := &SignedTerms{Terms: *terms}
	if err := st.Verify(kr); !errors.Is(err, ErrMissingSignature) {
		t.Fatal("expected missing signature")
	}
	_ = st.SignAs(RoleSeller, signer("vendedora"))
	_ = st.SignAs(RoleBuyer, signer("intruso"))
	if st.Verify(kr) == nil {
		t.Fatal("intruder signature accepted")
	}
	_ = st.SignAs(RoleBuyer, signer("banco"))
	if err := st.Verify(kr); err != nil {
		t.Fatal(err)
	}
	st.Terms.Price.Value++ // adulteração depois da assinatura
	if st.Verify(kr) == nil {
		t.Fatal("tampered terms accepted")
	}
}

func lockedViews(terms *SwapTerms) (FabricView, BesuView) {
	return FabricView{
		State: FabricLocked, AssetID: terms.AssetID, Sender: terms.Seller.FSCIdentity,
		Recipient: terms.Buyer.FSCIdentity, Hashlock: terms.Hashlock, HashOK: true, Deadline: terms.T1Time(),
	}, BesuView{
		State: BesuLocked, Contract: terms.Networks.PaymentContract, Token: terms.Price.Token,
		Sender: terms.Buyer.EVM, Recipient: terms.Seller.EVM, Amount: terms.Price.Value,
		Hashlock: terms.Hashlock, Timelock: terms.T2Time(),
	}
}

func TestCheckFabricLock(t *testing.T) {
	terms, p, _ := validTerms(t)
	f, _ := lockedViews(terms)
	if err := CheckFabricLock(terms, p, f); err != nil {
		t.Fatal(err)
	}
	muts := map[string]func(*FabricView){
		"state":     func(v *FabricView) { v.State = FabricNone },
		"asset":     func(v *FabricView) { v.AssetID = "outro" },
		"recipient": func(v *FabricView) { v.Recipient = []byte("laranja") },
		"sender":    func(v *FabricView) { v.Sender = []byte("x") },
		"hashFunc":  func(v *FabricView) { v.HashOK = false },
		"hashlock":  func(v *FabricView) { v.Hashlock[0] ^= 1 },
		"deadline-": func(v *FabricView) { v.Deadline = v.Deadline.Add(-p.DeadlineTolerance - time.Second) },
		"deadline+": func(v *FabricView) { v.Deadline = v.Deadline.Add(p.DeadlineTolerance + time.Second) },
	}
	for name, m := range muts {
		v := f
		m(&v)
		if CheckFabricLock(terms, p, v) == nil {
			t.Fatalf("%s: expected error", name)
		}
	}
	v := f
	v.Deadline = v.Deadline.Add(p.DeadlineTolerance) // exatamente na tolerância: aceito
	if err := CheckFabricLock(terms, p, v); err != nil {
		t.Fatal(err)
	}
}

func TestCheckBesuLock(t *testing.T) {
	terms, p, _ := validTerms(t)
	_, b := lockedViews(terms)
	if err := CheckBesuLock(terms, p, b); err != nil {
		t.Fatal(err)
	}
	muts := map[string]func(*BesuView){
		"state":     func(v *BesuView) { v.State = BesuNone },
		"contract":  func(v *BesuView) { v.Contract[0] ^= 1 },
		"token":     func(v *BesuView) { v.Token[0] ^= 1 },
		"amount":    func(v *BesuView) { v.Amount-- },
		"recipient": func(v *BesuView) { v.Recipient[0] ^= 1 },
		"sender":    func(v *BesuView) { v.Sender[0] ^= 1 },
		"hashlock":  func(v *BesuView) { v.Hashlock[0] ^= 1 },
		"timelock>": func(v *BesuView) { v.Timelock = v.Timelock.Add(time.Second) },
	}
	for name, m := range muts {
		v := b
		m(&v)
		if CheckBesuLock(terms, p, v) == nil {
			t.Fatalf("%s: expected error", name)
		}
	}
	// B: minFabricDeadline acima de T1 torna o settle impossível
	tb := *terms
	tb.Protocol = ProtocolCacti
	v := b
	v.MinFabricDeadline = tb.T1Time().Add(time.Second)
	if CheckBesuLock(&tb, p, v) == nil {
		t.Fatal("minFabricDeadline > T1 must be rejected")
	}
}

func TestDerive(t *testing.T) {
	cases := []struct {
		f    FabricState
		b    BesuState
		want Status
	}{
		{FabricNone, BesuNone, StatusNegotiated},
		{FabricLocked, BesuNone, StatusAssetLocked},
		{FabricNone, BesuLocked, StatusPaymentLocked},
		{FabricLocked, BesuLocked, StatusBothLocked},
		{FabricLocked, BesuClaimed, StatusRevealed},
		{FabricClaimed, BesuClaimed, StatusCompleted},
		{FabricReclaimed, BesuRefunded, StatusReverted},
		{FabricReclaimed, BesuNone, StatusReverted},
		{FabricNone, BesuRefunded, StatusReverted},
		{FabricReclaimed, BesuClaimed, StatusSellerHasBoth},
		{FabricClaimed, BesuRefunded, StatusBuyerHasBoth},
	}
	for _, c := range cases {
		got := Derive(Observation{Fabric: FabricView{State: c.f}, Besu: BesuView{State: c.b}})
		if got != c.want {
			t.Fatalf("(%s,%s) = %s, want %s", c.f, c.b, got, c.want)
		}
	}
	if !StatusCompleted.Terminal() || StatusRevealed.Terminal() {
		t.Fatal("terminal")
	}
}

func TestNextActionSellerNeverRevealsInsideMargin(t *testing.T) {
	terms, p, _ := validTerms(t)
	f, b := lockedViews(terms)
	s := sha256.Sum256([]byte("sim-secret"))
	for dt := -5 * time.Second; dt <= 5*time.Second; dt += time.Second {
		now := terms.T2Time().Add(-p.RevealMargin + dt)
		d := NextAction(RoleSeller, terms, p, Observation{Now: now, Fabric: f, Besu: b, Preimage: &s})
		reveals := d.Kind == ActClaimPayment
		if reveals != (dt < 0) {
			t.Fatalf("dt=%s: kind=%s", dt, d.Kind)
		}
	}
}

func TestNextActionBuyerRejectsBadLock(t *testing.T) {
	terms, p, now := validTerms(t)
	f, _ := lockedViews(terms)
	f.Recipient = []byte("laranja")
	d := NextAction(RoleBuyer, terms, p, Observation{Now: now, Fabric: f})
	if d.Kind != ActAbort {
		t.Fatalf("kind=%s", d.Kind)
	}
	if NextAction("x", terms, p, Observation{}).Kind != ActAlert {
		t.Fatal("unknown role")
	}
}
