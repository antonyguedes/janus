// Command vectors gera os vetores de teste compartilhados entre `go test` e `forge test`:
//
//	testdata/vectors/htlc.json         pares (S, H = SHA256(S))
//	testdata/vectors/attestation.json  digest + assinaturas P-256 de uma FabricLockAttestation
//
// As chaves P-256 são derivadas de escalares fixos e servem SÓ para teste.
// Uso: go run ./cmd/vectors -out ../testdata/vectors
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"janus/node/internal/attest"
	"janus/node/internal/hashlock"
)

func hx(b []byte) string { return "0x" + hex.EncodeToString(b) }

// VerifierAddress e ChainID usados no vetor de atestação. O teste Foundry faz o deploy do
// verificador exatamente neste endereço (deployCodeTo) e fixa o chainId com vm.chainId.
const (
	ChainID         = 31337
	VerifierAddress = "0x00000000000000000000000000000000000a77e5"
)

type htlcVector struct {
	Preimage string `json:"preimage"`
	Hashlock string `json:"hashlock"`
}

type orgVector struct {
	MSPID     string `json:"mspId"`
	MSPIDHash string `json:"mspIdHash"`
	Qx        string `json:"qx"`
	Qy        string `json:"qy"`
	R         string `json:"r"`
	S         string `json:"s"`
}

type attVector struct {
	ChainID         uint64      `json:"chainId"`
	Verifier        string      `json:"verifier"`
	DomainTag       string      `json:"domainTag"`
	SwapID          string      `json:"swapId"`
	Hashlock        string      `json:"hashlock"`
	Preimage        string      `json:"preimage"`
	AssetID         string      `json:"assetId"`
	AssetIDHash     string      `json:"assetIdHash"`
	RecipientIDHash string      `json:"recipientIdHash"`
	Deadline        uint64      `json:"deadline"`
	Channel         string      `json:"channel"`
	ChannelHash     string      `json:"channelHash"`
	Namespace       string      `json:"namespace"`
	NamespaceHash   string      `json:"namespaceHash"`
	FabricTxID      string      `json:"fabricTxId"`
	BlockNumber     uint64      `json:"blockNumber"`
	Digest          string      `json:"digest"`
	Orgs            []orgVector `json:"orgs"` // ordenadas por mspIdHash
}

// TestKey deriva uma chave P-256 determinística de um rótulo (SÓ para vetores).
func TestKey(label string) *ecdsa.PrivateKey {
	seed := sha256.Sum256([]byte("janus-test-key:" + label))
	priv, err := ecdsa.ParseRawPrivateKey(elliptic.P256(), seed[:])
	if err != nil {
		log.Fatal(err)
	}
	return priv
}

func main() {
	out := flag.String("out", "../testdata/vectors", "diretório de saída")
	flag.Parse()
	if err := os.MkdirAll(*out, 0o755); err != nil {
		log.Fatal(err)
	}

	// ---------------------------------------------------------------- HTLC
	var preimages []hashlock.Preimage
	var zero, ones, seq hashlock.Preimage
	for i := range ones {
		ones[i] = 0xff
		seq[i] = byte(i)
	}
	preimages = append(preimages, zero, ones, seq)
	for i := range 5 {
		preimages = append(preimages, hashlock.Preimage(sha256.Sum256(fmt.Appendf(nil, "janus-preimage-%d", i))))
	}
	var hv []htlcVector
	for _, s := range preimages {
		h := s.Hash()
		hv = append(hv, htlcVector{Preimage: s.Hex(), Hashlock: h.Hex()})
	}
	write(filepath.Join(*out, "htlc.json"), map[string]any{
		"description": "S de 32 bytes; H = SHA256(S) (Go crypto/sha256 == Solidity sha256(abi.encodePacked(bytes32)) == Panurus HTLC SHA256/None)",
		"vectors":     hv,
	})

	// ---------------------------------------------------------------- Atestação
	verifier, err := hashlock.ParseHash32("0x" + fmt.Sprintf("%064s", VerifierAddress[2:]))
	if err != nil {
		log.Fatal(err)
	}
	var vaddr attest.Address
	copy(vaddr[:], verifier[12:])
	dom := attest.Domain{ChainID: ChainID, Verifier: vaddr}

	s := preimages[3]
	h := s.Hash()
	assetID := "recebivel-0001"
	channel, namespace := "recebiveis", "tcc"
	txid := sha256.Sum256([]byte("fabric-tx-lock-0001"))
	a := attest.FabricLockAttestation{
		SwapID:          sha256.Sum256([]byte("swap-0001")),
		Hashlock:        h,
		AssetIDHash:     attest.HashString(assetID),
		RecipientIDHash: attest.HashString("BancoMSP:banco-owner"),
		Deadline:        1_800_000_000,
		ChannelHash:     attest.HashString(channel),
		NamespaceHash:   attest.HashString(namespace),
		FabricTxID:      txid,
		BlockNumber:     42,
	}
	digest := attest.Digest(dom, a)

	var orgs []orgVector
	var sigs []attest.OrgSignature
	keys := map[[32]byte]*ecdsa.PrivateKey{}
	for _, msp := range []string{"RegistradoraMSP", "BancoMSP"} {
		k := TestKey(msp)
		sig, err := attest.Signer{MSPID: msp, Key: k}.Sign(dom, a)
		if err != nil {
			log.Fatal(err)
		}
		if !attest.Verify(&k.PublicKey, dom, a, sig) {
			log.Fatalf("self-verify failed for %s", msp)
		}
		sigs = append(sigs, sig)
		keys[sig.MSPIDHash] = k
	}
	attest.SortSignatures(sigs)
	mspByHash := map[[32]byte]string{}
	for _, msp := range []string{"RegistradoraMSP", "BancoMSP"} {
		mspByHash[attest.MSPIDHash(msp)] = msp
	}
	for _, sig := range sigs {
		qx, qy, _ := attest.PublicKeyXY(&keys[sig.MSPIDHash].PublicKey)
		orgs = append(orgs, orgVector{
			MSPID: mspByHash[sig.MSPIDHash], MSPIDHash: hx(sig.MSPIDHash[:]),
			Qx: hx(qx[:]), Qy: hx(qy[:]), R: hx(sig.R[:]), S: hx(sig.S[:]),
		})
	}
	write(filepath.Join(*out, "attestation.json"), attVector{
		ChainID: ChainID, Verifier: VerifierAddress, DomainTag: hx(attest.DomainTag[:]),
		SwapID: hx(a.SwapID[:]), Hashlock: hx(a.Hashlock[:]), Preimage: s.Hex(),
		AssetID: assetID, AssetIDHash: hx(a.AssetIDHash[:]), RecipientIDHash: hx(a.RecipientIDHash[:]),
		Deadline: a.Deadline, Channel: channel, ChannelHash: hx(a.ChannelHash[:]),
		Namespace: namespace, NamespaceHash: hx(a.NamespaceHash[:]),
		FabricTxID: hx(a.FabricTxID[:]), BlockNumber: a.BlockNumber, Digest: hx(digest[:]), Orgs: orgs,
	})
	fmt.Println("vetores gerados em", *out)
}

func write(path string, v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(path, append(b, '\n'), 0o644); err != nil {
		log.Fatal(err)
	}
}
