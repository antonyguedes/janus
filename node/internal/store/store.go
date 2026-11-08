// Package store persiste os swaps do nó (termos assinados, papel e último status) em SQLite.
// A fonte da verdade do andamento são as cadeias. O store guarda o que só o nó conhece
// (termos assinados, papel) e permite a reconciliação no boot.
package store

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite"

	"janus/node/internal/swap"
)

// Record é um swap conhecido pelo nó.
type Record struct {
	SwapID    [32]byte         `json:"swapId"`
	Role      swap.Role        `json:"role"`
	Signed    swap.SignedTerms `json:"signed"`
	Status    swap.Status      `json:"status"`
	Terminal  bool             `json:"terminal"`
	UpdatedAt time.Time        `json:"updatedAt"`
}

// Store em SQLite.
type Store struct{ db *sql.DB }

var ErrNotFound = errors.New("store: swap not found")

// Open abre (ou cria) o banco em path. ":memory:" funciona para testes.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1) // SQLite: um escritor
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS swaps (
		id TEXT PRIMARY KEY,
		role TEXT NOT NULL,
		signed TEXT NOT NULL,
		status TEXT NOT NULL,
		terminal INTEGER NOT NULL DEFAULT 0,
		updated_at INTEGER NOT NULL
	)`); err != nil {
		return nil, fmt.Errorf("store: migrate: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

func idHex(id [32]byte) string { return hex.EncodeToString(id[:]) }

// Put insere ou atualiza o registro.
func (s *Store) Put(ctx context.Context, r Record) error {
	raw, err := json.Marshal(r.Signed)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, `INSERT INTO swaps (id, role, signed, status, terminal, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET signed=excluded.signed, status=excluded.status,
			terminal=excluded.terminal, updated_at=excluded.updated_at`,
		idHex(r.SwapID), string(r.Role), string(raw), string(r.Status), boolInt(r.Terminal), time.Now().UnixMilli())
	return err
}

// SetStatus atualiza só o status, chamado a cada passo do motor.
func (s *Store) SetStatus(ctx context.Context, id [32]byte, st swap.Status) error {
	res, err := s.db.ExecContext(ctx, `UPDATE swaps SET status=?, terminal=?, updated_at=? WHERE id=?`,
		string(st), boolInt(st.Terminal()), time.Now().UnixMilli(), idHex(id))
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) Get(ctx context.Context, id [32]byte) (Record, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, role, signed, status, terminal, updated_at FROM swaps WHERE id=?`, idHex(id))
	r, err := scan(row)
	if errors.Is(err, sql.ErrNoRows) {
		return r, ErrNotFound
	}
	return r, err
}

// List devolve os swaps, mais recentes primeiro. activeOnly filtra os não terminais
// (usado na reconciliação de boot).
func (s *Store) List(ctx context.Context, activeOnly bool) ([]Record, error) {
	q := `SELECT id, role, signed, status, terminal, updated_at FROM swaps`
	if activeOnly {
		q += ` WHERE terminal = 0`
	}
	rows, err := s.db.QueryContext(ctx, q+` ORDER BY updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Record
	for rows.Next() {
		r, err := scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

type scanner interface{ Scan(...any) error }

func scan(sc scanner) (Record, error) {
	var (
		r                   Record
		id, role, raw, stat string
		term                int
		upd                 int64
	)
	if err := sc.Scan(&id, &role, &raw, &stat, &term, &upd); err != nil {
		return r, err
	}
	b, err := hex.DecodeString(id)
	if err != nil || len(b) != 32 {
		return r, fmt.Errorf("store: bad id %q", id)
	}
	copy(r.SwapID[:], b)
	r.Role, r.Status, r.Terminal, r.UpdatedAt = swap.Role(role), swap.Status(stat), term == 1, time.UnixMilli(upd)
	return r, json.Unmarshal([]byte(raw), &r.Signed)
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
