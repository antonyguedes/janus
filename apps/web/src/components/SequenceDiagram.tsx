import type { Protocol, SimEvent } from "../types";
import { actionLabel, targetOf } from "../labels";

interface Props {
  events: SimEvent[];
  now: number;
  protocol: Protocol;
}

type Lane = "vendedora" | "fabric" | "relay" | "besu" | "banco";

const laneTitle: Record<Lane, string> = {
  vendedora: "Nó Vendedora",
  fabric: "Fabric · Panurus HTLC",
  relay: "Relay (atestações)",
  besu: "Besu · Drex",
  banco: "Nó Banco",
};

/** Diagrama de sequência: uma seta por transação comitada (ou rejeitada) até o instante atual. */
export function SequenceDiagram({ events, now, protocol }: Props) {
  const lanes: Lane[] = protocol === "cacti"
    ? ["vendedora", "fabric", "relay", "besu", "banco"]
    : ["vendedora", "fabric", "besu", "banco"];
  const width = 760;
  const gap = width / lanes.length;
  const x = (l: Lane) => gap / 2 + lanes.indexOf(l) * gap;
  const rowH = 46;
  const top = 72;
  const visible = events.filter((e) => e.t <= now);
  const height = top + Math.max(events.length, 4) * rowH + 24;

  return (
    <svg className="seq" viewBox={`0 0 ${width} ${height}`} role="img" aria-label="Diagrama de sequência do swap">
      <defs>
        <marker id="arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse">
          <path d="M0,0 L10,5 L0,10 z" className="arrowhead" />
        </marker>
        <marker id="arrow-bad" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse">
          <path d="M0,0 L10,5 L0,10 z" className="arrowhead bad" />
        </marker>
      </defs>
      {lanes.map((l) => (
        <g key={l}>
          <rect x={x(l) - 64} y={8} width={128} height={30} rx={8} className={`lane-head lane-${l}`} />
          <text x={x(l)} y={28} textAnchor="middle" className="lane-title">{laneTitle[l]}</text>
          <line x1={x(l)} x2={x(l)} y1={40} y2={height - 8} className="lifeline" />
        </g>
      ))}
      {visible.map((e, i) => {
        const target = targetOf(e.action);
        if (!target) return null;
        const from = x(e.actor as Lane);
        const to = x(target);
        const y = top + i * rowH;
        const label = `t=${e.t}s · ${actionLabel[e.action]}${e.ok ? "" : " ✕"}`;
        return (
          <g key={`${e.t}-${e.action}-${i}`} className={`msg ${e.ok ? "ok" : "bad"} appear`}>
            <line x1={from} x2={to + (to > from ? -6 : 6)} y1={y} y2={y} markerEnd={e.ok ? "url(#arrow)" : "url(#arrow-bad)"} />
            <text x={(from + to) / 2} y={y - 7} textAnchor="middle" className="msg-label">{label}</text>
            {e.note && (
              <text x={(from + to) / 2} y={y + 15} textAnchor="middle" className="msg-note">{e.note}</text>
            )}
          </g>
        );
      })}
    </svg>
  );
}
