package chaincode

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Transaction struct {
	Commitment  string `json:"commitment"`
	MerkleRoot  string `json:"merkle_root"`
	MerkleProof string `json:"merkle_proof"`
}

func NewTransaction(commitment, merkleRoot, merkleProof string) *Transaction {
	return &Transaction{
		Commitment:  commitment,
		MerkleRoot:  merkleRoot,
		MerkleProof: merkleProof,
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
