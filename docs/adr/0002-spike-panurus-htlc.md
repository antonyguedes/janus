# ADR-0002 — Spike Panurus: HTLC, hash e fonte de tempo

- **Status:** aceito (Fase 0)
- **Data:** 2026-11-08
- **Método:** leitura do código-fonte do Panurus **v0.27.0** (commit `632947b`). As referências abaixo são caminhos dentro do repositório.

## Evidências

| Pergunta | Resposta | Onde |
|---|---|---|
| Função de hash do HTLC | `crypto.SHA256` por padrão, configurável pelo atributo `htlc.hashFunc` | `token/services/interop/htlc/transaction.go` (`Lock`) |
| Encoding da imagem | `encoding.None` (zero value) por padrão, ou seja, os **bytes crus** do digest. Também há `Base64` e `Hex` | `token/services/interop/encoding/encoding.go` |
| Imagem | `HashEncoding(HashFunc(preimage))` | `HashInfo.Image` em `script.go` |
| Preimagem própria | `htlc.WithHash(H)` deixa o chamador fornecer `H`. Sem ele, o Panurus sorteia um nonce de **24 bytes** | `recipientAsScript`, `getRandomNonce` |
| Deadline | **relativo** no `Lock` (`time.Duration`) e convertido para `time.Now().Add(d)` **no relógio de quem trava** | `recipientAsScript` |
| Fonte de tempo do validador | **`time.Now()` do peer que executa o Token Chaincode**, ou seja, o endossante. **Não** é o timestamp declarado pelo cliente | `token/core/fabtoken/v1/validator/validator_transfer.go` (`TransferHTLCValidate`) |
| Fronteira | claim se `now.Before(deadline)`, reclaim caso contrário (`now >= deadline`) | `token/services/identity/interop/htlc/validator.go` (`VerifyOwner`) |
| Lock com deadline no passado | rejeitado (`Script.Validate(now)`) | `script.go` |
| Claim | exige a preimagem e a assinatura do *recipient* (`ClaimSignature{RecipientSignature, Preimage}`). **Só o recipient (Banco) consegue fazer o claim** | `validator.go` (`MetadataClaimKeyCheck`) |
| Forma do gasto | 1 entrada → 1 saída (transferência integral de propriedade), compatível com um NFT | `TransferHTLCValidate` |
| NFT | serviço `token/services/nfttx` (UUID único por token, quantidade 1) | `token/services/nfttx` |

## Decisões

1. **Hash interoperável:** `S` com 32 bytes, `H = SHA256(S)` em bytes crus, passado ao Panurus com `htlc.WithHash(H)`, `WithHashFunc(crypto.SHA256)` e `WithHashEncoding(encoding.None)`, **sempre explícitos**, sem depender de defaults. No Solidity: `sha256(abi.encodePacked(bytes32 S))`. Os vetores ficam compartilhados em `testdata/vectors/htlc.json`.
2. **T1 absoluto → deadline relativo:** o nó da Vendedora calcula `d = T1 − agora` no instante do lock. O Banco **não confia** nesse cálculo e confere no vault que `|script.Deadline − T1| ≤ tolerância de skew`, e que `script.HashInfo` tem exatamente (SHA256, None, H). Ver o ADR-0004.
3. **Claim no Fabric é do Banco:** como a assinatura do recipient é exigida, nas duas abordagens quem faz o `htlc.Claim` no Fabric é o nó do Banco (watchtower). Um "claim permissionless" no Fabric não é possível sem mudar o Panurus, e isso entra no modelo de ameaças como requisito de disponibilidade do Banco.
4. **O modelo de ameaça de tempo muda** em relação à seção 6 do plano. O risco não é o timestamp do cliente, e sim um **endossante com relógio manipulado**. Detalhes e mitigação no ADR-0004.

## Adendo (Fase 4): dono dos outputs na auditoria

Nos outputs de **transferência** devolvidos por `ttx.Auditor.Audit` (Panurus v0.27.0, `token/request.go`), o campo `token.Output.Owner` é o **recipient dos metadados** (no HTLC, o Banco), e o dono real no ledger (a identidade tipada do script HTLC) fica em `token.Output.Token.Owner`, embutido. O `htlc.ToOutput` lê o campo externo. Por isso, uma auditoria que o usa para detectar locks **nunca enxerga o script**. A view de auditoria de exemplo do upstream (`integration/token/interop/views/auditor.go`) tem essa lacuna latente: a validação do script dos novos locks fica sem efeito.

**Decisão:** a auditoria do Janus usa o dono no ledger (`ledgerOwner`) para detectar scripts HTLC, validar novos locks (`CheckLockScript`), classificar claim/reclaim e registrar a evidência de atestação (B2). Há teste de regressão em `TestAuditSeesHTLCLockByLedgerOwner`. Candidato a reporte upstream.

O mesmo vale para `Type`: `token.Output` embute `token.Token` e declara um `Type` próprio. A auditoria lê o primeiro não vazio.
