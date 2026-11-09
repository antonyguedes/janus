# ADR-0003 — Spike Cacti e decisão B1/B2

- **Status:** aceito (Fase 0). **Decisão: B2, "HTLC atestado".**
- **Data:** 2026-11-08
- **Método:** inventário do Hyperledger Cacti **v3.0.1** (árvore do repositório no tag) e testes do precompile P-256 no Besu 26.11.0.

## Evidências

| Componente Cacti v3.0.1 | O que cobre | Serve ao Janus? |
|---|---|---|
| `weaver/core/drivers` | drivers de relay **só para `fabric-driver` e `corda-driver`** | Fabric → *: **sim** (coleta a view com os endossos dos peers). Besu → Fabric: **não existe driver Besu** |
| `weaver/core/network/besu`, `weaver/sdks/besu` | contratos e SDK de **asset exchange (HTLC)** no Besu | não serve para provas. O HTLC do Besu é nosso (`HTLCPayment.sol`) |
| `cactus-plugin-ledger-connector-fabric` / `-besu` | conectores (consulta, submissão de tx, eventos) | **sim**: transporte e consulta |
| `cactus-plugin-satp-hermes` | gateway SATP de **transferência** de ativos | **não**: move o ativo, o que contraria a premissa jurídica |
| `cactus-plugin-htlc-eth-besu(-erc20)` | HTLC no Besu | redundante com o `HTLCPayment.sol` |
| Besu `P256VERIFY` (EIP-7951, `0x100`) | ECDSA secp256r1 nativo | **sim**: com `osakaTime: 0` no genesis, assinatura válida → `1`, inválida → vazio, cerca de 6,9k gas |

## Decisão

- **B1 rejeitada.** Liberar o NFT no Fabric com uma *prova de evento do Besu* exige: (a) um componente que produza provas do Besu, que o Cacti não tem (não há driver Besu no Weaver); e (b) uma condição de gasto nova no validador do Panurus, ou seja, um fork ou contribuição no núcleo do token. Isso é invasivo demais para o MVP.
- **B2 adotada, "HTLC atestado":**
  - Fabric: o NFT continua travado pelo **HTLC nativo do Panurus** (hashlock `H` e deadline `T1`).
  - Besu: o `AttestedSettlement` só libera o Drex para a Vendedora se receber **(i)** a preimagem `S` com `sha256(S) = H` **e (ii)** uma **atestação do lock no Fabric** válida para a política (k de n orgs), verificada pelo `FabricProofVerifier` com P-256 nativo.
  - A revelação de `S` no Besu libera o NFT: o Banco (watchtower) faz o `htlc.Claim(S)` no Fabric. **Não é preciso prova Besu → Fabric** e não existe cadastro de validadores QBFT no lado Fabric. Esse item da Fase 4 cai.
  - **Ganho em relação à A:** o Banco pode depositar **antes** do lock no Fabric (o contrato garante que o pagamento só sai com o lock provado) e deixa de depender da leitura do próprio vault como única garantia.

## Formato da atestação (MVP)

Os peers do Fabric não assinam mensagens arbitrárias, e a view do Weaver é protobuf (caro de decodificar no EVM). Por isso, no MVP, **cada org atestadora roda um atestador** (dentro do nó Janus dela) que:

1. lê o lock **já comitado** no ledger, pelo próprio peer/vault ou pelo conector Fabric do Cacti;
2. assina com uma chave **P-256 emitida pela CA da própria org** (X.509) o digest canônico `FabricLockAttestation` (ver `contracts/src/FabricProofVerifier.sol` e `node/internal/attest`);
3. o relay (componente Cacti ou nó Janus) agrega as assinaturas e submete ao Besu.

No `FabricProofVerifier`, a cadeia X.509 é validada **no cadastro** (feito por governança, fora da cadeia), e no contrato fica a chave pública do atestador vinculada ao `mspId`. A política padrão é **2 de 2: Registradora (auditora) + Banco**. A Vendedora não atesta, porque é a parte que se beneficiaria de uma atestação falsa.

**Evolução prevista (pós-MVP):** substituir o atestador pela **view do Weaver com os endossos reais dos peers**, verificando no contrato a assinatura do endossante sobre o `ProposalResponsePayload`. A interface `verify(attestation, sigs)` e o digest P-256 foram desenhados para comportar essa troca.

## Consequências

- O SATP-Hermes fica documentado como **não usado** (ele transfere o ativo).
- A disponibilidade do relay e dos atestadores passa a ser requisito de vivacidade, não de segurança: sem atestação não há pagamento, e o Drex volta depois de `T2`.

## Atualização (2026-11-09)

Os conectores do Cacti passaram a ser usados na coleta das atestações: cada atestador prova o lock pelo conector Fabric do Cacti da própria org. Ver o **ADR-0005**.
