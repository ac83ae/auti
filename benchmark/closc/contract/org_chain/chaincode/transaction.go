package chaincode

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Transaction struct {
	MerkleRoot string `json:"merkle_root"`
}

func NewTransaction(merkleRoot string) *Transaction {
	return &Transaction{
		MerkleRoot: merkleRoot,
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
