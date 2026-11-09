package cacti

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// fake reproduz o formato real das respostas do conector Fabric do Cacti v3.0.1 (capturado no spike).
func fake(t *testing.T, filter string, endorsers []string, txInBlock string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req map[string]any
		_ = json.Unmarshal(body, &req)
		switch {
		case strings.HasSuffix(r.URL.Path, "/get-transaction-receipt-by-txid"):
			if req["methodName"] != "GetBlockByTxID" || req["contractName"] != "qscc" {
				t.Errorf("receipt request: %s", body)
			}
			var es []map[string]string
			for _, e := range endorsers {
				es = append(es, map[string]string{"mspid": e, "endorserID": "cert", "signature": "sig"})
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"blockNumber":   map[string]any{"low": 8, "high": 0, "unsigned": true},
				"chainCodeName": "tns", "responseStatus": 200, "transactionEndorsement": es,
			})
		case strings.HasSuffix(r.URL.Path, "/get-block"):
			q := req["query"].(map[string]any)
			if q["blockNumber"] != "8" || req["responseType"] != "full" {
				t.Errorf("get-block request: %s", body)
			}
			_, _ = io.WriteString(w, `{"decodedBlock":{"data":{"data":[
				{"payload":{"header":{"channel_header":{"tx_id":"other"}}}},
				{"payload":{"header":{"channel_header":{"tx_id":"`+txInBlock+`"}}}}]},
				"metadata":{"metadata":[{},{},`+filter+`,{}]}}}`)
		default:
			http.NotFound(w, r)
		}
	}))
}

func connector(url string) *FabricConnector {
	return &FabricConnector{BaseURL: url, KeychainID: "janus", KeychainRef: "registradora", Channel: "testchannel"}
}

func TestTxProofValid(t *testing.T) {
	srv := fake(t, `[0,0]`, []string{"RegistradoraMSP", "BancoMSP"}, "tx-lock")
	defer srv.Close()
	p, err := connector(srv.URL).TxProof(context.Background(), "tx-lock")
	if err != nil {
		t.Fatal(err)
	}
	if p.BlockNumber != 8 || !p.Valid || p.Chaincode != "tns" || len(p.Endorsers) != 2 {
		t.Fatalf("%+v", p)
	}
	if err := p.Satisfies("tns", []string{"RegistradoraMSP", "BancoMSP"}); err != nil {
		t.Fatal(err)
	}
	if p.Satisfies("outro", nil) == nil {
		t.Fatal("chaincode divergente deveria falhar")
	}
}

func TestTxProofInvalidatedInBlock(t *testing.T) {
	// MVCC_READ_CONFLICT (11) para a segunda transação, em formato Buffer
	srv := fake(t, `{"type":"Buffer","data":[0,11]}`, []string{"RegistradoraMSP", "BancoMSP"}, "tx-lock")
	defer srv.Close()
	p, err := connector(srv.URL).TxProof(context.Background(), "tx-lock")
	if err != nil {
		t.Fatal(err)
	}
	if p.Valid || p.Code != 11 || p.Satisfies("tns", nil) == nil {
		t.Fatalf("transação invalidada aceita: %+v", p)
	}
}

func TestTxProofMissingEndorser(t *testing.T) {
	srv := fake(t, `[0,0]`, []string{"VendedoraMSP"}, "tx-lock")
	defer srv.Close()
	p, err := connector(srv.URL).TxProof(context.Background(), "tx-lock")
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Satisfies("tns", []string{"RegistradoraMSP", "BancoMSP"}); err == nil {
		t.Fatal("endosso só da Vendedora não cumpre a política")
	}
}

func TestTxProofTxNotInBlock(t *testing.T) {
	srv := fake(t, `[0,0]`, []string{"RegistradoraMSP"}, "outra-tx")
	defer srv.Close()
	if _, err := connector(srv.URL).TxProof(context.Background(), "tx-lock"); err == nil {
		t.Fatal("tx ausente do bloco")
	}
}

func TestTxProofHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"error":"boom"}`, http.StatusInternalServerError)
	}))
	defer srv.Close()
	if _, err := connector(srv.URL).TxProof(context.Background(), "x"); err == nil || !strings.Contains(err.Error(), "500") {
		t.Fatalf("err=%v", err)
	}
}

func TestFlexUint(t *testing.T) {
	for in, want := range map[string]uint64{`8`: 8, `"42"`: 42, `{"low":5,"high":1,"unsigned":true}`: 1<<32 + 5} {
		var f flexUint
		if err := json.Unmarshal([]byte(in), &f); err != nil || uint64(f) != want {
			t.Fatalf("%s -> %d (%v)", in, f, err)
		}
	}
}
