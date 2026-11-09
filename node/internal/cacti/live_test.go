package cacti

import (
	"context"
	"os"
	"testing"
)

// TestLiveTxProof roda contra um servidor Cacti real (JANUS_CACTI_URL + JANUS_CACTI_TX).
func TestLiveTxProof(t *testing.T) {
	url, tx := os.Getenv("JANUS_CACTI_URL"), os.Getenv("JANUS_CACTI_TX")
	if url == "" || tx == "" {
		t.Skip("JANUS_CACTI_URL/JANUS_CACTI_TX não definidos")
	}
	c := &FabricConnector{BaseURL: url, KeychainID: "janus", KeychainRef: os.Getenv("JANUS_CACTI_REF"), Channel: "testchannel"}
	p, err := c.TxProof(context.Background(), tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", p)
	if err := p.Satisfies("tns", []string{"RegistradoraMSP", "BancoMSP"}); err != nil {
		t.Fatal(err)
	}
}
