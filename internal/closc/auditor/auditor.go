package auditor

import (
	mt "github.com/txaty/go-merkletree"
	"go.dedis.ch/kyber/v3"

	"github.com/ac83ae/auti/internal/closc/organization"
	"github.com/ac83ae/auti/internal/closc/transaction"
	"github.com/ac83ae/auti/internal/crypto"
)

type TypeID string

type TypeEpochID kyber.Point

type Auditor struct {
	ID            TypeID
	AuditedOrgIDs []organization.TypeID
	EpochID       TypeEpochID
}

func New(id string, organizations []*organization.Organization) *Auditor {
	aud := &Auditor{
		ID: TypeID(id),
	}
	aud.AuditedOrgIDs = make([]organization.TypeID, len(organizations))
	for idx, org := range organizations {
		aud.AuditedOrgIDs[idx] = org.ID
	}
	return aud
}

func (a *Auditor) SetEpochID(id kyber.Point) {
	a.EpochID = crypto.KyberSuite.Point().Set(id)
}

func (a *Auditor) VerifyMerkleProof(tx transaction.LocalOnChain) (uint, error) {
	txPlain, err := tx.ToPlain()
	if err != nil {
		return 0, err
	}
	merkleProof, err := crypto.MerkleProofUnmarshal(txPlain.MerkleProof)
	if err != nil {
		return 0, err
	}
	ok, err := crypto.VerifyMerkleProof(txPlain, merkleProof, txPlain.MerkleRoot)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, nil
	}
	return 1, nil
}

func (a *Auditor) SummarizeMerkleProofVerificationResults(verificationResults []uint) bool {
	if len(verificationResults) == 0 {
		return false
	}
	result := verificationResults[0]
	for i := 1; i < len(verificationResults); i++ {
		result *= verificationResults[i]
		if result == 0 {
			return false
		}
	}
	return true
}

func (a *Auditor) VerifyCommitments(commitmentList1, commitmentList2 [][]byte,
	hashPoints1, hashPoints2 []kyber.Point) (bool, error) {
	if len(commitmentList1) != len(commitmentList2) {
		return false, nil
	}
	if len(commitmentList1) != len(hashPoints1) {
		return false, nil
	}
	if len(commitmentList2) != len(hashPoints2) {
		return false, nil
	}
	sum := crypto.KyberSuite.Point().Null()
	for i := 0; i < len(commitmentList1); i++ {
		// commitment 1
		commitPoint := crypto.KyberSuite.Point()
		if err := commitPoint.UnmarshalBinary(commitmentList1[i]); err != nil {
			return false, err
		}
		sum = sum.Add(sum, commitPoint)
		// commitment 2
		commitPoint = crypto.KyberSuite.Point()
		if err := commitPoint.UnmarshalBinary(commitmentList2[i]); err != nil {
			return false, err
		}
		sum = sum.Add(sum, commitPoint)
		// hash point 1
		sum = sum.Sub(sum, hashPoints1[i])
		// hash point 2
		sum = sum.Sub(sum, hashPoints2[i])
	}
	// TODO: test this
	return sum.Equal(crypto.KyberSuite.Point().Null()), nil
}

func (a *Auditor) AccumulateCommitments(commitments []kyber.Point) (kyber.Point, error) {
	sum := crypto.KyberSuite.Point().Set(a.EpochID)
	for _, commitment := range commitments {
		sum = sum.Add(sum, commitment)
	}
	return sum, nil
}

func (a *Auditor) MergeProof(commitments []mt.DataBlock, proofs []*mt.Proof) ([]byte, error) {
	batchProof, err := crypto.NewMerkleBatchProof(commitments, proofs)
	if err != nil {
		return nil, err
	}
	return crypto.MerkleBatchProofMarshal(batchProof)
}
