package token

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// KV é o armazenamento durável do nó (store SQLite). Os registros deste pacote (gravames, locks
// auditados e expectativas de lock) fazem write-through nele e são recarregados no boot. Sem KV
// anexado, o comportamento é o do MVP (só memória).
type KV interface {
	Put(ctx context.Context, bucket, key string, value []byte) error
	Delete(ctx context.Context, bucket, key string) error
	All(ctx context.Context, bucket string) (map[string][]byte, error)
}

const (
	bucketGravames     = "gravames"
	bucketLocks        = "audited_locks"
	bucketExpectations = "lock_expectations"
)

// AttachKV anexa o armazenamento durável aos registros padrão e recarrega o estado salvo.
// Chamado pela SDK do Janus no boot, antes de qualquer view atender pedidos.
func AttachKV(ctx context.Context, kv KV) error {
	if err := DefaultGravames.attach(ctx, kv); err != nil {
		return fmt.Errorf("gravames: %w", err)
	}
	if err := DefaultLocks.attach(ctx, kv); err != nil {
		return fmt.Errorf("locks auditados: %w", err)
	}
	if err := DefaultExpectation.attach(ctx, kv); err != nil {
		return fmt.Errorf("expectativas: %w", err)
	}
	return nil
}

func hashKey(h []byte) string { return hex.EncodeToString(h) }

// ------------------------------------------------------------------ gravames

func (g *Gravames) attach(ctx context.Context, kv KV) error {
	all, err := kv.All(ctx, bucketGravames)
	if err != nil {
		return err
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.kv = kv
	for id, motivo := range all {
		g.set[id] = string(motivo)
	}
	return nil
}

// ------------------------------------------------------------------ locks auditados

func (r *LockRegistry) attach(ctx context.Context, kv KV) error {
	all, err := kv.All(ctx, bucketLocks)
	if err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.kv = kv
	for _, raw := range all {
		var x LockRecord
		if err := json.Unmarshal(raw, &x); err != nil {
			return err
		}
		r.m[string(x.Hash)] = x
	}
	return nil
}

// ------------------------------------------------------------------ expectativas

func (e *Expectations) attach(ctx context.Context, kv KV) error {
	all, err := kv.All(ctx, bucketExpectations)
	if err != nil {
		return err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.kv = kv
	for _, raw := range all {
		var x Expectation
		if err := json.Unmarshal(raw, &x); err != nil {
			return err
		}
		e.m[string(x.Hash)] = x
	}
	return nil
}
