import { useEffect, useMemo, useState } from "react";
import type { Protocol, ScenarioInfo, SimResult } from "./types";
import { loadCore, runScenario, scenarios as listScenarios } from "./wasm";
import { endOf, keyTimes, nextKey, prevKey, snapshotAt } from "./playback";
import { SequenceDiagram } from "./components/SequenceDiagram";
import { Timeline } from "./components/Timeline";
import { Inspector } from "./components/Inspector";
import { LivePanel } from "./components/LivePanel";
import { JanusMark } from "./components/JanusMark";
import { PortalBridge } from "./components/PortalBridge";

type Theme = "auto" | "dark" | "light";

function useTheme(): [Theme, () => void] {
  const [theme, setTheme] = useState<Theme>(() => {
    const fromURL = new URLSearchParams(globalThis.location?.search ?? "").get("theme");
    if (fromURL === "dark" || fromURL === "light") return fromURL;
    try {
      return (localStorage.getItem("janus.theme") as Theme) || "auto";
    } catch {
      return "auto";
    }
  });
  useEffect(() => {
    const root = document.documentElement;
    if (theme === "auto") delete root.dataset.theme;
    else root.dataset.theme = theme;
    try { localStorage.setItem("janus.theme", theme); } catch { /* sem storage */ }
  }, [theme]);
  const next = () => setTheme((t) => (t === "auto" ? "dark" : t === "dark" ? "light" : "auto"));
  return [theme, next];
}

type Mode = "offline" | "live";
type Tab = "success" | "failure";

const onlyA = new Set(["reveal-t2-minus", "late-reveal-leak", "buyer-offline-past-t1"]);

// Parâmetros de URL (links compartilháveis): ?protocol=cacti&scenario=seller-no-reveal&t=end&mode=live
const qs = new URLSearchParams(globalThis.location?.search ?? "");
const initialT = qs.get("t");

