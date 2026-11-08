package it

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"janus/node/internal/api"
	"janus/node/internal/besu"
	"janus/node/internal/janus"
	"janus/node/internal/swap"
	jt "janus/node/internal/token"
)

// Chaves PÚBLICAS de desenvolvimento do anvil.
const (
	keyAdmin  = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	keySeller = "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"
	keyBuyer  = "5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a"
)

func freePort(t *testing.T) int {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

// besuForE2E sobe um anvil (EVM Osaka, gas zero, chainId 1337) e faz o deploy dos contratos.
func besuForE2E(t *testing.T) (rpc string, depPath string, dep *besu.Deployment) {
	t.Helper()
	bin, err := exec.LookPath("anvil")
	if err != nil {
		t.Skip("anvil não encontrado")
	}
	port := freePort(t)
	cmd := exec.Command(bin, "--port", fmt.Sprint(port), "--hardfork", "osaka", "--block-time", "1",
		"--block-base-fee-per-gas", "0", "--gas-price", "0", "--chain-id", "1337", "--silent")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill(); _ = cmd.Wait() })
	rpc = fmt.Sprintf("http://127.0.0.1:%d", port)
	var eth *ethclient.Client
	for range 100 {
		if eth, err = ethclient.Dial(rpc); err == nil {
			if _, err = eth.ChainID(context.Background()); err == nil {
				break
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	if err != nil {
		t.Fatal(err)
	}
	defer eth.Close()
	addr := func(h string) common.Address {
		k, _ := crypto.HexToECDSA(h)
		return crypto.PubkeyToAddress(k.PublicKey)
	}
	admin, _ := crypto.HexToECDSA(keyAdmin)
	dep, err = besu.DeployAll(context.Background(), eth, admin, besu.DeploySpec{
		Channel: "testchannel", Namespace: "tns",
		Allow: []common.Address{addr(keySeller), addr(keyBuyer)},
		Mint:  map[common.Address]*big.Int{addr(keyBuyer): big.NewInt(100_000_000_00)},
	})
	if err != nil {
		t.Fatal(err)
	}
	depPath = filepath.Join(t.TempDir(), "1337.local.json")
	raw, _ := json.Marshal(dep)
	if err := os.WriteFile(depPath, raw, 0o600); err != nil {
		t.Fatal(err)
	}
	return rpc, depPath, dep
}

func writeJanusConfigs(t *testing.T, rpc, depPath string, faults map[string]janus.Config) map[string]string {
	dir := t.TempDir()
	t.Setenv("JANUS_CONFIG_DIR", dir)
	urls := map[string]string{}
	write := func(node string, c janus.Config) {
		port := freePort(t)
		c.API.Listen = fmt.Sprintf("127.0.0.1:%d", port)
		c.DataDir = filepath.Join(dir, node+"-data")
		c.Policy = "dev"
		if f, ok := faults[node]; ok {
			c.Faults = f.Faults
		}
		raw, _ := json.MarshalIndent(c, "", "  ")
		for _, id := range []string{node, node + ".0"} {
			if err := os.WriteFile(filepath.Join(dir, id+".json"), raw, 0o600); err != nil {
				t.Fatal(err)
			}
		}
		urls[node] = "http://" + c.API.Listen
	}
	write(jt.AuditorNode, janus.Config{Role: "registradora", MSPID: "RegistradoraMSP"})
	write("vendedora", janus.Config{Role: "vendedora", MSPID: "VendedoraMSP", Counterparty: "banco",
		Besu: &janus.BesuConfig{RPC: rpc, KeyHex: keySeller, Deployment: depPath}})
	write("banco", janus.Config{Role: "banco", MSPID: "BancoMSP",
		Besu: &janus.BesuConfig{RPC: rpc, KeyHex: keyBuyer, Deployment: depPath}})
	return urls
}

func post(t *testing.T, url string, body, out any) int {
	t.Helper()
	raw, _ := json.Marshal(body)
	res, err := http.Post(url, "application/json", bytes.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if out != nil {
		_ = json.NewDecoder(res.Body).Decode(out)
	}
	return res.StatusCode
}

func getSwap(url, id string) (api.SwapDTO, error) {
	var d api.SwapDTO
	res, err := http.Get(url + "/api/v1/swaps/" + id)
	if err != nil {
		return d, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return d, fmt.Errorf("status %d", res.StatusCode)
	}
	return d, json.NewDecoder(res.Body).Decode(&d)
}

func waitAPI(t *testing.T, url string) {
	t.Helper()
	deadline := time.Now().Add(90 * time.Second)
	for time.Now().Before(deadline) {
		if res, err := http.Get(url + "/api/v1/health"); err == nil {
			res.Body.Close()
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatalf("api %s não subiu", url)
}

func waitStatus(t *testing.T, url, id string, want swap.Status, timeout time.Duration) api.SwapDTO {
	t.Helper()
	var last api.SwapDTO
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if d, err := getSwap(url, id); err == nil {
			last = d
			if d.Status == want {
				return d
			}
		}
		time.Sleep(time.Second)
	}
	t.Fatalf("%s: swap %s status=%s, esperado %s", url, id[:12], last.Status, want)
	return last
}

// TestE2EHTLC: Abordagem A ponta a ponta com os nós Janus reais (SDK do Janus dentro dos nós
// FSC), Fabric 3.1.5 + Panurus via NWO e Besu (anvil, EVM Osaka). Tudo dirigido pelas APIs REST,
// como o simulador faria.
func TestE2EHTLC(t *testing.T) {
	if os.Getenv("FAB_BINS") == "" {
		t.Skip("FAB_BINS não definido")
	}
	rpc, depPath, dep := besuForE2E(t)
	urls := writeJanusConfigs(t, rpc, depPath, nil)
	ii := startNetwork(t, Options{})
	_ = ii
	for _, u := range urls {
		waitAPI(t, u)
	}

	// 1) Registradora emite o recebível para a Vendedora
	var issued map[string]string
	var code int
	for i := 0; i < 30; i++ { // o auditor é registrado em background
		if code = post(t, urls[jt.AuditorNode]+"/api/v1/recebiveis", api.IssueRequest{Owner: "vendedora",
			Registro: "REG-E2E-1", Cedente: "Vendedora S.A.", Sacado: "Sacado Ltda", ValorFace: 1_250_000_00,
			Vencimento: "2027-03-31"}, &issued); code == http.StatusCreated {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if code != http.StatusCreated || issued["linearId"] == "" {
		t.Fatalf("emissão: %d %v", code, issued)
	}

	// 2) Vendedora propõe o swap ao Banco (termos assinados por sessão P2P do FSC)
	var proposed api.SwapDTO
	if code := post(t, urls["vendedora"]+"/api/v1/swaps", api.ProposeRequest{BuyerNode: "banco",
		AssetID: issued["linearId"], Price: 1_200_000_00, Protocol: swap.ProtocolHTLC}, &proposed); code != http.StatusCreated {
		t.Fatalf("proposta: %d %+v", code, proposed)
	}
	if proposed.Terms == nil || proposed.Terms.T1-proposed.Terms.T2 < 90 {
		t.Fatalf("termos: %+v", proposed)
	}

	// 3) Os motores dos dois nós levam o swap a "completed" sozinhos
	waitStatus(t, urls["banco"], proposed.SwapID, swap.StatusCompleted, 2*time.Minute)
	waitStatus(t, urls["vendedora"], proposed.SwapID, swap.StatusCompleted, time.Minute)

	// 4) Conferência nas duas cadeias
	k, _ := crypto.HexToECDSA(keySeller)
	c, err := besu.Dial(context.Background(), rpc, k, dep)
	if err != nil {
		t.Fatal(err)
	}
	bal, _ := c.Balance(context.Background(), c.From())
	if bal.Uint64() != 1_200_000_00 {
		t.Fatalf("saldo Drex da Vendedora: %s", bal)
	}
	ok(t, ii, "banco", "get", jt.Get{LinearID: issued["linearId"]})
	if _, err := call(t, ii, "vendedora", "get", jt.Get{LinearID: issued["linearId"]}); err == nil {
		t.Fatal("o NFT deveria estar com o Banco")
	}
}
