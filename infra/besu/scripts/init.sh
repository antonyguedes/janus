#!/usr/bin/env bash
# Gera (uma vez) genesis QBFT + chaves dos 4 validadores e distribui em data/nodeN.
# Idempotente: não regenera se data/ já existir (use `make besu-clean` para zerar).
set -euo pipefail
cd "$(dirname "$0")/.."
BESU_IMAGE="${BESU_IMAGE:-hyperledger/besu:26.11.0}"

if [[ -f data/genesis.json ]]; then
  echo "besu: data/ já existe, mantendo"; exit 0
fi

rm -rf generated
docker run --rm -u "$(id -u):$(id -g)" -v "$PWD:/work" "$BESU_IMAGE" \
  operator generate-blockchain-config \
  --config-file=/work/config/qbft-config.json \
  --to=/work/generated --private-key-file-name=key >/dev/null

mkdir -p data
cp generated/genesis.json data/genesis.json
i=1
for k in generated/keys/*; do
  mkdir -p "data/node$i"
  cp "$k/key" "$k/key.pub" "data/node$i/"
  cp config/permissions_config.toml "data/node$i/"
  i=$((i+1))
done

# enode do node1 (bootnode) — IP fixo da rede docker
pub=$(sed 's/^0x//' data/node1/key.pub)
echo "enode://${pub}@172.30.0.11:30303" > data/bootnode
echo "besu: rede QBFT gerada em infra/besu/data"
