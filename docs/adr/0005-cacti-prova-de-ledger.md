# ADR-0005 — Conectores do Hyperledger Cacti na coleta de atestações (B2)

- **Status:** aceito
- **Data:** 2026-11-09
- **Contexto:** no MVP da B2 (ADR-0003), cada atestador assinava o lock com base na evidência local: o registro da auditoria (Registradora, em memória) ou o vault (Banco), mais `Ledger.Status(txID)` do próprio nó FSC. O número do bloco ia zerado. O plano previa o Cacti coletando a prova do lock no Fabric.

## Spike (evidência)

| Pergunta | Resposta |
|---|---|
| Versão | O release **v3.0.1** do Cacti publicou os pacotes JS sob o escopo **`@hyperledger-cacti/*` 3.0.1**. Sob `@hyperledger/*`, a última versão é a 2.0.0, de 2024 |
| Conector Fabric × Fabric 3.1.5 | **funciona.** O `cactus-plugin-ledger-connector-fabric` 3.0.1 (SDK `fabric-network` 2.5) leu blocos e transações da rede NWO, com TLS e discovery `asLocalhost` |
| `get-transaction-receipt-by-txid` | devolve `blockNumber` (como `Long` `{low, high}`), `chainCodeName`, `responseStatus` e os **endossantes** (`mspid`, certificado e assinatura). Para o lock real: bloco 8, `tns`, 200, `[RegistradoraMSP, BancoMSP]` |
| validade da transação | o recibo **não** traz o código de validação (ele usa `GetBlockByTxID`), e uma transação pode estar no bloco e ter sido invalidada. É preciso ler o bloco (`get-block`, `responseType: full`) e conferir o `TRANSACTIONS_FILTER` (`metadata[2]`) |

## Decisão

1. **Cada org atestadora (Registradora, Banco) roda um servidor Cacti** (`infra/cacti/relay`: `cactus-cmd-api-server` v3.0.1 com o conector Fabric e um keychain) usando a **identidade da própria org**.
2. Ao receber um pedido de atestação, o atestador do nó Janus obtém pelo Cacti a **prova da transação de lock** (`node/internal/cacti`):
   - recibo → bloco e endossantes;
   - bloco → `TRANSACTIONS_FILTER == VALID`;
   - exige o chaincode do token (`tns`) e **endossos de todas as orgs da política** (`AND(Registradora, Banco)`).
3. **Falha fechada:** com o Cacti configurado no `janus.json` (`cacti`), a ausência ou a reprovação da prova impede a assinatura.
4. O **número do bloco** da prova entra na `FabricLockAttestation`. Como o relay (`CollectView`) exige que todas as orgs atestem fatos idênticos, cada org confirma de forma independente, pelo próprio conector, o mesmo bloco, a mesma transação e o mesmo resultado de validação.
5. Os fatos do script HTLC (hash, recipient, deadline, ativo) continuam vindo da evidência local (auditoria/vault), porque o recibo do Cacti só expõe a primeira escrita do rwset. O Cacti fornece a prova de **comitado, válido e endossado pela política, no bloco N**.

## Fora deste passo

- **Submissão do `settle` pelo conector Besu do Cacti:** feita depois, no **ADR-0006**.
- **Keychain em memória:** só para DEV. Em produção, keychain com Vault/KMS e autorização no ApiServer (`authorizationProtocol`), hoje `NONE` e escutando apenas em `127.0.0.1`.

## Consequências

- Verificação: `node/internal/cacti` (testes com respostas no formato real, incluindo transação invalidada no bloco, endossante faltante e `Long`) e `TestE2EAttested`. No E2E, o settle aceito no Besu traz o bloco que o conector do Cacti reporta para a transação de lock, com os dois endossantes.
- Operação: `make demo` sobe os servidores Cacti em `:4100` (Registradora) e `:4101` (Banco) quando as dependências npm estão instaladas.
