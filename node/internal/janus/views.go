package janus

import (
	"errors"

	fscview "github.com/hyperledger-labs/fabric-smart-client/platform/view/services/view"

	"janus/node/internal/token"
)

// RegisterViews registra todas as views do Janus. Cada nó registra todas: o papel efetivo vem
// do janus.json e das identidades do nó, não do conjunto de views.
func RegisterViews(r *fscview.Registry) error {
	return errors.Join(
		r.RegisterFactory("issue", &token.IssueViewFactory{}),
		r.RegisterFactory("gravame", &token.GravameViewFactory{}),
		r.RegisterFactory("registerAuditor", &token.RegisterAuditorViewFactory{}),
		r.RegisterFactory("get", &token.GetViewFactory{}),
		r.RegisterFactory("lock", &token.LockViewFactory{}),
		r.RegisterFactory("reclaim", &token.ReclaimViewFactory{}),
		r.RegisterFactory("claim", &token.ClaimViewFactory{}),
		r.RegisterFactory("expect", &token.ExpectLockViewFactory{}),
		r.RegisterFactory("status", &token.StatusViewFactory{}),
		r.RegisterFactory("propose", &token.ProposeViewFactory{}),
		r.RegisterFactory("collect", &token.CollectViewFactory{}),
		r.RegisterResponder(&token.AcceptIssuedView{}, &token.IssueView{}),
		r.RegisterResponder(&token.LockAcceptView{}, &token.LockView{}),
		r.RegisterResponder(&token.AcceptProposalView{}, &token.ProposeView{}),
		r.RegisterResponder(&token.AttestResponderView{}, &token.CollectView{}),
	)
}
