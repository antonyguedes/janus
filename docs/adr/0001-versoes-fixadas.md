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
| Hyperledger Cacti | **v3.0.1** | GitHub releases (2026-08-21) |
| Node.js (simulador) | **≥ 24** (local v24.19.0) | |

## Consequências

- Panurus e FSC andam juntos. O **upgrade do FSC só acontece quando o Panurus o adotar** no `go.mod` dele, para evitar uma combinação que ninguém testou.
- O Besu 26.11.0 tem só 6 dias. Se surgir uma regressão, o fallback documentado é a **26.10.2**, porque o genesis não usa nenhum recurso exclusivo da 26.11.
- O Panurus avisa que **não foi auditado**. Ver o ADR-0002 e o modelo de ameaças.
