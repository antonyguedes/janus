# Modelo de ameaças — Janus

**Escopo:** liquidação DvP entre o recebível (NFT do Panurus no Fabric) e o Drex (ERC-20 no Besu), pelas abordagens A (HTLC) e B2 (HTLC atestado).
**Partes:** a Registradora (emissora e auditora do NFT, atestadora em B2), a Vendedora (dona do NFT) e o Banco (comprador, watchtower e atestador em B2).
**Premissas de confiança:** os consensos das duas redes são seguros (Raft/BFT no Fabric, QBFT com ≥ 3f+1 no Besu). Cada instituição protege as próprias chaves. Os relógios dos nós honestos ficam dentro de `ε_skew` (NTP).

## Propriedade de segurança

> **Atomicidade para a parte honesta:** uma parte que segue o protocolo e fica disponível dentro da janela exigida termina com **o próprio ativo de volta** ou **com o ativo da contraparte**, nunca sem nenhum dos dois.

Verificação: o simulador (`node/internal/swap/sim`) cobre os 9 cenários nas duas abordagens e uma varredura de quedas e latências, sempre sem perda para os honestos. O motor real é verificado sobre cadeias falsas (`internal/engine`), e os E2E rodam em redes reais (`node/topology/it`).

## Ameaças e mitigações

| # | Ameaça | Abordagem | Mitigação | Evidência |
|---|---|---|---|---|
| T1 | A Vendedora recebe o Drex e faz o reclaim do NFT com um **endossante de relógio adiantado** | A, B2 | (1) política de endosso `AND(Registradora, Banco)`; (2) a auditora confere o tempo com o próprio relógio (reclaim só em `deadline + 30 s`); (3) margem Δ | ADR-0004; integração: reclaim recusado em T1+2 s pela auditoria; `TestCheckSpendTimeBoundaries` |
| T2 | A Vendedora **revela S tarde**: a tx falha depois de T2 e S vaza no mempool | A | o nó honesto não revela a menos de `RevealMargin` de T2 | `NextAction`; cenário `late-reveal-leak` (mostra a perda da Vendedora imprudente) e `TestNextActionSellerNeverRevealsInsideMargin` |
| T3 | A Vendedora segura S até **T2 − ε** para encurtar a janela do Banco | A | `T1 − T2 ≥ Δ` validado nos termos; o watchtower age por evento | cenário `reveal-t2-minus`, `TermsValidate` |
| T4 | NFT travado com **deadline, recipient, hash ou ativo divergentes** | A, B2 | o Banco só aceita o lock esperado (`LockAcceptView`) e só paga depois de `CheckFabricLock`; em B2, o contrato exige a igualdade dos campos atestados | integração "lock sem termos"; cenário `wrong-deadline`; `test_revert_attestationMismatches` |
| T5 | **Pagamento travado errado** (valor, recipient, timelock > T2, contrato) | A, B2 | a Vendedora só revela depois de `CheckBesuLock` | `TestCheckBesuLock`; teste anvil `CheckBesuLock` |
| T6 | **Reuso de hashlock ou preimagem** em outro swap | A, B2 | `hashlockUsed` nos contratos; `swapId = sha256(H‖partes‖nonce)` | `test_revert_hashlockReuse*` |
| T7 | **Atestação forjada, adulterada, insuficiente, maleável ou de outro deployment** | B2 | P-256 com low-s, k-de-n orgs distintas (ordem estrita), domínio com chainId + verificador, chaves revogáveis, Vendedora fora do conjunto de atestadores | `AttestedSettlement.t.sol` (11 casos); `TestAttestedSettlementWithGoAttestations` |
| T8 | **Atestador assina um lock não comitado** | B2 | evidência da auditoria ou do vault **e** `Ledger.Status(txID) == Valid` antes de assinar | `buildAttestation` |
| T9 | **Watchtower do Banco fora do ar** por mais de `T1 − T2` depois da revelação | A, B2 | claim permissionless no Besu; reconciliação no boot; S persistido assim que visto; alertas; watchtowers redundantes (operacional) | cenários `buyer-crash-recover` (ok) e `buyer-offline-past-t1` (**risco residual documentado**) |
| T10 | **Relay ou atestadores parados** | B2 | é questão de vivacidade, não de segurança: sem atestação não há pagamento, e as duas pernas revertem | cenário `seller-no-reveal` em B (relay parado) |
| T11 | **Recebível com gravame** ou travado duas vezes | A, B2 | registro de gravames na auditoria; UTXO do Panurus (um token só se gasta uma vez) | integração "gravame bloqueia o lock" |
| T12 | **Emissão duplicada** do mesmo recebível | — | ID único derivado do número de registro (`uniqueness`) | `IssueView` |
| T13 | **Termos adulterados** em trânsito ou repudiados | A, B2 | sessão FSC autenticada (X.509); as duas partes assinam o hash canônico com a identidade de owner | `TestSignedTerms`, E2E (negociação real) |
| T14 | **Vazamento de S em repouso** | A, B2 | AES-256-GCM com swapId como AAD; chave mestra 0600 (KMS/HSM fora do escopo do MVP) | `secrets_test.go` |
| T15 | **Conta não autorizada** movimenta Drex | — | allowlist no token **e** permissionamento de contas do Besu | `DrexToken.t.sol`, `permissions_config.toml` |
| T16 | **Reentrância ou token anômalo** nos contratos | A, B2 | `ReentrancyGuard`, `SafeERC20`, checagem do delta de saldo, CEI | invariantes `HTLCPayment.invariant.t.sol` |

## Riscos residuais

1. **Disponibilidade do Banco (T9):** é inerente a HTLC. A janela é `T1 − T2` (2 h em produção).
2. **Panurus não auditado:** versão fixada (ADR-0001), mudanças monitoradas e dependência isolada atrás de portas.
3. **Privacidade no Besu:** valor e partes do lock são visíveis no Besu. O ZKAT-DLog cobre só o lado Fabric (fase opcional).
4. **Chaves em arquivo (DEV):** em produção, as chaves Besu, de atestador e mestra vão para KMS/HSM, e as chaves de atestador são emitidas pela CA da org.
5. **Estado em memória:** registro de gravames, expectativas de lock e locks auditados ficam em memória no MVP. A expectativa é reinstalada na reconciliação de boot. Gravames e locks auditados precisam de persistência antes de produção.
