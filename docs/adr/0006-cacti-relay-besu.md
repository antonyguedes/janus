# ADR-0006 — Conector Besu do Hyperledger Cacti na submissão do settle (B2)

- **Status:** aceito
- **Data:** 2026-11-09
- **Contexto:** na B2, a Vendedora submetia `AttestedSettlement.settle(S, atestação, assinaturas)` pelo próprio cliente Go. O plano previa o Cacti transportando a prova até o Besu (`R->>B: submete a prova`). O ADR-0005 integrou o conector Fabric do Cacti na coleta das atestações, e este ADR completa o relay.

## Spike (evidência)

| Pergunta | Resposta |
|---|---|
| Conector | `@hyperledger-cacti/cactus-plugin-ledger-connector-besu` 3.0.1 (`web3` 1.6.1), no mesmo servidor Cacti (`infra/cacti/relay`) |
| `invoke-contract` sem pré-registro | aceita `contractAbi` + `contractAddress` na requisição. A ABI vem dos bindings gerados do Foundry |
| Credencial | `CACTUS_KEYCHAIN_REF`: a chave da conta de relay fica no keychain do servidor, e a requisição só leva a referência |
| Gas zero | com `gasPrice: "0"` e `gas` explícito, o `web3` 1.6 assina uma transação legacy (`type 0x0`, EIP-155) aceita pelo anvil Osaka **e pelo Besu 26.11 QBFT** |
| Tuplas | `FabricLockAttestation` e `OrgSignature[]` como objetos com os nomes da ABI, `bytes32` em hex e `uint64` em string decimal |
| Formato da resposta | o servidor v3.0.1 responde `{"success", "out": {"transactionReceipt": {...}}}`, enquanto o OpenAPI descreve o recibo no topo. O cliente Go aceita os dois |
| Limites do Besu (achados) | `eth_getLogs` com `toBlock` além do topo → `Invalid params`. Intervalo > 5000 blocos → `exceeds maximum RPC range limit`. A busca do settle passou a usar janelas de 2000 blocos, do topo para trás |

## Decisão

1. A Vendedora opera um **relay**: um servidor Cacti com o conector Besu e a **conta de relay** (`0x15d3…`, já prevista na allowlist de contas do Besu).
2. Com `relay` no `janus.json`, o `settle` da B2 sai pelo `invoke-contract` (`node/internal/cacti.BesuConnector`). O nó **não confia só na resposta do relay**: espera a inclusão pelo próprio cliente Besu e exige `status == 1`.
3. **Sem fallback por padrão.** Se o relay estiver fora, o motor tenta de novo no próximo ciclo. Se nada for liquidado até T2, as duas pernas revertem (seguro). `fallbackDirect: true` habilita a submissão direta como alternativa de vivacidade.

## Confiança no relay

O relay é **transporte, não custódia**:
- `settle` é permissionless e paga sempre o `seller` fixado no depósito. O relay não consegue desviar o pagamento.
- A atestação e as assinaturas são verificadas on-chain (`FabricProofVerifier`), e o relay não consegue forjá-las.
- O relay vê S antes da inclusão, como qualquer nó do mempool. S está destinado a ser público, e o Banco é quem se beneficia de conhecê-lo.
- Se a chave do relay vazar, o atacante só consegue submeter `settle` válidos (de graça, numa rede com gas zero) ou nada. O risco é de vivacidade (recusar ou atrasar), mitigado pelo fallback opcional e pelos timelocks.

## Consequências

- Verificação: `TestCactiBesuSettle` (anvil e Besu QBFT real via `JANUS_BESU_RPC`/`JANUS_BESU_WS`); `TestE2EAttested`, com o settle aceito on-chain **enviado pela conta de relay**.
- Operação: `make demo` sobe o relay em `:4102`. O `make demo-besu` usa o WebSocket do Besu em `:8546`.
- Keychain em memória e ApiServer sem autorização continuam aceitáveis só em DEV (ADR-0005).
