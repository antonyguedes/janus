// Command janus-node é o binário único de cada instituição (Registradora, Vendedora ou Banco):
// um nó FSC com o Panurus (FabToken), as views do Janus e a SDK do Janus (Besu, motor, API).
// O papel vem do janus.json ao lado do core.yaml do FSC (FSCNODE_CFG_PATH).
package main

import (
	fscnode "github.com/hyperledger-labs/fabric-smart-client/node"
	viewregistry "github.com/hyperledger-labs/fabric-smart-client/platform/view/services/view"

	"janus/node/internal/janus"
)

func main() {
	n := fscnode.New()
	n.InstallSDK(janus.NewTokenSDK(n))
	n.InstallSDK(janus.NewSDK(n))
	n.Execute(func() error {
		return janus.RegisterViews(viewregistry.GetRegistry(n))
	})
}
