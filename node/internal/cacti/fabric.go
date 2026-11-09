// Package cacti é o cliente Go do conector Fabric do Hyperledger Cacti (v3.0.1, REST do
// cactus-cmd-api-server). O atestador do Janus usa o conector como fonte de prova do lock no
// ledger comitado (ADR-0005): bloco, endossantes e validade da transação.
package cacti

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"slices"
	"strconv"
	"time"
)

const fabricPluginPath = "/api/v1/plugins/@hyperledger-cacti/cactus-plugin-ledger-connector-fabric"

// FabricConnector fala com o conector Fabric de um servidor Cacti.
type FabricConnector struct {
	BaseURL     string // ex.: http://127.0.0.1:4100
	KeychainID  string
	KeychainRef string // identidade da org no keychain do servidor Cacti
	Channel     string
	HTTP        *http.Client
}

// TxProof é o que o ledger diz sobre uma transação, visto pelo conector do Cacti.
type TxProof struct {
	TxID        string   `json:"txId"`
	BlockNumber uint64   `json:"blockNumber"`
	Chaincode   string   `json:"chaincode"`
	Endorsers   []string `json:"endorsers"` // MSP IDs dos peers endossantes
	Valid       bool     `json:"valid"`     // TRANSACTIONS_FILTER == VALID (0)
	Code        int      `json:"code"`      // código de validação do Fabric
}

// Satisfies verifica a prova contra o chaincode e as orgs exigidas (política AND).
func (p TxProof) Satisfies(chaincode string, required []string) error {
	if !p.Valid {
		return fmt.Errorf("cacti: tx %s invalidada no bloco %d (código %d)", p.TxID, p.BlockNumber, p.Code)
	}
	if chaincode != "" && p.Chaincode != chaincode {
		return fmt.Errorf("cacti: tx %s é do chaincode %q, esperado %q", p.TxID, p.Chaincode, chaincode)
	}
	for _, msp := range required {
		if !slices.Contains(p.Endorsers, msp) {
			return fmt.Errorf("cacti: tx %s sem endosso de %s (endossantes %v)", p.TxID, msp, p.Endorsers)
		}
	}
	return nil
}

func (c *FabricConnector) client() *http.Client {
	if c.HTTP != nil {
		return c.HTTP
	}
	return &http.Client{Timeout: 30 * time.Second}
}

