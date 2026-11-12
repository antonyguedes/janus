<p align="center">
  <img src="docs/assets/janus-banner.svg" alt="Janus: dois portais, o dourado do Fabric e o ciano do Besu, guardados por Jano" width="100%">
</p>

<p align="center">
  <img alt="Hyperledger Fabric 3.1.5" src="https://img.shields.io/badge/Fabric-3.1.5-E7A83A?style=for-the-badge&labelColor=0B0920">
  <img alt="Panurus v0.27.0" src="https://img.shields.io/badge/Panurus-v0.27.0-FFD98A?style=for-the-badge&labelColor=0B0920">
  <img alt="Hyperledger Besu 26.11" src="https://img.shields.io/badge/Besu-26.11_QBFT-19C3E6?style=for-the-badge&labelColor=0B0920">
  <img alt="Hyperledger Cacti v3.0.1" src="https://img.shields.io/badge/Cacti-v3.0.1-E05CF2?style=for-the-badge&labelColor=0B0920">
  <br>
  <img alt="Go 1.27" src="https://img.shields.io/badge/Go-1.27-B49BFF?style=flat-square&labelColor=151233">
  <img alt="Solidity 0.8.30" src="https://img.shields.io/badge/Solidity-0.8.30_osaka-B49BFF?style=flat-square&labelColor=151233">
  <img alt="React 19" src="https://img.shields.io/badge/React-19_+_WASM-B49BFF?style=flat-square&labelColor=151233">
  <img alt="8 cenários em rede real" src="https://img.shields.io/badge/cenários_do_plano-8%2F8-3DDC97?style=flat-square&labelColor=151233">
</p>

<p align="center">
  <b>Liquidação atômica interbancária cross-chain (DvP)</b><br>
  o recebível tokenizado vive no <b>Hyperledger Fabric</b>, o Drex vive no <b>Hyperledger Besu</b>,<br>
  e <b>nenhum dos dois ativos sai da própria rede</b>.
</p>

---

## ⛩ As duas faces de Jano

Na mitologia romana, **Jano** é o deus das portas e das passagens. Ele tem duas faces, uma voltada para cada lado do umbral. O Janus é exatamente isso: um guardião entre dois mundos que nunca se misturam.

<table>
<tr>
<td width="50%" valign="top">

### 🟡 Portal do Fabric
**o recebível**

- NFT emitido pela **Registradora** (Panurus, FabToken)
- Travado no **HTLC nativo do Panurus**
- Política de endosso `AND(Registradora, Banco)`
- Auditoria com relógio próprio e gravames persistidos

</td>
<td width="50%" valign="top">

### 🔵 Portal do Besu
**o pagamento**

- **Drex** ERC-20 permissionado (2 casas, allowlist)
- `HTLCPayment` (A) ou `AttestedSettlement` (B2)
- QBFT com 4 validadores e gas zero
- P-256 nativo (EIP-7951) para verificar atestações

</td>
</tr>
</table>

A **chave** que abre os dois portais é a preimagem **S**: quem a revela recebe o Drex, e com ela o Banco reivindica o recebível. Os prazos `T1 > T2` garantem que, se algo der errado, **cada um volta com o que trouxe**.

## 🌀 Duas formas de atravessar

| | **A · HTLC** | **B2 · HTLC atestado** |
|---|---|---|
| Fabric | HTLC do Panurus (SHA-256, deadline `T1`) | o mesmo |
| Besu | `HTLCPayment`: claim com `S` antes de `T2` | `AttestedSettlement`: libera com `S` **e** uma atestação 2-de-2 do lock |
| Prova do lock | o Banco confere no próprio vault | **Hyperledger Cacti**: conector Fabric por org (bloco, validade, endossos), atestadores P-256, relay Besu do Cacti submete o `settle` |
| Ordem | NFT trava → Drex trava | o Drex pode travar primeiro (o contrato exige a prova) |
| ADRs | [0002](docs/adr/0002-spike-panurus-htlc.md), [0004](docs/adr/0004-tempo-no-fabric.md) | [0003](docs/adr/0003-spike-cacti-abordagem-b.md), [0005](docs/adr/0005-cacti-prova-de-ledger.md), [0006](docs/adr/0006-cacti-relay-besu.md) |

