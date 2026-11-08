package sim

import (
	"testing"
	"time"

	"janus/node/internal/swap"
)

func protocols() []swap.Protocol { return []swap.Protocol{swap.ProtocolHTLC, swap.ProtocolCacti} }

func run(t *testing.T, proto swap.Protocol, sc Scenario) Result {
	t.Helper()
	r, err := Run(Config{Protocol: proto, Scenario: sc})
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Os cenários do plano (Fase 6) com a expectativa de resultado de cada um.
func TestScenarios(t *testing.T) {
	cases := []struct {
		sc     Scenario
		protos []swap.Protocol
		check  func(Result) bool
	}{
		{Happy, protocols(), func(r Result) bool { return r.Outcome.Status == swap.StatusCompleted }},
		{BuyerNoLock, protocols(), func(r Result) bool {
			return r.Outcome.Status == swap.StatusReverted && !r.Outcome.SellerLost()
		}},
		{SellerNoReveal, protocols(), func(r Result) bool {
			return r.Outcome.Status == swap.StatusReverted && !r.Outcome.SellerLost() && !r.Outcome.BuyerLost()
		}},
		{RevealAtT2Minus, []swap.Protocol{swap.ProtocolHTLC}, func(r Result) bool {
			return r.Outcome.Status == swap.StatusCompleted
		}},
		{BuyerCrashRecover, protocols(), func(r Result) bool { return r.Outcome.Status == swap.StatusCompleted }},
		{BuyerOfflinePastT1, []swap.Protocol{swap.ProtocolHTLC}, func(r Result) bool {
			// risco documentado: o Banco perde, e a Vendedora honesta NÃO faz reclaim (o NFT segue travado)
			return r.Outcome.BuyerLost() && r.Snapshots[len(r.Snapshots)-1].Fabric == swap.FabricLocked
		}},
		{WrongDeadline, protocols(), func(r Result) bool {
			return !r.Outcome.BuyerLost() && r.Outcome.Status != swap.StatusCompleted
		}},
		{EarlyReclaim, protocols(), func(r Result) bool {
			return r.Outcome.Status == swap.StatusCompleted && !r.Outcome.BuyerLost()
		}},
		{LateRevealLeak, []swap.Protocol{swap.ProtocolHTLC}, func(r Result) bool {
			return r.Outcome.SellerLost() && r.Outcome.Status == swap.StatusBuyerHasBoth
		}},
	}
	for _, c := range cases {
		for _, p := range c.protos {
			t.Run(string(c.sc)+"/"+string(p), func(t *testing.T) {
				r := run(t, p, c.sc)
				if !c.check(r) {
					t.Fatalf("outcome=%+v\nevents=%+v", r.Outcome, r.Events)
				}
				if len(r.Snapshots) == 0 || r.T1-r.T2 < r.Delta {
					t.Fatalf("result malformed: %+v", r)
				}
			})
		}
	}
	if _, err := Run(Config{Scenario: "nope"}); err == nil {
		t.Fatal("cenário desconhecido")
	}
	if _, err := Run(Config{Protocol: swap.ProtocolCacti, Scenario: RevealAtT2Minus}); err == nil {
		t.Fatal("cenário só de A")
	}
	if len(Scenarios) != 9 {
		t.Fatal("catálogo de cenários")
	}
}

// Banco honesto, no cenário da revelação tardia, nunca fica com as duas pernas.
func TestHonestBuyerNeverTakesBoth(t *testing.T) {
	w, s, _, _ := Build(Config{Protocol: swap.ProtocolHTLC, Scenario: LateRevealLeak})
	if r := RunAgents(w, s, Honest(swap.RoleBuyer)); r.Outcome.Status == swap.StatusBuyerHasBoth {
		t.Fatalf("%+v", r.Events)
	}
}

// Varredura: dois honestos, quedas curtas (< Δ) e latências variadas, nunca terminam em anomalia/perda.
func TestSweepHonestWithOutages(t *testing.T) {
	for _, proto := range protocols() {
		for lat := int64(1); lat <= 10; lat += 3 {
			for start := int64(0); start < 360; start += 17 {
				for _, who := range []swap.Role{swap.RoleSeller, swap.RoleBuyer} {
					w := NewWorld(proto)
					w.Latency = lat
					outage := int64(w.P.Delta()/time.Second) - lat - 5
					a := map[swap.Role]Agent{swap.RoleSeller: Honest(swap.RoleSeller), swap.RoleBuyer: Honest(swap.RoleBuyer)}
					s := start
					a[who] = Agent{Role: who, Online: func(rel int64) bool { return rel < s || rel >= s+outage }}
					o := RunAgents(w, a[swap.RoleSeller], a[swap.RoleBuyer]).Outcome
					if o.SellerLost() || o.BuyerLost() || o.Status == swap.StatusSellerHasBoth || o.Status == swap.StatusBuyerHasBoth {
						t.Fatalf("%s lat=%d outage[%s]=%d+%d: %+v", proto, lat, who, s, outage, o)
					}
				}
			}
		}
	}
}
