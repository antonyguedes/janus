# infra/fabric

A rede Hyperledger Fabric 3.1.5 do Janus é **gerada programaticamente** pelo NWO (Network Orchestrator do Fabric Smart Client), seguindo o caminho recomendado pelo Panurus. Não há `docker-compose` aqui:

| O quê | Onde |
|---|---|
| Topologia (3 orgs, canal, Token Chaincode do Panurus, política `AND(Registradora, Banco)`, nós FSC) | [`node/topology/it/topology.go`](../../node/topology/it/topology.go) |
| Subida da rede para testes e demonstração | [`node/topology/it/network.go`](../../node/topology/it/network.go), [`node/topology/cmd/janus-net`](../../node/topology/cmd/janus-net) |
| External builder do chaincode (Go 1.27 no host, ADR-0001) | [`node/ci/external-builders/golang`](../../node/ci/external-builders/golang) |
| Binários `peer`, `orderer`, `configtxgen`, … | `make fabric-bins` baixa para `infra/fabric/bin` (não versionado) |

Os testes ponta a ponta ficam em [`node/topology/it`](../../node/topology/it) (`make e2e`).
