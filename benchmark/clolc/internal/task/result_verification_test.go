package task

import (
	"crypto/rand"
	"testing"

	"go.dedis.ch/kyber/v3"

	"github.com/ac83ae/auti/internal/clolc/organization"
	"github.com/ac83ae/auti/internal/clolc/transaction"
	"github.com/ac83ae/auti/internal/crypto"
)

func TestRVOrgAndAudResult(t *testing.T) {
	// entity setup
	com, auditors, organizations := generateEntities(2)
	publicKeyMap, err := com.InitializeEpoch(auditors, organizations)
	if err != nil {
		t.Fatal(err)
	}
	// compose local transactions
	localTXs1, localTXs2 := generateLocalTXPairList(organizations[0].ID, organizations[1].ID)
	// hide the local transactions
	var (
		hiddenTXs1   = make([]*transaction.LocalHidden, testNumTXs)
		hiddenTXs2   = make([]*transaction.LocalHidden, testNumTXs)
		points1      = make([]kyber.Point, testNumTXs)
		points2      = make([]kyber.Point, testNumTXs)
		randScalars1 = make([]kyber.Scalar, testNumTXs)
		randScalars2 = make([]kyber.Scalar, testNumTXs)
	)
	for i := 0; i < testNumTXs; i++ {
		hiddenTX1, point1, scalar1, err := localTXs1[i].Hide()
		if err != nil {
			t.Fatal(err)
		}
		hiddenTXs1[i] = hiddenTX1
		points1[i] = point1
		randScalars1[i] = scalar1
		organizations[0].Accumulate(organizations[1].ID, point1)
		hiddenTX2, point2, scalar2, err := localTXs2[i].Hide()
		if err != nil {
			t.Fatal(err)
		}
		hiddenTXs2[i] = hiddenTX2
		points2[i] = point2
		randScalars2[i] = scalar2
		organizations[1].Accumulate(organizations[0].ID, point2)
	}
	// compose org chain transactions
	orgTX1, err := organizations[0].ComposeTXOrgChain(organizations[1].ID)
	if err != nil {
		t.Fatal(err)
	}
	orgTX2, err := organizations[1].ComposeTXOrgChain(organizations[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	// compute res
	res1, err := auditors[0].AccumulateCommitments(organizations[0].ID, hiddenTXs1)
	if err != nil {
		t.Fatal(err)
	}
	res2, err := auditors[1].AccumulateCommitments(organizations[1].ID, hiddenTXs2)
	if err != nil {
		t.Fatal(err)
	}
	// compute b
	txRandList := auditors[0].GetEpochTXRandomness(organizations[0].ID, organizations[1].ID)
	if txRandList == nil {
		t.Fatal("txRandList is nil")
	}
	b1, err := auditors[0].ComputeB(randScalars1, txRandList)
	if err != nil {
		t.Fatal(err)
	}
	b2, err := auditors[1].ComputeB(randScalars2, txRandList)
	if err != nil {
		t.Fatal(err)
	}
	// compute c
	orgIDPoint1 := organization.EpochIDHashPoint(organizations[0].EpochID)
	orgIDPoint2 := organization.EpochIDHashPoint(organizations[1].EpochID)
	acc1 := crypto.KyberSuite.Point()
	if err := acc1.UnmarshalBinary(orgTX1.Accumulator); err != nil {
		t.Fatal(err)
	}
	acc1.Sub(acc1, orgIDPoint1)
	acc2 := crypto.KyberSuite.Point()
	if err := acc2.UnmarshalBinary(orgTX2.Accumulator); err != nil {
		t.Fatal(err)
	}
	acc2.Sub(acc2, orgIDPoint2)
	c1 := auditors[0].ComputeC(res1, acc1)
	c2 := auditors[1].ComputeC(res2, acc2)
	// compute d
	d1 := auditors[0].ComputeD(acc1, b1)
	d2 := auditors[1].ComputeD(acc2, b2)
	randID1 := make([]byte, 32)
	if _, err = rand.Read(randID1); err != nil {
		t.Fatal(err)
	}
	randID2 := make([]byte, 32)
	if _, err := rand.Read(randID2); err != nil {
		t.Fatal(err)
	}
	idHash1 := organization.IDHashString(organizations[0].ID)
	idHash2 := organization.IDHashString(organizations[1].ID)
	publicKey1 := publicKeyMap[idHash1]
	publicKey2 := publicKeyMap[idHash2]
	audTX1, err := auditors[0].EncryptConsistencyExamResult(organizations[0].ID, idHash2, res1, b1, c1, d1, publicKey1)
	if err != nil {
		t.Fatal(err)
	}
	audTX2, err := auditors[1].EncryptConsistencyExamResult(organizations[1].ID, idHash1, res2, b2, c2, d2, publicKey2)
	if err != nil {
		t.Fatal(err)
	}
	// verify
	res, err := com.VerifyOrgAndAudResult(organizations[0].ID, auditors[0].ID, orgTX1.ToOnChain(), audTX1.ToOnChain())
	if err != nil {
		t.Fatal(err)
	}
	if !res {
		t.Fatal("verify failed")
	}
	res, err = com.VerifyOrgAndAudResult(organizations[1].ID, auditors[1].ID, orgTX2.ToOnChain(), audTX2.ToOnChain())
	if err != nil {
		t.Fatal(err)
	}
	if !res {
		t.Fatal("verify failed")
	}
}

func TestRVVerifyAuditPairResult(t *testing.T) {
	// entity setup
	com, auditors, organizations := generateEntities(2)
	publicKeyMap, err := com.InitializeEpoch(auditors, organizations)
	if err != nil {
		t.Fatal(err)
	}
	// compose local transactions
	localTXs1, localTXs2 := generateLocalTXPairList(organizations[0].ID, organizations[1].ID)
	// hide the local transactions
	var (
		hiddenTXs1   = make([]*transaction.LocalHidden, testNumTXs)
		hiddenTXs2   = make([]*transaction.LocalHidden, testNumTXs)
		points1      = make([]kyber.Point, testNumTXs)
		points2      = make([]kyber.Point, testNumTXs)
		randScalars1 = make([]kyber.Scalar, testNumTXs)
		randScalars2 = make([]kyber.Scalar, testNumTXs)
	)
	for i := 0; i < testNumTXs; i++ {
		hiddenTX1, point1, scalar1, err := localTXs1[i].Hide()
		if err != nil {
			t.Fatal(err)
		}
		hiddenTXs1[i] = hiddenTX1
		points1[i] = point1
		randScalars1[i] = scalar1
		organizations[0].Accumulate(organizations[1].ID, point1)
		hiddenTX2, point2, scalar2, err := localTXs2[i].Hide()
		if err != nil {
			t.Fatal(err)
		}
		hiddenTXs2[i] = hiddenTX2
		points2[i] = point2
		randScalars2[i] = scalar2
		organizations[1].Accumulate(organizations[0].ID, point2)
	}
	// compose org chain transactions
	orgTX1, err := organizations[0].ComposeTXOrgChain(organizations[1].ID)
	if err != nil {
		t.Fatal(err)
	}
	orgTX2, err := organizations[1].ComposeTXOrgChain(organizations[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	// compute res
	res1, err := auditors[0].AccumulateCommitments(organizations[0].ID, hiddenTXs1)
	if err != nil {
		t.Fatal(err)
	}
	res2, err := auditors[1].AccumulateCommitments(organizations[1].ID, hiddenTXs2)
	if err != nil {
		t.Fatal(err)
	}
	// compute b
	txRandList := auditors[0].GetEpochTXRandomness(organizations[0].ID, organizations[1].ID)
	if txRandList == nil {
		t.Fatal("txRandList is nil")
	}
	b1, err := auditors[0].ComputeB(randScalars1, txRandList)
	if err != nil {
		t.Fatal(err)
	}
	b2, err := auditors[1].ComputeB(randScalars2, txRandList)
	if err != nil {
		t.Fatal(err)
	}
	// compute c
	orgIDPoint1 := organization.EpochIDHashPoint(organizations[0].EpochID)
	orgIDPoint2 := organization.EpochIDHashPoint(organizations[1].EpochID)
	acc1 := crypto.KyberSuite.Point()
	if err := acc1.UnmarshalBinary(orgTX1.Accumulator); err != nil {
		t.Fatal(err)
	}
	acc1.Sub(acc1, orgIDPoint1)
	acc2 := crypto.KyberSuite.Point()
	if err := acc2.UnmarshalBinary(orgTX2.Accumulator); err != nil {
		t.Fatal(err)
	}
	acc2.Sub(acc2, orgIDPoint2)
	c1 := auditors[0].ComputeC(res1, acc1)
	c2 := auditors[1].ComputeC(res2, acc2)
	// compute d
	d1 := auditors[0].ComputeD(acc1, b1)
	d2 := auditors[1].ComputeD(acc2, b2)
	randID1 := make([]byte, 32)
	if _, err = rand.Read(randID1); err != nil {
		t.Fatal(err)
	}
	randID2 := make([]byte, 32)
	if _, err := rand.Read(randID2); err != nil {
		t.Fatal(err)
	}
	idHash1 := organization.IDHashString(organizations[0].ID)
	idHash2 := organization.IDHashString(organizations[1].ID)
	publicKey1 := publicKeyMap[idHash1]
	publicKey2 := publicKeyMap[idHash2]
	audTX1, err := auditors[0].EncryptConsistencyExamResult(organizations[0].ID, idHash2, res1, b1, c1, d1, publicKey1)
	if err != nil {
		t.Fatal(err)
	}
	audTX2, err := auditors[1].EncryptConsistencyExamResult(organizations[1].ID, idHash1, res2, b2, c2, d2, publicKey2)
	if err != nil {
		t.Fatal(err)
	}
	// verify
	res, err := com.VerifyAuditPairResult(
		organizations[0].ID, organizations[1].ID,
		auditors[0].ID, auditors[1].ID,
		audTX1.ToOnChain(), audTX2.ToOnChain(),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !res {
		t.Fatal("verify failed")
	}

}
