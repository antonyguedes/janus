package it

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	ptoken "github.com/LFDT-Panurus/panurus/integration/token"
	"github.com/LFDT-Panurus/panurus/integration/token/common/sdk/ffabtoken"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hyperledger-labs/fabric-smart-client/integration"
	"github.com/hyperledger-labs/fabric-smart-client/integration/nwo/fabric"
	"github.com/hyperledger-labs/fabric-smart-client/integration/nwo/fsc"

	"janus/node/internal/besu"
	"janus/node/internal/engine"
	"janus/node/internal/janus"
	jt "janus/node/internal/token"
)

// Chaves PÚBLICAS de desenvolvimento do Foundry/anvil (também na allowlist de infra/besu).
const (
	KeyAdmin  = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	KeySeller = "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"
	KeyBuyer  = "5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a"
	// KeyRelay: conta de relay que submete o settle via conector Besu do Cacti (allowlist do Besu).
	KeyRelay = "47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a"
)

// Network é uma rede Janus completa: Fabric 3.x + Panurus (NWO) com os três nós FSC rodando a
// SDK do Janus. Usada pelos testes e pelo runner de demonstração (cmd/janus-net).
type Network struct {
	II    *integration.Infrastructure
	suite *ptoken.TestSuite
}

// StartNetwork gera e sobe a topologia. O chamador precisa ter registrado um fail handler do
// gomega (os testes usam RegisterTestingT, o runner usa um handler que entra em pânico).
func StartNetwork(o Options) (n *Network, err error) {
	topos := Topology(o)
	for _, tp := range topos {
		if f, ok := tp.(*fsc.Topology); ok {
			f.AddSDK(&ffabtoken.SDK{})
			f.AddSDK(&janus.SDK{}) // sem janus.json (JANUS_CONFIG_DIR) é um no-op
			f.P2PCommunicationType = fsc.WebSocket
		}
	}
	n = &Network{suite: ptoken.NewTestSuite(func() int { return 21000 }, topos)}
	defer func() {
		if r := recover(); r != nil {
			n.Stop()
			err = fmt.Errorf("setup da rede falhou: %v", r)
		}
	}()
	n.suite.Setup()
	n.II = n.suite.II
	return n, nil
}

// Stop derruba a rede (idempotente) e remove os binários compilados pelo NWO.
func (n *Network) Stop() {
	if n == nil || n.suite == nil {
		return
	}
	if n.suite.II != nil {
		n.suite.TearDown()
		n.suite.II = nil
	}
	ptoken.CleanupBuildArtifacts()
}

// FreePort devolve uma porta TCP livre em 127.0.0.1.
func FreePort() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func addrOf(hexKey string) common.Address {
	k, _ := crypto.HexToECDSA(hexKey)
	return crypto.PubkeyToAddress(k.PublicKey)
}

