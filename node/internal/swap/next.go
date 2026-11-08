package swap

import (
	"fmt"
	"time"
)

// ActionKind é a próxima ação que o nó deve executar.
type ActionKind string

const (
	ActWait ActionKind = "wait" // nada a fazer agora. Reavaliar em WakeAt ou no próximo evento

	ActLockAsset     ActionKind = "lock_asset"     // Vendedora: htlc.Lock do NFT (H, recipient=Banco, T1)
	ActLockPayment   ActionKind = "lock_payment"   // Banco: HTLCPayment.lock (A) / AttestedSettlement.deposit (B)
	ActClaimPayment  ActionKind = "claim_payment"  // Vendedora (A): HTLCPayment.claim(S), revela S
	ActRequestAttest ActionKind = "request_attest" // Vendedora (B): coletar atestações do lock via relay
	ActSettlePayment ActionKind = "settle_payment" // Vendedora (B): AttestedSettlement.settle(S, att), revela S
	ActClaimAsset    ActionKind = "claim_asset"    // Banco: htlc.Claim(S) no Fabric
	ActRefundPayment ActionKind = "refund_payment" // Banco: refund depois de T2
	ActReclaimAsset  ActionKind = "reclaim_asset"  // Vendedora: htlc.Reclaim depois de T1
	ActDone          ActionKind = "done"           // terminal para este nó
	ActAbort         ActionKind = "abort"          // não prosseguir. O swap não chegou a comprometer este nó
	ActAlert         ActionKind = "alert"          // situação anômala: exige operador
)

// Decision é o resultado de NextAction.
type Decision struct {
	Kind   ActionKind `json:"kind"`
	Reason string     `json:"reason"`
	// WakeAt: quando reavaliar mesmo sem eventos (o watchdog usa). Zero = só por evento.
	WakeAt time.Time `json:"wakeAt,omitzero"`
}

func wait(reason string, at time.Time) Decision {
	return Decision{Kind: ActWait, Reason: reason, WakeAt: at}
}
func do(k ActionKind, reason string) Decision { return Decision{Kind: k, Reason: reason} }

// NextAction concentra TODA a lógica de decisão do protocolo para um nó honesto.
// Ela é pura: a mesma entrada sempre produz a mesma decisão.
//
// Princípios:
//   - a fronteira de tempo é idêntica nas duas cadeias: ação "antes de T" se now < T, reversão se now ≥ T;
//   - a Vendedora só revela S se o pagamento estiver travado corretamente e faltar mais que RevealMargin para T2;
//   - o Banco só trava o pagamento depois de validar o lock do NFT (A) e, assim que vê S, faz o claim no Fabric;
//   - um nó honesto nunca fica com as duas pernas: a Vendedora paga não faz reclaim.
func NextAction(role Role, t *SwapTerms, p TimePolicy, o Observation) Decision {
	switch role {
	case RoleSeller:
		return sellerNext(t, p, o)
	case RoleBuyer:
		return buyerNext(t, p, o)
	}
	return do(ActAlert, fmt.Sprintf("papel desconhecido %q", role))
}

func sellerNext(t *SwapTerms, p TimePolicy, o Observation) Decision {
	now, T1, T2 := o.Now, t.T1Time(), t.T2Time()
	f, b := o.Fabric, o.Besu

	// 1) Já recebi o pagamento.
	if b.State == BesuClaimed {
		switch f.State {
		case FabricClaimed:
			return do(ActDone, "swap concluído: pagamento recebido e NFT entregue")
		case FabricLocked:
			if now.Before(T1) {
				return wait("pago; aguardando o Banco fazer o claim do NFT", T1)
			}
			return do(ActAlert, "pago, mas o Banco não fez o claim até T1. O nó honesto NÃO faz reclaim: resolver com o operador")
		default:
			return do(ActAlert, fmt.Sprintf("pagamento recebido com o NFT em estado %s", f.State))
		}
	}

	switch f.State {
	case FabricReclaimed:
		return do(ActDone, "swap revertido: NFT reavido")

	case FabricClaimed:
		// O Banco tem o NFT. Se S ainda não está no Besu, receber agora.
		if b.State == BesuLocked && o.Preimage != nil && now.Before(b.Timelock) {
			return do(pickPaymentAction(t, o), "NFT já foi reivindicado: receber o pagamento")
		}
		return do(ActAlert, "NFT reivindicado sem o pagamento recebido")

	case FabricNone:
		if b.State == BesuRefunded {
			return do(ActDone, "swap abortado antes do lock do NFT")
		}
		if !now.Before(T2.Add(-(p.BuyerLockMargin + p.RevealMargin))) {
			return do(ActAbort, "janela para o lock do NFT expirou")
		}
		if o.Preimage == nil {
			return do(ActAlert, "preimagem S não encontrada no armazenamento seguro")
		}
		return do(ActLockAsset, "termos assinados: travar o NFT no HTLC do Panurus")
	}

	// f.State == FabricLocked e pagamento não recebido.
	if err := CheckFabricLock(t, p, f); err != nil {
		return do(ActAlert, "o próprio lock do NFT diverge dos termos: "+err.Error())
	}
	reclaimOr := func(reason string) Decision {
		if !now.Before(T1) {
			return do(ActReclaimAsset, "T1 atingido: reaver o NFT ("+reason+")")
		}
		return wait(reason+"; reclaim em T1", T1)
	}

	switch b.State {
	case BesuNone:
		return reclaimOr("aguardando o Banco travar o pagamento")
	case BesuRefunded:
		return reclaimOr("pagamento devolvido ao Banco")
	}

	// b.State == BesuLocked
	if err := CheckBesuLock(t, p, b); err != nil {
		return reclaimOr("lock do pagamento inválido, NÃO revelar S: " + err.Error())
	}
	revealDeadline := b.Timelock.Add(-p.RevealMargin)
	if !now.Before(revealDeadline) {
		return reclaimOr("tarde demais para revelar S com segurança")
	}
	if o.Preimage == nil {
		return do(ActAlert, "preimagem S não encontrada no armazenamento seguro")
	}
	if t.Protocol == ProtocolCacti && !o.AttestationReady {
		return Decision{Kind: ActRequestAttest, Reason: "coletar atestações do lock no Fabric", WakeAt: revealDeadline}
	}
	return do(pickPaymentAction(t, o), "pagamento travado e válido: revelar S e receber")
}

