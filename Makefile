# Janus — orquestração do monorepo. `make help` lista os alvos.
SHELL := /usr/bin/env bash
.DEFAULT_GOAL := help

BESU_DIR := infra/besu
COMPOSE_BESU := docker compose -f $(BESU_DIR)/docker-compose.yml
export UID_GID := $(shell id -u):$(shell id -g)

.PHONY: help
help: ## lista os alvos
	@grep -hE '^[a-zA-Z0-9_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-16s\033[0m %s\n",$$1,$$2}'

# ------------------------------------------------------------------ dependências
.PHONY: deps
deps: ## instala dependências (Foundry libs, módulos Go)
	cd contracts && [ -d lib/openzeppelin-contracts ] || forge install --no-git OpenZeppelin/openzeppelin-contracts@v5.7.0
	cd contracts && [ -d lib/forge-std ] || forge install --no-git foundry-rs/forge-std
	cd node && go mod download

# ------------------------------------------------------------------ redes
.PHONY: up down besu-up besu-down besu-clean
up: besu-up ## sobe as redes locais
down: besu-down ## derruba as redes locais

besu-up: ## sobe o Besu QBFT (4 validadores, gas zero)
	$(BESU_DIR)/scripts/init.sh
	UID=$(shell id -u) GID=$(shell id -g) $(COMPOSE_BESU) up -d
	@echo "RPC: http://localhost:8545"

besu-down: ## para o Besu
	$(COMPOSE_BESU) down

besu-clean: besu-down ## apaga o estado e as chaves do Besu
	rm -rf $(BESU_DIR)/data $(BESU_DIR)/generated

# ------------------------------------------------------------------ geração
.PHONY: bindings
bindings: ## regenera os bindings Go (abigen v2) a partir dos contratos
	./scripts/gen-bindings.sh

# ------------------------------------------------------------------ testes
.PHONY: test test-contracts test-node test-vectors vectors
test: test-vectors test-contracts test-node ## roda todos os testes offline

vectors: ## regenera os vetores compartilhados Go <-> Solidity
	cd node && go run ./cmd/vectors -out ../testdata/vectors

test-vectors: ## vetores de hash/atestação nos dois lados
	cd node && go test ./internal/hashlock/... ./internal/attest/... -run Vector -count=1
	cd contracts && forge test --match-contract VectorsTest

test-contracts: ## forge test (unit + fuzz + invariantes)
	cd contracts && forge test

test-node: ## go test do nó Janus
	cd node && go test ./... -cover
