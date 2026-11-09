package it

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"janus/node/internal/attest"
	"janus/node/internal/besu"
	"janus/node/internal/besu/bindings"
	"janus/node/internal/cacti"
	"janus/node/internal/hashlock"
	"janus/node/internal/swap"
)

// TestCactiBesuSettle: o relay da Vendedora submete AttestedSettlement.settle pelo conector Besu do
// Cacti (web3 1.6.1, gasPrice 0) num anvil Osaka. Não precisa de Fabric: as atestações são assinadas
// em Go com chaves cadastradas no verificador.
func TestCactiBesuSettle(t *testing.T) {
	if err := CactiAvailable(); err != nil {
		t.Skip(err)
	}
	ctx := context.Background()
	// JANUS_BESU_RPC/JANUS_BESU_WS: usa um Besu já rodando (ex.: make besu-up); senão, anvil.
	rpc, ws := os.Getenv("JANUS_BESU_RPC"), os.Getenv("JANUS_BESU_WS")
	if rpc == "" {
		port, _ := FreePort()
		var stop func()
		var err error
		if rpc, stop, err = StartAnvil(port); err != nil {
			t.Skip(err)
		}
		t.Cleanup(stop)
		ws = strings.Replace(rpc, "http://", "ws://", 1)
	}
	dep, _, err := DeployBesu(ctx, rpc, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("besu %s chainId %d", rpc, dep.ChainID)
	rport, _ := FreePort()
	relaySrv, err := StartCactiRelay(t.TempDir(), rport, rpc, ws, KeyRelay, os.Stderr)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(relaySrv.Stop)

	dial := func(h string) *besu.Client {
		k, _ := crypto.HexToECDSA(h)
		c, err := besu.Dial(ctx, rpc, k, dep)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(c.Close)
		return c
	}
	admin, buyer := dial(KeyAdmin), dial(KeyBuyer)

	// atestadores 2-de-2 (chaves Go)
	keys := map[string]*ecdsa.PrivateKey{}
	var atts []besu.Attestor
	for _, msp := range RequiredEndorsers {
		k, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		keys[msp] = k
		x, y, _ := attest.PublicKeyXY(&k.PublicKey)
		atts = append(atts, besu.Attestor{MSPID: msp, Qx: x, Qy: y})
	}
	if err := admin.RegisterAttestors(ctx, 2, atts...); err != nil {
		t.Fatal(err)
	}

	// depósito do Banco
	s, h, _ := hashlock.New()
	hdr, _ := buyer.Eth().HeaderByNumber(ctx, nil)
	seller := swap.Party{MSPID: "VendedoraMSP", FSCIdentity: []byte("vendedora"), EVM: swap.Address(addrOf(KeySeller))}
	buyerP := swap.Party{MSPID: "BancoMSP", FSCIdentity: []byte("banco"), EVM: swap.Address(addrOf(KeyBuyer))}
	tt := &swap.SwapTerms{Protocol: swap.ProtocolCacti, AssetID: "nft-cacti", Seller: seller, Buyer: buyerP,
		Price: swap.Amount{Token: swap.Address(dep.DrexToken), Value: 777_00}, Hashlock: h,
		T2: int64(hdr.Time) + 180, T1: int64(hdr.Time) + 360,
		Networks: swap.NetworkRefs{FabricChannel: "testchannel", FabricNamespace: "tns", BesuChainID: dep.ChainID,
			PaymentContract: swap.Address(dep.AttestedSettlement)}}
	tt.SwapID = swap.ComputeSwapID(h, seller, buyerP, tt.Nonce)
	asset, recip := attest.HashString(tt.AssetID), attest.HashBytes(tt.Buyer.FSCIdentity)
	if _, err := buyer.Deposit(ctx, tt, asset, recip, tt.T1-30); err != nil {
		t.Fatal(err)
	}
	a := attest.FabricLockAttestation{SwapID: tt.SwapID, Hashlock: h, AssetIDHash: asset, RecipientIDHash: recip,
		Deadline: uint64(tt.T1), ChannelHash: attest.HashString("testchannel"), NamespaceHash: attest.HashString("tns"),
		FabricTxID: [32]byte{0xfa}, BlockNumber: 7}
	var sigs []attest.OrgSignature
	for _, msp := range RequiredEndorsers {
		sig, err := attest.Signer{MSPID: msp, Key: keys[msp]}.Sign(buyer.AttestationDomain(), a)
		if err != nil {
			t.Fatal(err)
		}
		sigs = append(sigs, sig)
	}
	attest.SortSignatures(sigs)

	// settle pelo conector Besu do Cacti, com a conta de relay
	relay := &cacti.BesuConnector{BaseURL: relaySrv.URL, KeychainID: "janus", KeychainRef: "relay",
		Account: addrOf(KeyRelay).Hex()}
	abi := bindings.AttestedSettlementMetaData.ABI
	res, err := relay.Send(ctx, abi, dep.AttestedSettlement.Hex(), "settle", cacti.SettleParams(tt.SwapID, s, a, sigs))
	if err != nil {
		t.Fatal(err)
	}
	rc, err := buyer.WaitReceipt(ctx, common.HexToHash(res.TxHash), 30*time.Second)
	if err != nil || rc.Status != 1 {
		t.Fatalf("settle via Cacti: status=%v err=%v", rc, err)
	}
	tx, _, err := buyer.Eth().TransactionByHash(ctx, common.HexToHash(res.TxHash))
	if err != nil {
		t.Fatal(err)
	}
	from, err := buyer.Sender(tx)
	if err != nil || from != addrOf(KeyRelay) {
		t.Fatalf("remetente %s != relay %s (%v)", from, addrOf(KeyRelay), err)
	}
	if tx.GasPrice().Sign() != 0 {
		t.Fatalf("gasPrice %s != 0", tx.GasPrice())
	}
	v, _ := buyer.SettlementView(ctx, tt.SwapID)
	if v.State != swap.BesuClaimed || v.Preimage == nil || *v.Preimage != s {
		t.Fatalf("depósito não liquidado: %+v", v)
	}
	got, err := buyer.SettledAttestation(ctx, tt.SwapID)
	if err != nil || got != a {
		t.Fatalf("atestação on-chain: %+v %v", got, err)
	}
	bal, _ := buyer.Balance(ctx, addrOf(KeySeller))
	if bal.Uint64() != 777_00 {
		t.Fatalf("Drex da Vendedora: %s", bal)
	}
	t.Logf("settle %s submetido pelo relay %s via Cacti; Vendedora recebeu R$ 777,00", res.TxHash[:12], from.Hex()[:10])

	// uma segunda tentativa (swap já liquidado) é rejeitada pela cadeia
	if res2, err := relay.Send(ctx, abi, dep.AttestedSettlement.Hex(), "settle", cacti.SettleParams(tt.SwapID, s, a, sigs)); err == nil {
		if rc, err := buyer.WaitReceipt(ctx, common.HexToHash(res2.TxHash), 30*time.Second); err == nil && rc.Status == 1 {
			t.Fatal("settle duplicado deveria reverter")
		}
	}
}
