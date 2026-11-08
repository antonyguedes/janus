// Package janus é a SDK do FSC que transforma um nó FSC/Panurus no nó Janus de uma instituição:
// cliente Besu, motor de liquidação (watchtower/watchdog), negociação de termos, API REST + SSE,
// store SQLite e segredos cifrados. É instalada tanto pelo binário cmd/janus-node quanto pelos
// nós gerados pelo NWO nos testes. Assim o código testado é o mesmo que roda em produção.
//
// Configuração: janus.json no mesmo diretório do core.yaml do nó (ou JANUS_CONFIG).
package janus

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/config"
	fscview "github.com/hyperledger-labs/fabric-smart-client/platform/view/services/view"

	"janus/node/internal/api"
	"janus/node/internal/besu"
	"janus/node/internal/engine"
	"janus/node/internal/hashlock"
	"janus/node/internal/protocol/htlc"
	"janus/node/internal/secrets"
	"janus/node/internal/store"
	"janus/node/internal/swap"
	"janus/node/internal/token"
)

// Config é o conteúdo de janus.json.
type Config struct {
	Role         string        `json:"role"` // vendedora | banco | registradora
	MSPID        string        `json:"mspId"`
	Counterparty string        `json:"counterparty,omitempty"` // nó FSC do Banco (só a Vendedora)
	DataDir      string        `json:"dataDir"`
	Policy       string        `json:"policy"` // dev | prod
	API          APIConfig     `json:"api"`
	Besu         *BesuConfig   `json:"besu,omitempty"`
	Faults       engine.Faults `json:"faults"`
}

type APIConfig struct {
	Listen      string `json:"listen"`
	AllowOrigin string `json:"allowOrigin,omitempty"`
}

type BesuConfig struct {
	RPC        string `json:"rpc"`
	KeyHex     string `json:"keyHex"`     // DEV: em produção, KMS/HSM
	Deployment string `json:"deployment"` // deployments/<chainId>.local.json
	FromBlock  uint64 `json:"fromBlock"`
}

// SDK implementa node.SDK (+ PostStart) do FSC.
type SDK struct {
	registry services.Registry
	log      *slog.Logger
	cancel   context.CancelFunc
	srv      *http.Server
}

// NewSDK segue a convenção usada pelo NWO no main gerado: <pkg>.NewSDK(n).
func NewSDK(registry services.Registry) *SDK {
	return &SDK{registry: registry, log: slog.Default().With("component", "janus")}
}

func (s *SDK) Install() error              { return nil }
func (s *SDK) Start(context.Context) error { return nil }

// configPath: JANUS_CONFIG (arquivo), ou JANUS_CONFIG_DIR/<id do nó>.json (vários nós num host,
// como no NWO), ou janus.json ao lado do core.yaml.
func (s *SDK) configPath() string {
	if p := os.Getenv("JANUS_CONFIG"); p != "" {
		return p
	}
	cp := config.GetProvider(s.registry)
	if dir := os.Getenv("JANUS_CONFIG_DIR"); dir != "" {
		id := cp.ID()
		for _, cand := range []string{id, strings.TrimSuffix(id, filepath.Ext(id))} {
			if p := filepath.Join(dir, cand+".json"); fileExists(p) {
				return p
			}
		}
		return filepath.Join(dir, id+".json")
	}
	return filepath.Join(filepath.Dir(cp.ConfigFileUsed()), "janus.json")
}

func fileExists(p string) bool { _, err := os.Stat(p); return err == nil }

// LoadConfig lê e valida janus.json.
func LoadConfig(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	if err := json.Unmarshal(raw, c); err != nil {
		return nil, fmt.Errorf("janus: %s: %w", path, err)
	}
	switch c.Role {
	case "vendedora", "banco", "registradora":
	default:
		return nil, fmt.Errorf("janus: invalid role %q", c.Role)
	}
	if c.DataDir == "" {
		c.DataDir = filepath.Join(filepath.Dir(path), "janus-data")
	}
	if c.Role != "registradora" && c.Besu == nil {
		return nil, errors.New("janus: besu config required for vendedora/banco")
	}
	if c.Role == "vendedora" && c.Counterparty == "" {
		return nil, errors.New("janus: vendedora needs counterparty (banco node id)")
	}
	return c, nil
}

