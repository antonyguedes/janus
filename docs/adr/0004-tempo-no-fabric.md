# ADR-0004 — Tempo no Fabric e no Besu

- **Status:** aceito (Fase 0). O spike de reprodução do ataque roda na Fase 2 (topologia NWO). Este ADR registra o modelo e a mitigação.
- **Data:** 2026-11-08

## Correção do modelo de ameaça

O plano (seção 6) supunha que o validador HTLC usava o **timestamp do cliente**. O spike (ADR-0002) mostrou que, no Panurus v0.27.0, o validador usa **`time.Now()` do peer endossante** (`TransferHTLCValidate`). As consequências são estas:

| Vetor | Efeito | Mitigação |
|---|---|---|
| **Endossante com relógio adiantado, controlado pela Vendedora** | aceita um *reclaim* antes de `T1`, depois que a Vendedora já recebeu o Drex | **Camada 1:** a política de endosso do Token Chaincode exige **AND(Registradora, Banco)** em transações de transferência. Um peer só da Vendedora não basta |
| Relógio do Banco atrasado | o Banco acredita que ainda pode fazer o claim, mas os endossantes já veem `now ≥ T1` | margem `Δ` com `ε_skew` e monitoramento NTP. O watchtower age em `T2 + t_detecção`, muito antes de `T1` |
| Vendedora trava com o deadline errado (`Deadline` vem do relógio dela) | o Banco recebe um prazo menor que o combinado | **Camada 2:** o Banco confere `|script.Deadline − T1| ≤ tol` antes de travar o Drex (A) ou de assinar a atestação (B). A Registradora (auditora) recusa no `AuditView` um lock cujo deadline divirja dos termos registrados além de `tol` |
| Claim endossado antes de `T1` e comitado depois | aceito pelo Fabric. A validação de tempo só acontece no endosso | favorável ao Banco, sem prejuízo à Vendedora, que já tem `S` revelado e foi paga. Um reclaim concorrente gera conflito MVCC: só um dos dois comita |
| Timestamp do bloco Besu | o QBFT limita o desvio do proposer em relação ao bloco anterior. O desvio vem da configuração (`blockperiodseconds`) | `ε_skew` inclui o período do bloco |

## Regra

$$
T_1 - T_2 \ge \Delta = t_{det} + t_{fin,Fabric} + t_{relay} + \varepsilon_{skew} + m_{op}
$$

- `ε_skew` = tolerância de relógio (padrão de 30 s) + o período de bloco do Besu.
- Valores padrão em `node/internal/swap/timepolicy.go`: em produção, `T1 = 4h` e `T2 = 2h`. Em dev, `T1 = 6 min`, `T2 = 3 min` e `Δ_min = 90 s`.
- A fronteira é a mesma nas duas redes: **claim se `agora < T`, reversão se `agora ≥ T`**. O `HTLCPayment.sol` segue a mesma regra do `VerifyOwner` do Panurus.

## Spike de reprodução (Fase 2)

Na topologia NWO, rodar um peer da Vendedora com `faketime +5h` e tentar o reclaim antecipado:
- com a política `OR(...)`, o ataque **passa** (reproduz);
- com a política `AND(Registradora, Banco)`, ele **falha** no endosso (mitigado).

Se a camada 2 (auditor) não for viável na API de auditoria, a alternativa é **contribuir upstream** uma checagem de deadline máximo, nunca fazer um fork local.
