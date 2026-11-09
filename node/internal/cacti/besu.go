package cacti

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"janus/node/internal/attest"
)

const besuPluginPath = "/api/v1/plugins/@hyperledger-cacti/cactus-plugin-ledger-connector-besu"

// BesuConnector submete transações pelo conector Besu do Cacti (relay da Vendedora, ADR-0006).
// A chave da conta de relay fica no keychain do servidor Cacti. A requisição só carrega a
// referência (CACTUS_KEYCHAIN_REF).
type BesuConnector struct {
	BaseURL     string
	KeychainID  string
	KeychainRef string // entrada do keychain com a chave privada da conta de relay
	Account     string // endereço da conta de relay (0x...)
	Gas         uint64 // limite de gas explícito (evita estimateGas sem "from" no conector)
	HTTP        *http.Client
}

func (b *BesuConnector) client() *http.Client {
	if b.HTTP != nil {
		return b.HTTP
	}
	return &http.Client{Timeout: 90 * time.Second}
}

// InvokeResult é o retorno de invoke-contract (SEND).
type InvokeResult struct {
	Success bool   `json:"success"`
	TxHash  string `json:"txHash"`
}

// Send chama invoke-contract com invocationType SEND, gasPrice 0 (rede com gas zero) e a ABI
// passada na própria requisição (sem pré-registro do contrato no keychain).
func (b *BesuConnector) Send(ctx context.Context, abiJSON, address, method string, params []any) (InvokeResult, error) {
	var raw json.RawMessage
	gas := b.Gas
	if gas == 0 {
		gas = 3_000_000
	}
	err := postJSON(ctx, b.client(), b.BaseURL+besuPluginPath+"/invoke-contract", "invoke-contract", map[string]any{
		"contractName":    "janus-" + method,
		"contractAbi":     json.RawMessage(abiJSON),
		"contractAddress": address,
		"invocationType":  "SEND",
		"methodName":      method,
		"params":          params,
		"gas":             gas,
		"gasPrice":        "0",
		"signingCredential": map[string]any{
			"type": "CACTUS_KEYCHAIN_REF", "ethAccount": b.Account,
			"keychainId": b.KeychainID, "keychainEntryKey": b.KeychainRef,
		},
	}, &raw)
	if err != nil {
		return InvokeResult{}, err
	}
	// O servidor v3.0.1 responde {"success", "out": {"transactionReceipt": {...}}}; o OpenAPI
	// descreve o recibo no topo. Aceita os dois formatos.
	type rcpt struct {
		TransactionHash string `json:"transactionHash"`
		Status          any    `json:"status"`
	}
	var out struct {
		Success            bool `json:"success"`
		TransactionReceipt rcpt `json:"transactionReceipt"`
		Out                struct {
			TransactionReceipt rcpt `json:"transactionReceipt"`
		} `json:"out"`
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return InvokeResult{}, err
	}
	r := out.TransactionReceipt
	if r.TransactionHash == "" {
		r = out.Out.TransactionReceipt
	}
	if !out.Success || r.TransactionHash == "" {
		return InvokeResult{}, fmt.Errorf("cacti besu: invoke-contract sem sucesso ou sem hash: %.400s", raw)
	}
	if st, ok := r.Status.(bool); ok && !st {
		return InvokeResult{TxHash: r.TransactionHash}, fmt.Errorf("cacti besu: transação %s revertida", r.TransactionHash)
	}
	return InvokeResult{Success: true, TxHash: r.TransactionHash}, nil
}

func hex32(b [32]byte) string { return "0x" + hex.EncodeToString(b[:]) }

// SettleParams monta os parâmetros de AttestedSettlement.settle no formato do web3 1.6:
// bytes32 em hex, tuplas como objetos com os nomes da ABI e uint64 como string decimal.
func SettleParams(swapID, preimage [32]byte, a attest.FabricLockAttestation, sigs []attest.OrgSignature) []any {
	att := map[string]any{
		"swapId": hex32(a.SwapID), "hashlock": hex32(a.Hashlock), "assetIdHash": hex32(a.AssetIDHash),
		"recipientIdHash": hex32(a.RecipientIDHash), "deadline": strconv.FormatUint(a.Deadline, 10),
		"channelHash": hex32(a.ChannelHash), "namespaceHash": hex32(a.NamespaceHash),
		"fabricTxId": hex32(a.FabricTxID), "blockNumber": strconv.FormatUint(a.BlockNumber, 10),
	}
	ss := make([]any, len(sigs))
	for i, s := range sigs {
		ss[i] = map[string]any{"mspIdHash": hex32(s.MSPIDHash), "keyIndex": s.KeyIndex, "r": hex32(s.R), "s": hex32(s.S)}
	}
	return []any{hex32(swapID), hex32(preimage), att, ss}
}

// NormalizeAccount devolve o endereço em minúsculas com 0x.
func NormalizeAccount(a string) string {
	a = strings.ToLower(strings.TrimSpace(a))
	if !strings.HasPrefix(a, "0x") {
		a = "0x" + a
	}
	return a
}
