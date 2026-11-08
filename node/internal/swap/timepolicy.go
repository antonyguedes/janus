package swap

import (
	"fmt"
	"time"
)

// TimePolicy define os prazos do swap e a margem Δ (ADR-0004):
//
//	T1 − T2 ≥ Δ = t_detecção + t_finalidade_Fabric + t_relay + ε_skew + m_operacional
type TimePolicy struct {
	T1Offset time.Duration // prazo do lock do ativo (Fabric), contado a partir da proposta
	T2Offset time.Duration // prazo do lock do pagamento (Besu)

	Detection      time.Duration // tempo para o watchtower detectar a revelação no Besu
	FabricFinality time.Duration // endosso + ordenação + commit no Fabric
	Relay          time.Duration // coleta e submissão de prova (Abordagem B)
	Skew           time.Duration // tolerância de relógio entre nós/peers + período de bloco do Besu
	Operational    time.Duration // margem operacional (retries, filas)

	// RevealMargin: a Vendedora NÃO revela S se faltar menos que isso para T2. Assim, um claim que
	// chegue depois de T2 nunca expõe S enquanto o Banco ainda pode fazer o refund.
	RevealMargin time.Duration
	// BuyerLockMargin: o Banco só trava o pagamento se ainda houver pelo menos isso até T2,
	// para dar tempo de a Vendedora conferir e revelar.
	BuyerLockMargin time.Duration
	// DeadlineTolerance: divergência máxima aceita entre o deadline do script HTLC no Fabric
	// (calculado com o relógio de quem trava) e o T1 combinado.
	DeadlineTolerance time.Duration
}

// Delta é a margem mínima exigida entre T1 e T2.
func (p TimePolicy) Delta() time.Duration {
	return p.Detection + p.FabricFinality + p.Relay + p.Skew + p.Operational
}

// Validate confere a consistência interna da política.
func (p TimePolicy) Validate() error {
	if p.T2Offset <= 0 || p.T1Offset <= p.T2Offset {
		return fmt.Errorf("timepolicy: need 0 < T2 < T1 (T1=%s, T2=%s)", p.T1Offset, p.T2Offset)
	}
	if gap := p.T1Offset - p.T2Offset; gap < p.Delta() {
		return fmt.Errorf("timepolicy: T1-T2=%s < Δ=%s", gap, p.Delta())
	}
	if p.T2Offset <= p.RevealMargin+p.BuyerLockMargin {
		return fmt.Errorf("timepolicy: T2=%s leaves no window for lock+reveal", p.T2Offset)
	}
	if p.DeadlineTolerance > p.Skew {
		return fmt.Errorf("timepolicy: deadline tolerance %s must be <= skew %s", p.DeadlineTolerance, p.Skew)
	}
	return nil
}

// ProductionPolicy: T1 = 4h, T2 = 2h (plano, seção 3.4).
func ProductionPolicy() TimePolicy {
	return TimePolicy{
		T1Offset: 4 * time.Hour, T2Offset: 2 * time.Hour,
		Detection: 5 * time.Minute, FabricFinality: 2 * time.Minute, Relay: 5 * time.Minute,
		Skew: 30*time.Second + 2*time.Second, Operational: 30 * time.Minute,
		RevealMargin: 30 * time.Minute, BuyerLockMargin: 45 * time.Minute,
		DeadlineTolerance: 30 * time.Second,
	}
}

// DevPolicy: prazos em minutos para dev/E2E (T1 = 6 min, T2 = 3 min, Δ ≥ 90 s).
func DevPolicy() TimePolicy {
	return TimePolicy{
		T1Offset: 6 * time.Minute, T2Offset: 3 * time.Minute,
		Detection: 10 * time.Second, FabricFinality: 10 * time.Second, Relay: 15 * time.Second,
		Skew: 32 * time.Second, Operational: 23 * time.Second,
		RevealMargin: 45 * time.Second, BuyerLockMargin: 60 * time.Second,
		DeadlineTolerance: 30 * time.Second,
	}
}

// E2EPolicy: prazos curtos para os testes ponta a ponta em rede real (T2 = 75 s, T1 = 170 s, Δ = 90 s).
func E2EPolicy() TimePolicy {
	p := DevPolicy()
	p.T1Offset, p.T2Offset = 170*time.Second, 75*time.Second
	p.RevealMargin, p.BuyerLockMargin = 20*time.Second, 25*time.Second
	return p
}
