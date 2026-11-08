import type { SimResult, Snapshot } from "./types";

/** Instantes relevantes (eventos e mudanças de estado/decisão), para a navegação passo a passo. */
export function keyTimes(r: SimResult): number[] {
  const s = new Set<number>([0]);
  for (const x of r.snapshots) s.add(x.t);
  for (const e of r.events) s.add(e.t);
  return [...s].sort((a, b) => a - b);
}

/** Último snapshot com t <= now. */
export function snapshotAt(r: SimResult, now: number): Snapshot {
  let cur = r.snapshots[0]!;
  for (const x of r.snapshots) {
    if (x.t > now) break;
    cur = x;
  }
  return cur;
}

export function nextKey(keys: number[], now: number): number {
  return keys.find((k) => k > now) ?? keys[keys.length - 1] ?? now;
}

export function prevKey(keys: number[], now: number): number {
  const before = keys.filter((k) => k < now);
  return before[before.length - 1] ?? 0;
}

export function endOf(r: SimResult): number {
  return r.t1 + 120;
}
