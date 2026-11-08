// Package token contém as views do FSC/Panurus do Janus: emissão do recebível (NFT), consulta,
// lock/claim/reclaim HTLC e a auditoria da Registradora.
//
// No Panurus, um NFT é um token de quantidade 1 cujo Type é base64(json(estado)) (nfttx). Por isso o
// lock HTLC de um recebível é um htlc.Lock comum desse Type com valor 1, sem fork e sem código no
// chaincode.
package token

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/LFDT-Panurus/panurus/token/services/nfttx/marshaller"
	token2 "github.com/LFDT-Panurus/panurus/token/token"
)

// Recebivel é o estado do NFT emitido pela Registradora. Ele nunca sai do Fabric.
type Recebivel struct {
	LinearID   string `json:"LinearID"`   // atribuído na emissão (UUID)
	Registro   string `json:"Registro"`   // número do registro na Registradora (único)
	Cedente    string `json:"Cedente"`    // titular original
	Sacado     string `json:"Sacado"`     // devedor
	ValorFace  uint64 `json:"ValorFace"`  // centavos
	Vencimento string `json:"Vencimento"` // AAAA-MM-DD
}

// SetLinearID implementa nfttx.LinearState.
func (r *Recebivel) SetLinearID(id string) string {
	if r.LinearID == "" {
		r.LinearID = id
	}
	return r.LinearID
}

// Validate confere os campos obrigatórios.
func (r *Recebivel) Validate() error {
	switch {
	case r.Registro == "":
		return errors.New("recebivel: registro obrigatório")
	case r.Cedente == "" || r.Sacado == "":
		return errors.New("recebivel: cedente e sacado obrigatórios")
	case r.ValorFace == 0:
		return errors.New("recebivel: valor de face obrigatório")
	case len(r.Vencimento) != len("2006-01-02"):
		return errors.New("recebivel: vencimento AAAA-MM-DD")
	}
	return nil
}

// TokenType devolve o Type do Panurus que representa este NFT, idêntico ao gerado por nfttx.
func TokenType(r *Recebivel) (token2.Type, error) {
	raw, err := marshaller.Marshal(r)
	if err != nil {
		return "", fmt.Errorf("recebivel: marshal: %w", err)
	}
	return token2.Type(base64.StdEncoding.EncodeToString(raw)), nil
}

// FromTokenType decodifica o estado de um Type de NFT. Devolve erro se o Type não for um recebível.
func FromTokenType(t token2.Type) (*Recebivel, error) {
	raw, err := base64.StdEncoding.DecodeString(string(t))
	if err != nil {
		return nil, fmt.Errorf("recebivel: not an nft type: %w", err)
	}
	r := &Recebivel{}
	if err := marshaller.Unmarshal(raw, r); err != nil {
		return nil, fmt.Errorf("recebivel: unmarshal: %w", err)
	}
	if r.LinearID == "" {
		return nil, errors.New("recebivel: missing LinearID")
	}
	return r, nil
}
