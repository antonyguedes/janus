package swap

import (
	"crypto/sha256"
	"fmt"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------------------------
// Mundo simulado: regras idênticas às dos contratos (HTLCPayment/AttestedSettlement) e às do
// validador HTLC do Panurus (claim se now < deadline, reclaim se now >= deadline). Cada ação
// leva `latency` segundos para comitar e é revalidada NO INSTANTE DO COMMIT, como numa rede real.
// ---------------------------------------------------------------------------------------------

const base = 1_900_000_000 // epoch arbitrária

type pending struct {
	at   int64
	kind ActionKind
	by   Role
}

type world struct {
	t     *SwapTerms
	p     TimePolicy
	now   int64
	fab   FabricView
	besu  BesuView
	queue []pending
	attAt int64 // instante em que as atestações ficam prontas (B); 0 = não solicitadas
	log   []string

	secret   [32]byte
	sellerS  *[32]byte // o segredo guardado pela Vendedora
	buyerS   *[32]byte // o que o Banco aprendeu
	latency  int64
	relayOK  bool
	relayLat int64
	// lockDeadlineSkew: desvio do relógio da Vendedora ao calcular script.Deadline.
	lockDeadlineSkew time.Duration
}

func newWorld(proto Protocol) *world {
	p := DevPolicy()
	s := sha256.Sum256([]byte("sim-secret"))
	h := sha256.Sum256(s[:])
	seller := Party{MSPID: "VendedoraMSP", FSCIdentity: []byte("vendedora"), EVM: Address{1}}
	buyer := Party{MSPID: "BancoMSP", FSCIdentity: []byte("banco"), EVM: Address{2}}
	nonce := [32]byte{9}
	t := &SwapTerms{
		Protocol: proto, AssetID: "recebivel-0001", Seller: seller, Buyer: buyer,
		Price: Amount{Token: Address{0xd}, Value: 1_000_00}, Hashlock: h, Nonce: nonce,
		T1: base + int64(p.T1Offset/time.Second), T2: base + int64(p.T2Offset/time.Second),
		Networks: NetworkRefs{FabricChannel: "recebiveis", FabricNamespace: "tcc", BesuChainID: 1337, PaymentContract: Address{0xc}},
	}
	t.SwapID = ComputeSwapID(h, seller, buyer, nonce)
	return &world{t: t, p: p, now: base, secret: s, sellerS: &s, latency: 4, relayOK: true, relayLat: 8}
}

func (w *world) ts(x int64) time.Time { return time.Unix(x, 0) }

func (w *world) observe(role Role) Observation {
	o := Observation{Now: w.ts(w.now), Fabric: w.fab, Besu: w.besu}
	if role == RoleSeller {
		o.Preimage = w.sellerS
		o.AttestationReady = w.attAt != 0 && w.now >= w.attAt
	} else {
		o.Preimage = w.buyerS
	}
	return o
}

func (w *world) inFlight(k ActionKind) bool {
	for _, q := range w.queue {
		if q.kind == k {
			return true
		}
	}
	return false
}

func (w *world) submit(by Role, k ActionKind) {
	switch k {
	case ActWait, ActDone, ActAbort, ActAlert:
		return
	case ActRequestAttest:
		if w.attAt == 0 && w.relayOK && w.fab.State != FabricNone {
			w.attAt = w.now + w.relayLat
		}
		return
	}
	if w.inFlight(k) {
		return
	}
	w.queue = append(w.queue, pending{at: w.now + w.latency, kind: k, by: by})
}

// commit aplica as transações maduras com as regras on-chain.
func (w *world) commit() {
	var rest []pending
	for _, q := range w.queue {
		if q.at > w.now {
			rest = append(rest, q)
			continue
		}
		now := w.ts(w.now)
		ok := false
		switch q.kind {
		case ActLockAsset:
			if w.fab.State == FabricNone {
				w.fab = FabricView{
					State: FabricLocked, AssetID: w.t.AssetID, Sender: w.t.Seller.FSCIdentity,
					Recipient: w.t.Buyer.FSCIdentity, Hashlock: w.t.Hashlock, HashOK: true,
					Deadline: w.t.T1Time().Add(w.lockDeadlineSkew), TxID: "tx-lock", Block: 7,
				}
				ok = true
			}
		case ActLockPayment:
			if w.besu.State == BesuNone && now.Before(w.t.T2Time()) {
				w.besu = BesuView{
					State: BesuLocked, Contract: w.t.Networks.PaymentContract, Token: w.t.Price.Token,
					Sender: w.t.Buyer.EVM, Recipient: w.t.Seller.EVM, Amount: w.t.Price.Value,
					Hashlock: w.t.Hashlock, Timelock: w.t.T2Time(),
				}
				if w.t.Protocol == ProtocolCacti {
					w.besu.MinFabricDeadline = w.t.T1Time().Add(-w.p.DeadlineTolerance)
				}
				ok = true
			}
		case ActClaimPayment, ActSettlePayment:
			valid := w.besu.State == BesuLocked && now.Before(w.besu.Timelock)
			if q.kind == ActSettlePayment {
				// AttestedSettlement: exige uma atestação de lock com deadline >= minFabricDeadline
				valid = valid && w.attAt != 0 && w.now >= w.attAt && w.fab.State != FabricNone &&
					!w.fab.Deadline.Before(w.besu.MinFabricDeadline)
			}
			// A preimagem fica pública no mempool mesmo que a tx falhe (pior caso para a Vendedora).
			s := w.secret
			w.buyerS = &s
			if valid {
				w.besu.State, w.besu.Preimage = BesuClaimed, &s
				ok = true
			}
		case ActClaimAsset:
			if w.fab.State == FabricLocked && now.Before(w.fab.Deadline) && w.buyerS != nil {
				w.fab.State = FabricClaimed
				ok = true
			}
		case ActRefundPayment:
			if w.besu.State == BesuLocked && !now.Before(w.besu.Timelock) {
				w.besu.State = BesuRefunded
				ok = true
			}
		case ActReclaimAsset:
			if w.fab.State == FabricLocked && !now.Before(w.fab.Deadline) {
				w.fab.State = FabricReclaimed
				ok = true
			}
		}
		w.log = append(w.log, fmt.Sprintf("t=%+4d %-9s %-15s ok=%v", w.now-base, q.by, q.kind, ok))
	}
	w.queue = rest
}

// agent é um comportamento: honesto (NextAction) ou adversarial.
type agent struct {
	role   Role
	online func(rel int64) bool
	// override, quando não nil, substitui a decisão honesta (adversário).
	override func(w *world, d Decision) Decision
}

func honest(r Role) agent { return agent{role: r, online: func(int64) bool { return true }} }

type outcome struct {
	status        Status
	sellerHasNFT  bool
	sellerHasCash bool
	buyerHasNFT   bool
	buyerHasCash  bool
	log           []string
}

func run(w *world, seller, buyer agent) outcome {
	end := w.t.T1 + 120
	for ; w.now <= end; w.now++ {
		w.commit()
		for _, a := range []agent{seller, buyer} {
			if !a.online(w.now - base) {
				continue
			}
			d := NextAction(a.role, w.t, w.p, w.observe(a.role))
			if a.override != nil {
				d = a.override(w, d)
			}
			w.submit(a.role, d.Kind)
		}
	}
	f, b := w.fab.State, w.besu.State
	return outcome{
		status:        Derive(Observation{Fabric: w.fab, Besu: w.besu}),
		sellerHasNFT:  f == FabricNone || f == FabricReclaimed || f == FabricLocked,
		sellerHasCash: b == BesuClaimed,
		buyerHasNFT:   f == FabricClaimed,
		buyerHasCash:  b == BesuNone || b == BesuRefunded || b == BesuLocked,
		log:           w.log,
	}
}

func (o outcome) sellerLost() bool { return !o.sellerHasNFT && !o.sellerHasCash }
func (o outcome) buyerLost() bool  { return !o.buyerHasNFT && !o.buyerHasCash }

// ---------------------------------------------------------------------------------------------
// Cenários (seção 4, Fase 6, do plano), nas duas abordagens.
// ---------------------------------------------------------------------------------------------

func protocols() []Protocol { return []Protocol{ProtocolHTLC, ProtocolCacti} }

func TestSim_1_HappyPath(t *testing.T) {
	for _, proto := range protocols() {
		t.Run(string(proto), func(t *testing.T) {
			o := run(newWorld(proto), honest(RoleSeller), honest(RoleBuyer))
			if o.status != StatusCompleted {
				t.Fatalf("status=%s\n%v", o.status, o.log)
			}
		})
	}
}

func TestSim_2_BuyerNeverLocks(t *testing.T) {
	for _, proto := range protocols() {
		t.Run(string(proto), func(t *testing.T) {
			buyer := honest(RoleBuyer)
			buyer.override = func(_ *world, d Decision) Decision {
				if d.Kind == ActLockPayment {
					return wait("adversário: não trava", time.Time{})
				}
				return d
			}
			o := run(newWorld(proto), honest(RoleSeller), buyer)
			if o.status != StatusReverted || o.sellerLost() {
				t.Fatalf("status=%s\n%v", o.status, o.log)
			}
		})
	}
}

func TestSim_3_SellerNeverReveals_OrRelayDown(t *testing.T) {
	for _, proto := range protocols() {
		t.Run(string(proto), func(t *testing.T) {
			w := newWorld(proto)
			seller := honest(RoleSeller)
			if proto == ProtocolCacti {
				w.relayOK = false // relay parado: a Vendedora honesta nunca consegue a atestação
			} else {
				seller.override = func(_ *world, d Decision) Decision {
					if d.Kind == ActClaimPayment {
						return wait("adversário: não revela", time.Time{})
					}
					return d
				}
			}
			o := run(w, seller, honest(RoleBuyer))
			if o.status != StatusReverted || o.buyerLost() || o.sellerLost() {
				t.Fatalf("status=%s\n%v", o.status, o.log)
			}
		})
	}
}

func TestSim_4_RevealAtT2MinusEpsilon(t *testing.T) {
	// Vendedora maliciosa segura S e só revela no último segundo antes de T2.
	// O Banco honesto ainda precisa conseguir fazer o claim antes de T1.
	w := newWorld(ProtocolHTLC)
	seller := honest(RoleSeller)
	seller.override = func(w *world, d Decision) Decision {
		if w.besu.State == BesuLocked && w.now == w.t.T2-w.latency-1 {
			return do(ActClaimPayment, "adversário: revelação tardia")
		}
		if d.Kind == ActClaimPayment {
			return wait("adversário: segura S", time.Time{})
		}
		return d
	}
	o := run(w, seller, honest(RoleBuyer))
	if o.status != StatusCompleted || o.buyerLost() {
		t.Fatalf("status=%s\n%v", o.status, o.log)
	}
}

func TestSim_4b_RevealTooLateFailsButLeaksS(t *testing.T) {
	// A revelação chega depois de T2: a tx falha, mas S vaza no mempool. O Banco faz o refund
	// e o claim do NFT. Por isso a Vendedora honesta respeita RevealMargin.
	w := newWorld(ProtocolHTLC)
	seller := honest(RoleSeller)
	seller.override = func(w *world, d Decision) Decision {
		if w.besu.State == BesuLocked && w.now == w.t.T2-1 {
			return do(ActClaimPayment, "imprudente: revelação após a margem")
		}
		if d.Kind == ActClaimPayment {
			return wait("segura S", time.Time{})
		}
		return d
	}
	// Banco ganancioso: depois do refund, usa o S vazado para pegar também o NFT.
	greedy := honest(RoleBuyer)
	greedy.override = func(w *world, d Decision) Decision {
		if w.buyerS != nil && w.fab.State == FabricLocked && w.ts(w.now).Before(w.fab.Deadline) &&
			w.besu.State == BesuRefunded {
			return do(ActClaimAsset, "adversário: claim com S vazado")
		}
		return d
	}
	o := run(w, seller, greedy)
	if !o.sellerLost() || o.status != StatusBuyerHasBoth {
		t.Fatalf("esperado: a Vendedora imprudente perde\nstatus=%s\n%v", o.status, o.log)
	}
	// Já o Banco HONESTO, no mesmo cenário, não fica com as duas pernas.
	w2 := newWorld(ProtocolHTLC)
	if o2 := run(w2, seller, honest(RoleBuyer)); o2.status == StatusBuyerHasBoth {
		t.Fatalf("Banco honesto nunca fica com as duas pernas\n%v", o2.log)
	}
}

func TestSim_5_BuyerCrashAfterRevealAndRecovers(t *testing.T) {
	for _, proto := range protocols() {
		t.Run(string(proto), func(t *testing.T) {
			w := newWorld(proto)
			down := false
			var downAt int64
			buyer := honest(RoleBuyer)
			// cai logo depois de a Vendedora revelar e volta dentro de (T1 − T2)
			buyer.online = func(rel int64) bool {
				if !down && w.besu.State == BesuClaimed && downAt == 0 {
					down, downAt = true, rel
				}
				if down && rel > downAt+int64(w.p.Delta()/time.Second) {
					down = false
				}
				return !down
			}
			o := run(w, honest(RoleSeller), buyer)
			if o.status != StatusCompleted {
				t.Fatalf("status=%s\n%v", o.status, o.log)
			}
		})
	}
}

func TestSim_5b_BuyerOfflinePastT1_KnownRisk(t *testing.T) {
	// Risco documentado (plano, seção 6): watchtower fora do ar por mais de T1 − T2 depois da revelação.
	w := newWorld(ProtocolHTLC)
	buyer := honest(RoleBuyer)
	buyer.online = func(rel int64) bool { return w.besu.State != BesuClaimed }
	o := run(w, honest(RoleSeller), buyer)
	if !o.buyerLost() {
		t.Fatalf("esperado: perda do Banco ausente\nstatus=%s\n%v", o.status, o.log)
	}
	// e a Vendedora honesta não aproveitou para fazer o reclaim (o NFT continua travado)
	if w.fab.State != FabricLocked {
		t.Fatalf("Vendedora honesta não deveria fazer reclaim depois de paga: %s", w.fab.State)
	}
}

func TestSim_6_WrongTermsRejected(t *testing.T) {
	// NFT travado com deadline bem menor que T1 (relógio da Vendedora adiantado ou má-fé).
	for _, proto := range protocols() {
		t.Run(string(proto), func(t *testing.T) {
			w := newWorld(proto)
			w.lockDeadlineSkew = -5 * time.Minute
			o := run(w, honest(RoleSeller), honest(RoleBuyer))
			if o.buyerLost() || o.status == StatusCompleted {
				t.Fatalf("status=%s\n%v", o.status, o.log)
			}
		})
	}
}

func TestSim_8_EarlyReclaimRejected(t *testing.T) {
	// Vendedora tenta o reclaim antecipado (peer com relógio forjado). Com a política
	// AND(Registradora, Banco), o endosso usa relógios honestos: no mundo, a regra é a do commit.
	for _, proto := range protocols() {
		t.Run(string(proto), func(t *testing.T) {
			w := newWorld(proto)
			seller := honest(RoleSeller)
			seller.override = func(w *world, d Decision) Decision {
				if w.besu.State == BesuClaimed && w.fab.State == FabricLocked {
					return do(ActReclaimAsset, "adversário: reclaim antecipado")
				}
				return d
			}
			o := run(w, seller, honest(RoleBuyer))
			if o.status != StatusCompleted || o.buyerLost() {
				t.Fatalf("status=%s\n%v", o.status, o.log)
			}
		})
	}
}

// Varredura: os dois honestos, com quedas aleatórias curtas (< Δ) e latências variadas, nunca
// terminam em estado anômalo nem com perda.
func TestSim_SweepHonestWithOutages(t *testing.T) {
	for _, proto := range protocols() {
		for lat := int64(1); lat <= 10; lat += 3 {
			for start := int64(0); start < 360; start += 17 {
				for _, who := range []Role{RoleSeller, RoleBuyer} {
					w := newWorld(proto)
					w.latency = lat
					outage := int64(w.p.Delta()/time.Second) - lat - 5
					a := map[Role]agent{RoleSeller: honest(RoleSeller), RoleBuyer: honest(RoleBuyer)}
					s := start
					a[who] = agent{role: who, online: func(rel int64) bool { return rel < s || rel >= s+outage }}
					o := run(w, a[RoleSeller], a[RoleBuyer])
					if o.sellerLost() || o.buyerLost() || o.status == StatusSellerHasBoth || o.status == StatusBuyerHasBoth {
						t.Fatalf("%s lat=%d outage[%s]=%d+%d: status=%s\n%v", proto, lat, who, s, outage, o.status, o.log)
					}
				}
			}
		}
	}
}
