package auditor

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"go.dedis.ch/kyber/v3"

	clolcorg "github.com/ac83ae/auti/internal/clolc/organization"
	"github.com/ac83ae/auti/internal/clolc/transaction"
	"github.com/ac83ae/auti/internal/constants"
	"github.com/ac83ae/auti/internal/crypto"
)

var sha256Func = sha256.New()

type TypeID string
type TypeEpochID []byte

type Auditor struct {
	ID                   TypeID
	AuditedOrgIDs        []clolcorg.TypeID
	epochTXRandMap       map[[2]string][]kyber.Scalar
	EpochID              TypeEpochID
	epochOrgSecretKeyMap map[string]crypto.TypePrivateKey
	epochOrgIDMap        map[clolcorg.TypeID]clolcorg.TypeEpochID
}

func New(id string, organizations []*clolcorg.Organization) *Auditor {
	aud := &Auditor{
		ID: TypeID(id),
	}
	aud.AuditedOrgIDs = make([]clolcorg.TypeID, len(organizations))
	for idx, org := range organizations {
		aud.AuditedOrgIDs[idx] = org.ID
	}
	return aud
}

func (a *Auditor) SetEpochTXRandomness(txRandMap map[[2]string][]kyber.Scalar) {
	a.epochTXRandMap = txRandMap
}

func (a *Auditor) GetEpochTXRandomness(orgID1, orgID2 clolcorg.TypeID) []kyber.Scalar {
	key := clolcorg.IDHashKey(clolcorg.IDHashString(orgID1), clolcorg.IDHashString(orgID2))
	if txRand, ok := a.epochTXRandMap[key]; ok {
		return txRand
	}
	return nil
}

func (a *Auditor) SetEpochSecretKey(orgSecretKeyMap map[string]crypto.TypePrivateKey) {
	a.epochOrgSecretKeyMap = orgSecretKeyMap
}

func (a *Auditor) SetEpochID(id []byte) {
	a.EpochID = id
}

func (a *Auditor) SetEpochOrgIDMap(idMap map[clolcorg.TypeID]clolcorg.TypeEpochID) {
	a.epochOrgIDMap = idMap
}

func (a *Auditor) AccumulateCommitments(
	orgID clolcorg.TypeID, txList []*transaction.LocalHidden,
) (kyber.Point, error) {
	if len(txList) == 0 {
		return nil, fmt.Errorf("empty transaction list")
	}
	if constants.MaxNumTXInEpoch < len(txList) {
		return nil, fmt.Errorf("too many transactions in the epoch: %d", len(txList))
	}
	orgIDHashStr := clolcorg.IDHashString(orgID)
	counterPartyIDHashStr := hex.EncodeToString(txList[0].CounterParty)
	orgKey := clolcorg.IDHashKey(orgIDHashStr, counterPartyIDHashStr)
	randomScalars := a.epochTXRandMap[orgKey]
	result := crypto.KyberSuite.Point().Null()
	for idx, tx := range txList {
		commitmentBytes := tx.Commitment
		commitmentPoint := crypto.KyberSuite.Point()
		if err := commitmentPoint.UnmarshalBinary(commitmentBytes); err != nil {
			return nil, err
		}
		commitmentPoint.Mul(randomScalars[idx], commitmentPoint)
		result.Add(result, commitmentPoint)
	}
	return result, nil
}

func (a *Auditor) ComputeB(orgTXRandList, comTXRandList []kyber.Scalar) (kyber.Point, error) {
	if len(orgTXRandList) != len(comTXRandList) {
		return nil, fmt.Errorf("length of two lists are not equal")
	}
	scalar := crypto.KyberSuite.Scalar().Zero()
	for idx := range orgTXRandList {
		tmp := crypto.KyberSuite.Scalar().Mul(orgTXRandList[idx], comTXRandList[idx])
		scalar.Sub(scalar, tmp)
	}
	result := crypto.KyberSuite.Point().Mul(scalar, crypto.PointH)
	return result, nil
}

func (a *Auditor) ComputeC(res, A kyber.Point) kyber.Point {
	result := crypto.KyberSuite.Point().Sub(A, res)
	return result
}

