// Command janus-net sobe um ambiente Janus completo para demonstração e uso do simulador web em
// modo Live: Fabric 3.x + Panurus (NWO), Besu (anvil Osaka ou um Besu QBFT já rodando), os
// contratos, e os três nós Janus (Registradora, Vendedora, Banco) com as APIs em portas fixas.
// A governança B2 (atestadores 2-de-2) é cadastrada automaticamente.
//
//	FAB_BINS=infra/fabric/bin go run ./cmd/janus-net            # anvil
//	FAB_BINS=infra/fabric/bin go run ./cmd/janus-net -besu http://localhost:8545
//
// Ctrl-C derruba tudo.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/onsi/gomega"

	"janus/node/internal/api"
	"janus/node/internal/besu"
	jt "janus/node/internal/token"
	"janus/node/topology/it"
)

func main() {
	besuRPC := flag.String("besu", "anvil", `"anvil" para subir um anvil, ou a URL de um Besu QBFT (make besu-up)`)
	besuWS := flag.String("besu-ws", "ws://localhost:8546", "URL WebSocket do Besu (relay Cacti); ignorado com anvil")
	policy := flag.String("policy", "dev", "política de tempo: dev (T2=3min, T1=6min) | e2e (T2=75s, T1=170s)")
	dir := flag.String("dir", ".janus/net", "diretório de estado (configs, SQLite, segredos)")
	origin := flag.String("origin", "http://localhost:5173", "origem do simulador web (CORS)")
	useCacti := flag.Bool("cacti", true, "atestadores provam o lock via conector Fabric do Hyperledger Cacti (ADR-0005)")
	flag.Parse()

	if os.Getenv("FAB_BINS") == "" {
		log.Fatal("defina FAB_BINS com os binários do Fabric (make fabric-bins)")
	}
	gomega.RegisterFailHandler(func(msg string, _ ...int) { panic(msg) })
	abs, err := filepath.Abs(*dir)
	must(err)
	must(os.RemoveAll(abs)) // estado sempre novo: a rede NWO também é nova a cada execução
	must(os.MkdirAll(abs, 0o700))
	ctx := context.Background()

	// 1) Besu
	rpc := *besuRPC
	if rpc == "anvil" {
		port, err := it.FreePort()
		must(err)
		var stop func()
		rpc, stop, err = it.StartAnvil(port)
		must(err)
		defer stop()
		*besuWS = strings.Replace(rpc, "http://", "ws://", 1)
	}
	dep, depPath, err := it.DeployBesu(ctx, rpc, abs)
	must(err)
	log.Printf("besu: %s (chainId %d) | HTLCPayment %s | AttestedSettlement %s", rpc, dep.ChainID, dep.HTLCPayment, dep.AttestedSettlement)

	// 2) Configuração dos nós Janus (APIs nas portas padrão do simulador)
	cactiPorts := map[string]int{jt.AuditorNode: 4100, "banco": 4101}
	const relayPort = 4102
	if *useCacti {
		if err := it.CactiAvailable(); err != nil {
			log.Printf("Cacti desabilitado: %v", err)
			*useCacti = false
		}
	}
	relay := relayPort
	if !*useCacti {
		cactiPorts, relay = nil, 0
	}
	urls, err := it.WriteJanusConfigs(it.NodeConfig{Dir: abs, RPC: rpc, Deployment: depPath, Policy: *policy,
		AllowOrigin: *origin, Ports: map[string]int{jt.AuditorNode: 8080, "vendedora": 8081, "banco": 8082},
		Cacti: cactiPorts, Relay: relay})
	must(err)
	must(os.Setenv("JANUS_CONFIG_DIR", abs))

	// 3) Rede Fabric + Panurus + nós FSC com a SDK do Janus
	log.Print("subindo Fabric 3.x + Panurus (NWO)... (1 a 2 minutos)")
	n, err := it.StartNetwork(it.Options{})
	must(err)
	defer n.Stop()
	for node, u := range urls {
		must(waitHealth(u, 2*time.Minute))
		log.Printf("nó %-12s %s", node, u)
	}
	// connection profile + identidades das orgs para os conectores do Cacti
	ccp, ids, err := n.FabricConnectionProfile()
	must(err)
	cdir := filepath.Join(abs, "cacti")
	must(os.MkdirAll(cdir, 0o700))
	must(writeJSON(filepath.Join(cdir, "ccp.json"), ccp))
	must(writeJSON(filepath.Join(cdir, "identities.json"), ids))
	msp := map[string]string{jt.AuditorNode: "RegistradoraMSP", "banco": "BancoMSP"}
	for node, port := range cactiPorts {
		logf, err := os.Create(filepath.Join(cdir, node+".log"))
		must(err)
		srv, err := it.StartCacti(cdir, node, port, ccp, ids[msp[node]], logf)
		must(err)
		defer srv.Stop()
		log.Printf("cacti %-12s %s (conector Fabric, %s)", node, srv.URL, msp[node])
	}
	if relay > 0 {
		logf, err := os.Create(filepath.Join(cdir, "relay.log"))
		must(err)
		srv, err := it.StartCactiRelay(cdir, relay, rpc, *besuWS, it.KeyRelay, logf)
		must(err)
		defer srv.Stop()
		log.Printf("cacti %-12s %s (conector Besu, relay da Vendedora)", "relay", srv.URL)
	}

	// 4) Governança B2: cadastra as chaves de atestador (Registradora + Banco), política 2-de-2
	admin, _ := crypto.HexToECDSA(it.KeyAdmin)
	bc, err := besu.Dial(ctx, rpc, admin, dep)
	must(err)
	var atts []besu.Attestor
	for _, node := range []string{jt.AuditorNode, "banco"} {
		in, err := info(urls[node])
		must(err)
		if in.Attestor == nil {
			must(fmt.Errorf("%s sem chave de atestador", node))
		}
		atts = append(atts, besu.Attestor{MSPID: in.Attestor.MSPID,
			Qx: common.HexToHash(in.Attestor.Qx), Qy: common.HexToHash(in.Attestor.Qy)})
	}
	must(bc.RegisterAttestors(ctx, 2, atts...))
	bc.Close()

	fmt.Printf(`
Janus no ar.
  Registradora  %s
  Vendedora     %s
  Banco         %s
  Besu          %s
  Cacti         %s

Simulador web: npm run web  →  http://localhost:5173/?mode=live
Ctrl-C para derrubar tudo.
`, urls[jt.AuditorNode], urls["vendedora"], urls["banco"], rpc, cactiSummary(cactiPorts))

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	log.Print("derrubando a rede...")
}

func cactiSummary(ports map[string]int) string {
	if len(ports) == 0 {
		return "desabilitado (atestação sem prova de ledger via Cacti)"
	}
	return fmt.Sprintf("Registradora :%d, Banco :%d (conector Fabric), relay :4102 (conector Besu) — v3.0.1",
		ports[jt.AuditorNode], ports["banco"])
}

func writeJSON(path string, v any) error {
	raw, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, raw, 0o600)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func waitHealth(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if r, err := http.Get(url + "/api/v1/health"); err == nil {
			r.Body.Close()
			if r.StatusCode == http.StatusOK {
				return nil
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return errors.New("api não respondeu: " + url)
}

func info(url string) (api.Info, error) {
	var in api.Info
	r, err := http.Get(url + "/api/v1/info")
	if err != nil {
		return in, err
	}
	defer r.Body.Close()
	return in, json.NewDecoder(r.Body).Decode(&in)
}
