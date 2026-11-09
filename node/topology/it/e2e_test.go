package it

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hyperledger-labs/fabric-smart-client/integration"

	"janus/node/internal/api"
	"janus/node/internal/besu"
	"janus/node/internal/cacti"
	"janus/node/internal/engine"
	"janus/node/internal/swap"
	jt "janus/node/internal/token"
)

const (
	keyAdmin  = KeyAdmin
	keySeller = KeySeller
	keyBuyer  = KeyBuyer
)

// besuForE2E sobe um anvil (EVM Osaka, gas zero, chainId 1337) e faz o deploy dos contratos.
func besuForE2E(t *testing.T) (rpc string, depPath string, dep *besu.Deployment) {
	t.Helper()
	if _, err := exec.LookPath("anvil"); err != nil {
		t.Skip("anvil não encontrado")
	}
	port, err := FreePort()
	if err != nil {
		t.Fatal(err)
	}
	rpc, stop, err := StartAnvil(port)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(stop)
	dep, depPath, err = DeployBesu(context.Background(), rpc, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return rpc, depPath, dep
}

func writeJanusConfigs(t *testing.T, rpc, depPath, policy string, faults map[string]engine.Faults) map[string]string {
	dir := t.TempDir()
	t.Setenv("JANUS_CONFIG_DIR", dir)
	urls, err := WriteJanusConfigs(NodeConfig{Dir: dir, RPC: rpc, Deployment: depPath, Policy: policy, Faults: faults})
	if err != nil {
		t.Fatal(err)
	}
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

func getInfo(t *testing.T, url string) api.Info {
	t.Helper()
	res, err := http.Get(url + "/api/v1/info")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var in api.Info
	if err := json.NewDecoder(res.Body).Decode(&in); err != nil {
		t.Fatal(err)
	}
	return in
}

func hex32(t *testing.T, s string) [32]byte {
	var out [32]byte
	b := common.FromHex(s)
	if len(b) != 32 {
		t.Fatalf("bad hex32 %q", s)
	}
	copy(out[:], b)
	return out
}

type e2eOpts struct {
	proto    swap.Protocol
	registro string
	policy   string // dev | e2e
	faults   map[string]engine.Faults
	cacti    bool // atestadores provam o lock via conector Fabric do Cacti (ADR-0005)
	setup    func(urls map[string]string, rpc string, dep *besu.Deployment)
}

type e2eEnv struct {
	ii    *integration.Infrastructure
	urls  map[string]string
	rpc   string
	dep   *besu.Deployment
	cacti map[string]*CactiServer
}

func bootE2E(t *testing.T, o e2eOpts) *e2eEnv {
	if os.Getenv("FAB_BINS") == "" {
		t.Skip("FAB_BINS não definido")
	}
	if o.policy == "" {
		o.policy = "dev"
	}
	var cactiPorts map[string]int
	if o.cacti {
		if err := CactiAvailable(); err != nil {
			t.Skip(err)
		}
		cactiPorts = map[string]int{}
		for _, node := range []string{jt.AuditorNode, "banco"} {
			p, err := FreePort()
			if err != nil {
				t.Fatal(err)
			}
			cactiPorts[node] = p
		}
	}
	rpc, depPath, dep := besuForE2E(t)
	dir := t.TempDir()
	t.Setenv("JANUS_CONFIG_DIR", dir)
	urls, err := WriteJanusConfigs(NodeConfig{Dir: dir, RPC: rpc, Deployment: depPath, Policy: o.policy,
		Faults: o.faults, Cacti: cactiPorts})
	if err != nil {
		t.Fatal(err)
	}
	n, err := startNetworkN(t, Options{})
	if err != nil {
		t.Fatal(err)
	}
	env := &e2eEnv{ii: n.II, urls: urls, rpc: rpc, dep: dep, cacti: map[string]*CactiServer{}}
	if o.cacti {
		ccp, ids, err := n.FabricConnectionProfile()
		if err != nil {
			t.Fatal(err)
		}
		msp := map[string]string{jt.AuditorNode: "RegistradoraMSP", "banco": "BancoMSP"}
		for node, port := range cactiPorts {
			srv, err := StartCacti(t.TempDir(), node, port, ccp, ids[msp[node]], nil)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(srv.Stop)
			env.cacti[node] = srv
		}
	}
	for _, u := range urls {
		waitAPI(t, u)
	}
	if o.setup != nil {
		o.setup(urls, rpc, dep)
	}
	return env
}

// issueAndPropose: a Registradora emite o recebível e a Vendedora propõe o swap ao Banco.
func issueAndPropose(t *testing.T, env *e2eEnv, o e2eOpts) (string, api.SwapDTO) {
	var issued map[string]string
	var code int
	for i := 0; i < 30; i++ { // o auditor é registrado em background
		if code = post(t, env.urls[jt.AuditorNode]+"/api/v1/recebiveis", api.IssueRequest{Owner: "vendedora",
			Registro: o.registro, Cedente: "Vendedora S.A.", Sacado: "Sacado Ltda", ValorFace: 1_250_000_00,
			Vencimento: "2027-03-31"}, &issued); code == http.StatusCreated {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if code != http.StatusCreated || issued["linearId"] == "" {
		t.Fatalf("emissão: %d %v", code, issued)
	}
	var proposed api.SwapDTO
	if code := post(t, env.urls["vendedora"]+"/api/v1/swaps", api.ProposeRequest{BuyerNode: "banco",
		AssetID: issued["linearId"], Price: 1_200_000_00, Protocol: o.proto}, &proposed); code != http.StatusCreated {
		t.Fatalf("proposta: %d %+v", code, proposed)
	}
	if proposed.Terms == nil || proposed.Terms.T1-proposed.Terms.T2 < 90 || proposed.Protocol != o.proto {
		t.Fatalf("termos: %+v", proposed)
	}
	return issued["linearId"], proposed
}

func drexOf(t *testing.T, env *e2eEnv, key string) uint64 {
	k, _ := crypto.HexToECDSA(key)
	c, err := besu.Dial(context.Background(), env.rpc, k, env.dep)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	bal, err := c.Balance(context.Background(), c.From())
	if err != nil {
		t.Fatal(err)
	}
	return bal.Uint64()
}

func nftWith(t *testing.T, env *e2eEnv, owner, other, linearID string) {
	t.Helper()
	ok(t, env.ii, owner, "get", jt.Get{LinearID: linearID})
	if _, err := call(t, env.ii, other, "get", jt.Get{LinearID: linearID}); err == nil {
		t.Fatalf("o NFT não deveria estar com %s", other)
	}
}

// e2e executa o caminho feliz completo pelas APIs e confere as duas cadeias.
func e2e(t *testing.T, o e2eOpts) (*e2eEnv, api.SwapDTO) {
	env := bootE2E(t, o)
	id, proposed := issueAndPropose(t, env, o)
	waitStatus(t, env.urls["banco"], proposed.SwapID, swap.StatusCompleted, 2*time.Minute)
	waitStatus(t, env.urls["vendedora"], proposed.SwapID, swap.StatusCompleted, time.Minute)
	if got := drexOf(t, env, keySeller); got != 1_200_000_00 {
		t.Fatalf("saldo Drex da Vendedora: %d", got)
	}
	nftWith(t, env, "banco", "vendedora", id)
	return env, proposed
}

// TestE2EHTLC: Abordagem A ponta a ponta com os nós Janus reais (SDK do Janus dentro dos nós
// FSC), Fabric 3.1.5 + Panurus via NWO e Besu (anvil, EVM Osaka). Tudo dirigido pelas APIs REST,
// como o simulador faria.
func TestE2EHTLC(t *testing.T) { e2e(t, e2eOpts{proto: swap.ProtocolHTLC, registro: "REG-E2E-A"}) }

// TestE2EAttested: Abordagem B2 ponta a ponta. O Banco deposita antes do lock, a Vendedora coleta
// as atestações da Registradora (auditoria) e do Banco (vault) por sessões FSC, e o
// AttestedSettlement verifica a política 2-de-2 com P-256 nativo antes de liberar o Drex.
func TestE2EAttested(t *testing.T) {
	env, proposed := e2e(t, e2eOpts{proto: swap.ProtocolCacti, registro: "REG-E2E-B", cacti: true, setup: func(urls map[string]string, rpc string, dep *besu.Deployment) {
		admin, _ := crypto.HexToECDSA(keyAdmin)
		c, err := besu.Dial(context.Background(), rpc, admin, dep)
		if err != nil {
			t.Fatal(err)
		}
		var atts []besu.Attestor
		for _, node := range []string{jt.AuditorNode, "banco"} {
			in := getInfo(t, urls[node])
			if in.Attestor == nil {
				t.Fatalf("%s sem chave de atestador", node)
			}
			atts = append(atts, besu.Attestor{MSPID: in.Attestor.MSPID, Qx: hex32(t, in.Attestor.Qx), Qy: hex32(t, in.Attestor.Qy)})
		}
		if err := c.RegisterAttestors(context.Background(), 2, atts...); err != nil {
			t.Fatal(err)
		}
	}})

	// A atestação aceita on-chain traz o bloco que o conector Fabric do Cacti reporta para a
	// transação de lock: a prova veio do ledger comitado, não da memória dos nós.
	k, _ := crypto.HexToECDSA(keyBuyer)
	bc, err := besu.Dial(context.Background(), env.rpc, k, env.dep)
	if err != nil {
		t.Fatal(err)
	}
	defer bc.Close()
	id := hex32(t, "0x"+proposed.SwapID)
	att, err := bc.SettledAttestation(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}
	if att.BlockNumber == 0 {
		t.Fatal("atestação sem número de bloco: a prova via Cacti não foi usada")
	}
	srv := env.cacti["banco"]
	conn := &cacti.FabricConnector{BaseURL: srv.URL, KeychainID: "janus", KeychainRef: "banco", Channel: "testchannel"}
	proof, err := conn.TxProof(context.Background(), hex.EncodeToString(att.FabricTxID[:]))
	if err != nil {
		t.Fatal(err)
	}
	if proof.BlockNumber != att.BlockNumber || proof.Satisfies("tns", RequiredEndorsers) != nil {
		t.Fatalf("prova Cacti %+v != atestação on-chain (bloco %d)", proof, att.BlockNumber)
	}
	t.Logf("lock %x comitado no bloco %d, endossado por %v; atestação 2-de-2 aceita on-chain",
		att.FabricTxID[:6], proof.BlockNumber, proof.Endorsers)
}

// TestE2ESellerNoReveal (cenário 3, abordagem A): a Vendedora nunca revela S (falha injetada).
// Os motores conduzem as duas reversões em rede real: refund no Besu depois de T2 (Banco) e
// reclaim no Fabric depois de T1 + tolerância da auditoria (Vendedora).
func TestE2ESellerNoReveal(t *testing.T) {
	o := e2eOpts{proto: swap.ProtocolHTLC, registro: "REG-E2E-3", policy: "e2e",
		faults: map[string]engine.Faults{"vendedora": {NoReveal: true}}}
	env := bootE2E(t, o)
	before := drexOf(t, env, keyBuyer)
	id, proposed := issueAndPropose(t, env, o)
	waitStatus(t, env.urls["vendedora"], proposed.SwapID, swap.StatusBothLocked, time.Minute)
	waitStatus(t, env.urls["vendedora"], proposed.SwapID, swap.StatusReverted, 5*time.Minute)
	if got := drexOf(t, env, keyBuyer); got != before {
		t.Fatalf("o Banco deveria ter reavido o Drex: %d != %d", got, before)
	}
	if got := drexOf(t, env, keySeller); got != 0 {
		t.Fatalf("a Vendedora não deveria ter recebido: %d", got)
	}
	nftWith(t, env, "vendedora", "banco", id)
}

// TestE2EBuyerCrashRecover (cenário 5): o nó do Banco cai com as duas pernas travadas, a
// Vendedora revela S com o Banco fora do ar, e o Banco volta. A reconciliação de boot (store
// SQLite + backfill do watcher a partir do checkpoint) faz o claim do NFT antes de T1.
func TestE2EBuyerCrashRecover(t *testing.T) {
	o := e2eOpts{proto: swap.ProtocolHTLC, registro: "REG-E2E-5", policy: "e2e"}
	env := bootE2E(t, o)
	id, proposed := issueAndPropose(t, env, o)
	waitStatus(t, env.urls["banco"], proposed.SwapID, swap.StatusBothLocked, time.Minute)
	env.ii.StopFSCNode("banco")
	waitStatus(t, env.urls["vendedora"], proposed.SwapID, swap.StatusRevealed, time.Minute)
	time.Sleep(3 * time.Second)
	env.ii.StartFSCNode("banco")
	waitAPI(t, env.urls["banco"])
	waitStatus(t, env.urls["banco"], proposed.SwapID, swap.StatusCompleted, 2*time.Minute)
	waitStatus(t, env.urls["vendedora"], proposed.SwapID, swap.StatusCompleted, time.Minute)
	nftWith(t, env, "banco", "vendedora", id)
}
