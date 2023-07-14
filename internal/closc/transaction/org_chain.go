package transaction

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type OrgPlain struct {
	MerkleRoot []byte
}

func NewOrgPlain(merkleRoot []byte) *OrgPlain {
	return &OrgPlain{
		MerkleRoot: merkleRoot,
	}
}

func (o *OrgPlain) ToOnChain() *OrgOnChain {
	return NewOrgOnChain(hex.EncodeToString(o.MerkleRoot))
}

type OrgOnChain struct {
	MerkleRoot string `json:"merkle_root"`
}

func NewOrgOnChain(merkleRoot string) *OrgOnChain {
	return &OrgOnChain{
		MerkleRoot: merkleRoot,
	}
}

func (o *OrgOnChain) ToPlain() (*OrgPlain, error) {
	merkleRoot, err := hex.DecodeString(o.MerkleRoot)
	if err != nil {
		return nil, err
	}
	return NewOrgPlain(merkleRoot), nil
}

func (o *OrgOnChain) KeyVal() (string, []byte, error) {
	txJSON, err := json.Marshal(o)
	if err != nil {
		return "", nil, err
	}
	sha256Func := sha256.New()
	sha256Func.Write(txJSON)
	return hex.EncodeToString(sha256Func.Sum(nil)), txJSON, nil
}
