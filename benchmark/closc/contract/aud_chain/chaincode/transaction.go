package chaincode

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Transaction struct {
	Commitment string `json:"commitment"`
	Hash       string `json:"hash"`
}

func NewTransaction(commitment, hash string) *Transaction {
	return &Transaction{
		Commitment: commitment,
		Hash:       hash,
	}
}

func (t *Transaction) KeyVal() (string, []byte, error) {
	txJSON, err := json.Marshal(t)
	if err != nil {
		return "", nil, err
	}
	sha256Func := sha256.New()
	sha256Func.Write(txJSON)
	return hex.EncodeToString(sha256Func.Sum(nil)), txJSON, nil
}
