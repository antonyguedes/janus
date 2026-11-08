package api

import (
	"bufio"
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"janus/node/internal/engine"
	"janus/node/internal/store"
	"janus/node/internal/swap"
)

// stuck é um protocolo que nunca sai de "wait", usado só para gerar passos.
type stuck struct{}

func (stuck) Name() swap.Protocol { return swap.ProtocolHTLC }
func (stuck) Observe(context.Context, swap.Role, *swap.SwapTerms) (swap.Observation, error) {
	return swap.Observation{Now: time.Now()}, nil
}
func (stuck) Execute(context.Context, swap.Role, *swap.SwapTerms, swap.Decision) error { return nil }

func TestAPI(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	eng := engine.New(engine.Config{Role: swap.RoleBuyer, Policy: swap.DevPolicy(), Poll: 30 * time.Millisecond,
		Protocol: map[swap.Protocol]swap.SettlementProtocol{swap.ProtocolHTLC: stuck{}}})
	srv := httptest.NewServer((&Server{Info: Info{Node: "banco", Role: swap.RoleBuyer}, Store: st, Engine: eng}).Handler())
	defer srv.Close()

	terms := swap.SwapTerms{SwapID: [32]byte{0xab}, Protocol: swap.ProtocolHTLC, AssetID: "r1",
		Price: swap.Amount{Value: 100}, T1: time.Now().Add(time.Hour).Unix(), T2: time.Now().Add(30 * time.Minute).Unix()}
	ctx := context.Background()
	if err := st.Put(ctx, store.Record{SwapID: terms.SwapID, Role: swap.RoleBuyer,
		Signed: swap.SignedTerms{Terms: terms}, Status: swap.StatusNegotiated}); err != nil {
		t.Fatal(err)
	}
	if err := st.SetStatus(ctx, terms.SwapID, swap.StatusBothLocked); err != nil {
		t.Fatal(err)
	}

	// REST
	res, err := http.Get(srv.URL + "/api/v1/swaps")
	if err != nil {
		t.Fatal(err)
	}
	var list []SwapDTO
	_ = json.NewDecoder(res.Body).Decode(&list)
	res.Body.Close()
	if len(list) != 1 || list[0].Status != swap.StatusBothLocked || list[0].AssetID != "r1" {
		t.Fatalf("list: %+v", list)
	}
	id := hex.EncodeToString(terms.SwapID[:])
	res, _ = http.Get(srv.URL + "/api/v1/swaps/" + id)
	var one SwapDTO
	_ = json.NewDecoder(res.Body).Decode(&one)
	res.Body.Close()
	if one.Terms == nil || one.Terms.AssetID != "r1" {
		t.Fatalf("get: %+v", one)
	}
	if res, _ := http.Get(srv.URL + "/api/v1/swaps/zz"); res.StatusCode != http.StatusBadRequest {
		t.Fatalf("bad id: %d", res.StatusCode)
	}
	if res, _ := http.Get(srv.URL + "/api/v1/swaps/" + strings.Repeat("00", 32)); res.StatusCode != http.StatusNotFound {
		t.Fatalf("not found: %d", res.StatusCode)
	}

	// SSE: hello e depois passos do motor
	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/api/v1/events", nil)
	sctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	res, err = http.DefaultClient.Do(req.WithContext(sctx))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if ct := res.Header.Get("Content-Type"); ct != "text/event-stream" {
		t.Fatalf("content-type %q", ct)
	}
	go func() { _, _ = eng.Run(sctx, &terms) }()
	sc := bufio.NewScanner(res.Body)
	var gotHello, gotStep bool
	for sc.Scan() && !(gotHello && gotStep) {
		line := sc.Text()
		gotHello = gotHello || line == "event: hello"
		if strings.HasPrefix(line, "data: ") && strings.Contains(line, `"swapId":"`+id) {
			gotStep = true
		}
	}
	if !gotHello || !gotStep {
		t.Fatalf("sse hello=%v step=%v", gotHello, gotStep)
	}

	// reconciliação: o swap ativo aparece em List(activeOnly)
	act, _ := st.List(ctx, true)
	if len(act) != 1 {
		t.Fatalf("ativos: %d", len(act))
	}
	_ = st.SetStatus(ctx, terms.SwapID, swap.StatusCompleted)
	if act, _ := st.List(ctx, true); len(act) != 0 {
		t.Fatal("terminal não deveria ser ativo")
	}
}
