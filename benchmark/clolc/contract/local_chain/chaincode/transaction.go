package chaincode

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Transaction struct {
	CounterParty string `json:"counter_party"`
	Commitment   string `json:"commitment"`
	Timestamp    string `json:"timestamp"`
}

func NewTransaction(counterParty, commitment, timestamp string) *Transaction {
	return &Transaction{
		CounterParty: counterParty,
		Commitment:   commitment,
		Timestamp:    timestamp,
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
