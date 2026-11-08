# Janus

Liquidação atômica interbancária **cross-chain (DvP)**: o recebível tokenizado (NFT do Panurus) fica no **Hyperledger Fabric** e o pagamento em **Drex** (ERC-20 permissionado) fica no **Hyperledger Besu**. Nenhum dos dois ativos sai da própria rede. A troca é atômica por:

- **Abordagem A, HTLC:** HTLC nativo do Panurus no Fabric + `HTLCPayment.sol` no Besu. Usa SHA-256 e timelocks escalonados (`T1 > T2`).
- **Abordagem B2, HTLC atestado:** o Besu só libera o pagamento com a preimagem **e** uma atestação do lock no Fabric, assinada por uma política k-de-n de orgs (P-256 nativo). Ver o ADR-0003.

O plano completo está em [`docs/implementation_plan.md`](docs/implementation_plan.md), as decisões em [`docs/adr/`](docs/adr), o [modelo de ameaças](docs/threat-model.md) e o [runbook](docs/runbook.md).

## Estrutura

| Caminho | Conteúdo |
|---|---|
| `contracts/` | Foundry: `DrexToken`, `HTLCPayment`, `FabricProofVerifier`, `AttestedSettlement` |
| `node/internal/swap` | núcleo puro: termos, política de tempo, máquina de estados, `NextAction`, simulador |
| `node/internal/token` | views FSC/Panurus: emissão, lock/claim/reclaim, status, auditoria da Registradora |
| `node/internal/besu` | cliente Besu (bindings abigen v2, watcher com checkpoint) |
| `node/internal/engine` | laço observar → decidir → executar (watchtower/watchdog), com falhas injetáveis |
| `node/internal/protocol/htlc`, `cacti` | Abordagem A e Abordagem B2 (HTLC atestado) sobre portas Fabric/Besu |
| `node/internal/janus`, `node/cmd/janus-node` | SDK do FSC que monta o nó Janus e o binário único por instituição |
| `node/internal/swap/sim`, `node/cmd/simwasm` | simulador determinístico (9 cenários), também compilado para WebAssembly |
| `apps/web` | simulador web (React + Vite + TS): modo offline (WASM) e modo live (SSE dos nós) |
| `node/internal/api`, `store`, `secrets` | REST + SSE, SQLite, S cifrado (AES-GCM) |
| `node/topology` | topologia NWO e testes de integração com rede Fabric real |
| `infra/besu` | QBFT com 4 validadores, gas zero, allowlist |
| `testdata/vectors` | vetores compartilhados Go ↔ Solidity |

## Pré-requisitos

Go 1.27.3, Foundry ≥ 1.8, Docker, `jq` e Node ≥ 24 (simulador). As versões estão fixadas no [ADR-0001](docs/adr/0001-versoes-fixadas.md).

```sh
make deps          # libs do Foundry + módulos Go
make fabric-bins   # binários do Fabric 3.1.5 / CA 1.5.22 (só para a integração)
```

## Testes

| Nível | Comando |
|---|---|
| Vetores Go ↔ Solidity | `make test-vectors` |
| Contratos (unit, fuzz, invariantes) | `make test-contracts` |
| Nó (lógica, motor, API, cliente Besu via anvil) | `make test-node` |
| Integração Fabric/Panurus (rede NWO real) | `make it` |
| E2E (Fabric NWO + anvil + 3 nós Janus): A, B2, cenário 3, cenário 5 | `make it` (os testes `TestE2E*`) |
| Simulador web | `npm run test -w @janus/web` |

## Simulador web

```sh
npm install
npm run web        # compila o núcleo para WASM e sobe em http://localhost:5173
```

Os links aceitam parâmetros: `?protocol=cacti&scenario=seller-no-reveal&t=end`. O modo **Live** conversa com as APIs dos nós (`/api/v1/...`).

## Redes locais

```sh
make besu-up       # Besu QBFT em http://localhost:8545
cd contracts && forge script script/Deploy.s.sol --rpc-url http://localhost:8545 --broadcast --legacy --with-gas-price 0
make besu-down
```

> As chaves em `infra/besu/config/permissions_config.toml` e em `contracts/script/Deploy.s.sol` são as contas **públicas de desenvolvimento** do Foundry/anvil. Nunca use essas chaves fora do ambiente local.