export function App() {
  const [theme, cycleTheme] = useTheme();
  const [ready, setReady] = useState(false);
  const [err, setErr] = useState("");
  const [mode, setMode] = useState<Mode>(qs.get("mode") === "live" ? "live" : "offline");
  const [protocol, setProtocol] = useState<Protocol>(qs.get("protocol") === "cacti" ? "cacti" : "htlc");
  const [tab, setTab] = useState<Tab>("success");
  const [catalog, setCatalog] = useState<ScenarioInfo[]>([]);
  const [scenario, setScenario] = useState(qs.get("scenario") ?? "happy");
  const [result, setResult] = useState<SimResult | null>(null);
  const [now, setNow] = useState(0);
  const [playing, setPlaying] = useState(false);
  const [speed, setSpeed] = useState(8);

  useEffect(() => {
    loadCore()
      .then(() => {
        const cat = listScenarios();
        setCatalog(cat);
        const sc = cat.find((s) => s.id === scenario);
        if (sc) setTab(sc.success ? "success" : "failure");
        setReady(true);
      })
      .catch((e: Error) => setErr(e.message));
  }, []);

  const visible = catalog.filter((s) => (tab === "success") === s.success && !(protocol === "cacti" && onlyA.has(s.id)));

  useEffect(() => {
    if (!ready) return;
    if (!visible.some((s) => s.id === scenario) && visible[0]) {
      setScenario(visible[0].id);
      return;
    }
    try {
      const r = runScenario(protocol, scenario);
      setResult(r);
      setNow(initialT === "end" ? endOf(r) : Number(initialT ?? 0) || 0);
      setPlaying(false);
      setErr("");
    } catch (e) { setErr((e as Error).message); }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [ready, protocol, scenario, tab]);

  const keys = useMemo(() => (result ? keyTimes(result) : [0]), [result]);
  const end = result ? endOf(result) : 1;

  useEffect(() => {
    if (!playing) return;
    const id = setInterval(() => {
      setNow((t) => {
        const n = Math.min(end, t + 1);
        if (n >= end) setPlaying(false);
        return n;
      });
    }, 1000 / speed);
    return () => clearInterval(id);
  }, [playing, speed, end]);

  const snap = result ? snapshotAt(result, now) : null;
  const info = catalog.find((s) => s.id === scenario);

  return (
    <div className="app">
      <header>
        <div className="brand">
          <JanusMark size={48} />
          <div>
            <h1>Janus</h1>
            <p>Guardião das passagens · recebível no <span className="t-gold">Fabric</span> ⇄ Drex no <span className="t-cyan">Besu</span></p>
          </div>
        </div>
        <div className="toggles">
          <div className="seg" role="tablist" aria-label="Abordagem">
            <button className={protocol === "htlc" ? "on" : ""} onClick={() => setProtocol("htlc")}>A · HTLC</button>
            <button className={protocol === "cacti" ? "on" : ""} onClick={() => setProtocol("cacti")}>B · HTLC atestado</button>
          </div>
          <div className="seg" role="tablist" aria-label="Modo">
            <button className={mode === "offline" ? "on" : ""} onClick={() => setMode("offline")}>Offline</button>
            <button className={mode === "live" ? "on" : ""} onClick={() => setMode("live")}>Live</button>
          </div>
          <button className="theme-btn" onClick={cycleTheme} title={`Tema: ${theme}`} aria-label={`Tema: ${theme}`}>
            {theme === "auto" ? "◐" : theme === "dark" ? "☾" : "☀"}
          </button>
        </div>
      </header>

      {err && <div className="error">{err}</div>}

      {mode === "live" ? (
        <LivePanel protocol={protocol} />
      ) : !ready || !result || !snap ? (
        <p className="hint">Carregando o núcleo do protocolo (WebAssembly)…</p>
      ) : (
        <main>
          <nav className="tabs">
            <button className={tab === "success" ? "on" : ""} onClick={() => setTab("success")}>Sucesso</button>
            <button className={tab === "failure" ? "on" : ""} onClick={() => setTab("failure")}>Timeout / Falha</button>
          </nav>
          <div className="chips">
            {visible.map((s) => (
              <button key={s.id} className={`chip ${s.id === scenario ? "on" : ""}`} onClick={() => setScenario(s.id)}>{s.title}</button>
            ))}
          </div>
          {info && <p className="scenario-desc">{info.description}</p>}

          <PortalBridge snap={snap} protocol={protocol} />

          <Timeline now={now} end={end} t1={result.t1} t2={result.t2} delta={result.delta}
            marks={result.events.map((e) => e.t)} onSeek={(t) => { setNow(t); setPlaying(false); }} />

          <div className="controls">
            <button onClick={() => { setNow(prevKey(keys, now)); setPlaying(false); }} aria-label="Passo anterior">⏮ passo</button>
            <button className="primary" onClick={() => { if (now >= end) setNow(0); setPlaying(!playing); }}>{playing ? "⏸ pausar" : "▶ reproduzir"}</button>
            <button onClick={() => { setNow(nextKey(keys, now)); setPlaying(false); }} aria-label="Próximo passo">passo ⏭</button>
            <button onClick={() => { setNow(end); setPlaying(false); }}>fim</button>
            <label className="speed">velocidade
              <select value={speed} onChange={(e) => setSpeed(Number(e.target.value))}>
                {[2, 8, 30, 90].map((v) => <option key={v} value={v}>{v}×</option>)}
              </select>
            </label>
          </div>

          <div className="grid">
            <section className="card diagram">
              <SequenceDiagram events={result.events} now={now} protocol={protocol} />
            </section>
            <Inspector snap={snap} protocol={protocol} t1={result.t1} t2={result.t2}
              outcome={result.outcome} finished={now >= end} />
          </div>
          <p className="footnote">
            Modo offline: a lógica de decisão é o mesmo <code>swap.NextAction</code> dos nós Janus, compilado para WebAssembly.
            Tema original inspirado em Jano, deus romano das passagens.
            Prazos de desenvolvimento: T2 = {result.t2}s, T1 = {result.t1}s, Δ ≥ {result.delta}s.
          </p>
        </main>
      )}
    </div>
  );
}
