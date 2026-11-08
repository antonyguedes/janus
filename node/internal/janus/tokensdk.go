package janus

import (
	"errors"

	fabtoken "github.com/LFDT-Panurus/panurus/token/core/fabtoken/v1/driver"
	"github.com/LFDT-Panurus/panurus/token/sdk"
	tokensdk "github.com/LFDT-Panurus/panurus/token/sdk/dig"
	"github.com/LFDT-Panurus/panurus/token/services/network/fabric"
	dig2 "github.com/hyperledger-labs/fabric-smart-client/platform/common/sdk/dig"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services"
	"go.uber.org/dig"
)

// TokenSDK é a SDK de produção do Panurus para o Janus: rede Fabric + driver FabToken.
// Equivale à ffabtoken usada nos testes do Panurus, sem a dependência do módulo de integração
// (comunicação P2P por websocket, o padrão do FSC).
type TokenSDK struct {
	dig2.SDK
}

func NewTokenSDK(registry services.Registry) *TokenSDK {
	return &TokenSDK{SDK: tokensdk.NewSDK(registry)}
}

func (p *TokenSDK) Install() error {
	err := errors.Join(
		sdk.RegisterTokenDriverDependencies(p.Container()),
		p.Container().Provide(fabric.NewGenericDriver, dig.Group("network-drivers")),
		p.Container().Provide(fabtoken.NewTokenDriver, dig.Group("token-drivers")),
		p.Container().Provide(fabtoken.NewValidatorDriver, dig.Group("validator-drivers")),
	)
	if err != nil {
		return err
	}
	return p.SDK.Install()
}
