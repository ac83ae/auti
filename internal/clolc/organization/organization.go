package organization

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"go.dedis.ch/kyber/v3"

	"github.com/ac83ae/auti/internal/clolc/transaction"
	"github.com/ac83ae/auti/internal/crypto"
)

var sha256Func = sha256.New()

type TypeID string
type TypeEpochID []byte

type Organization struct {
	ID                  TypeID
	IDHash              string
	EpochID             TypeEpochID
	epochAccumulatorMap map[[2]string]kyber.Point
	epochTXRandomness   map[[2]string]kyber.Scalar
}

func New(id string) *Organization {
	defer sha256Func.Reset()
	sha256Func.Write([]byte(id))
	idHash := hex.EncodeToString(sha256Func.Sum(nil))
	org := &Organization{
		ID:                  TypeID(id),
		IDHash:              idHash,
		epochAccumulatorMap: make(map[[2]string]kyber.Point),
		epochTXRandomness:   make(map[[2]string]kyber.Scalar),
	}
	return org
}

func (c *Organization) SetEpochID(randID []byte) {
	c.EpochID = randID
}

func (c *Organization) RecordTransaction(tx *transaction.LocalPlain) error {
	// Submit the transaction to the local chain
	sha256Func := sha256.New()
	sha256Func.Write([]byte(tx.CounterParty))
	counterPartyHash := sha256Func.Sum(nil)
	commitment, randScalar, err := crypto.PedersenCommit(tx.Amount)
	if err != nil {
		return err
	}
	commitmentBytes, err := commitment.MarshalBinary()
	if err != nil {
		return err
	}
	clolcHidden := &transaction.LocalHidden{
		CounterParty: counterPartyHash,
		Commitment:   commitmentBytes,
		Timestamp:    tx.Timestamp,
	}
	if err = c.SubmitTXLocalChain(clolcHidden); err != nil {
		return err
	}
	counterPartyHashStr := hex.EncodeToString(counterPartyHash)
	orgMapKey := IDHashKey(c.IDHash, counterPartyHashStr)
	// Accumulate the commitment to the corresponding accumulator
	if _, ok := c.epochAccumulatorMap[orgMapKey]; !ok {
		c.epochAccumulatorMap[orgMapKey] = commitment
	} else {
		c.epochAccumulatorMap[orgMapKey].Add(
			c.epochAccumulatorMap[orgMapKey],
			commitment,
		)
	}
	// Record the randomness used in the commitment
	c.epochTXRandomness[orgMapKey] = randScalar
	return nil
}

func (c *Organization) Accumulate(counterParty TypeID, commitment kyber.Point) {
	counterPartyHashStr := IDHashString(counterParty)
	orgMapKey := IDHashKey(c.IDHash, counterPartyHashStr)
	// Accumulate the commitment to the corresponding accumulator
	if _, ok := c.epochAccumulatorMap[orgMapKey]; !ok {
		c.epochAccumulatorMap[orgMapKey] = commitment
	} else {
		c.epochAccumulatorMap[orgMapKey].Add(
			c.epochAccumulatorMap[orgMapKey],
			commitment,
		)
	}
}

func (c *Organization) SubmitTXLocalChain(tx *transaction.LocalHidden) error {
	panic("not implemented")
}

func (c *Organization) ComposeTXOrgChain(counterParty TypeID) (*transaction.OrgPlain, error) {
	counterPartyHashStr := IDHashString(counterParty)
	orgMapKey := IDHashKey(c.IDHash, counterPartyHashStr)
	accumulator, ok := c.epochAccumulatorMap[orgMapKey]
	if !ok {
		return nil, fmt.Errorf("no transaction from %s to %s", c.ID, counterParty)
	}
	epochIDHashPoint := EpochIDHashPoint(c.EpochID)
	resultPoint := crypto.KyberSuite.Point().Add(accumulator, epochIDHashPoint)
	result, err := resultPoint.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return transaction.NewOrgPlain(result), nil
}

func IDHashBytes(id TypeID) []byte {
	defer sha256Func.Reset()
	sha256Func.Write([]byte(id))
	return sha256Func.Sum(nil)
}

func IDHashString(id TypeID) string {
	return hex.EncodeToString(IDHashBytes(id))
}

func IDHashScalar(id TypeID) kyber.Scalar {
	return crypto.KyberSuite.Scalar().SetBytes(IDHashBytes(id))
}

func IDHashPoint(id TypeID) kyber.Point {
	return crypto.KyberSuite.Point().Mul(IDHashScalar(id), nil)
}

func IDHashKey(orgIDHash1, orgIDHash2 string) [2]string {
	if orgIDHash1 < orgIDHash2 {
		return [2]string{orgIDHash1, orgIDHash2}
	}
	return [2]string{orgIDHash2, orgIDHash1}
}

func EpochIDHashBytes(epochID TypeEpochID) []byte {
	defer sha256Func.Reset()
	sha256Func.Write(epochID)
	return sha256Func.Sum(nil)
}

func EpochIDHashString(epochID TypeEpochID) string {
	return hex.EncodeToString(EpochIDHashBytes(epochID))
}

func EpochIDHashScalar(epochID TypeEpochID) kyber.Scalar {
	return crypto.KyberSuite.Scalar().SetBytes(EpochIDHashBytes(epochID))
}

func EpochIDHashPoint(epochID TypeEpochID) kyber.Point {
	return crypto.KyberSuite.Point().Mul(EpochIDHashScalar(epochID), nil)
}