```mermaid
%%{init: {"theme":"base","themeVariables":{"background":"#0b0920","primaryColor":"#151233","primaryTextColor":"#ece9ff","primaryBorderColor":"#7c4dff","lineColor":"#b49bff","secondaryColor":"#1d1945","tertiaryColor":"#151233","actorBkg":"#151233","actorBorder":"#b49bff","actorTextColor":"#ece9ff","signalColor":"#b49bff","signalTextColor":"#ece9ff","noteBkgColor":"#1d1945","noteTextColor":"#ffd98a","noteBorderColor":"#e7a83a","sequenceNumberColor":"#0b0920"}}}%%
sequenceDiagram
    autonumber
    participant V as Nó Vendedora
    participant F as 🟡 Fabric (Panurus HTLC)
    participant B as 🔵 Besu (Drex)
    participant C as Nó Banco
    V->>C: SwapTerms por sessão P2P do FSC (os dois assinam, S nunca trafega)
    V->>F: trava o NFT com H = SHA256(S) até T1
    C->>F: confere o lock (recipient, H, T1)
    C->>B: trava o Drex com H até T2
    V->>B: claim(S), revela S
    B-->>C: evento com S (watchtower)
    C->>F: htlc.Claim(S) antes de T1, o NFT atravessa
    Note over F,B: falhou? refund no Besu depois de T2 e reclaim no Fabric depois de T1
```

## 🏛 Arquitetura

```mermaid
%%{init: {"theme":"base","themeVariables":{"background":"#0b0920","primaryColor":"#151233","primaryTextColor":"#ece9ff","primaryBorderColor":"#7c4dff","lineColor":"#b49bff","clusterBkg":"#110e2a","clusterBorder":"#2d2860","edgeLabelBackground":"#151233"}}}%%
flowchart LR
    subgraph FAB["🟡 Hyperledger Fabric 3.1.5 · canal do recebível"]
        TCC["Token Chaincode do Panurus<br/>HTLC nativo · FabToken"]
    end
    subgraph BESU["🔵 Hyperledger Besu 26.11 · QBFT · gas zero"]
        TK["DrexToken"]
        HT["HTLCPayment · A"]
        VB["FabricProofVerifier +<br/>AttestedSettlement · B2"]
    end
    subgraph CACTI["🟣 Cacti v3.0.1"]
        CF1["conector Fabric<br/>Registradora"]
        CF2["conector Fabric<br/>Banco"]
        CB["conector Besu<br/>relay da Vendedora"]
    end
    REG["Nó Janus · Registradora<br/>emissora · auditora · atestadora"]
    VEN["Nó Janus · Vendedora"]
    BAN["Nó Janus · Banco<br/>watchtower · atestador"]
    UI["Simulador web<br/>React + WASM"]

    REG & VEN & BAN --> TCC
    VEN <-->|"sessão P2P FSC"| BAN
    VEN & BAN --> HT
    BAN --> VB
    REG -.prova do lock.-> CF1 --> TCC
    BAN -.prova do lock.-> CF2 --> TCC
    VEN -.settle.-> CB --> VB
    HT & VB --> TK
    UI -->|REST + SSE| VEN & BAN & REG

    style FAB fill:#1c1608,stroke:#e7a83a,color:#ffd98a
    style BESU fill:#071a22,stroke:#19c3e6,color:#8cf0ff
    style CACTI fill:#1d0d24,stroke:#e05cf2,color:#f2b8fb
```

