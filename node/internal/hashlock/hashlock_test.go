package hashlock

import (
	"encoding/json"
	"os"
	"testing"
)

func TestVectorsShared(t *testing.T) {
	raw, err := os.ReadFile("../../../testdata/vectors/htlc.json")
	if err != nil {
		t.Fatal(err)
	}
	var f struct {
		Vectors []struct{ Preimage, Hashlock string } `json:"vectors"`
	}
	if err := json.Unmarshal(raw, &f); err != nil {
		t.Fatal(err)
	}
	if len(f.Vectors) == 0 {
		t.Fatal("no vectors")
	}
	for i, v := range f.Vectors {
		s, err := ParseHash32(v.Preimage)
		if err != nil {
			t.Fatal(err)
		}
		h, err := ParseHash32(v.Hashlock)
		if err != nil {
			t.Fatal(err)
		}
		if err := Verify(Preimage(s), Hash(h)); err != nil {
			t.Fatalf("vector %d: %v", i, err)
		}
	}
}

func TestNewAndVerify(t *testing.T) {
	s, h, err := New()
	if err != nil {
		t.Fatal(err)
	}
	if err := Verify(s, h); err != nil {
		t.Fatal(err)
	}
	s[0] ^= 1
	if err := Verify(s, h); err != ErrMismatch {
		t.Fatalf("expected mismatch, got %v", err)
	}
}

func TestParseHash32Rejects(t *testing.T) {
	for _, in := range []string{"0x00", "zz", "0x" + string(make([]byte, 0))} {
		if _, err := ParseHash32(in); err == nil {
			t.Fatalf("expected error for %q", in)
		}
	}
}