func (a *Auditor) ComputeD(pointA, pointB kyber.Point) kyber.Point {
	result := crypto.KyberSuite.Point().Add(pointA, pointB)
	result.Neg(result)
	return result
}

func (a *Auditor) EncryptConsistencyExamResult(
	orgID clolcorg.TypeID, counterPartyIDHash string,
	res, pointB, pointC, pointD kyber.Point, publicKey kyber.Point,
) (*transaction.AudPlain, error) {
	txID, err := a.ComputeCETransactionID(orgID, counterPartyIDHash)
	if err != nil {
		return nil, err
	}
	cipherRes, err := crypto.EncryptPoint(publicKey, res)
	if err != nil {
		return nil, err
	}
	cipherResBytes, err := cipherRes.Serialize()
	if err != nil {
		return nil, err
	}
	cipherB, err := crypto.EncryptPoint(publicKey, pointB)
	if err != nil {
		return nil, err
	}
	cipherBBytes, err := cipherB.Serialize()
	if err != nil {
		return nil, err
	}
	cipherC, err := crypto.EncryptPoint(publicKey, pointC)
	if err != nil {
		return nil, err
	}
	cipherCBytes, err := cipherC.Serialize()
	if err != nil {
		return nil, err
	}
	epochIDHashPoint := EpochIDHashPoint(a.EpochID)
	idPointD := crypto.KyberSuite.Point().Add(epochIDHashPoint, pointD)
	cipherD, err := crypto.EncryptPoint(publicKey, idPointD)
	if err != nil {
		return nil, err
	}
	cipherDBytes, err := cipherD.Serialize()
	if err != nil {
		return nil, err
	}
	return transaction.NewAudPlain(
		txID, cipherResBytes, cipherBBytes, cipherCBytes, cipherDBytes,
	), nil
}

func (a *Auditor) ComputeCETransactionID(
	orgID clolcorg.TypeID, counterPartyIDHash string,
) ([]byte, error) {
	orgIDHashStr := clolcorg.IDHashString(orgID)
	orgKey := clolcorg.IDHashKey(orgIDHashStr, counterPartyIDHash)
	randomnesses := a.epochTXRandMap[orgKey]
	epochOrgID := a.epochOrgIDMap[orgID]
	epochOrgIDBytes := make([]byte, len(epochOrgID))
	copy(epochOrgIDBytes, epochOrgID)
	randAccumulator := crypto.KyberSuite.Scalar().Zero()
	for _, randScalar := range randomnesses {
		randAccumulator.Add(randAccumulator, randScalar)
	}
	randAccumulatorBytes, err := randAccumulator.MarshalBinary()
	if err != nil {
		return nil, err
	}
	concatBytes := append(epochOrgIDBytes, randAccumulatorBytes...)
	sha256Func := sha256.New()
	sha256Func.Write(concatBytes)
	result := sha256Func.Sum(nil)
	return result, nil
}

func (a *Auditor) DecryptResAndB(orgIDHash string,
	tx *transaction.AudOnChain) (kyber.Point, kyber.Point, error) {
	plainTX, err := tx.ToPlain()
	if err != nil {
		return nil, nil, err
	}
	privateKey, ok := a.epochOrgSecretKeyMap[orgIDHash]
	if !ok {
		return nil, nil, fmt.Errorf("no private key for organization %s", orgIDHash)
	}
	res, err := crypto.DecryptPoint(privateKey, plainTX.CipherRes)
	if err != nil {
		return nil, nil, err
	}
	pointB, err := crypto.DecryptPoint(privateKey, plainTX.CipherB)
	if err != nil {
		return nil, nil, err
	}
	return res, pointB, nil
}

func (a *Auditor) CheckResultConsistency(res, B, txRes, txB kyber.Point) bool {
	result := crypto.KyberSuite.Point().Null()
	result.Add(result, res)
	result.Add(result, B)
	result.Add(result, txRes)
	result.Add(result, txB)
	return result.Equal(crypto.KyberSuite.Point().Null())
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
