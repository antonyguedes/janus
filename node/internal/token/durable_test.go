package token

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"janus/node/internal/store"
)

func freshRegistries() (*Gravames, *LockRegistry, *Expectations) {
	return NewGravames(), NewLockRegistry(), NewExpectations()
}

// Estado operacional sobrevive ao reinício do nó (mesmo arquivo SQLite).
func TestDurableRegistriesSurviveRestart(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "janus.db")
	st, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	g, l, e := freshRegistries()
	for _, f := range []func(context.Context, KV) error{g.attach, l.attach, e.attach} {
		if err := f(ctx, st.KV()); err != nil {
			t.Fatal(err)
		}
	}
	if err := g.Set("nft-1", "penhor judicial"); err != nil {
		t.Fatal(err)
	}
	_ = g.Set("nft-2", "alienação")
	_ = g.Clear("nft-2")
	h := []byte("0123456789abcdef0123456789abcdef")
	dl := time.Unix(1_900_000_000, 0).UTC()
	if err := l.Put(LockRecord{Hash: h, LinearID: "nft-3", Sender: []byte("v"), Recipient: []byte("b"), Deadline: dl, TxID: "tx"}); err != nil {
		t.Fatal(err)
	}
	if err := e.Put(Expectation{Hash: h, LinearID: "nft-3", Deadline: dl, Tolerance: 30 * time.Second}); err != nil {
		t.Fatal(err)
	}
	_ = st.Close()

	// "reinício": registros novos, mesmo banco
	st2, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer st2.Close()
	g2, l2, e2 := freshRegistries()
	for _, f := range []func(context.Context, KV) error{g2.attach, l2.attach, e2.attach} {
		if err := f(ctx, st2.KV()); err != nil {
			t.Fatal(err)
		}
	}
	if g2.Check("nft-1") == nil {
		t.Fatal("gravame perdido no reinício")
	}
	if g2.Check("nft-2") != nil {
		t.Fatal("baixa de gravame perdida no reinício")
	}
	if rec, ok := l2.Get(h); !ok || rec.LinearID != "nft-3" || rec.TxID != "tx" || !rec.Deadline.Equal(dl) {
		t.Fatalf("lock auditado perdido: %+v %v", rec, ok)
	}
	if x, ok := e2.Get(h); !ok || x.LinearID != "nft-3" || x.Tolerance != 30*time.Second {
		t.Fatalf("expectativa perdida: %+v %v", x, ok)
	}
}

type failingKV struct{}

func (failingKV) Put(context.Context, string, string, []byte) error { return errors.New("disco cheio") }
func (failingKV) Delete(context.Context, string, string) error      { return errors.New("disco cheio") }
func (failingKV) All(context.Context, string) (map[string][]byte, error) {
	return map[string][]byte{}, nil
}

// Falha fechada: sem persistir, a operação não vale (nem em memória).
func TestDurableFailClosed(t *testing.T) {
	ctx := context.Background()
	g, l, e := freshRegistries()
	_ = g.attach(ctx, failingKV{})
	_ = l.attach(ctx, failingKV{})
	_ = e.attach(ctx, failingKV{})
	if g.Set("nft", "x") == nil || g.Check("nft") != nil {
		t.Fatal("gravame não persistido não pode valer em memória")
	}
	h := []byte("h")
	if l.Put(LockRecord{Hash: h}) == nil {
		t.Fatal("lock")
	}
	if _, ok := l.Get(h); ok {
		t.Fatal("lock não persistido ficou em memória")
	}
	if e.Put(Expectation{Hash: h}) == nil {
		t.Fatal("expectativa")
	}
	if _, ok := e.Get(h); ok {
		t.Fatal("expectativa não persistida ficou em memória")
	}
}
