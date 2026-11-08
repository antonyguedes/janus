package sim

import (
	"fmt"
	"time"

	"janus/node/internal/swap"
)

func wait(reason string) swap.Decision { return swap.Decision{Kind: swap.ActWait, Reason: reason} }
func act(k swap.ActionKind, reason string) swap.Decision {
	return swap.Decision{Kind: k, Reason: reason}
}

// Build monta o mundo e os agentes de um cenário.
func Build(cfg Config) (*World, Agent, Agent, error) {
	if cfg.Protocol == "" {
		cfg.Protocol = swap.ProtocolHTLC
	}
	w := NewWorld(cfg.Protocol)
	if cfg.Latency > 0 {
		w.Latency = cfg.Latency
	}
	seller, buyer := Honest(swap.RoleSeller), Honest(swap.RoleBuyer)
	switch cfg.Scenario {
	case Happy, "":
	case BuyerNoLock:
		buyer.Override = func(_ *World, d swap.Decision) swap.Decision {
			if d.Kind == swap.ActLockPayment {
				return wait("adversário: não trava o pagamento")
			}
			return d
		}
	case SellerNoReveal:
		if cfg.Protocol == swap.ProtocolCacti {
			w.RelayOK = false // relay parado: a Vendedora honesta nunca obtém a atestação
		} else {
			seller.Override = func(_ *World, d swap.Decision) swap.Decision {
				if d.Kind == swap.ActClaimPayment {
					return wait("adversário: não revela S")
				}
				return d
			}
		}
	case RevealAtT2Minus:
		if cfg.Protocol != swap.ProtocolHTLC {
			return nil, Agent{}, Agent{}, fmt.Errorf("cenário %s só se aplica à abordagem A", cfg.Scenario)
		}
		seller.Override = func(w *World, d swap.Decision) swap.Decision {
			if w.Besu.State == swap.BesuLocked && w.Now == w.T.T2-w.Latency-1 {
				return act(swap.ActClaimPayment, "adversário: revelação no último instante")
			}
			if d.Kind == swap.ActClaimPayment {
				return wait("adversário: segura S")
			}
			return d
		}
	case BuyerCrashRecover:
		var downAt int64 = -1
		buyer.Online = func(rel int64) bool {
			if downAt < 0 && w.Besu.State == swap.BesuClaimed {
				downAt = rel
			}
			return downAt < 0 || rel > downAt+int64(w.P.Delta()/time.Second)
		}
	case BuyerOfflinePastT1:
		buyer.Online = func(int64) bool { return w.Besu.State != swap.BesuClaimed }
	case WrongDeadline:
		w.LockDeadlineSkew = -5 * time.Minute
	case EarlyReclaim:
		seller.Override = func(w *World, d swap.Decision) swap.Decision {
			if w.Besu.State == swap.BesuClaimed && w.Fab.State == swap.FabricLocked {
				return act(swap.ActReclaimAsset, "adversário: reclaim antecipado depois de paga")
			}
			return d
		}
	case LateRevealLeak:
		if cfg.Protocol != swap.ProtocolHTLC {
			return nil, Agent{}, Agent{}, fmt.Errorf("cenário %s só se aplica à abordagem A", cfg.Scenario)
		}
		seller.Override = func(w *World, d swap.Decision) swap.Decision {
			if w.Besu.State == swap.BesuLocked && w.Now == w.T.T2-1 {
				return act(swap.ActClaimPayment, "imprudente: revela depois da margem")
			}
			if d.Kind == swap.ActClaimPayment {
				return wait("segura S")
			}
			return d
		}
		buyer.Override = func(w *World, d swap.Decision) swap.Decision {
			if w.BuyerS != nil && w.Fab.State == swap.FabricLocked && ts(w.Now).Before(w.Fab.Deadline) &&
				w.Besu.State == swap.BesuRefunded {
				return act(swap.ActClaimAsset, "adversário: usa o S vazado")
			}
			return d
		}
	default:
		return nil, Agent{}, Agent{}, fmt.Errorf("cenário desconhecido %q", cfg.Scenario)
	}
	return w, seller, buyer, nil
}

// Run executa um cenário pré-definido.
func Run(cfg Config) (Result, error) {
	w, s, b, err := Build(cfg)
	if err != nil {
		return Result{}, err
	}
	res := RunAgents(w, s, b)
	res.Config = cfg
	return res, nil
}
