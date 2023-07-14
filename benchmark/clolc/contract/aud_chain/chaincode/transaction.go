package chaincode

import "encoding/json"

type Transaction struct {
	ID        string
	CipherRes string `json:"cipher_res"`
	CipherB   string `json:"cipher_b"`
	CipherC   string `json:"cipher_c"`
	CipherD   string `json:"cipher_d"`
}

func NewTransaction(id, cipherRes, cipherB, cipherC, cipherD string) *Transaction {
	return &Transaction{
		ID:        id,
		CipherRes: cipherRes,
		CipherB:   cipherB,
		CipherC:   cipherC,
		CipherD:   cipherD,
	}
}

func (t *Transaction) KeyVal() (string, []byte, error) {
	txJSON, err := json.Marshal(t)
	if err != nil {
		return "", nil, err
	}
	return t.ID, txJSON, nil
}
