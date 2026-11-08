package token

import (
	"fmt"
	"time"

	"github.com/LFDT-Panurus/panurus/token"
	"github.com/LFDT-Panurus/panurus/token/services/interop/htlc"
	"github.com/LFDT-Panurus/panurus/token/services/ttx"
	token2 "github.com/LFDT-Panurus/panurus/token/token"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/view"
)

// AuditView: a Registradora audita TODA transação de token (emissão, lock, claim, reclaim).
//   - recusa movimentar recebível com gravame;
//   - valida o script de novos locks HTLC (SHA256/None, limites de deadline);
//   - confere claim/reclaim com o PRÓPRIO relógio (ADR-0004, camada 2). Reclaim só depois de
//     deadline + ReclaimGrace, mesmo que um endossante com relógio adiantado já tenha aprovado.
type AuditView struct{}

func (a *AuditView) Call(ctx view.Context) (any, error) {
	tx, err := ttx.ReceiveTransaction(ctx)
	if err != nil {
		return nil, err
	}
	if err := tx.IsValid(ctx.Context()); err != nil {
		return nil, fmt.Errorf("audit: invalid tx: %w", err)
	}
	w := ttx.MyAuditorWallet(ctx, token.WithTMSID(tx.TokenService().ID()))
	if w == nil {
		return nil, fmt.Errorf("audit: auditor wallet not found")
	}
	auditor, err := ttx.NewAuditor(ctx, w)
	if err != nil {
		return nil, err
	}
	if err := auditor.Validate(ctx.Context(), tx); err != nil {
		return nil, fmt.Errorf("audit: validate: %w", err)
	}
	inputs, outputs, err := auditor.Audit(ctx.Context(), tx)
	if err != nil {
		return nil, fmt.Errorf("audit: %w", err)
	}
	defer auditor.Release(ctx.Context(), tx)

	if err := CheckAudit(DefaultAudit, DefaultGravames, inputs, outputs, time.Now()); err != nil {
		return nil, err
	}
	return ctx.RunView(ttx.NewAuditApproveView(w, tx))
}

// CheckAudit aplica a política a uma transação já decodificada.
func CheckAudit(p AuditPolicy, g *Gravames, inputs *token.InputStream, outputs *token.OutputStream, now time.Time) error {
	checkType := func(t token2.Type) error {
		if r, err := FromTokenType(t); err == nil {
			return g.Check(r.LinearID)
		}
		return nil // tokens que não são recebíveis (não usados no MVP) passam
	}
	for i := range inputs.Count() {
		in := inputs.At(i)
		if err := checkType(in.Type); err != nil {
			return err
		}
		hin, err := htlc.ToInput(in)
		if err != nil {
			return fmt.Errorf("audit: input %d: %w", i, err)
		}
		if !hin.IsHTLC() {
			continue
		}
		s, err := hin.Script()
		if err != nil {
			return fmt.Errorf("audit: input %d script: %w", i, err)
		}
		out, err := outputForAction(outputs, in.ActionIndex)
		if err != nil {
			return err
		}
		op, err := ClassifySpend(s, out.Owner)
		if err != nil {
			return err
		}
		if err := p.CheckSpend(s, op, now); err != nil {
			return err
		}
	}
	for i := range outputs.Count() {
		out := outputs.At(i)
		if err := checkType(out.Type); err != nil {
			return err
		}
		ho, err := htlc.ToOutput(out)
		if err != nil {
			return fmt.Errorf("audit: output %d: %w", i, err)
		}
		if !ho.IsHTLC() {
			continue
		}
		s, err := ho.Script()
		if err != nil {
			return fmt.Errorf("audit: output %d script: %w", i, err)
		}
		if err := p.CheckLockScript(s, now); err != nil {
			return err
		}
	}
	return nil
}

func outputForAction(outputs *token.OutputStream, action int) (*token.Output, error) {
	var found *token.Output
	for i := range outputs.Count() {
		o := outputs.At(i)
		if o.ActionIndex == action {
			if found != nil {
				return nil, fmt.Errorf("audit: htlc spend must have exactly one output")
			}
			found = o
		}
	}
	if found == nil {
		return nil, fmt.Errorf("audit: htlc spend without output")
	}
	return found, nil
}

// RegisterAuditorView registra AuditView como responder de auditoria do nó.
type RegisterAuditorView struct{}

func (r *RegisterAuditorView) Call(ctx view.Context) (any, error) {
	return ctx.RunView(ttx.NewRegisterAuditorView(&AuditView{}))
}

type RegisterAuditorViewFactory struct{}

func (f *RegisterAuditorViewFactory) NewView([]byte) (view.View, error) {
	return &RegisterAuditorView{}, nil
}
