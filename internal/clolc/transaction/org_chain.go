package transaction

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type OrgPlain struct {
	Accumulator []byte
}

func NewOrgPlain(accumulator []byte) *OrgPlain {
	return &OrgPlain{
		Accumulator: accumulator,
	}
}

func (o *OrgPlain) ToOnChain() *OrgOnChain {
	accumulatorString := hex.EncodeToString(o.Accumulator)
	return NewOrgOnChain(accumulatorString)
}

type OrgOnChain struct {
	Accumulator string `json:"accumulator"`
}

func NewOrgOnChain(accumulator string) *OrgOnChain {
	return &OrgOnChain{
		Accumulator: accumulator,
	}
}

func (o *OrgOnChain) ToPlain() (*OrgPlain, error) {
	accumulatorBytes, err := hex.DecodeString(o.Accumulator)
	if err != nil {
		return nil, err
	}
	return NewOrgPlain(accumulatorBytes), nil
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
