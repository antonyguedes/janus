import type { Protocol, Snapshot } from "../types";
import { statusLabel } from "../labels";

const fabricState: Record<string, { label: string; holder: string }> = {
  none: { label: "Recebível livre", holder: "com a Vendedora" },
  locked: { label: "Travado no portal", holder: "HTLC do Panurus" },
  claimed: { label: "Atravessou", holder: "com o Banco" },
  reclaimed: { label: "Voltou", holder: "com a Vendedora" },
};
const besuState: Record<string, { label: string; holder: string }> = {
  none: { label: "Drex livre", holder: "com o Banco" },
  locked: { label: "Travado no portal", holder: "contrato de liquidação" },
  claimed: { label: "Atravessou", holder: "com a Vendedora" },
  refunded: { label: "Voltou", holder: "com o Banco" },
};

interface Props {
  snap: Snapshot;
  protocol: Protocol;
}

/** Os dois portais de Jano: o Fabric (dourado) e o Besu (ciano), com a chave S entre eles. */
export function PortalBridge({ snap, protocol }: Props) {
  const f = fabricState[snap.fabricLabel] ?? fabricState.none!;
  const b = besuState[snap.besuLabel] ?? besuState.none!;
  const key = snap.revealed ? "revealed" : snap.buyerKnowsS ? "leaked" : "secret";
  const done = snap.status === "completed";
  const bad = snap.status.startsWith("anomaly");
  return (
    <section className={`bridge ${done ? "done" : ""} ${bad ? "bad" : ""}`} aria-label="Estado das duas redes">
      <Portal side="fabric" title="Fabric" subtitle="Recebível · NFT" state={snap.fabricLabel} label={f.label} holder={f.holder} />
      <div className="bridge-core">
        <div className={`key key-${key}`} title="Preimagem S">
          <span className="key-glyph" aria-hidden>⚿</span>
          <span>{key === "revealed" ? "S revelada" : key === "leaked" ? "S exposta" : "S secreta"}</span>
        </div>
        <div className="bridge-beam" aria-hidden />
        <div className="bridge-status">{statusLabel[snap.status]}</div>
        <div className="bridge-mode">{protocol === "htlc" ? "A · HTLC" : "B · HTLC atestado (Cacti)"}</div>
        {protocol === "cacti" && (
          <div className={`attest ${snap.attested ? "on" : ""}`}>{snap.attested ? "atestação 2-de-2 ✓" : "atestação pendente"}</div>
        )}
      </div>
      <Portal side="besu" title="Besu" subtitle="Pagamento · Drex" state={snap.besuLabel} label={b.label} holder={b.holder} />
    </section>
  );
}

function Portal(p: { side: "fabric" | "besu"; title: string; subtitle: string; state: string; label: string; holder: string }) {
  return (
    <div className={`portal portal-${p.side} state-${p.state}`}>
      <div className="portal-ring" aria-hidden>
        <span className="ring r1" />
        <span className="ring r2" />
        <span className="portal-glyph">{p.side === "fabric" ? "◈" : "◎"}</span>
      </div>
      <div className="portal-text">
        <strong>{p.title}</strong>
        <small>{p.subtitle}</small>
        <span className="portal-state">{p.label}</span>
        <span className="portal-holder">{p.holder}</span>
      </div>
    </div>
  );
}
