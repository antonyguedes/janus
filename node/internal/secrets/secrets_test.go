package secrets

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestRoundTripAndBinding(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	key, err := LoadOrCreateKey(filepath.Join(dir, "master.key"))
	if err != nil {
		t.Fatal(err)
	}
	again, _ := LoadOrCreateKey(filepath.Join(dir, "master.key"))
	if !bytes.Equal(key, again) {
		t.Fatal("chave deve ser estável")
	}
	st, err := NewFileStore(filepath.Join(dir, "s"), key)
	if err != nil {
		t.Fatal(err)
	}
	id1, id2, s := [32]byte{1}, [32]byte{2}, [32]byte{0xaa, 0xbb}
	if got, _ := st.Get(ctx, id1); got != nil {
		t.Fatal("vazio")
	}
	if err := st.Put(ctx, id1, s); err != nil {
		t.Fatal(err)
	}
	got, err := st.Get(ctx, id1)
	if err != nil || got == nil || *got != s {
		t.Fatalf("roundtrip: %v %v", got, err)
	}
	// em repouso, o S não aparece em claro
	raw, _ := os.ReadFile(st.path(id1))
	if bytes.Contains(raw, s[:2]) && bytes.Contains(raw, s[:]) {
		t.Fatal("preimagem em claro no disco")
	}
	// arquivo copiado para outro swap não decifra (AAD = swapId)
	_ = os.WriteFile(st.path(id2), raw, 0o600)
	if _, err := st.Get(ctx, id2); err == nil {
		t.Fatal("AAD deveria impedir reuso do ciphertext")
	}
	// chave errada não decifra
	other, _ := NewFileStore(st.Dir, bytes.Repeat([]byte{7}, KeySize))
	if _, err := other.Get(ctx, id1); err == nil {
		t.Fatal("chave errada deveria falhar")
	}
	if _, err := NewFileStore(dir, []byte("curta")); err == nil {
		t.Fatal("chave curta")
	}
}
