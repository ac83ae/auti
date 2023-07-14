package task

import (
	"math/rand"
	"testing"
	"time"

	"go.dedis.ch/kyber/v3"

	"github.com/ac83ae/auti/internal/clolc/organization"
	"github.com/ac83ae/auti/internal/clolc/transaction"
	"github.com/ac83ae/auti/internal/constants"
)

const testNumTXs = constants.MaxNumTXInEpoch

func TestCECheck(t *testing.T) {
	// setup
	com, auditors, organizations := generateEntities(2)
	_, err := com.InitializeEpoch(auditors, organizations)
	if err != nil {
		t.Fatal(err)
	}
	// test
	txList1, txList2 := generateLocalTXPairList(organizations[0].ID, organizations[1].ID)
	// compute hidden transactions
	var (
		hiddenTXs1   = make([]*transaction.LocalHidden, testNumTXs)
		hiddenTXs2   = make([]*transaction.LocalHidden, testNumTXs)
		points1      = make([]kyber.Point, testNumTXs)
		points2      = make([]kyber.Point, testNumTXs)
		randScalars1 = make([]kyber.Scalar, testNumTXs)
		randScalars2 = make([]kyber.Scalar, testNumTXs)
	)
	for i := 0; i < testNumTXs; i++ {
		hiddenTX1, point1, scalar1, err := txList1[i].Hide()
		if err != nil {
			t.Fatal(err)
		}
		hiddenTXs1[i] = hiddenTX1
		points1[i] = point1
		randScalars1[i] = scalar1
		hiddenTX2, point2, scalar2, err := txList2[i].Hide()
		if err != nil {
			t.Fatal(err)
		}
		hiddenTXs2[i] = hiddenTX2
		points2[i] = point2
		randScalars2[i] = scalar2
	}
	res1, err := auditors[0].AccumulateCommitments(organizations[0].ID, hiddenTXs1)
	if err != nil {
		t.Fatal(err)
	}
	res2, err := auditors[1].AccumulateCommitments(organizations[1].ID, hiddenTXs2)
	if err != nil {
		t.Fatal(err)
	}
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
	if !auditors[0].CheckResultConsistency(res1, b1, res2, b2) {
		t.Fatal("result consistency check failed")
	}
}

func randAmount() float64 {
	amount := rand.Float64()
	integerPart := rand.Int()%10000 - 5000
	amount += float64(integerPart)
	return amount
}

func generateLocalTXPairList(
	orgID1, orgID2 organization.TypeID,
) ([]*transaction.LocalPlain, []*transaction.LocalPlain) {
	txList1 := make([]*transaction.LocalPlain, testNumTXs)
	txList2 := make([]*transaction.LocalPlain, testNumTXs)
	for i := 0; i < testNumTXs; i++ {
		currTimestamp := time.Now().UnixNano()
		tx1, tx2 := transaction.NewPairLocalPlain(
			string(orgID1),
			string(orgID2),
			randAmount(),
			currTimestamp,
		)
		txList1[i] = tx1
		txList2[i] = tx2
	}
	return txList1, txList2
}
