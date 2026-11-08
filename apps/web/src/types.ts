// Tipos espelhados do Go (node/internal/swap, node/internal/swap/sim, node/internal/api).

export type Protocol = "htlc" | "cacti";
export type Role = "vendedora" | "banco";

export type ActionKind =
  | "wait"
  | "lock_asset"
  | "lock_payment"
  | "claim_payment"
  | "request_attest"
  | "settle_payment"
  | "claim_asset"
  | "refund_payment"
  | "reclaim_asset"
  | "done"
  | "abort"
  | "alert";

export type Status =
  | "negotiated"
  | "asset_locked"
  | "payment_locked"
  | "both_locked"
  | "revealed"
  | "completed"
  | "reverted"
  | "anomaly_seller_has_both"
  | "anomaly_buyer_has_both";

export interface Decision {
  kind: ActionKind;
  reason: string;
}

export interface ScenarioInfo {
  id: string;
  title: string;
  description: string;
  success: boolean;
}

export interface SimEvent {
  t: number;
  actor: Role;
  action: ActionKind;
  chain: "fabric" | "besu" | "relay";
  ok: boolean;
  note?: string;
}

export interface Snapshot {
  t: number;
  fabric: number;
  fabricLabel: string;
  besu: number;
  besuLabel: string;
  status: Status;
  seller: Decision;
  buyer: Decision;
  sellerOnline: boolean;
  buyerOnline: boolean;
  revealed: boolean;
  attested: boolean;
  buyerKnowsS: boolean;
  fabricDeadline: number;
}

export interface Outcome {
  status: Status;
  sellerHasNft: boolean;
  sellerHasCash: boolean;
  buyerHasNft: boolean;
  buyerHasCash: boolean;
}

export interface SimResult {
  config: { protocol: Protocol; scenario: string; latency: number };
  t1: number;
  t2: number;
  delta: number;
  snapshots: Snapshot[];
  events: SimEvent[];
  outcome: Outcome;
}

// --- modo live (API dos nós) ---

export interface NodeInfo {
  node: string;
  role: string;
  mspId: string;
  evm?: string;
  protocols?: Protocol[];
  attestor?: { mspId: string; qx: string; qy: string };
}

export interface SwapDTO {
  swapId: string;
  role: Role;
  protocol: Protocol;
  assetId: string;
  price: number;
  hashlock: string;
  t1: number;
  t2: number;
  status: Status;
  terminal: boolean;
  updatedAt: string;
}

export interface StepDTO {
  swapId: string;
  role: Role;
  at: string;
  status: Status;
  decision: Decision;
  fabric: string;
  besu: string;
  error?: string;
  executed: boolean;
}