Cada instituição roda **um binário** (`janus-node`): um nó do Fabric Smart Client com o Panurus, as views do Janus, o cliente Besu, o motor de liquidação, a API REST + SSE, um store SQLite e S cifrado (AES-GCM). Não existe custodiante central.

## 🔑 Atravessando: a demonstração completa

```sh
make deps          # Foundry, módulos Go e npm (simulador + Cacti)
make fabric-bins   # binários do Fabric 3.1.5 / CA 1.5.22 (uma vez)

make besu-up && make demo-besu    # Fabric + Panurus + Besu QBFT + Cacti + 3 nós Janus
# ou: make demo                    # o mesmo, com anvil no lugar do Besu QBFT

npm run web                        # outro terminal → http://localhost:5173
```

| Portal | Endereço |
|---|---|
| 🌐 Simulador web | http://localhost:5173 (`?mode=live` para os nós reais) |
| 🏛 Registradora · Vendedora · Banco | `:8080` · `:8081` · `:8082` (`/api/v1/...`) |
| 🟣 Cacti | `:4100` Registradora · `:4101` Banco (conector Fabric) · `:4102` relay (conector Besu) |
| 🔵 Besu QBFT | http://localhost:8545 |

**Pelo terminal:**

```sh
NFT=$(curl -s -X POST localhost:8080/api/v1/recebiveis -H 'content-type: application/json' \
  -d '{"owner":"vendedora","registro":"REG-001","cedente":"V","sacado":"S","valorFace":100000000,"vencimento":"2027-03-31"}' | jq -r .linearId)

SWAP=$(curl -s -X POST localhost:8081/api/v1/swaps -H 'content-type: application/json' \
  -d "{\"buyerNode\":\"banco\",\"assetId\":\"$NFT\",\"price\":90000000,\"protocol\":\"cacti\"}" | jq -r .swapId)

curl -s localhost:8082/api/v1/swaps/$SWAP | jq .status     # → "completed"
curl -sN localhost:8081/api/v1/events                        # passos do motor ao vivo
```

## 🖥 O simulador

<p align="center">
  <img src="docs/assets/simulador-escuro.png" alt="Simulador do Janus, tema escuro, abordagem B2 concluída" width="49%">
  <img src="docs/assets/simulador-claro.png" alt="Simulador do Janus, tema claro, cenário de reversão" width="49%">
</p>

- **Offline:** os 9 cenários (sucesso e falha) rodam com o **mesmo `swap.NextAction` dos nós**, compilado para WebAssembly.
- **Live:** emite, propõe e acompanha swaps reais pelos nós, via SSE.
- Os dois portais mostram o estado de cada perna e da chave S. Há relógio com `T2`/`T1`, diagrama de sequência animado, inspetor por rede, tema escuro ("noite dos portais") e claro ("aurora"), e links compartilháveis: `?theme=dark&protocol=cacti&scenario=seller-no-reveal&t=end`.

## 🧪 Provas

| Nível | Comando | O que cobre |
|---|---|---|
| Vetores Go ↔ Solidity | `make test-vectors` | SHA-256 e digest/assinaturas P-256 idênticos nos dois lados |
| Contratos | `make test-contracts` | 51 testes, fuzz nas fronteiras de tempo e 4 invariantes |
| Nó | `make test-node` | núcleo, simulador, motor, API, cliente Besu (anvil), Cacti |
| Simulador | `make test-web` | reprodução e typecheck |
| Rede real | `make e2e` | Fabric 3.1.5 + Panurus (NWO) + Besu + Cacti + 3 nós Janus |

**Os 8 cenários do plano:**