func (c *Config) policy() swap.TimePolicy {
	if c.Policy == "prod" {
		return swap.ProductionPolicy()
	}
	return swap.DevPolicy()
}

// PostStart sobe o nó Janus depois que o FSC e o Panurus estão de pé.
func (s *SDK) PostStart(ctx context.Context) error {
	path := s.configPath()
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		s.log.Info("janus.json não encontrado: nó FSC sem papel Janus", "path", path)
		return nil
	}
	cfg, err := LoadConfig(path)
	if err != nil {
		return err
	}
	s.log = s.log.With("role", cfg.Role)
	views, err := fscview.GetManager(s.registry)
	if err != nil {
		return fmt.Errorf("janus: view manager: %w", err)
	}
	if err := os.MkdirAll(cfg.DataDir, 0o700); err != nil {
		return err
	}
	runCtx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	if cfg.Role == "registradora" {
		return s.startRegistradora(runCtx, cfg, views)
	}
	return s.startParty(runCtx, cfg, views)
}

// Stop encerra API e motores (chamado pelo binário no shutdown).
func (s *SDK) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
	if s.srv != nil {
		ctx, c := context.WithTimeout(context.Background(), 3*time.Second)
		defer c()
		_ = s.srv.Shutdown(ctx)
	}
}

func (s *SDK) serve(cfg *Config, h http.Handler) error {
	if cfg.API.Listen == "" {
		return nil
	}
	ln, err := net.Listen("tcp", cfg.API.Listen)
	if err != nil {
		return fmt.Errorf("janus: api listen: %w", err)
	}
	s.srv = &http.Server{Handler: h, ReadHeaderTimeout: 5 * time.Second}
	go func() { _ = s.srv.Serve(ln) }()
	s.log.Info("api no ar", "addr", ln.Addr().String())
	return nil
}

// --------------------------------------------------------------------- Registradora

func (s *SDK) startRegistradora(ctx context.Context, cfg *Config, views *fscview.Manager) error {
	// Os parâmetros públicos do Panurus podem chegar depois do PostStart: registra com retry.
	go func() {
		for i := 0; i < 60; i++ {
			if _, err := views.InitiateView(ctx, &token.RegisterAuditorView{}); err == nil {
				s.log.Info("auditor registrado")
				return
			} else if i%10 == 0 {
				s.log.Warn("registro do auditor pendente", "err", err)
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(2 * time.Second):
			}
		}
	}()
	st, err := store.Open(filepath.Join(cfg.DataDir, "janus.db"))
	if err != nil {
		return err
	}
	eng := engine.New(engine.Config{Role: "registradora", Policy: cfg.policy()})
	srv := &api.Server{Info: api.Info{Node: cfg.Role, Role: "registradora", MSPID: cfg.MSPID},
		Store: st, Engine: eng, AllowOrigin: cfg.API.AllowOrigin,
		Issue: func(r *http.Request, req api.IssueRequest) (string, error) {
			res, err := views.InitiateView(r.Context(), &token.IssueView{Issue: &token.Issue{
				Recipient: req.Owner,
				Recebivel: token.Recebivel{Registro: req.Registro, Cedente: req.Cedente, Sacado: req.Sacado,
					ValorFace: req.ValorFace, Vencimento: req.Vencimento},
			}})
			if err != nil {
				return "", err
			}
			id, _ := res.(string)
			return id, nil
		}}
	return s.serve(cfg, srv.Handler())
}

// --------------------------------------------------------------------- Vendedora / Banco

type besuAdapter struct{ c *besu.Client }

