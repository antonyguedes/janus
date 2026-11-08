// Package secrets guarda as preimagens S cifradas em repouso (AES-256-GCM).
// A Vendedora gera S e o guarda aqui antes de qualquer lock. O Banco grava aqui o S que aprende
// no Besu (watchtower), para sobreviver a quedas. O swapId é o dado associado (AAD) do GCM: um
// arquivo copiado para outro swap não decifra.
package secrets

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// KeySize é o tamanho da chave mestra (AES-256).
const KeySize = 32

// FileStore guarda um arquivo por swap em Dir.
type FileStore struct {
	Dir  string
	aead cipher.AEAD
}

// NewFileStore cria o diretório (0700) e prepara o AEAD com a chave mestra.
func NewFileStore(dir string, key []byte) (*FileStore, error) {
	if len(key) != KeySize {
		return nil, fmt.Errorf("secrets: key must be %d bytes", KeySize)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	return &FileStore{Dir: dir, aead: aead}, nil
}

// LoadOrCreateKey lê a chave mestra de path ou cria uma nova (0600). Em produção, a chave
// viria de um KMS/HSM (fora do escopo do MVP).
func LoadOrCreateKey(path string) ([]byte, error) {
	raw, err := os.ReadFile(path)
	if err == nil {
		k, err := hex.DecodeString(string(raw))
		if err != nil || len(k) != KeySize {
			return nil, fmt.Errorf("secrets: invalid key file %s", path)
		}
		return k, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	k := make([]byte, KeySize)
	if _, err := rand.Read(k); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, err
	}
	return k, os.WriteFile(path, []byte(hex.EncodeToString(k)), 0o600)
}

func (s *FileStore) path(id [32]byte) string {
	return filepath.Join(s.Dir, hex.EncodeToString(id[:])+".s")
}

// Put cifra e grava S de forma atômica. Gravar o mesmo S de novo é idempotente.
func (s *FileStore) Put(_ context.Context, id, preimage [32]byte) error {
	nonce := make([]byte, s.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return err
	}
	ct := s.aead.Seal(nonce, nonce, preimage[:], id[:])
	tmp := s.path(id) + ".tmp"
	if err := os.WriteFile(tmp, ct, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path(id))
}

// Get devolve nil se não houver S para o swap.
func (s *FileStore) Get(_ context.Context, id [32]byte) (*[32]byte, error) {
	ct, err := os.ReadFile(s.path(id))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	n := s.aead.NonceSize()
	if len(ct) < n {
		return nil, errors.New("secrets: corrupted file")
	}
	pt, err := s.aead.Open(nil, ct[:n], ct[n:], id[:])
	if err != nil {
		return nil, fmt.Errorf("secrets: decrypt: %w", err)
	}
	if len(pt) != 32 {
		return nil, errors.New("secrets: invalid preimage length")
	}
	var out [32]byte
	copy(out[:], pt)
	return &out, nil
}
