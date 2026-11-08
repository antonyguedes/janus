package token

import (
	"crypto"
	"crypto/sha256"
	"strings"
	"testing"
	"time"

	"github.com/LFDT-Panurus/panurus/token/services/interop/encoding"
	"github.com/LFDT-Panurus/panurus/token/services/interop/htlc"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/view"
)

func script(deadline time.Time) *htlc.Script {
	h := sha256.Sum256([]byte("s"))
	return &htlc.Script{
		Sender: view.Identity("vendedora"), Recipient: view.Identity("banco"), Deadline: deadline,
		HashInfo: htlc.HashInfo{Hash: h[:], HashFunc: crypto.SHA256, HashEncoding: encoding.None},
	}
}

func TestCheckLockScript(t *testing.T) {
	p, now := DefaultAuditPolicy(), time.Now()
	if err := p.CheckLockScript(script(now.Add(6*time.Minute)), now); err != nil {
		t.Fatal(err)
	}
	bad := map[string]*htlc.Script{}
	s := script(now.Add(10 * time.Second))
	bad["too short"] = s
	s = script(now.Add(30 * 24 * time.Hour))
	bad["too long"] = s
	s = script(now.Add(time.Hour))
	s.HashInfo.HashEncoding = encoding.Base64
	bad["encoding"] = s
	s = script(now.Add(time.Hour))
	s.HashInfo.HashFunc = crypto.SHA512
	bad["hashfunc"] = s
	s = script(now.Add(time.Hour))
	s.HashInfo.Hash = []byte("short")
	bad["hashlen"] = s
	s = script(now.Add(time.Hour))
	s.Recipient = s.Sender
	bad["parties"] = s
	for name, s := range bad {
		if p.CheckLockScript(s, now) == nil {
			t.Fatalf("%s: expected rejection", name)
		}
	}
}

func TestClassifySpend(t *testing.T) {
	s := script(time.Now())
	if op, err := ClassifySpend(s, []byte("banco")); err != nil || op != OpClaim {
		t.Fatal("claim")
	}
	if op, err := ClassifySpend(s, []byte("vendedora")); err != nil || op != OpReclaim {
		t.Fatal("reclaim")
	}
	if _, err := ClassifySpend(s, []byte("laranja")); err == nil {
		t.Fatal("third party")
	}
}

// Camada 2 do ADR-0004: o auditor recusa reclaim antes de deadline+grace com o próprio relógio,
// mesmo que os endossantes já tenham aprovado.
func TestCheckSpendTimeBoundaries(t *testing.T) {
	p := DefaultAuditPolicy()
	T1 := time.Unix(1_900_000_000, 0)
	s := script(T1)
	cases := []struct {
		op   HTLCOp
		now  time.Time
		want bool
	}{
		{OpClaim, T1.Add(-time.Second), true},
		{OpClaim, T1, false},
		{OpReclaim, T1, false},                                   // endossante diria ok; auditor recusa
		{OpReclaim, T1.Add(p.ReclaimGrace - time.Second), false}, // ainda na tolerância
		{OpReclaim, T1.Add(p.ReclaimGrace), true},
		{OpReclaim, T1.Add(-5 * time.Hour), false}, // ataque do relógio forjado
	}
	for i, c := range cases {
		err := p.CheckSpend(s, c.op, c.now)
		if (err == nil) != c.want {
			t.Fatalf("case %d: err=%v", i, err)
		}
	}
}

func TestGravames(t *testing.T) {
	g := NewGravames()
	if g.Check("r1") != nil {
		t.Fatal("clean")
	}
	g.Set("r1", "penhor")
	if err := g.Check("r1"); err == nil || !strings.Contains(err.Error(), "penhor") {
		t.Fatal("gravame")
	}
	g.Clear("r1")
	if g.Check("r1") != nil {
		t.Fatal("cleared")
	}
}

func TestExpectationVerify(t *testing.T) {
	T1 := time.Unix(1_900_000_000, 0)
	s := script(T1.Add(10 * time.Second))
	x := Expectation{Hash: s.HashInfo.Hash, LinearID: "abc", Deadline: T1, Tolerance: 30 * time.Second}
	if err := x.Verify(s, "abc"); err != nil {
		t.Fatal(err)
	}
	if x.Verify(s, "outro") == nil {
		t.Fatal("asset")
	}
	if x.Verify(script(T1.Add(-time.Minute)), "abc") == nil {
		t.Fatal("deadline")
	}
}

func TestTokenTypeRoundTrip(t *testing.T) {
	r := &Recebivel{LinearID: "id-1", Registro: "REG-1", Cedente: "Vendedora", Sacado: "Sacado SA",
		ValorFace: 1_250_000_00, Vencimento: "2027-03-31"}
	if err := r.Validate(); err != nil {
		t.Fatal(err)
	}
	typ, err := TokenType(r)
	if err != nil {
		t.Fatal(err)
	}
	back, err := FromTokenType(typ)
	if err != nil || *back != *r {
		t.Fatalf("roundtrip: %v %+v", err, back)
	}
	if _, err := FromTokenType("REAL"); err == nil {
		t.Fatal("non-nft type must fail")
	}
	if (&Recebivel{}).Validate() == nil {
		t.Fatal("empty recebivel")
	}
}

func TestToFabricView(t *testing.T) {
	h := [32]byte{1}
	cases := map[string]int{"none": 0, "locked": 1, "claimed": 2, "reclaimed": 3, "": 0}
	for in, want := range cases {
		v := toFabricView(h, LockStatus{State: in, LinearID: "x", HashOK: true})
		if int(v.State) != want || v.Hashlock != h || v.AssetID != "x" {
			t.Fatalf("%q -> %+v", in, v)
		}
	}
}
