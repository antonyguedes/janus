// Package sim é um simulador determinístico do protocolo: duas cadeias com as regras dos
// contratos (HTLCPayment/AttestedSettlement) e do validador HTLC do Panurus, e dois nós que
// decidem com swap.NextAction, o mesmo código dos nós reais. Cada transação leva Latency
// segundos e é revalidada no instante do commit.
//
// Ele alimenta o modo offline do simulador web (compilado para WebAssembly em cmd/simwasm) e
// os testes de propriedade do protocolo.
package sim

import (
	"crypto/sha256"
	"fmt"
	"time"

	"janus/node/internal/swap"
)

// Base é a epoch arbitrária do simulador (t=0).
const Base = 1_900_000_000

// Scenario identifica um cenário pré-definido (seção 4, Fase 6, do plano).
type Scenario string

const (
	Happy              Scenario = "happy"                 // 1
	BuyerNoLock        Scenario = "buyer-no-lock"         // 2
	SellerNoReveal     Scenario = "seller-no-reveal"      // 3 (A) / relay parado (B)
	RevealAtT2Minus    Scenario = "reveal-t2-minus"       // 4
	BuyerCrashRecover  Scenario = "buyer-crash-recover"   // 5
	BuyerOfflinePastT1 Scenario = "buyer-offline-past-t1" // risco documentado
	WrongDeadline      Scenario = "wrong-deadline"        // 6 (termos divergentes)
	EarlyReclaim       Scenario = "early-reclaim"         // 8
	LateRevealLeak     Scenario = "late-reveal-leak"      // revelação após a margem (didático)
)