func (b besuAdapter) View(ctx context.Context, t *swap.SwapTerms) (swap.BesuView, error) {
	return b.c.View(ctx, t)
}
func (b besuAdapter) LockPayment(ctx context.Context, t *swap.SwapTerms) error {
	_, err := b.c.LockPayment(ctx, t)
	return err
}
func (b besuAdapter) ClaimPayment(ctx context.Context, id, s [32]byte) error {
	_, err := b.c.ClaimPayment(ctx, id, s)
	return err
}
func (b besuAdapter) RefundPayment(ctx context.Context, id [32]byte) error {
	_, err := b.c.RefundPayment(ctx, id)
	return err
}

func (s *SDK) startParty(ctx context.Context, cfg *Config, views *fscview.Manager) error {
	role := swap.RoleSeller
	if cfg.Role == "banco" {
		role = swap.RoleBuyer
	}
	policy := cfg.policy()

	// Besu
	dep, err := besu.LoadDeployment(cfg.Besu.Deployment)
	if err != nil {
		return fmt.Errorf("janus: deployment: %w", err)
	}
	key, err := crypto.HexToECDSA(trim0x(cfg.Besu.KeyHex))
	if err != nil {
		return fmt.Errorf("janus: besu key: %w", err)
	}
	bc, err := besu.Dial(ctx, cfg.Besu.RPC, key, dep)
	if err != nil {
		return err
	}

	// Store e segredos
	st, err := store.Open(filepath.Join(cfg.DataDir, "janus.db"))
	if err != nil {
		return err
	}
	mk, err := secrets.LoadOrCreateKey(filepath.Join(cfg.DataDir, "master.key"))
	if err != nil {
		return err
	}
	sec, err := secrets.NewFileStore(filepath.Join(cfg.DataDir, "secrets"), mk)
	if err != nil {
		return err
	}

	proto := &htlc.Protocol{
		Fabric:  &token.FSCPort{Views: views, Counterparty: cfg.Counterparty},
		Besu:    besuAdapter{bc},
		Secrets: sec,
		Policy:  policy,
	}
	eng := engine.New(engine.Config{Role: role, Policy: policy, Faults: cfg.Faults, Log: s.log,
		Protocol: map[swap.Protocol]swap.SettlementProtocol{swap.ProtocolHTLC: proto}})

	// Status de cada passo vai para o store (a API lista a partir dele).
	steps, _ := eng.Subscribe()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case stp := <-steps:
				if stp.Status != "" {
					_ = st.SetStatus(ctx, stp.SwapID, stp.Status)
				}
			}
		}
	}()

	run := func(t swap.SwapTerms) {
		go func() {
			d, err := eng.Run(ctx, &t)
			s.log.Info("swap encerrado neste nó", "swap", hex.EncodeToString(t.SwapID[:8]), "decision", d.Kind, "err", err)
			if errors.Is(err, engine.ErrCrash) {
				s.log.Warn("falha injetada: crash depois da revelação. Reinicie o nó para a reconciliação")
				os.Exit(3)
			}
		}()
	}

	// Negociação (hooks das views de proposta/aceite)
	evm := swap.Address(bc.From())
	token.SetNegotiationHooks(token.NegotiationHooks{
		Policy:     policy,
		LocalParty: func() swap.Party { return swap.Party{MSPID: cfg.MSPID, EVM: evm} },
		Accept: func(_ context.Context, t *swap.SwapTerms) error {
			if t.Protocol != swap.ProtocolHTLC {
				return fmt.Errorf("protocol %s not enabled on this node", t.Protocol)
			}
			if t.Price.Token != swap.Address(dep.DrexToken) || t.Networks.PaymentContract != swap.Address(dep.HTLCPayment) ||
				t.Networks.BesuChainID != dep.ChainID {
				return errors.New("terms reference unknown Besu deployment")
			}
			return nil
		},
		OnSigned: func(c context.Context, r swap.Role, sig *swap.SignedTerms) error {
			if r == swap.RoleSeller {
				// S foi guardado sob o hashlock antes da negociação. Indexa pelo swapId antes de
				// iniciar o motor, que procura o segredo pelo swapId.
				pre, err := sec.Get(c, sig.Terms.Hashlock)
				if err != nil || pre == nil {
					return fmt.Errorf("preimage for hashlock not found: %v", err)
				}
				if err := sec.Put(c, sig.Terms.SwapID, *pre); err != nil {
					return err
				}
			}
			if err := st.Put(c, store.Record{SwapID: sig.Terms.SwapID, Role: r, Signed: *sig, Status: swap.StatusNegotiated}); err != nil {
				return err
			}
			run(sig.Terms)
			return nil
		},
	})

	// Reconciliação no boot: retoma todo swap não terminal. A cadeia é a fonte da verdade.
	active, err := st.List(ctx, true)
	if err != nil {
		return err
	}
	for _, rec := range active {
		if rec.Role == swap.RoleBuyer {
			_ = proto.Accept(ctx, &rec.Signed.Terms) // reinstala a expectativa de lock (memória)
		}
		s.log.Info("reconciliando swap", "swap", hex.EncodeToString(rec.SwapID[:8]), "status", rec.Status)
		run(rec.Signed.Terms)
	}

	// Watcher do Besu: cada evento acorda o motor do swap correspondente.
	events := make(chan besu.Event, 64)
	go func() {
		_ = bc.Watch(ctx, besu.FileCheckpoint{Path: filepath.Join(cfg.DataDir, "besu.checkpoint")},
			besu.WatchOptions{FromBlock: cfg.Besu.FromBlock, Poll: time.Second}, events)
	}()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case e := <-events:
				eng.Kick(e.SwapID)
			}
		}
	}()

	srv := &api.Server{
		Info: api.Info{Node: cfg.Role, Role: role, MSPID: cfg.MSPID, EVM: bc.From().Hex(),
			Protocols: []swap.Protocol{swap.ProtocolHTLC}, Faults: cfg.Faults},
		Store: st, Engine: eng, Log: s.log, AllowOrigin: cfg.API.AllowOrigin,
	}
	if role == swap.RoleSeller {
		srv.Propose = func(r *http.Request, req api.ProposeRequest) (api.SwapDTO, error) {
			return s.propose(r.Context(), views, sec, st, dep, req)
		}
	}
	return s.serve(cfg, srv.Handler())
}

