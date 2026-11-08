import { describe, expect, it } from "vitest";
import { keyTimes, nextKey, prevKey, snapshotAt } from "./playback";
import type { SimResult, Snapshot } from "./types";

const snap = (t: number, fabricLabel: string): Snapshot => ({
  t, fabric: 0, fabricLabel, besu: 0, besuLabel: "none", status: "negotiated",
  seller: { kind: "wait", reason: "" }, buyer: { kind: "wait", reason: "" },
  sellerOnline: true, buyerOnline: true, revealed: false, attested: false, buyerKnowsS: false, fabricDeadline: 0,
});

const r: SimResult = {
  config: { protocol: "htlc", scenario: "happy", latency: 4 },
  t1: 360, t2: 180, delta: 90,
  snapshots: [snap(0, "none"), snap(4, "locked"), snap(12, "claimed")],
  events: [{ t: 4, actor: "vendedora", action: "lock_asset", chain: "fabric", ok: true }],
  outcome: { status: "completed", sellerHasNft: false, sellerHasCash: true, buyerHasNft: true, buyerHasCash: false },
};

describe("playback", () => {
  it("key times merge snapshots and events", () => {
    expect(keyTimes(r)).toEqual([0, 4, 12]);
  });
  it("snapshotAt picks the last snapshot not after now", () => {
    expect(snapshotAt(r, 0).fabricLabel).toBe("none");
    expect(snapshotAt(r, 5).fabricLabel).toBe("locked");
    expect(snapshotAt(r, 999).fabricLabel).toBe("claimed");
  });
  it("step navigation", () => {
    const k = keyTimes(r);
    expect(nextKey(k, 0)).toBe(4);
    expect(nextKey(k, 12)).toBe(12);
    expect(prevKey(k, 12)).toBe(4);
    expect(prevKey(k, 0)).toBe(0);
  });
});
