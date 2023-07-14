package transaction

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"go.dedis.ch/kyber/v3"
)

type AudPlain struct {
	Commitment []byte
	Hash       []byte
}

func NewAudPlain(commitment, hash []byte) *AudPlain {
	return &AudPlain{
		Commitment: commitment,
		Hash:       hash,
	}
}

func NewAudPlainFromPoint(commitment kyber.Point, hash []byte) (*AudPlain, error) {
	commitmentBytes, err := commitment.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return NewAudPlain(commitmentBytes, hash), nil
}

func (a *AudPlain) ToOnChain() *AudOnChain {
	return NewAudOnChain(hex.EncodeToString(a.Commitment), hex.EncodeToString(a.Hash))
}

type AudOnChain struct {
	Commitment string `json:"commitment"`
	Hash       string `json:"hash"`
}

func NewAudOnChain(commitment, hash string) *AudOnChain {
	return &AudOnChain{
		Commitment: commitment,
		Hash:       hash,
	}
}

func (a *AudOnChain) ToPlain() (*AudPlain, error) {
	commitment, err := hex.DecodeString(a.Commitment)
	if err != nil {
		return nil, err
	}
	hash, err := hex.DecodeString(a.Hash)
	if err != nil {
		return nil, err
	}
	return NewAudPlain(commitment, hash), nil
}

func (a *AudOnChain) KeyVal() (string, []byte, error) {
	txJSON, err := json.Marshal(a)
	if err != nil {
		return "", nil, err
	}
	sha256Func := sha256.New()
	sha256Func.Write(txJSON)
	return hex.EncodeToString(sha256Func.Sum(nil)), txJSON, nil
}