// propose: gera S, guarda cifrado ANTES de qualquer mensagem e negocia os termos com o Banco.
func (s *SDK) propose(ctx context.Context, views *fscview.Manager, sec *secrets.FileStore, st *store.Store,
	dep *besu.Deployment, req api.ProposeRequest) (api.SwapDTO, error) {
	if req.Protocol != swap.ProtocolHTLC {
		return api.SwapDTO{}, fmt.Errorf("protocol %s not enabled", req.Protocol)
	}
	S, H, err := hashlock.New()
	if err != nil {
		return api.SwapDTO{}, err
	}
	var nonce [32]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return api.SwapDTO{}, err
	}
	// O swapId depende das identidades, que só a view conhece. Por isso o segredo é indexado
	// pelo hashlock até a negociação terminar. O hook OnSigned copia S para o swapId.
	if err := sec.Put(ctx, [32]byte(H), S); err != nil {
		return api.SwapDTO{}, err
	}
	res, err := views.InitiateView(ctx, &token.ProposeView{Propose: &token.Propose{
		BuyerNode: req.BuyerNode, Protocol: req.Protocol, AssetID: req.AssetID,
		Price:    swap.Amount{Token: swap.Address(dep.DrexToken), Value: req.Price},
		Hashlock: H, Nonce: nonce,
		Networks: swap.NetworkRefs{FabricChannel: dep.FabricChannel, FabricNamespace: dep.FabricNamespace,
			BesuChainID: dep.ChainID, PaymentContract: swap.Address(dep.HTLCPayment)},
	}})
	if err != nil {
		return api.SwapDTO{}, err
	}
	raw, _ := res.([]byte)
	var signed swap.SignedTerms
	if err := json.Unmarshal(raw, &signed); err != nil {
		return api.SwapDTO{}, err
	}
	rec, err := st.Get(ctx, signed.Terms.SwapID)
	if err != nil {
		return api.SwapDTO{}, err
	}
	return api.ToDTO(rec), nil
}

func trim0x(s string) string {
	if len(s) > 1 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
		return s[2:]
	}
	return s
}
