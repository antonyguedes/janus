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

## Spike de reprodução (Fase 2): resultado

**Ambiente:** topologia NWO (`node/topology/it`), Fabric 3.1.5, Panurus v0.27.0 (FabToken), 3 orgs, política `AND('RegistradoraMSP.member','BancoMSP.member')`.

| Tentativa de reclaim (Vendedora) | Validador (relógio do endossante) | Auditor (relógio da Registradora) | Resultado |
|---|---|---|---|
| antes de T1 | recusa (`now < deadline`) | recusa | ❌ rejeitado |
| **T1 + 2 s** | **aceitaria** (`now ≥ deadline`) | **recusa**: `audit: early reclaim refused (now … < deadline+grace …)` | ❌ rejeitado (camada 2) |
| T1 + 30 s + 2 s | aceita | aceita | ✅ reclaim comitado, status `reclaimed` |

A segunda linha é a evidência de que **a camada 2 bloqueia uma transação que os endossantes aprovariam**: o auditor decide com o próprio relógio, independentemente do relógio dos peers. Com a política `AND`, o relógio do peer da Vendedora nem participa. Um reclaim antecipado exigiria comprometer um peer da Registradora **e** um do Banco (camada 1) **e** a auditora (camada 2).

**Limite do spike:** um peer com relógio adiantado não pode ser simulado. O validador roda em Go, que lê o relógio via vDSO (o `libfaketime` não intercepta), e os *time namespaces* do Linux não deslocam o `CLOCK_REALTIME`. A reprodução literal do "peer forjado" exigiria um binário de chaincode adulterado em uma única org. Fica registrada como teste manual opcional.

### Roteiro original

Na topologia NWO, rodar um peer da Vendedora com `faketime +5h` e tentar o reclaim antecipado:
- com a política `OR(...)`, o ataque **passa** (reproduz);
- com a política `AND(Registradora, Banco)`, ele **falha** no endosso (mitigado).

Se a camada 2 (auditor) não for viável na API de auditoria, a alternativa é **contribuir upstream** uma checagem de deadline máximo, nunca fazer um fork local.
