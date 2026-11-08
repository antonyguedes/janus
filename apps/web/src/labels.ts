import type { ActionKind, Status } from "./types";

export const actionLabel: Record<ActionKind, string> = {
  wait: "aguardar",
  lock_asset: "htlc.Lock (NFT)",
  lock_payment: "lock / deposit (Drex)",
  claim_payment: "claim(S)",
  request_attest: "coletar atestações",
  settle_payment: "settle(S, atestação)",
  claim_asset: "htlc.Claim(S)",
  refund_payment: "refund",
  reclaim_asset: "htlc.Reclaim",
  done: "concluído",
  abort: "abortado",
  alert: "alerta",
};

export const statusLabel: Record<Status, string> = {
  negotiated: "Termos assinados",
  asset_locked: "NFT travado",
  payment_locked: "Drex travado",
  both_locked: "Duas pernas travadas",
  revealed: "S revelado",
  completed: "Concluído",
  reverted: "Revertido",
  anomaly_seller_has_both: "Anomalia: Vendedora com tudo",
  anomaly_buyer_has_both: "Anomalia: Banco com tudo",
};

export function statusTone(s: Status): "ok" | "bad" | "warn" | "neutral" {
  if (s === "completed") return "ok";
  if (s.startsWith("anomaly")) return "bad";
  if (s === "reverted") return "warn";
  return "neutral";
}

/** Lifeline de destino de cada ação no diagrama de sequência. */
export function targetOf(a: ActionKind): "fabric" | "besu" | "relay" | null {
  switch (a) {
    case "lock_asset":
    case "claim_asset":
    case "reclaim_asset":
      return "fabric";
    case "lock_payment":
    case "claim_payment":
    case "settle_payment":
    case "refund_payment":
      return "besu";
    case "request_attest":
      return "relay";
    default:
      return null;
  }
}

export const brl = (cents: number) =>
  (cents / 100).toLocaleString("pt-BR", { style: "currency", currency: "BRL" });
