import type { Outcome, Protocol, Snapshot } from "../types";
import { actionLabel, statusLabel, statusTone } from "../labels";

const fabricText: Record<string, string> = {
  none: "Livre na wallet da Vendedora",
  locked: "Travado no script HTLC",
  claimed: "Reivindicado pelo Banco",
  reclaimed: "Reavido pela Vendedora",
};
const besuText: Record<string, string> = {
  none: "Sem pagamento travado",
  locked: "Drex travado no contrato",
  claimed: "Pago à Vendedora (S público)",
  refunded: "Devolvido ao Banco",
};

interface Props {
  snap: Snapshot;
  protocol: Protocol;
  t1: number;
  t2: number;
  outcome?: Outcome;
  finished: boolean;
}

/** Inspetor de estado por rede e por nó. */
export function Inspector({ snap, protocol, t1, t2, outcome, finished }: Props) {
  const tone = statusTone(snap.status);
  return (
    <div className="inspector">
      <div className={`status-pill tone-${tone}`}>{statusLabel[snap.status]}</div>
      <div className="cards">
        <section className="card">
          <h3>Fabric · NFT do recebível</h3>
          <p className={`state state-${snap.fabricLabel}`}>{fabricText[snap.fabricLabel]}</p>
          <dl>
            <dt>Deadline do lock</dt>
            <dd>{snap.fabricDeadline ? `t=${snap.fabricDeadline}s` : "—"} <small>(T1 = {t1}s)</small></dd>
            <dt>Hash</dt><dd>SHA-256 · encoding None</dd>
          </dl>
        </section>
        <section className="card">
          <h3>Besu · {protocol === "cacti" ? "AttestedSettlement" : "HTLCPayment"}</h3>
          <p className={`state state-${snap.besuLabel}`}>{besuText[snap.besuLabel]}</p>
          <dl>
            <dt>Timelock</dt><dd>T2 = {t2}s</dd>
            <dt>Preimagem S</dt><dd>{snap.revealed ? "revelada on-chain" : snap.buyerKnowsS ? "exposta no mempool (tx revertida)" : "secreta"}</dd>
            {protocol === "cacti" && (<><dt>Atestação 2-de-2</dt><dd>{snap.attested ? "coletada (Registradora + Banco)" : "pendente"}</dd></>)}
          </dl>
        </section>
        <section className={`card node ${snap.sellerOnline ? "" : "offline"}`}>
          <h3>Nó Vendedora {snap.sellerOnline ? "" : "· offline"}</h3>
          <p className="decision"><b>{actionLabel[snap.seller.kind] ?? snap.seller.kind}</b></p>
          <p className="reason">{snap.seller.reason}</p>
        </section>
        <section className={`card node ${snap.buyerOnline ? "" : "offline"}`}>
          <h3>Nó Banco (watchtower) {snap.buyerOnline ? "" : "· offline"}</h3>
          <p className="decision"><b>{actionLabel[snap.buyer.kind] ?? snap.buyer.kind}</b></p>
          <p className="reason">{snap.buyer.reason}</p>
          <p className="reason">Conhece S: {snap.buyerKnowsS ? "sim" : "não"}</p>
        </section>
      </div>
      {finished && outcome && (
        <section className={`card outcome tone-${statusTone(outcome.status)}`}>
          <h3>Resultado</h3>
          <ul>
            <li>Vendedora: {outcome.sellerHasCash ? "recebeu o Drex" : outcome.sellerHasNft ? "manteve o NFT" : "perdeu as duas pernas"}</li>
            <li>Banco: {outcome.buyerHasNft ? "recebeu o NFT" : outcome.buyerHasCash ? "manteve o Drex" : "perdeu as duas pernas"}</li>
          </ul>
        </section>
      )}
    </div>
  );
}