func pickPaymentAction(t *SwapTerms, o Observation) ActionKind {
	if t.Protocol == ProtocolCacti {
		if !o.AttestationReady {
			return ActRequestAttest
		}
		return ActSettlePayment
	}
	return ActClaimPayment
}

func buyerNext(t *SwapTerms, p TimePolicy, o Observation) Decision {
	now, T1, T2 := o.Now, t.T1Time(), t.T2Time()
	f, b := o.Fabric, o.Besu

	// A preimagem vem do evento do Besu ou do conhecimento local.
	preimage := b.Preimage
	if preimage == nil {
		preimage = o.Preimage
	}

	// 1) O NFT já é meu.
	if f.State == FabricClaimed {
		if b.State == BesuClaimed {
			return do(ActDone, "swap concluído: NFT recebido e pagamento entregue")
		}
		if b.State == BesuLocked && !now.Before(b.Timelock) {
			return do(ActAlert, "NFT recebido e pagamento ainda travado depois de T2. O nó honesto NÃO faz refund: resolver com o operador")
		}
		return wait("NFT recebido; aguardando a Vendedora receber o pagamento", b.Timelock)
	}

	// 2) A Vendedora foi paga (S está público): fazer o claim do NFT antes de T1, com prioridade máxima.
	if b.State == BesuClaimed {
		switch f.State {
		case FabricLocked:
			if preimage == nil {
				return do(ActAlert, "pagamento reivindicado mas a preimagem não foi capturada: varrer eventos")
			}
			if now.Before(T1) {
				return do(ActClaimAsset, "S revelado no Besu: fazer o claim do NFT antes de T1")
			}
			return do(ActAlert, "PERDA: S foi revelado mas o claim do NFT não aconteceu antes de T1")
		case FabricReclaimed:
			return do(ActAlert, "PERDA: pagamento entregue e NFT reavido pela Vendedora")
		default:
			return do(ActAlert, fmt.Sprintf("pagamento entregue com o NFT em estado %s", f.State))
		}
	}

	// 3) Pagamento travado e ainda não reivindicado.
	if b.State == BesuLocked {
		if !now.Before(b.Timelock) {
			return do(ActRefundPayment, "T2 atingido: reaver o pagamento")
		}
		if preimage != nil && f.State == FabricLocked && now.Before(T1) && CheckFabricLock(t, p, f) == nil {
			return do(ActClaimAsset, "S conhecido: fazer o claim do NFT")
		}
		return wait("pagamento travado; aguardando a revelação de S", b.Timelock)
	}

	if b.State == BesuRefunded {
		return do(ActDone, "swap revertido: pagamento reavido")
	}

	// 4) b.State == BesuNone: decidir se trava o pagamento.
	lockDeadline := T2.Add(-p.BuyerLockMargin)
	if f.State == FabricReclaimed {
		return do(ActDone, "NFT reavido pela Vendedora antes do pagamento: nada a fazer")
	}
	if !now.Before(lockDeadline) {
		return do(ActAbort, "janela para travar o pagamento expirou")
	}
	if f.State == FabricLocked {
		if err := CheckFabricLock(t, p, f); err != nil {
			return do(ActAbort, "lock do NFT diverge dos termos, NÃO pagar: "+err.Error())
		}
		return do(ActLockPayment, "NFT travado e válido: travar o pagamento")
	}
	// f.State == FabricNone
	if t.Protocol == ProtocolCacti {
		// B2: o contrato só libera o pagamento com prova do lock, então o depósito antecipado é seguro.
		return do(ActLockPayment, "abordagem B: depositar (liberação condicionada à atestação do lock)")
	}
	return wait("aguardando a Vendedora travar o NFT", lockDeadline)
}
