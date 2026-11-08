// Package api expõe o nó Janus para o simulador web e para operadores: REST para consultar
// swaps e SSE (/api/v1/events) para acompanhar cada passo do motor em tempo real.
package api

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"janus/node/internal/engine"
	"janus/node/internal/store"
	"janus/node/internal/swap"
)

// Info identifica o nó (papel, MSP, endereço EVM, protocolos habilitados).
type Info struct {
	Node      string          `json:"node"`
	Role      swap.Role       `json:"role"`
	MSPID     string          `json:"mspId"`
	EVM       string          `json:"evm"`
	Protocols []swap.Protocol `json:"protocols"`
	Faults    engine.Faults   `json:"faults"`
}

// Server é o handler HTTP.
type Server struct {
	Info   Info
	Store  *store.Store
	Engine *engine.Engine
	Log    *slog.Logger
	// AllowOrigin habilita CORS para o simulador (ex.: http://localhost:5173).
	AllowOrigin string

	// Propose (só na Vendedora) inicia a negociação de um swap com o Banco.
	Propose func(r *http.Request, req ProposeRequest) (SwapDTO, error)
	// Issue (só na Registradora) emite um recebível para um owner.
	Issue func(r *http.Request, req IssueRequest) (string, error)
}

// ProposeRequest é o corpo de POST /api/v1/swaps.
type ProposeRequest struct {
	BuyerNode string        `json:"buyerNode"`
	AssetID   string        `json:"assetId"`
	Price     uint64        `json:"price"` // centavos de Drex
	Protocol  swap.Protocol `json:"protocol"`
}

// IssueRequest é o corpo de POST /api/v1/recebiveis.
type IssueRequest struct {
	Owner      string `json:"owner"` // nó FSC do titular
	Registro   string `json:"registro"`
	Cedente    string `json:"cedente"`
	Sacado     string `json:"sacado"`
	ValorFace  uint64 `json:"valorFace"`
	Vencimento string `json:"vencimento"`
}

// Handler monta as rotas.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", s.health)
	mux.HandleFunc("GET /api/v1/info", func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, http.StatusOK, s.Info) })
	mux.HandleFunc("GET /api/v1/swaps", s.list)
	mux.HandleFunc("GET /api/v1/swaps/{id}", s.get)
	mux.HandleFunc("POST /api/v1/swaps/{id}/kick", s.kick)
	mux.HandleFunc("GET /api/v1/events", s.events)
	mux.HandleFunc("POST /api/v1/swaps", s.propose)
	mux.HandleFunc("POST /api/v1/recebiveis", s.issue)
	return s.cors(mux)
}

func (s *Server) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.AllowOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", s.AllowOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "time": time.Now().UTC()})
}

// SwapDTO é a representação pública de um swap.
type SwapDTO struct {
	SwapID    string          `json:"swapId"`
	Role      swap.Role       `json:"role"`
	Protocol  swap.Protocol   `json:"protocol"`
	AssetID   string          `json:"assetId"`
	Price     uint64          `json:"price"`
	Hashlock  string          `json:"hashlock"`
	T1        int64           `json:"t1"`
	T2        int64           `json:"t2"`
	Status    swap.Status     `json:"status"`
	Terminal  bool            `json:"terminal"`
	UpdatedAt time.Time       `json:"updatedAt"`
	Terms     *swap.SwapTerms `json:"terms,omitempty"`
}

func toDTO(r store.Record, full bool) SwapDTO {
	t := r.Signed.Terms
	d := SwapDTO{
		SwapID: hex.EncodeToString(r.SwapID[:]), Role: r.Role, Protocol: t.Protocol, AssetID: t.AssetID,
		Price: t.Price.Value, Hashlock: hex.EncodeToString(t.Hashlock[:]), T1: t.T1, T2: t.T2,
		Status: r.Status, Terminal: r.Terminal, UpdatedAt: r.UpdatedAt,
	}
	if full {
		d.Terms = &t
	}
	return d
}

func (s *Server) list(w http.ResponseWriter, r *http.Request) {
	recs, err := s.Store.List(r.Context(), r.URL.Query().Get("active") == "true")
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	out := make([]SwapDTO, 0, len(recs))
	for _, rec := range recs {
		out = append(out, toDTO(rec, false))
	}
	writeJSON(w, http.StatusOK, out)
}

func parseID(raw string) ([32]byte, error) {
	var id [32]byte
	b, err := hex.DecodeString(strings.TrimPrefix(raw, "0x"))
	if err != nil || len(b) != 32 {
		return id, errors.New("swapId must be 32 bytes hex")
	}
	copy(id[:], b)
	return id, nil
}

func (s *Server) get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	rec, err := s.Store.Get(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		writeErr(w, http.StatusNotFound, err)
		return
	} else if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, toDTO(rec, true))
}

func (s *Server) propose(w http.ResponseWriter, r *http.Request) {
	if s.Propose == nil {
		writeErr(w, http.StatusMethodNotAllowed, errors.New("only the seller node proposes swaps"))
		return
	}
	var req ProposeRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	if req.Protocol == "" {
		req.Protocol = swap.ProtocolHTLC
	}
	if req.BuyerNode == "" || req.AssetID == "" || req.Price == 0 {
		writeErr(w, http.StatusBadRequest, errors.New("buyerNode, assetId and price are required"))
		return
	}
	dto, err := s.Propose(r, req)
	if err != nil {
		writeErr(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusCreated, dto)
}

func (s *Server) issue(w http.ResponseWriter, r *http.Request) {
	if s.Issue == nil {
		writeErr(w, http.StatusMethodNotAllowed, errors.New("only the registradora issues recebiveis"))
		return
	}
	var req IssueRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	id, err := s.Issue(r, req)
	if err != nil {
		writeErr(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"linearId": id})
}

// ToDTO converte um registro para a representação pública.
func ToDTO(r store.Record) SwapDTO { return toDTO(r, true) }

func (s *Server) kick(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	s.Engine.Kick(id)
	w.WriteHeader(http.StatusAccepted)
}

// StepDTO é um passo do motor no formato do SSE.
type StepDTO struct {
	engine.Step
	SwapID string `json:"swapId"`
	Fabric string `json:"fabric"`
	Besu   string `json:"besu"`
}

// events transmite os passos do motor como Server-Sent Events (evento "step"), com heartbeat.
func (s *Server) events(w http.ResponseWriter, r *http.Request) {
	fl, ok := w.(http.Flusher)
	if !ok {
		writeErr(w, http.StatusInternalServerError, errors.New("streaming unsupported"))
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	steps, unsub := s.Engine.Subscribe()
	defer unsub()
	fmt.Fprintf(w, "event: hello\ndata: %s\n\n", mustJSON(s.Info))
	fl.Flush()
	hb := time.NewTicker(15 * time.Second)
	defer hb.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case <-hb.C:
			fmt.Fprint(w, ": heartbeat\n\n")
			fl.Flush()
		case st := <-steps:
			dto := StepDTO{Step: st, SwapID: hex.EncodeToString(st.SwapID[:]), Fabric: st.Fabric.String(), Besu: st.Besu.String()}
			fmt.Fprintf(w, "event: step\ndata: %s\n\n", mustJSON(dto))
			fl.Flush()
		}
	}
}

func mustJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, err error) {
	writeJSON(w, code, map[string]string{"error": err.Error()})
}
