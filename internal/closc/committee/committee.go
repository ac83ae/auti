package committee

import (
	mt "github.com/txaty/go-merkletree"
	"go.dedis.ch/kyber/v3"

	"github.com/ac83ae/auti/internal/closc/auditor"
	closcorg "github.com/ac83ae/auti/internal/closc/organization"
	"github.com/ac83ae/auti/internal/crypto"
)

type TypeID string

type Committee struct {
	ID                TypeID
	managedEntityMap  map[auditor.TypeID][]closcorg.TypeID
	managedAuditorIDs []auditor.TypeID
	managedOrgIDs     []closcorg.TypeID
	epochAuditorIDMap map[auditor.TypeID]auditor.TypeEpochID
}

func New(id string, auditors []*auditor.Auditor) *Committee {
	com := &Committee{
		ID:               TypeID(id),
		managedEntityMap: make(map[auditor.TypeID][]closcorg.TypeID),
	}
	com.managedAuditorIDs = make([]auditor.TypeID, len(auditors))
	for idx, aud := range auditors {
		com.managedEntityMap[aud.ID] = aud.AuditedOrgIDs
		com.managedAuditorIDs[idx] = aud.ID
		com.managedOrgIDs = append(com.managedOrgIDs, aud.AuditedOrgIDs...)
	}
	return com
}

func (c *Committee) reinitializeMaps() {
	c.epochAuditorIDMap = make(map[auditor.TypeID]auditor.TypeEpochID)
}

func GenerateAuditorEpochID() kyber.Point {
	randScalar := crypto.KyberSuite.Scalar().Pick(crypto.KyberSuite.RandomStream())
	randPoint := crypto.KyberSuite.Point().Mul(randScalar, crypto.PointG)
	return randPoint
}

func (c *Committee) InitializeEpoch(auditors []*auditor.Auditor) error {
	c.reinitializeMaps()
	for _, aud := range auditors {
		// Generate epoch ID for each auditor
		epochID := GenerateAuditorEpochID()
		c.epochAuditorIDMap[aud.ID] = epochID
		// Distribute epoch auditor IDs
		aud.SetEpochID(epochID)
	}
	return nil
}

func (c *Committee) VerifyMerkleBatchProof(commitments []mt.DataBlock,
	batchProof *crypto.MerkleBatchProof, merkleRoot []byte) (uint, error) {
	ok, err := crypto.VerifyMerkleBatchProof(commitments, batchProof, merkleRoot)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, nil
	}
	return 1, nil
}

func (c *Committee) SummarizeMerkleBatchProofVerificationResults(verificationResults []uint) bool {
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

func (c *Committee) VerifyCommitment(commitments []kyber.Point) bool {
	sum := crypto.KyberSuite.Point().Null()
	for _, commitment := range commitments {
		sum = crypto.KyberSuite.Point().Add(sum, commitment)
	}
	for _, auditorID := range c.managedAuditorIDs {
		auditorEpochID := c.epochAuditorIDMap[auditorID]
		sum = crypto.KyberSuite.Point().Sub(sum, auditorEpochID)
	}
	// TODO: check this
	return sum.Equal(crypto.KyberSuite.Point().Null())
}
