#!/usr/bin/env bash
# Gera node/internal/besu/bindings a partir dos artefatos do Foundry (abigen v2, um único pacote,
# com structs compartilhadas deduplicadas). Uso: make bindings
set -euo pipefail
cd "$(dirname "$0")/.."
GETH_VERSION="${GETH_VERSION:-v1.17.10}"
(cd contracts && forge build -q)
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
python3 - "$tmp/combined.json" <<'PY'
import json, sys
out = {"contracts": {}, "version": "0.8.30"}
for c in ["DrexToken", "HTLCPayment", "FabricProofVerifier", "AttestedSettlement"]:
    a = json.load(open(f"contracts/out/{c}.sol/{c}.json"))
    out["contracts"][f"src/{c}.sol:{c}"] = {"abi": a["abi"], "bin": a["bytecode"]["object"].removeprefix("0x")}
json.dump(out, open(sys.argv[1], "w"))
PY
(cd node && go run "github.com/ethereum/go-ethereum/cmd/abigen@${GETH_VERSION}" --v2 \
  --combined-json "$tmp/combined.json" --pkg bindings --out internal/besu/bindings/bindings.go)
echo "bindings gerados em node/internal/besu/bindings/bindings.go"
