# ADR-0001 — Versões fixadas

- **Status:** aceito (Fase 0)
- **Data:** 2026-11-08
- **Contexto:** decisão 5 do plano ("sempre o mais atualizado"). As versões abaixo são as últimas estáveis em 2026-11-08, conferidas no proxy de módulos Go, nos releases do GitHub e no Docker Hub. Qualquer mudança aqui exige rodar de novo `make test` e os testes de integração.

| Componente | Versão | Fonte / evidência |
|---|---|---|
| Go | **1.27.3** | toolchain local. O `go.mod` do Panurus v0.27.0 já exige `go 1.27.3` |
| Panurus (ex-Fabric Token SDK) | **v0.27.0** (`github.com/LFDT-Panurus/panurus`, commit `632947b`) | proxy.golang.org, publicado em 2026-11-02 |
| Fabric Smart Client | **v0.27.0** (a versão que o Panurus v0.27.0 fixa) | `go.mod` do Panurus. A v0.28.0 saiu em 2026-11-06, mas o Panurus ainda não a adotou. Usamos a combinação testada upstream |
| Hyperledger Fabric (peer/orderer) | **3.1.5** | GitHub releases, Docker Hub `hyperledger/fabric-peer:3.1.5` |
| Fabric CA | **1.5.22** | GitHub releases |
| Hyperledger Besu | **26.11.0** | GitHub (`besu-eth/besu`), Docker Hub `hyperledger/besu:26.11.0`. Validado localmente: 4 validadores QBFT, gas zero, precompile P256VERIFY ativo |
| Solidity | **0.8.30** no código (`pragma ^0.8.30`) e EVM `osaka` | a 0.8.30 é a primeira com `osaka` como alvo padrão. O compilador vem do Foundry |
| Foundry | **≥ 1.8.1** (local). A última é a v1.9.0 | `forge --version`. Recomendado `foundryup` para a 1.9.0 |
| OpenZeppelin Contracts | **v5.7.0** | GitHub releases |
| go-ethereum (cliente + abigen) | **v1.17.10** | proxy.golang.org |
| Hyperledger Cacti | **v3.0.1**: pacotes JS `@hyperledger-cacti/*` 3.0.1 (`cactus-cmd-api-server`, `cactus-plugin-ledger-connector-fabric`, `cactus-plugin-keychain-memory`) | GitHub releases (2026-08-21) e npm. Atenção: em `@hyperledger/*` a última versão é a 2.0.0 (2024) |
| Node.js (simulador) | **≥ 24** (local v24.19.0) | |

## Consequências

- Panurus e FSC andam juntos. O **upgrade do FSC só acontece quando o Panurus o adotar** no `go.mod` dele, para evitar uma combinação que ninguém testou.
- O Besu 26.11.0 tem só 6 dias. Se surgir uma regressão, o fallback documentado é a **26.10.2**, porque o genesis não usa nenhum recurso exclusivo da 26.11.
- O Panurus avisa que **não foi auditado**. Ver o ADR-0002 e o modelo de ameaças.

## Adendo (Fase 2): build do Token Chaincode

- A imagem `hyperledger/fabric-ccenv:3.1` traz **Go 1.26.4**, mas o Token Chaincode do Panurus v0.27.0 exige **Go 1.27.3**. Com o build padrão via Docker, a instalação falha. O erro aparece mascarado como `failed to marshal response: string field contains invalid UTF-8`, porque o log do build truncado vai dentro da resposta gRPC. O mesmo acontece com Fabric 3.1.4 e 3.1.5, então não é regressão do Fabric.
- **Decisão:** o chaincode é compilado e executado por um **external builder do Janus** (`node/ci/external-builders/golang`: `detect`/`build`/`release`/`run`). Ele compila no host com o toolchain do projeto (`CGO_ENABLED=0`, `-trimpath`) e roda o processo com TLS mútuo a partir do `chaincode.json` fornecido pelo peer. O NWO do Panurus registra esse builder automaticamente (`<GOMOD>/../ci/external-builders/golang`).
- O NWO empacota o chaincode no layout GOPATH quando o `main` está num módulo de dependência. As dependências são copiadas do module cache, já verificado pelo `go.sum` do projeto.
- Os nós FSC usam **persistência SQLite** (`WithDefaultSqlitePersistence`) para vault, tokendb, auditdb e ttxdb.