// Scenarios lista os cenários com descrição (para a UI).
var Scenarios = []struct {
	ID          Scenario `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Success     bool     `json:"success"`
}{
	{Happy, "Caminho feliz", "Os dois nós seguem o protocolo: lock, lock, revelação, claim.", true},
	{BuyerNoLock, "Banco não trava", "O Banco nunca trava o Drex: o NFT volta à Vendedora depois de T1.", false},
	{SellerNoReveal, "Vendedora não revela / relay parado", "Sem revelação (A) ou sem atestação (B): as duas pernas revertem.", false},
	{RevealAtT2Minus, "Revelação em T2 − ε", "A Vendedora segura S até o último segundo. O Banco ainda faz o claim antes de T1.", true},
	{BuyerCrashRecover, "Banco cai e volta", "O Banco cai logo após a revelação e reconcilia dentro de T1 − T2.", true},
	{BuyerOfflinePastT1, "Banco fora além de T1", "Risco documentado: watchtower ausente por mais de T1 − T2.", false},
	{WrongDeadline, "Deadline divergente", "O NFT é travado com deadline menor que o combinado: o Banco não paga.", false},
	{EarlyReclaim, "Reclaim antecipado", "A Vendedora tenta o reclaim depois de paga: rejeitado pela regra de tempo.", true},
	{LateRevealLeak, "Revelação imprudente", "Revelação depois da margem: a tx falha, S vaza e um Banco ganancioso fica com tudo.", false},
}

// Config do simulador.
type Config struct {
	Protocol swap.Protocol `json:"protocol"`
	Scenario Scenario      `json:"scenario"`
	Latency  int64         `json:"latency"` // segundos por transação (padrão 4)
}

// Event é uma transação comitada (ou rejeitada) numa das cadeias.
type Event struct {
	T      int64           `json:"t"` // segundos desde o início
	Actor  swap.Role       `json:"actor"`
	Action swap.ActionKind `json:"action"`
	Chain  string          `json:"chain"` // fabric | besu | relay
	OK     bool            `json:"ok"`
	Note   string          `json:"note,omitempty"`
}

// Snapshot é o estado completo num instante (o inspetor da UI mostra isso).
type Snapshot struct {
	T            int64            `json:"t"`
	Fabric       swap.FabricState `json:"fabric"`
	FabricLabel  string           `json:"fabricLabel"`
	Besu         swap.BesuState   `json:"besu"`
	BesuLabel    string           `json:"besuLabel"`
	Status       swap.Status      `json:"status"`
	Seller       swap.Decision    `json:"seller"`
	Buyer        swap.Decision    `json:"buyer"`
	SellerOnline bool             `json:"sellerOnline"`
	BuyerOnline  bool             `json:"buyerOnline"`
	Revealed     bool             `json:"revealed"`       // S público no Besu
	Attested     bool             `json:"attested"`       // B: atestações prontas
	BuyerKnowsS  bool             `json:"buyerKnowsS"`    // o Banco tem S
	Deadline     int64            `json:"fabricDeadline"` // deadline do lock no Fabric (relativo)
}

// Result é a execução completa.
type Result struct {
	Config    Config     `json:"config"`
	T1        int64      `json:"t1"`
	T2        int64      `json:"t2"`
	Delta     int64      `json:"delta"`
	Snapshots []Snapshot `json:"snapshots"` // só os instantes com mudança
	Events    []Event    `json:"events"`
	Outcome   Outcome    `json:"outcome"`
}

// Outcome resume quem ficou com o quê.
type Outcome struct {
	Status        swap.Status `json:"status"`
	SellerHasNFT  bool        `json:"sellerHasNft"`
	SellerHasCash bool        `json:"sellerHasCash"`
	BuyerHasNFT   bool        `json:"buyerHasNft"`
	BuyerHasCash  bool        `json:"buyerHasCash"`
}

func (o Outcome) SellerLost() bool { return !o.SellerHasNFT && !o.SellerHasCash }
func (o Outcome) BuyerLost() bool  { return !o.BuyerHasNFT && !o.BuyerHasCash }

// ------------------------------------------------------------------------------- mundo

type pending struct {
	at   int64
	kind swap.ActionKind
	by   swap.Role
}

// World é o estado das duas cadeias e o conhecimento dos nós.
type World struct {
	T      *swap.SwapTerms
	P      swap.TimePolicy
	Now    int64
	Fab    swap.FabricView
	Besu   swap.BesuView
	queue  []pending
	attAt  int64
	events []Event

	secret   [32]byte
	SellerS  *[32]byte
	BuyerS   *[32]byte
	Latency  int64
	RelayOK  bool
	RelayLat int64
	// LockDeadlineSkew: desvio do relógio da Vendedora ao calcular script.Deadline.
	LockDeadlineSkew time.Duration
}

// NewWorld cria termos válidos (DevPolicy) com t=0 em Base.
func NewWorld(proto swap.Protocol) *World {
	p := swap.DevPolicy()
	s := sha256.Sum256([]byte("sim-secret"))
	h := sha256.Sum256(s[:])
	seller := swap.Party{MSPID: "VendedoraMSP", FSCIdentity: []byte("vendedora"), EVM: swap.Address{1}}
	buyer := swap.Party{MSPID: "BancoMSP", FSCIdentity: []byte("banco"), EVM: swap.Address{2}}
	nonce := [32]byte{9}
	t := &swap.SwapTerms{
		Protocol: proto, AssetID: "recebivel-0001", Seller: seller, Buyer: buyer,
		Price: swap.Amount{Token: swap.Address{0xd}, Value: 1_000_00}, Hashlock: h, Nonce: nonce,
		T1: Base + int64(p.T1Offset/time.Second), T2: Base + int64(p.T2Offset/time.Second),
		Networks: swap.NetworkRefs{FabricChannel: "recebiveis", FabricNamespace: "tcc", BesuChainID: 1337, PaymentContract: swap.Address{0xc}},
	}
	t.SwapID = swap.ComputeSwapID(h, seller, buyer, nonce)
	return &World{T: t, P: p, Now: Base, secret: s, SellerS: &s, Latency: 4, RelayOK: true, RelayLat: 8}
}

func ts(x int64) time.Time { return time.Unix(x, 0) }

// Observe monta a observação de um papel (o que o nó veria).
func (w *World) Observe(role swap.Role) swap.Observation {
	o := swap.Observation{Now: ts(w.Now), Fabric: w.Fab, Besu: w.Besu}
	if role == swap.RoleSeller {
		o.Preimage = w.SellerS
		o.AttestationReady = w.attAt != 0 && w.Now >= w.attAt
	} else {
		o.Preimage = w.BuyerS
	}
	return o
}

func (w *World) inFlight(k swap.ActionKind) bool {
	for _, q := range w.queue {
		if q.kind == k {
			return true
		}
	}
	return false
}

// Submit envia a ação de um nó (vira transação com latência).
func (w *World) Submit(by swap.Role, k swap.ActionKind) {
	switch k {
	case swap.ActWait, swap.ActDone, swap.ActAbort, swap.ActAlert:
		return
	case swap.ActRequestAttest:
		if w.attAt == 0 && w.RelayOK && w.Fab.State != swap.FabricNone {
			w.attAt = w.Now + w.RelayLat
			w.events = append(w.events, Event{T: w.Now - Base, Actor: by, Action: k, Chain: "relay", OK: true,
				Note: fmt.Sprintf("atestações prontas em t=%d", w.attAt-Base)})
		}
		return
	}
	if w.inFlight(k) {
		return
	}
	w.queue = append(w.queue, pending{at: w.Now + w.Latency, kind: k, by: by})
}

func chainOf(k swap.ActionKind) string {
	switch k {
	case swap.ActLockAsset, swap.ActClaimAsset, swap.ActReclaimAsset:
		return "fabric"
	}
	return "besu"
}

// Commit aplica as transações maduras com as regras on-chain.
func (w *World) Commit() {
	var rest []pending
	for _, q := range w.queue {
		if q.at > w.Now {
			rest = append(rest, q)
			continue
		}
		now := ts(w.Now)
		ok, note := false, ""
		switch q.kind {
		case swap.ActLockAsset:
			if w.Fab.State == swap.FabricNone {
				w.Fab = swap.FabricView{
					State: swap.FabricLocked, AssetID: w.T.AssetID, Sender: w.T.Seller.FSCIdentity,
					Recipient: w.T.Buyer.FSCIdentity, Hashlock: w.T.Hashlock, HashOK: true,
					Deadline: w.T.T1Time().Add(w.LockDeadlineSkew), TxID: "tx-lock", Block: 7,
				}
				ok = true
			}
		case swap.ActLockPayment:
			if w.Besu.State == swap.BesuNone && now.Before(w.T.T2Time()) {
				w.Besu = swap.BesuView{
					State: swap.BesuLocked, Contract: w.T.Networks.PaymentContract, Token: w.T.Price.Token,
					Sender: w.T.Buyer.EVM, Recipient: w.T.Seller.EVM, Amount: w.T.Price.Value,
					Hashlock: w.T.Hashlock, Timelock: w.T.T2Time(),
				}
				if w.T.Protocol == swap.ProtocolCacti {
					w.Besu.MinFabricDeadline = w.T.T1Time().Add(-w.P.DeadlineTolerance)
				}
				ok = true
			}
		case swap.ActClaimPayment, swap.ActSettlePayment:
			valid := w.Besu.State == swap.BesuLocked && now.Before(w.Besu.Timelock)
			if q.kind == swap.ActSettlePayment {
				valid = valid && w.attAt != 0 && w.Now >= w.attAt && w.Fab.State != swap.FabricNone &&
					!w.Fab.Deadline.Before(w.Besu.MinFabricDeadline)
			}
			// A preimagem fica pública no mempool mesmo que a tx falhe (pior caso para a Vendedora).
			s := w.secret
			w.BuyerS = &s
			if valid {
				w.Besu.State, w.Besu.Preimage = swap.BesuClaimed, &s
				ok = true
			} else {
				note = "revertida; S exposto no mempool"
			}
		case swap.ActClaimAsset:
			if w.Fab.State == swap.FabricLocked && now.Before(w.Fab.Deadline) && w.BuyerS != nil {
				w.Fab.State = swap.FabricClaimed
				ok = true
			}
		case swap.ActRefundPayment:
			if w.Besu.State == swap.BesuLocked && !now.Before(w.Besu.Timelock) {
				w.Besu.State = swap.BesuRefunded
				ok = true
			}
		case swap.ActReclaimAsset:
			if w.Fab.State == swap.FabricLocked && !now.Before(w.Fab.Deadline) {
				w.Fab.State = swap.FabricReclaimed
				ok = true
			} else if w.Fab.State == swap.FabricLocked {
				note = "rejeitada: deadline ainda não atingido"
			}
		}
		w.events = append(w.events, Event{T: w.Now - Base, Actor: q.by, Action: q.kind, Chain: chainOf(q.kind), OK: ok, Note: note})
	}
	w.queue = rest
}

// Events devolve as transações processadas até agora.
func (w *World) Events() []Event { return w.events }

// Agent é um comportamento de nó: honesto (NextAction) ou adversarial (Override).
type Agent struct {
	Role     swap.Role
	Online   func(rel int64) bool
	Override func(w *World, d swap.Decision) swap.Decision
}

// Honest devolve um agente que segue NextAction e está sempre online.
func Honest(r swap.Role) Agent { return Agent{Role: r, Online: func(int64) bool { return true }} }

// Outcome calcula quem ficou com o quê.
func (w *World) Outcome() Outcome {
	f, b := w.Fab.State, w.Besu.State
	return Outcome{
		Status:        swap.Derive(swap.Observation{Fabric: w.Fab, Besu: w.Besu}),
		SellerHasNFT:  f == swap.FabricNone || f == swap.FabricReclaimed || f == swap.FabricLocked,
		SellerHasCash: b == swap.BesuClaimed,
		BuyerHasNFT:   f == swap.FabricClaimed,
		BuyerHasCash:  b == swap.BesuNone || b == swap.BesuRefunded || b == swap.BesuLocked,
	}
}

// RunAgents executa até T1 + 120 s e devolve o resultado com snapshots só nas mudanças.
func RunAgents(w *World, seller, buyer Agent) Result {
	res := Result{T1: w.T.T1 - Base, T2: w.T.T2 - Base, Delta: int64(w.P.Delta() / time.Second)}
	end := w.T.T1 + 120
	var prev *Snapshot
	for ; w.Now <= end; w.Now++ {
		w.Commit()
		snap := Snapshot{T: w.Now - Base}
		for _, a := range []Agent{seller, buyer} {
			online := a.Online(w.Now - Base)
			if a.Role == swap.RoleSeller {
				snap.SellerOnline = online
			} else {
				snap.BuyerOnline = online
			}
			if !online {
				continue
			}
			d := swap.NextAction(a.Role, w.T, w.P, w.Observe(a.Role))
			if a.Override != nil {
				d = a.Override(w, d)
			}
			d.WakeAt = time.Time{} // não serializar epochs na UI
			if a.Role == swap.RoleSeller {
				snap.Seller = d
			} else {
				snap.Buyer = d
			}
			w.Submit(a.Role, d.Kind)
		}
		snap.Fabric, snap.Besu = w.Fab.State, w.Besu.State
		snap.FabricLabel, snap.BesuLabel = w.Fab.State.String(), w.Besu.State.String()
		snap.Status = swap.Derive(swap.Observation{Fabric: w.Fab, Besu: w.Besu})
		snap.Revealed = w.Besu.Preimage != nil
		snap.BuyerKnowsS = w.BuyerS != nil
		snap.Attested = w.attAt != 0 && w.Now >= w.attAt
		if !w.Fab.Deadline.IsZero() {
			snap.Deadline = w.Fab.Deadline.Unix() - Base
		}
		if prev == nil || changed(*prev, snap) {
			res.Snapshots = append(res.Snapshots, snap)
			s := snap
			prev = &s
		}
	}
	res.Events = w.events
	res.Outcome = w.Outcome()
	return res
}

func changed(a, b Snapshot) bool {
	a.T, b.T = 0, 0
	return a != b
}
