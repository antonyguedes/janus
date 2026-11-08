# Runbook — operação do nó Janus

## Componentes por instituição

Um processo `janus-node` (nó FSC + Panurus + SDK do Janus) por instituição. Ele lê:

- `core.yaml` do FSC (identidades, peers, persistência SQLite) em `FSCNODE_CFG_PATH`;
- `janus.json` no mesmo diretório (ou `JANUS_CONFIG`):

```json
{
  "role": "vendedora | banco | registradora",
  "mspId": "VendedoraMSP",
  "counterparty": "banco",
  "attestors": ["registradora", "banco"],
  "policy": "prod",
  "dataDir": "/var/lib/janus",
  "api": { "listen": "127.0.0.1:8081", "allowOrigin": "https://simulador.interno" },
  "besu": { "rpc": "https://besu.interno:8545", "keyHex": "<KMS em produção>", "deployment": "/etc/janus/1337.json" }
}
```

`dataDir` guarda `janus.db` (SQLite), `secrets/` (S cifrado), `master.key`, `attestor.key` e `besu.checkpoint`. **Faça backup de `dataDir` e `master.key` separadamente.**

## Subir

1. Redes: Fabric 3.1.5 com o Token Chaincode do Panurus (política `AND(Registradora, Banco)`), e Besu (`make besu-up` em dev).
2. Contratos: `forge script contracts/script/Deploy.s.sol --broadcast`, que grava o deployment JSON.
3. Governança B2: cadastre as chaves de atestador da Registradora e do Banco (`GET /api/v1/info` → `attestor`) com `besu.Client.RegisterAttestors(threshold=2)`.
4. Inicie os nós. A Registradora registra o auditor sozinha (com retry). Confira com `GET /api/v1/health`.

## Operar

| Ação | Chamada |
|---|---|
| Emitir recebível | `POST /api/v1/recebiveis` na Registradora |
| Propor swap | `POST /api/v1/swaps {buyerNode, assetId, price, protocol}` na Vendedora |
| Acompanhar | `GET /api/v1/swaps[?active=true]`, `GET /api/v1/swaps/{id}`, SSE `GET /api/v1/events` |
| Forçar reavaliação | `POST /api/v1/swaps/{id}/kick` |

## Incidentes

| Sintoma | Causa provável | Ação |
|---|---|---|
| Nó caiu com swaps ativos | — | **Reinicie.** A reconciliação de boot retoma todo swap não terminal a partir das cadeias e do S persistido |
| `alert: PERDA: S foi revelado mas o claim do NFT não aconteceu antes de T1` | watchtower do Banco ausente além de `T1 − T2` | risco residual (T9). Escalar para tratamento bilateral fora da cadeia |
| `alert: pago, mas o Banco não fez o claim até T1` (na Vendedora) | Banco ausente | o nó honesto **não** faz reclaim. Coordene a transferência manual do NFT ao Banco |
| Swap parado em `both_locked` (B2) | relay ou atestador fora do ar, ou chave de atestador não cadastrada | confira `attestors` no `janus.json`, `GET /api/v1/info` dos atestadores e o cadastro no `FabricProofVerifier`. Depois de T2, as duas pernas revertem sozinhas |
| `audit: early reclaim refused` | reclaim antes de `deadline + 30 s` | esperado: aguarde |
| `lock accept: no signed terms for this hashlock` | lock sem negociação prévia | todo lock passa por `POST /api/v1/swaps` |
| Instalação do chaincode falha com `invalid UTF-8` | build no `ccenv` com Go antigo | use o external builder (`node/ci/external-builders/golang`), ver o ADR-0001 |

## Falhas injetadas (E2E e demonstração)

No `janus.json`: `"faults": {"noLock": true}`, `{"noReveal": true}` ou `{"crashAfterReveal": true}`. A última encerra o processo com código 3 depois da revelação, e o reinício faz a reconciliação.
