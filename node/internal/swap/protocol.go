package swap

import "context"

// SettlementProtocol é implementado por cada abordagem (protocol/htlc, protocol/cacti).
// Toda a decisão está em NextAction. A implementação só observa as cadeias e executa ações.
type SettlementProtocol interface {
	Name() Protocol
	// Observe monta a observação combinada das duas cadeias para este nó.
	Observe(ctx context.Context, role Role, t *SwapTerms) (Observation, error)
	// Execute realiza a ação decidida. Precisa ser idempotente: o motor pode repetir uma ação
	// cujo efeito ainda não apareceu na observação.
	Execute(ctx context.Context, role Role, t *SwapTerms, d Decision) error
}
