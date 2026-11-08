package it

import (
	"github.com/LFDT-Panurus/panurus/integration/nwo/token"
	fabric2 "github.com/LFDT-Panurus/panurus/integration/nwo/token/fabric"
	"github.com/LFDT-Panurus/panurus/integration/nwo/token/generators/crypto/fabtokenv1"
	views2 "github.com/LFDT-Panurus/panurus/integration/token/common/views"
	"github.com/hyperledger-labs/fabric-smart-client/integration/nwo/api"
	"github.com/hyperledger-labs/fabric-smart-client/integration/nwo/fabric"
	"github.com/hyperledger-labs/fabric-smart-client/integration/nwo/fsc"

	jt "janus/node/internal/token"
)

// Orgs do consórcio. O NWO deriva o MSP ID como <Nome>MSP.
const (
	OrgRegistradora = "Registradora"
	OrgVendedora    = "Vendedora"
	OrgBanco        = "Banco"
)

// PolicyJanus é a política de endosso do Token Chaincode (ADR-0004, camada 1): toda
// transação precisa dos peers da Registradora E do Banco. Um peer da Vendedora, mesmo com
// relógio manipulado, não basta para aprovar um reclaim antecipado.
const PolicyJanus = "AND('RegistradoraMSP.member','BancoMSP.member')"

// Options parametriza a topologia (os spikes trocam a política).
type Options struct {
	NamespacePolicy   string // vazio = PolicyJanus
	EndorseOwnOrgOnly bool   // cada nó coleta endosso só dos peers da própria org (spike de tempo)
}

// Topology descreve a rede do Janus: Fabric 3.x com canal padrão, Token Chaincode do Panurus
// (driver FabToken) e três nós FSC, um por instituição.
func Topology(o Options) []api.Topology {
	ft := fabric.NewDefaultTopology()
	ft.EnableIdemix()
	ft.AddOrganizationsByName(OrgRegistradora, OrgVendedora, OrgBanco)
	ft.SetNamespaceApproverOrgs(OrgRegistradora)
	channel := ft.Channels[0].Name

	ft2 := fsc.NewTopology()

	// Registradora: emissora + auditora (ADR-0002)
	ft2.AddNodeByName(jt.AuditorNode).
		AddOptions(
			fabric.WithOrganization(OrgRegistradora),
			fabric.WithAnonymousIdentity(),
			fabric.WithDefaultSqlitePersistence(), // vault do FSC + tokendb/auditdb em SQLite
			token.WithDefaultIssuerIdentity(false),
			token.WithAuditorIdentity(false),
		).
		RegisterViewFactory("registerAuditor", &jt.RegisterAuditorViewFactory{}).
		RegisterViewFactory("issue", &jt.IssueViewFactory{}).
		RegisterViewFactory("gravame", &jt.GravameViewFactory{}).
		RegisterViewFactory("TxFinality", &views2.TxFinalityViewFactory{})

	ft2.AddNodeByName("vendedora").
		AddOptions(
			fabric.WithOrganization(OrgVendedora),
			fabric.WithAnonymousIdentity(),
			fabric.WithDefaultSqlitePersistence(),
			token.WithDefaultOwnerIdentity(),
		).
		RegisterResponder(&jt.AcceptIssuedView{}, &jt.IssueView{}).
		RegisterViewFactory("get", &jt.GetViewFactory{}).
		RegisterViewFactory("lock", &jt.LockViewFactory{}).
		RegisterViewFactory("reclaim", &jt.ReclaimViewFactory{}).
		RegisterViewFactory("status", &jt.StatusViewFactory{}).
		RegisterViewFactory("TxFinality", &views2.TxFinalityViewFactory{})

	ft2.AddNodeByName("banco").
		AddOptions(
			fabric.WithOrganization(OrgBanco),
			fabric.WithAnonymousIdentity(),
			fabric.WithDefaultSqlitePersistence(),
			token.WithDefaultOwnerIdentity(),
		).
		RegisterResponder(&jt.AcceptIssuedView{}, &jt.IssueView{}).
		RegisterResponder(&jt.LockAcceptView{}, &jt.LockView{}).
		RegisterViewFactory("get", &jt.GetViewFactory{}).
		RegisterViewFactory("expect", &jt.ExpectLockViewFactory{}).
		RegisterViewFactory("claim", &jt.ClaimViewFactory{}).
		RegisterViewFactory("status", &jt.StatusViewFactory{}).
		RegisterViewFactory("TxFinality", &views2.TxFinalityViewFactory{})

	tt := token.NewTopology()
	tms := tt.AddTMS(ft2.ListNodes(), ft, channel, fabtokenv1.DriverIdentifier)
	tms.SetTokenGenPublicParams("16")
	fabric2.SetOrgs(tms, OrgRegistradora, OrgVendedora, OrgBanco)
	policy := o.NamespacePolicy
	if policy == "" {
		policy = PolicyJanus
	}
	fabric2.WithNamespacePolicy(tms, policy)
	if o.EndorseOwnOrgOnly {
		fabric2.WithEndorsersFromOwnOrg(tms)
	}
	tms.AddAuditor(ft2.ListNodes(jt.AuditorNode)[0])
	tms.AddIssuerByID(jt.AuditorNode)

	return []api.Topology{ft, tt, ft2}
}
