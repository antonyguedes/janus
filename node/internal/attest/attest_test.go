package attest

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

type vecOrg struct{ MSPID, MSPIDHash, Qx, Qy, R, S string }

type vecFile struct {
	ChainID                                            uint64
	Verifier, DomainTag, SwapID, Hashlock, AssetID     string
	AssetIDHash, RecipientIDHash, Channel, ChannelHash string
	Namespace, NamespaceHash, FabricTxID, Digest       string
	Deadline, BlockNumber                              uint64
	Orgs                                               []vecOrg
}

func b32(t *testing.T, s string) [32]byte {
	t.Helper()
	b, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	if err != nil || len(b) != 32 {
		t.Fatalf("bad bytes32 %q", s)
	}
	var o [32]byte
	copy(o[:], b)
	return o
}

func load(t *testing.T) (vecFile, Domain, FabricLockAttestation) {
	raw, err := os.ReadFile("../../../testdata/vectors/attestation.json")
	if err != nil {
		t.Fatal(err)
	}
	var v vecFile
	if err := json.Unmarshal(raw, &v); err != nil {
		t.Fatal(err)
	}
	ab, _ := hex.DecodeString(strings.TrimPrefix(v.Verifier, "0x"))
	var addr Address
	copy(addr[:], ab)
	a := FabricLockAttestation{
		SwapID: b32(t, v.SwapID), Hashlock: b32(t, v.Hashlock), AssetIDHash: b32(t, v.AssetIDHash),
		RecipientIDHash: b32(t, v.RecipientIDHash), Deadline: v.Deadline, ChannelHash: b32(t, v.ChannelHash),
		NamespaceHash: b32(t, v.NamespaceHash), FabricTxID: b32(t, v.FabricTxID), BlockNumber: v.BlockNumber,
	}
	return v, Domain{ChainID: v.ChainID, Verifier: addr}, a
}

func TestVectorDigestAndSignatures(t *testing.T) {
	v, d, a := load(t)
	if DomainTag != b32(t, v.DomainTag) {
		t.Fatal("domain tag")
	}
	if HashString(v.AssetID) != a.AssetIDHash || HashString(v.Channel) != a.ChannelHash ||
		HashString(v.Namespace) != a.NamespaceHash {
		t.Fatal("string hashes")
	}
	if got := Digest(d, a); got != b32(t, v.Digest) {
		t.Fatalf("digest %x != %s", got, v.Digest)
	}
	for _, o := range v.Orgs {
		if MSPIDHash(o.MSPID) != b32(t, o.MSPIDHash) {
			t.Fatalf("mspIdHash %s", o.MSPID)
		}
		qx, qy := b32(t, o.Qx), b32(t, o.Qy)
		pub, err := ecdsa.ParseUncompressedPublicKey(elliptic.P256(), append(append([]byte{4}, qx[:]...), qy[:]...))
		if err != nil {
			t.Fatal(err)
		}
		sig := OrgSignature{MSPIDHash: b32(t, o.MSPIDHash), R: b32(t, o.R), S: b32(t, o.S)}
		if !Verify(pub, d, a, sig) {
			t.Fatalf("signature of %s does not verify", o.MSPID)
		}
		a2 := a
		a2.BlockNumber++
		if Verify(pub, d, a2, sig) {
			t.Fatal("tampered attestation verified")
		}
	}
}

func TestSignIsLowSAndSorted(t *testing.T) {
	_, d, a := load(t)
	var sigs []OrgSignature
	for _, msp := range []string{"BancoMSP", "RegistradoraMSP", "AuditorMSP"} {
		k, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		for i := range 32 { // várias assinaturas: cobre os casos high-s normalizados
			sig, err := Signer{MSPID: msp, Key: k}.Sign(d, a)
			if err != nil {
				t.Fatal(err)
			}
			if !Verify(&k.PublicKey, d, a, sig) {
				t.Fatal("own signature does not verify (low-s?)")
			}
			if i == 0 {
				sigs = append(sigs, sig) // uma por org
			}
		}
	}
	SortSignatures(sigs)
	for i := 1; i < len(sigs); i++ {
		if strings.Compare(string(sigs[i-1].MSPIDHash[:]), string(sigs[i].MSPIDHash[:])) >= 0 {
			t.Fatal("not sorted")
		}
	}
}

func TestDomainSeparation(t *testing.T) {
	_, d, a := load(t)
	d2 := d
	d2.ChainID++
	d3 := d
	d3.Verifier[19] ^= 1
	if Digest(d, a) == Digest(d2, a) || Digest(d, a) == Digest(d3, a) {
		t.Fatal("digest must bind chainId and verifier")
	}
}

func TestRejectsNonP256(t *testing.T) {
	k, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if _, err := (Signer{MSPID: "X", Key: k}).Sign(Domain{}, FabricLockAttestation{}); err != ErrNotP256 {
		t.Fatalf("expected ErrNotP256, got %v", err)
	}
}
