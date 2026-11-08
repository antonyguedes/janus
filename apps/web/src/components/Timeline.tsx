interface Props {
  now: number;
  end: number;
  t1: number;
  t2: number;
  delta: number;
  marks: number[];
  onSeek: (t: number) => void;
}

/** Relógio do swap: barra 0..end com T2, T1, a margem Δ e os instantes com eventos. */
export function Timeline({ now, end, t1, t2, delta, marks, onSeek }: Props) {
  const pct = (t: number) => `${Math.min(100, (t / end) * 100)}%`;
  const fmt = (t: number) => `${Math.floor(t / 60)}:${String(t % 60).padStart(2, "0")}`;
  return (
    <div className="timeline">
      <div className="clock" aria-live="polite">
        <span className="clock-label">t =</span> {fmt(now)}
        <span className="clock-sub">
          {now < t2 ? `T2 em ${t2 - now}s` : now < t1 ? `T1 em ${t1 - now}s` : "após T1"}
        </span>
      </div>
      <div className="track" onClick={(e) => {
        const r = (e.currentTarget as HTMLDivElement).getBoundingClientRect();
        onSeek(Math.round(((e.clientX - r.left) / r.width) * end));
      }}>
        <div className="zone zone-delta" style={{ left: pct(t2), width: `calc(${pct(t1)} - ${pct(t2)})` }}
          title={`T1 − T2 = ${t1 - t2}s (Δ mínimo ${delta}s)`} />
        <div className="progress" style={{ width: pct(now) }} />
        {marks.map((m, i) => <div key={i} className="mark" style={{ left: pct(m) }} />)}
        <div className="limit t2" style={{ left: pct(t2) }}><span>T2<em className="long"> · Besu</em></span></div>
        <div className="limit t1" style={{ left: pct(t1) }}><span>T1<em className="long"> · Fabric</em></span></div>
        <div className="cursor" style={{ left: pct(now) }} />
      </div>
    </div>
  );
}
