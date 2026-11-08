package janus

import (
	"os"
	"path/filepath"
	"testing"

	"janus/node/internal/swap"
)

func writeCfg(t *testing.T, body string) string {
	p := filepath.Join(t.TempDir(), "janus.json")
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoadConfig(t *testing.T) {
	ok := map[string]string{
		"registradora": `{"role":"registradora","mspId":"RegistradoraMSP"}`,
		"banco":        `{"role":"banco","mspId":"BancoMSP","besu":{"rpc":"http://x","keyHex":"00","deployment":"d.json"}}`,
		"vendedora":    `{"role":"vendedora","mspId":"VendedoraMSP","counterparty":"banco","besu":{"rpc":"http://x","keyHex":"00","deployment":"d.json"},"policy":"prod"}`,
	}
	for name, body := range ok {
		c, err := LoadConfig(writeCfg(t, body))
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if c.DataDir == "" {
			t.Fatalf("%s: dataDir default", name)
		}
	}
	c, _ := LoadConfig(writeCfg(t, ok["vendedora"]))
	if c.policy().T1Offset != swap.ProductionPolicy().T1Offset {
		t.Fatal("prod policy")
	}
	bad := map[string]string{
		"role":         `{"role":"x"}`,
		"besu":         `{"role":"banco"}`,
		"counterparty": `{"role":"vendedora","besu":{"rpc":"x"}}`,
		"json":         `{`,
	}
	for name, body := range bad {
		if _, err := LoadConfig(writeCfg(t, body)); err == nil {
			t.Fatalf("%s: expected error", name)
		}
	}
	if trim0x("0xab") != "ab" || trim0x("ab") != "ab" {
		t.Fatal("trim0x")
	}
}