// StartAnvil sobe um anvil com EVM Osaka, gas zero e chainId 1337 (as mesmas condições do Besu
// QBFT local). Devolve a URL e a função de parada.
func StartAnvil(port int) (string, func(), error) {
	bin, err := exec.LookPath("anvil")
	if err != nil {
		return "", nil, fmt.Errorf("anvil não encontrado no PATH: %w", err)
	}
	cmd := exec.Command(bin, "--port", fmt.Sprint(port), "--hardfork", "osaka", "--block-time", "1",
		"--block-base-fee-per-gas", "0", "--gas-price", "0", "--chain-id", "1337", "--silent")
	if err := cmd.Start(); err != nil {
		return "", nil, err
	}
	stop := func() { _ = cmd.Process.Kill(); _ = cmd.Wait() }
	url := fmt.Sprintf("http://127.0.0.1:%d", port)
	for range 100 {
		if c, err := ethclient.Dial(url); err == nil {
			_, err = c.ChainID(context.Background())
			c.Close()
			if err == nil {
				return url, stop, nil
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	stop()
	return "", nil, fmt.Errorf("anvil não respondeu em %s", url)
}

// DeployBesu faz o deploy dos contratos (canal/namespace da topologia NWO), coloca Vendedora e
// Banco na allowlist, emite Drex para o Banco e grava o deployment em dir.
func DeployBesu(ctx context.Context, rpc, dir string) (*besu.Deployment, string, error) {
	eth, err := ethclient.DialContext(ctx, rpc)
	if err != nil {
		return nil, "", err
	}
	defer eth.Close()
	admin, _ := crypto.HexToECDSA(KeyAdmin)
	dep, err := besu.DeployAll(ctx, eth, admin, besu.DeploySpec{
		Channel: "testchannel", Namespace: "tns",
		Allow: []common.Address{addrOf(KeySeller), addrOf(KeyBuyer)},
		Mint:  map[common.Address]*big.Int{addrOf(KeyBuyer): big.NewInt(100_000_000_00)},
	})
	if err != nil {
		return nil, "", err
	}
	path := filepath.Join(dir, fmt.Sprintf("%d.local.json", dep.ChainID))
	raw, _ := json.MarshalIndent(dep, "", "  ")
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		return nil, "", err
	}
	return dep, path, nil
}

// NodeConfig parametriza a escrita dos janus.json.
type NodeConfig struct {
	Dir         string // destino (vira JANUS_CONFIG_DIR)
	RPC         string
	Deployment  string
	Policy      string                   // dev | e2e | prod
	Ports       map[string]int           // porta da API por nó (0 = livre)
	Faults      map[string]engine.Faults // falhas injetadas por nó
	AllowOrigin string                   // CORS da API (simulador web)
	// Cacti: porta do servidor Cacti de cada nó atestador (Registradora, Banco). Quando
	// presente, o atestador exige a prova do ledger via conector Fabric do Cacti (ADR-0005).
	Cacti map[string]int
	// Relay: porta do servidor Cacti com o conector Besu (relay da Vendedora). Quando > 0, o
	// settle B2 da Vendedora sai pelo conector Besu do Cacti com a conta de relay (ADR-0006).
	Relay int
}

// RequiredEndorsers é a política de endosso do Token Chaincode (PolicyJanus) vista pelo atestador.
var RequiredEndorsers = []string{"RegistradoraMSP", "BancoMSP"}

// WriteJanusConfigs grava o janus.json de cada nó e devolve as URLs das APIs.
// O chamador exporta JANUS_CONFIG_DIR=Dir antes de subir a rede (os nós herdam o ambiente).
func WriteJanusConfigs(c NodeConfig) (map[string]string, error) {
	urls := map[string]string{}
	write := func(node string, jc janus.Config) error {
		port := c.Ports[node]
		if port == 0 {
			p, err := FreePort()
			if err != nil {
				return err
			}
			port = p
		}
		jc.API.Listen = fmt.Sprintf("127.0.0.1:%d", port)
		jc.API.AllowOrigin = c.AllowOrigin
		jc.DataDir = filepath.Join(c.Dir, node+"-data")
		jc.Policy = c.Policy
		jc.Faults = c.Faults[node]
		if c.Relay > 0 && jc.Role == "vendedora" {
			jc.Relay = &janus.RelayConfig{URL: fmt.Sprintf("http://127.0.0.1:%d", c.Relay), KeychainID: "janus",
				KeychainRef: "relay", Account: addrOf(KeyRelay).Hex()}
		}
		if p, ok := c.Cacti[node]; ok {
			jc.Cacti = &janus.CactiConfig{URL: fmt.Sprintf("http://127.0.0.1:%d", p), KeychainID: "janus",
				KeychainRef: node, RequiredEndorsers: RequiredEndorsers}
		}
		raw, _ := json.MarshalIndent(jc, "", "  ")
		for _, id := range []string{node, node + ".0"} { // o NWO pode sufixar a réplica
			if err := os.WriteFile(filepath.Join(c.Dir, id+".json"), raw, 0o600); err != nil {
				return err
			}
		}
		urls[node] = "http://" + jc.API.Listen
		return nil
	}
	besuCfg := func(key string) *janus.BesuConfig {
		return &janus.BesuConfig{RPC: c.RPC, KeyHex: key, Deployment: c.Deployment}
	}
	for node, jc := range map[string]janus.Config{
		jt.AuditorNode: {Role: "registradora", MSPID: "RegistradoraMSP", Channel: "testchannel", Namespace: "tns"},
		"vendedora": {Role: "vendedora", MSPID: "VendedoraMSP", Counterparty: "banco",
			Attestors: []string{jt.AuditorNode, "banco"}, Besu: besuCfg(KeySeller)},
		"banco": {Role: "banco", MSPID: "BancoMSP", Besu: besuCfg(KeyBuyer)},
	} {
		if err := write(node, jc); err != nil {
			return nil, err
		}
	}
	return urls, nil
}

// CactiIdentity é a credencial Fabric de uma org para o conector do Cacti (keychain).
type CactiIdentity struct {
	MSPID       string `json:"mspId"`
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"privateKey"`
}

// FabricConnectionProfile devolve o connection profile da rede (formato fabric-network/Cacti),
// com os endpoints em 127.0.0.1, e a identidade Admin de cada org (MSP ID -> credencial).
func (n *Network) FabricConnectionProfile() (map[string]any, map[string]CactiIdentity, error) {
	var fp *fabric.Platform
	for _, p := range n.II.NWO.Platforms {
		if f, ok := p.(*fabric.Platform); ok {
			fp = f
			break
		}
	}
	if fp == nil {
		return nil, nil, fmt.Errorf("plataforma Fabric não encontrada")
	}
	cp := fp.ConnectionProfile("janus", false)
	raw, err := json.Marshal(cp)
	if err != nil {
		return nil, nil, err
	}
	// o NWO aponta para o host "fabric" (uso em containers); o conector roda no host
	raw = bytes.ReplaceAll(raw, []byte("grpcs://fabric:"), []byte("grpcs://127.0.0.1:"))
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, nil, err
	}
	ids := map[string]CactiIdentity{}
	for _, org := range cp.Organizations {
		cert, _ := org.SignedCert["pem"].(string)
		key, _ := org.AdminPrivateKey["pem"].(string)
		ids[org.MSPID] = CactiIdentity{MSPID: org.MSPID, Certificate: cert, PrivateKey: key}
	}
	return out, ids, nil
}

// CactiServer é um servidor Hyperledger Cacti (infra/cacti/relay) de uma org atestadora.
type CactiServer struct {
	Node string
	URL  string
	cmd  *exec.Cmd
}

// Stop encerra o servidor (SIGTERM, depois SIGKILL).
func (c *CactiServer) Stop() {
	if c == nil || c.cmd == nil || c.cmd.Process == nil {
		return
	}
	_ = c.cmd.Process.Signal(syscall.SIGTERM)
	done := make(chan struct{})
	go func() { _ = c.cmd.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		_ = c.cmd.Process.Kill()
		<-done
	}
}

// CactiServerScript é o caminho de infra/cacti/relay/src/server.mjs.
func CactiServerScript() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "infra", "cacti", "relay", "src", "server.mjs")
}

// CactiAvailable informa se Node e as dependências do Cacti (npm install na raiz) estão presentes.
func CactiAvailable() error {
	if _, err := exec.LookPath("node"); err != nil {
		return fmt.Errorf("node não encontrado: %w", err)
	}
	mod := filepath.Join(filepath.Dir(CactiServerScript()), "..", "..", "..", "..", "node_modules",
		"@hyperledger-cacti", "cactus-plugin-ledger-connector-fabric")
	if _, err := os.Stat(mod); err != nil {
		return fmt.Errorf("dependências do Cacti ausentes (rode npm install na raiz): %w", err)
	}
	return nil
}

// StartCacti sobe o servidor Cacti de um nó atestador (conector Fabric) com a identidade da org.
func StartCacti(dir, node string, port int, ccp map[string]any, id CactiIdentity, logOut io.Writer) (*CactiServer, error) {
	ccpPath := filepath.Join(dir, "ccp.json")
	raw, _ := json.Marshal(ccp)
	if err := os.WriteFile(ccpPath, raw, 0o600); err != nil {
		return nil, err
	}
	return StartCactiConfig(dir, node, port, map[string]any{
		"fabric": map[string]any{"connectionProfile": ccpPath, "keychainRef": node, "identity": id},
	}, logOut)
}

// StartCactiRelay sobe o relay da Vendedora (conector Besu) com a chave da conta de relay.
func StartCactiRelay(dir string, port int, rpcHTTP, rpcWS, relayKey string, logOut io.Writer) (*CactiServer, error) {
	return StartCactiConfig(dir, "relay", port, map[string]any{
		"besu": map[string]any{"rpcApiHttpHost": rpcHTTP, "rpcApiWsHost": rpcWS, "keys": map[string]string{"relay": relayKey}},
	}, logOut)
}

// StartCactiConfig grava a configuração (seções "fabric" e/ou "besu") e espera o servidor ficar pronto.
func StartCactiConfig(dir, name string, port int, sections map[string]any, logOut io.Writer) (*CactiServer, error) {
	cfg := map[string]any{"apiHost": "127.0.0.1", "apiPort": port, "keychainId": "janus"}
	for k, v := range sections {
		cfg[k] = v
	}
	cfgPath := filepath.Join(dir, "cacti-"+name+".json")
	raw, _ := json.Marshal(cfg)
	if err := os.WriteFile(cfgPath, raw, 0o600); err != nil {
		return nil, err
	}
	cmd := exec.Command("node", CactiServerScript(), cfgPath)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if logOut != nil {
		cmd.Stderr = logOut
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	srv := &CactiServer{Node: name, URL: fmt.Sprintf("http://127.0.0.1:%d", port), cmd: cmd}
	ready := make(chan struct{})
	go func() {
		sc := bufio.NewScanner(out)
		signaled := false
		for sc.Scan() {
			if logOut != nil {
				fmt.Fprintln(logOut, sc.Text())
			}
			if !signaled && strings.Contains(sc.Text(), "janus-cacti pronto") {
				close(ready)
				signaled = true
			}
		}
	}()
	select {
	case <-ready:
		return srv, nil
	case <-time.After(90 * time.Second):
		srv.Stop()
		return nil, fmt.Errorf("servidor Cacti %s não ficou pronto", name)
	}
}
