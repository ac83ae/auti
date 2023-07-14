package transaction

import (
	"encoding/hex"
	"encoding/json"
)

type AudPlain struct {
	ID        []byte
	CipherRes []byte
	CipherB   []byte
	CipherC   []byte
	CipherD   []byte
}

func NewAudPlain(id, cipherRes, cipherB, cipherC, cipherD []byte) *AudPlain {
	return &AudPlain{
		ID:        id,
		CipherRes: cipherRes,
		CipherB:   cipherB,
		CipherC:   cipherC,
		CipherD:   cipherD,
	}
}

func (a *AudPlain) ToOnChain() *AudOnChain {
	return &AudOnChain{
		ID:        hex.EncodeToString(a.ID),
		CipherRes: hex.EncodeToString(a.CipherRes),
		CipherB:   hex.EncodeToString(a.CipherB),
		CipherC:   hex.EncodeToString(a.CipherC),
		CipherD:   hex.EncodeToString(a.CipherD),
	}
}

type AudOnChain struct {
	ID        string
	CipherRes string `json:"cipher_res"`
	CipherB   string `json:"cipher_b"`
	CipherC   string `json:"cipher_c"`
	CipherD   string `json:"cipher_d"`
}

func NewAudOnChain(id, cipherRes, cipherB, cipherC, cipherD string) *AudOnChain {
	return &AudOnChain{
		ID:        id,
		CipherRes: cipherRes,
		CipherB:   cipherB,
		CipherC:   cipherC,
		CipherD:   cipherD,
	}
}

func (a *AudOnChain) ToPlain() (*AudPlain, error) {
	id, err := hex.DecodeString(a.ID)
	if err != nil {
		return nil, err
	}
	cipherRes, err := hex.DecodeString(a.CipherRes)
	if err != nil {
		return nil, err
	}
	cipherB, err := hex.DecodeString(a.CipherB)
	if err != nil {
		return nil, err
	}
	cipherC, err := hex.DecodeString(a.CipherC)
	if err != nil {
		return nil, err
	}
	cipherD, err := hex.DecodeString(a.CipherD)
	if err != nil {
		return nil, err
	}
	return NewAudPlain(id, cipherRes, cipherB, cipherC, cipherD), nil
}

func (a *AudOnChain) KeyVal() (string, []byte, error) {
	txJSON, err := json.Marshal(a)
	if err != nil {
		return "", nil, err
	}
	return a.ID, txJSON, nil
}