func (c *FabricConnector) post(ctx context.Context, op string, body, out any) error {
	raw, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+fabricPluginPath+"/"+op, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.client().Do(req)
	if err != nil {
		return fmt.Errorf("cacti %s: %w", op, err)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(io.LimitReader(res.Body, 64<<20))
	if err != nil {
		return err
	}
	if res.StatusCode/100 != 2 {
		return fmt.Errorf("cacti %s: HTTP %d: %.300s", op, res.StatusCode, data)
	}
	return json.Unmarshal(data, out)
}

type credential struct {
	KeychainID  string `json:"keychainId"`
	KeychainRef string `json:"keychainRef"`
}

// receipt é o subconjunto usado de GetTransactionReceiptResponse.
type receipt struct {
	BlockNumber    flexUint `json:"blockNumber"`
	ChainCodeName  string   `json:"chainCodeName"`
	ResponseStatus flexUint `json:"responseStatus"`
	Endorsements   []struct {
		MSPID string `json:"mspid"`
	} `json:"transactionEndorsement"`
}

// fullBlock é o subconjunto usado do bloco decodificado (responseType "full").
type fullBlock struct {
	DecodedBlock struct {
		Data struct {
			Data []struct {
				Payload struct {
					Header struct {
						ChannelHeader struct {
							TxID string `json:"tx_id"`
						} `json:"channel_header"`
					} `json:"header"`
				} `json:"payload"`
			} `json:"data"`
		} `json:"data"`
		Metadata struct {
			Metadata []json.RawMessage `json:"metadata"`
		} `json:"metadata"`
	} `json:"decodedBlock"`
}

// ErrTxNotFound indica que o conector não achou a transação no ledger.
var ErrTxNotFound = errors.New("cacti: transação não encontrada no ledger")

// TxProof obtém do ledger, via Cacti, o bloco e os endossantes da transação (recibo) e a sua
// validade (TRANSACTIONS_FILTER do bloco). O recibo sozinho não basta: uma transação pode estar
// num bloco e ter sido invalidada (MVCC, política de endosso).
func (c *FabricConnector) TxProof(ctx context.Context, txID string) (TxProof, error) {
	p := TxProof{TxID: txID}
	cred := credential{c.KeychainID, c.KeychainRef}
	var rc receipt
	err := c.post(ctx, "get-transaction-receipt-by-txid", map[string]any{
		"signingCredential": cred, "channelName": c.Channel, "contractName": "qscc",
		"invocationType": "FabricContractInvocationType.CALL", "methodName": "GetBlockByTxID",
		"params": []string{c.Channel, txID},
	}, &rc)
	if err != nil {
		return p, err
	}
	p.BlockNumber, p.Chaincode = uint64(rc.BlockNumber), rc.ChainCodeName
	for _, e := range rc.Endorsements {
		p.Endorsers = append(p.Endorsers, e.MSPID)
	}
	if rc.ResponseStatus != 0 && rc.ResponseStatus != 200 {
		return p, fmt.Errorf("cacti: endosso da tx %s com status %d", txID, rc.ResponseStatus)
	}

	var blk fullBlock
	err = c.post(ctx, "get-block", map[string]any{
		"channelName":    c.Channel,
		"gatewayOptions": map[string]any{"identity": c.KeychainRef, "wallet": map[string]any{"keychain": cred}},
		"query":          map[string]any{"blockNumber": strconv.FormatUint(p.BlockNumber, 10)},
		"responseType":   "full",
	}, &blk)
	if err != nil {
		return p, err
	}
	idx := -1
	for i, d := range blk.DecodedBlock.Data.Data {
		if d.Payload.Header.ChannelHeader.TxID == txID {
			idx = i
			break
		}
	}
	if idx < 0 {
		return p, fmt.Errorf("%w: %s não está no bloco %d", ErrTxNotFound, txID, p.BlockNumber)
	}
	md := blk.DecodedBlock.Metadata.Metadata
	if len(md) <= 2 {
		return p, errors.New("cacti: bloco sem TRANSACTIONS_FILTER")
	}
	codes, err := decodeFilter(md[2])
	if err != nil {
		return p, err
	}
	if idx >= len(codes) {
		return p, errors.New("cacti: TRANSACTIONS_FILTER menor que o número de transações")
	}
	p.Code = codes[idx]
	p.Valid = p.Code == 0
	return p, nil
}

// decodeFilter aceita o TRANSACTIONS_FILTER como array de números ou como Buffer serializado
// ({"type":"Buffer","data":[...]}), formatos possíveis do decodificador do fabric-common.
func decodeFilter(raw json.RawMessage) ([]int, error) {
	var arr []int
	if err := json.Unmarshal(raw, &arr); err == nil {
		return arr, nil
	}
	var buf struct {
		Data []int `json:"data"`
	}
	if err := json.Unmarshal(raw, &buf); err == nil && buf.Data != nil {
		return buf.Data, nil
	}
	return nil, fmt.Errorf("cacti: TRANSACTIONS_FILTER em formato desconhecido: %.80s", raw)
}

// flexUint aceita número, string decimal ou Long do protobufjs ({low, high, unsigned}).
type flexUint uint64

func (f *flexUint) UnmarshalJSON(b []byte) error {
	var n float64
	if err := json.Unmarshal(b, &n); err == nil {
		if n < 0 || n > math.MaxUint64 {
			return fmt.Errorf("cacti: número fora do intervalo %v", n)
		}
		*f = flexUint(n)
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		v, err := strconv.ParseUint(s, 10, 64)
		*f = flexUint(v)
		return err
	}
	var l struct {
		Low  int64 `json:"low"`
		High int64 `json:"high"`
	}
	if err := json.Unmarshal(b, &l); err != nil {
		return fmt.Errorf("cacti: número em formato desconhecido: %s", b)
	}
	*f = flexUint(uint64(uint32(l.High))<<32 | uint64(uint32(l.Low)))
	return nil
}
