import { useEffect, useMemo, useState } from "react";
import type { Protocol, SimEvent, StepDTO, SwapDTO } from "../types";
import { brl, statusLabel, statusTone } from "../labels";
import { SequenceDiagram } from "./SequenceDiagram";

interface Nodes {
  registradora: string;
  vendedora: string;
  banco: string;
}

const defaults: Nodes = {
  registradora: "http://127.0.0.1:8080",
  vendedora: "http://127.0.0.1:8081",
  banco: "http://127.0.0.1:8082",
};

async function postJSON<T>(url: string, body: unknown): Promise<T> {
  const r = await fetch(url, { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(body) });
  const j = await r.json();
  if (!r.ok) throw new Error(j.error ?? r.statusText);
  return j as T;
}

/** Modo live: dirige os nós Janus pelas APIs e acompanha os passos dos motores via SSE. */
export function LivePanel({ protocol }: { protocol: Protocol }) {
  const [nodes, setNodes] = useState<Nodes>(() => {
    try {
      return { ...defaults, ...JSON.parse(localStorage.getItem("janus.nodes") ?? "{}") };
    } catch {
      return defaults;
    }
  });
  const [assetId, setAssetId] = useState("");
  const [price, setPrice] = useState(1_200_000_00);
  const [swaps, setSwaps] = useState<SwapDTO[]>([]);
  const [selected, setSelected] = useState<string>("");
  const [steps, setSteps] = useState<StepDTO[]>([]);
  const [msg, setMsg] = useState<string>("");

  useEffect(() => {
    try { localStorage.setItem("janus.nodes", JSON.stringify(nodes)); } catch { /* sem storage */ }
  }, [nodes]);

  // SSE dos dois nós participantes
  useEffect(() => {
    const sources = [nodes.vendedora, nodes.banco].map((base) => {
      const es = new EventSource(`${base}/api/v1/events`);
      es.addEventListener("step", (ev) => {
        const s = JSON.parse((ev as MessageEvent).data) as StepDTO;
        setSteps((prev) => [...prev.slice(-500), s]);
      });
      return es;
    });
    return () => sources.forEach((s) => s.close());
  }, [nodes.vendedora, nodes.banco]);

  // lista de swaps (polling leve)
  useEffect(() => {
    let stop = false;
    const tick = async () => {
      try {
        const r = await fetch(`${nodes.vendedora}/api/v1/swaps`);
        if (r.ok && !stop) setSwaps(await r.json());
      } catch { /* nó fora do ar */ }
    };
    void tick();
    const id = setInterval(tick, 2000);
    return () => { stop = true; clearInterval(id); };
  }, [nodes.vendedora]);

  const issue = async () => {
    try {
      const r = await postJSON<{ linearId: string }>(`${nodes.registradora}/api/v1/recebiveis`, {
        owner: "vendedora", registro: `REG-${Date.now()}`, cedente: "Vendedora S.A.", sacado: "Sacado Ltda",
        valorFace: 1_250_000_00, vencimento: "2027-03-31",
      });
      setAssetId(r.linearId);
      setMsg(`Recebível emitido: ${r.linearId}`);
    } catch (e) { setMsg(`Erro: ${(e as Error).message}`); }
  };

  const propose = async () => {
    try {
      const s = await postJSON<SwapDTO>(`${nodes.vendedora}/api/v1/swaps`, { buyerNode: "banco", assetId, price, protocol });
      setSelected(s.swapId);
      setMsg(`Swap ${s.swapId.slice(0, 12)}… proposto e assinado pelas duas partes`);
    } catch (e) { setMsg(`Erro: ${(e as Error).message}`); }
  };

  const events: SimEvent[] = useMemo(() => {
    const mine = steps.filter((s) => s.swapId === selected && s.executed);
    const t0 = mine[0] ? Date.parse(mine[0].at) : 0;
    return mine.map((s) => ({
      t: Math.round((Date.parse(s.at) - t0) / 1000), actor: s.role, action: s.decision.kind,
      chain: "fabric", ok: !s.error, note: s.error,
    }));
  }, [steps, selected]);

  const current = swaps.find((s) => s.swapId === selected);

  return (
    <div className="live">
      <section className="card">
        <h3>Nós Janus</h3>
        <div className="form-grid">
          {(Object.keys(nodes) as (keyof Nodes)[]).map((k) => (
            <label key={k}>{k}
              <input value={nodes[k]} onChange={(e) => setNodes({ ...nodes, [k]: e.target.value })} />
            </label>
          ))}
        </div>
      </section>
      <section className="card">
        <h3>Novo swap</h3>
        <div className="form-row">
          <button onClick={issue}>1 · Emitir recebível (Registradora)</button>
          <label>NFT <input value={assetId} onChange={(e) => setAssetId(e.target.value)} placeholder="LinearID" /></label>
          <label>Preço (centavos) <input type="number" value={price} onChange={(e) => setPrice(Number(e.target.value))} /></label>
          <button disabled={!assetId} onClick={propose}>2 · Propor ao Banco ({protocol === "htlc" ? "A" : "B"})</button>
        </div>
        {msg && <p className="hint">{msg}</p>}
      </section>
      <section className="card">
        <h3>Swaps (Vendedora)</h3>
        <table className="swaps">
          <thead><tr><th>swapId</th><th>abordagem</th><th>preço</th><th>status</th></tr></thead>
          <tbody>
            {swaps.map((s) => (
              <tr key={s.swapId} className={s.swapId === selected ? "sel" : ""} onClick={() => setSelected(s.swapId)}>
                <td><code>{s.swapId.slice(0, 16)}…</code></td>
                <td>{s.protocol === "htlc" ? "A · HTLC" : "B · atestado"}</td>
                <td>{brl(s.price)}</td>
                <td><span className={`badge tone-${statusTone(s.status)}`}>{statusLabel[s.status]}</span></td>
              </tr>
            ))}
            {swaps.length === 0 && <tr><td colSpan={4} className="hint">Nenhum swap ainda (ou nó fora do ar).</td></tr>}
          </tbody>
        </table>
      </section>
      {current && (
        <section className="card">
          <h3>Execução real · {statusLabel[current.status]}</h3>
          <SequenceDiagram events={events} now={Number.MAX_SAFE_INTEGER} protocol={current.protocol} />
        </section>
      )}
    </div>
  );
}
