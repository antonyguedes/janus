package it

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hyperledger-labs/fabric-smart-client/integration"
	"github.com/hyperledger-labs/fabric-smart-client/integration/nwo/common"
	"github.com/onsi/gomega"

	"janus/node/internal/hashlock"
	jt "janus/node/internal/token"
)

// Requer: FAB_BINS com os binários do Fabric 3.1.x (make fabric-bins) e Docker.
// Executar: make it   (ou: cd node/topology && FAB_BINS=... go test ./it -run . -timeout 30m)

func startNetwork(t *testing.T, o Options) *integration.Infrastructure {
	t.Helper()
	if os.Getenv("FAB_BINS") == "" {
		t.Skip("FAB_BINS não definido: teste de integração Fabric pulado")
	}
	gomega.RegisterTestingT(t)
	n, err := StartNetwork(o)
	t.Cleanup(n.Stop) // derruba a rede mesmo se o Setup falhar no meio
	if err != nil {
		t.Fatal(err)
	}
	return n.II
}

func call(t *testing.T, ii *integration.Infrastructure, node, fn string, in any) ([]byte, error) {
	t.Helper()
	res, err := ii.Client(node).CallView(fn, common.JSONMarshall(in))
	if err != nil {
		return nil, err
	}
	switch v := res.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	}
	return json.Marshal(res)
}

// ok chama a view e falha o teste em caso de erro.
func ok(t *testing.T, ii *integration.Infrastructure, node, fn string, in any) []byte {
	t.Helper()
	b, err := call(t, ii, node, fn, in)
	if err != nil {
		t.Fatalf("%s.%s: %v", node, fn, err)
	}
	return b
}

func status(t *testing.T, ii *integration.Infrastructure, node string, h hashlock.Hash) jt.LockStatus {
	t.Helper()
	raw := ok(t, ii, node, "status", jt.Status{Hash: h[:]})
	var s jt.LockStatus
	if err := json.Unmarshal(unquote(raw), &s); err != nil {
		t.Fatalf("status: %v (%s)", err, raw)
	}
	return s
}

// CallView devolve []byte serializado como JSON (base64) em alguns caminhos. Normaliza.
func unquote(b []byte) []byte {
	var inner []byte
	if json.Unmarshal(b, &inner) == nil && len(inner) > 0 && inner[0] == '{' {
		return inner
	}
	var s string
	if json.Unmarshal(b, &s) == nil {
		return []byte(s)
	}
	return b
}

// brief resume um LockStatus sem despejar as identidades (certificados) inteiras.
func brief(st jt.LockStatus) string {
	return fmt.Sprintf("state=%s asset=%s deadline=%s hashOk=%v preimage=%d bytes",
		st.State, st.LinearID, st.Deadline.Format(time.RFC3339), st.HashOK, len(st.Preimage))
}

func issue(t *testing.T, ii *integration.Infrastructure, registro string) string {
	t.Helper()
	raw := ok(t, ii, jt.AuditorNode, "issue", jt.Issue{
		Recipient: "vendedora",
		Recebivel: jt.Recebivel{Registro: registro, Cedente: "Vendedora S.A.", Sacado: "Sacado Ltda",
			ValorFace: 1_250_000_00, Vencimento: "2027-03-31"},
	})
	return strings.Trim(string(unquote(raw)), `"`)
}

