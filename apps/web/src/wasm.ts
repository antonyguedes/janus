// Carrega o núcleo do protocolo compilado para WebAssembly (node/cmd/simwasm).
// O modo offline usa exatamente o mesmo swap.NextAction dos nós Janus.
import type { Protocol, ScenarioInfo, SimResult } from "./types";

declare global {
  // definido por /wasm_exec.js (Go)
  // eslint-disable-next-line no-var
  var Go: new () => { importObject: WebAssembly.Imports; run(i: WebAssembly.Instance): Promise<void> };
  // eslint-disable-next-line no-var
  var janusScenarios: (() => string) | undefined;
  // eslint-disable-next-line no-var
  var janusRun: ((cfg: string) => string) | undefined;
  // eslint-disable-next-line no-var
  var janusReady: boolean | undefined;
}

let loading: Promise<void> | null = null;

export function loadCore(url = "/janus.wasm"): Promise<void> {
  if (loading) return loading;
  loading = (async () => {
    if (typeof globalThis.Go !== "function") throw new Error("wasm_exec.js não carregado");
    const go = new globalThis.Go();
    const res = await WebAssembly.instantiateStreaming(fetch(url), go.importObject);
    void go.run(res.instance);
    for (let i = 0; i < 500 && !globalThis.janusReady; i++) {
      await new Promise((r) => setTimeout(r, 10));
    }
    if (!globalThis.janusReady) throw new Error("núcleo WASM não inicializou");
  })();
  return loading;
}

export function scenarios(): ScenarioInfo[] {
  if (!globalThis.janusScenarios) throw new Error("núcleo não carregado");
  return JSON.parse(globalThis.janusScenarios()) as ScenarioInfo[];
}

export function runScenario(protocol: Protocol, scenario: string, latency = 4): SimResult {
  if (!globalThis.janusRun) throw new Error("núcleo não carregado");
  const out = JSON.parse(globalThis.janusRun(JSON.stringify({ protocol, scenario, latency })));
  if (out.error) throw new Error(out.error);
  return out as SimResult;
}
