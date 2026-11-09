package besu

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
	"time"

	"janus/node/internal/attest"
	"janus/node/internal/hashlock"
	"janus/node/internal/swap"
)

func TestAttestedSettlementWithGoAttestations(t *testing.T) {
	ctx, url, d := setup(t)
	admin, _ := Dial(ctx, url, devAdmin, d)
	buyer, _ := Dial(ctx, url, devBuyer, d)
	seller, _ := Dial(ctx, url, devSeller, d)

	// Governança: Registradora + Banco como atestadores, política 2-de-2.
	keys := map[string]*ecdsa.PrivateKey{}
	var atts []Attestor
	for _, msp := range []string{"RegistradoraMSP", "BancoMSP"} {
		k, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		keys[msp] = k
		x, y, _ := attest.PublicKeyXY(&k.PublicKey)
		atts = append(atts, Attestor{MSPID: msp, Qx: x, Qy: y})
	}
	if err := admin.RegisterAttestors(ctx, 2, atts...); err != nil {
		t.Fatal(err)
	}

	s, h, _ := hashlock.New()
	tt := terms(t, ctx, buyer, d, h, 3*time.Minute)
	tt.Protocol, tt.Networks.PaymentContract = swap.ProtocolCacti, swap.Address(d.AttestedSettlement)
	asset, recip := attest.HashString(tt.AssetID), attest.HashBytes(tt.Buyer.FSCIdentity)
	if _, err := buyer.Deposit(ctx, tt, asset, recip, tt.T1-30); err != nil {
		t.Fatal(err)
	}
	v, _ := seller.View(ctx, tt)
	if err := swap.CheckBesuLock(tt, swap.DevPolicy(), v); err != nil {
		t.Fatalf("a Vendedora deveria aceitar o depósito: %v", err)
	}

	a := attest.FabricLockAttestation{SwapID: tt.SwapID, Hashlock: h, AssetIDHash: asset, RecipientIDHash: recip,
		Deadline: uint64(tt.T1), ChannelHash: attest.HashString("recebiveis"), NamespaceHash: attest.HashString("tcc"),
		FabricTxID: [32]byte{0xfa}}
	dom := seller.AttestationDomain()
	sign := func(msp string) attest.OrgSignature {
		sig, err := attest.Signer{MSPID: msp, Key: keys[msp]}.Sign(dom, a)
		if err != nil {
			t.Fatal(err)
		}
		return sig
	}

	// Cenário 7: endossante insuficiente e prova adulterada são rejeitados.
	one := []attest.OrgSignature{sign("RegistradoraMSP")}
	if ok, _ := seller.VerifyAttestation(ctx, a, one); ok {
		t.Fatal("1 de 2 não deveria validar")
	}
	if _, err := seller.SettleAttested(ctx, tt.SwapID, s, a, one); err == nil {
		t.Fatal("settle com 1 assinatura deveria reverter")
	}
	sigs := []attest.OrgSignature{sign("RegistradoraMSP"), sign("BancoMSP")}
	attest.SortSignatures(sigs)
	tampered := a
	tampered.Deadline++
	if _, err := seller.SettleAttested(ctx, tt.SwapID, s, tampered, sigs); err == nil {
		t.Fatal("atestação adulterada deveria reverter")
	}

	if ok, err := seller.VerifyAttestation(ctx, a, sigs); err != nil || !ok {
		t.Fatalf("atestação válida rejeitada: %v", err)
	}
	if _, err := seller.SettleAttested(ctx, tt.SwapID, s, a, sigs); err != nil {
		t.Fatal(err)
	}
	v, _ = buyer.View(ctx, tt)
	if v.State != swap.BesuClaimed || v.Preimage == nil || *v.Preimage != s {
		t.Fatalf("view: %+v", v)
	}
	head, _ := buyer.Eth().BlockNumber(ctx)
	evs, err := buyer.Events(ctx, 0, head)
	if err != nil {
		t.Fatal(err)
	}
	var settled bool
	for _, e := range evs {
		if e.Contract == "settlement" && e.Kind == EvClaimed && e.Preimage != nil && *e.Preimage == s {
			settled = true
		}
	}
	if !settled {
		t.Fatalf("evento PaymentSettled com S não encontrado: %+v", evs)
	}
	got, err := buyer.SettledAttestation(ctx, tt.SwapID)
	if err != nil || got != a {
		t.Fatalf("atestação on-chain: %+v err=%v", got, err)
	}
}