func TestFabricHTLC(t *testing.T) {
	ii := startNetwork(t, Options{})
	time.Sleep(10 * time.Second) // nós carregam os parâmetros públicos
	ok(t, ii, jt.AuditorNode, "registerAuditor", nil)

	t.Run("emissão, lock, claim com S", func(t *testing.T) {
		id := issue(t, ii, "REG-0001")
		ok(t, ii, "vendedora", "get", jt.Get{LinearID: id})

		s, h, _ := hashlock.New()
		T1 := time.Now().Add(3 * time.Minute).Truncate(time.Second)
		ok(t, ii, "banco", "expect", jt.Expectation{Hash: h[:], LinearID: id, Deadline: T1, Tolerance: 30 * time.Second})
		ok(t, ii, "vendedora", "lock", jt.Lock{LinearID: id, Recipient: "banco", Hash: h[:], Deadline: T1})

		st := status(t, ii, "banco", h)
		if st.State != "locked" || !st.HashOK || st.LinearID != id || st.Deadline.Sub(T1).Abs() > 30*time.Second {
			t.Fatalf("status banco: %s", brief(st))
		}
		// preimagem errada é rejeitada
		wrong := sha256.Sum256([]byte("x"))
		if _, err := call(t, ii, "banco", "claim", jt.Claim{Preimage: wrong[:]}); err == nil {
			t.Fatal("claim com preimagem errada deveria falhar")
		}
		ok(t, ii, "banco", "claim", jt.Claim{Preimage: s[:]})
		ok(t, ii, "banco", "get", jt.Get{LinearID: id})
		if _, err := call(t, ii, "vendedora", "get", jt.Get{LinearID: id}); err == nil {
			t.Fatal("o NFT deveria ter saído da Vendedora")
		}
		if st := status(t, ii, "vendedora", h); st.State != "claimed" || string(st.Preimage) != string(s[:]) {
			t.Fatalf("status vendedora: %s", brief(st))
		}
	})

	t.Run("lock sem termos assinados é recusado pelo Banco", func(t *testing.T) {
		id := issue(t, ii, "REG-0002")
		_, h, _ := hashlock.New()
		_, err := call(t, ii, "vendedora", "lock", jt.Lock{LinearID: id, Recipient: "banco", Hash: h[:], Deadline: time.Now().Add(3 * time.Minute)})
		if err == nil {
			t.Fatal("lock sem expectativa deveria falhar")
		}
		ok(t, ii, "vendedora", "get", jt.Get{LinearID: id}) // continua com a Vendedora
	})

	t.Run("gravame bloqueia o lock (auditoria da Registradora)", func(t *testing.T) {
		id := issue(t, ii, "REG-0003")
		_, h, _ := hashlock.New()
		T1 := time.Now().Add(3 * time.Minute)
		ok(t, ii, "banco", "expect", jt.Expectation{Hash: h[:], LinearID: id, Deadline: T1, Tolerance: 30 * time.Second})
		ok(t, ii, jt.AuditorNode, "gravame", jt.Gravame{LinearID: id, Motivo: "penhor judicial"})
		if _, err := call(t, ii, "vendedora", "lock", jt.Lock{LinearID: id, Recipient: "banco", Hash: h[:], Deadline: T1}); err == nil {
			t.Fatal("lock de recebível com gravame deveria falhar")
		}
		ok(t, ii, jt.AuditorNode, "gravame", jt.Gravame{LinearID: id, Clear: true})
	})

	t.Run("reclaim só depois de T1 + tolerância do auditor", func(t *testing.T) {
		id := issue(t, ii, "REG-0004")
		_, h, _ := hashlock.New()
		T1 := time.Now().Add(70 * time.Second)
		ok(t, ii, "banco", "expect", jt.Expectation{Hash: h[:], LinearID: id, Deadline: T1, Tolerance: 30 * time.Second})
		ok(t, ii, "vendedora", "lock", jt.Lock{LinearID: id, Recipient: "banco", Hash: h[:], Deadline: T1})

		if _, err := call(t, ii, "vendedora", "reclaim", jt.Reclaim{Hash: h[:]}); err == nil {
			t.Fatal("reclaim antes de T1 deveria falhar")
		}
		time.Sleep(time.Until(T1) + 2*time.Second)
		if _, err := call(t, ii, "vendedora", "reclaim", jt.Reclaim{Hash: h[:]}); err == nil {
			t.Fatal("reclaim dentro da tolerância do auditor deveria falhar (ADR-0004 camada 2)")
		}
		// regressão: um lock expirado (depois de T1) ainda está travado, não "reclaimed"
		if st := status(t, ii, "vendedora", h); st.State != "locked" {
			t.Fatalf("lock expirado deveria continuar locked: %s", brief(st))
		}
		time.Sleep(time.Until(T1.Add(jt.DefaultAudit.ReclaimGrace)) + 2*time.Second)
		ok(t, ii, "vendedora", "reclaim", jt.Reclaim{Hash: h[:]})
		ok(t, ii, "vendedora", "get", jt.Get{LinearID: id})
		if st := status(t, ii, "vendedora", h); st.State != "reclaimed" {
			t.Fatalf("status: %s", brief(st))
		}
	})
}