| # | Cenário | A | B2 | Onde |
|---|---|:-:|:-:|---|
| 1 | Caminho feliz | ✅ | ✅ | A e B2 em rede real (também no Besu QBFT) |
| 2 | Banco não trava → o NFT volta depois de `T1` | ✅ | ✅ | A em rede real · B2 no simulador e no motor |
| 3 | Vendedora não revela / relay parado → as duas pernas revertem | ✅ | ✅ | A em rede real · B2 no simulador; relay e Cacti fora do ar verificados no demo |
| 4 | Revelação em `T2 − ε` → o Banco ainda faz o claim | ✅ | — | A em rede real |
| 5 | O nó cai e volta → reconciliação | ✅ | ✅ | A em rede real (Banco) · B2 em rede real (Registradora reiniciada) |
| 6 | Deadline divergente → o Banco recusa o lock | ✅ | ✅ | A em rede real · B2 no simulador |
| 7 | Prova adulterada ou endossante insuficiente → rejeitada | — | ✅ | contratos + EVM real · conector Cacti |
| 8 | Reclaim antecipado → recusado pela auditoria | ✅ | ✅ | rede real (mesmo HTLC do Panurus nas duas) |

## 📜 Pergaminhos

| Documento | |
|---|---|
| [Plano de implementação](docs/implementation_plan.md) | escopo, arquitetura e fases |
| [ADR-0001](docs/adr/0001-versoes-fixadas.md) · [0002](docs/adr/0002-spike-panurus-htlc.md) · [0003](docs/adr/0003-spike-cacti-abordagem-b.md) · [0004](docs/adr/0004-tempo-no-fabric.md) · [0005](docs/adr/0005-cacti-prova-de-ledger.md) · [0006](docs/adr/0006-cacti-relay-besu.md) | versões, Panurus, Cacti, tempo, prova de ledger, relay |
| [Modelo de ameaças](docs/threat-model.md) | 17 ameaças, mitigações, evidências e riscos residuais |
| [Runbook](docs/runbook.md) | `janus.json`, operação e incidentes |
| [infra/fabric](infra/fabric/README.md) | como a rede Fabric é gerada (NWO) |

<details>
<summary><b>🗺 Mapa do repositório</b></summary>

| Caminho | Conteúdo |
|---|---|
| `contracts/` | Foundry: `DrexToken`, `HTLCPayment`, `FabricProofVerifier`, `AttestedSettlement` |
| `node/internal/swap` · `swap/sim` | núcleo puro (termos, tempo, `NextAction`) e simulador determinístico |
| `node/internal/token` | views do Panurus: emissão, lock/claim/reclaim, status, auditoria, negociação, atestação |
| `node/internal/besu` · `cacti` | cliente Besu (abigen v2, watcher) e clientes dos conectores do Cacti |
| `node/internal/engine` · `protocol/{htlc,cacti}` | motor (watchtower/watchdog, falhas injetáveis) e as duas abordagens |
| `node/internal/janus` · `cmd/janus-node` | SDK do FSC que monta o nó e o binário único por instituição |
| `node/internal/{api,store,secrets}` | REST + SSE, SQLite (swaps e estado durável), S cifrado |
| `node/topology` | topologia NWO, testes de integração e E2E, runner `janus-net` |
| `node/ci/external-builders/golang` | external builder do chaincode (Go 1.27 no host) |
| `infra/besu` · `infra/cacti/relay` | QBFT local e servidores Cacti |
| `apps/web` · `node/cmd/simwasm` | simulador web e núcleo em WASM |
| `testdata/vectors` | vetores compartilhados Go ↔ Solidity |

</details>

---

<sub>
⚠️ As chaves em <code>infra/besu/config/permissions_config.toml</code>, <code>contracts/script/Deploy.s.sol</code> e nos testes são as contas <b>públicas de desenvolvimento</b> do Foundry/anvil. Nunca use essas chaves fora do ambiente local. Keychains em memória e servidores Cacti sem autorização servem apenas para desenvolvimento (ver o runbook).<br><br>
O tema visual é original, inspirado na figura mitológica de Jano, deus romano das passagens. Não usa arte, marcas nem personagens de terceiros, e não tem afiliação com jogos ou franquias que retratem Jano.
</sub>
