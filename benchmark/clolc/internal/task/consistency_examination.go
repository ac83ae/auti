package task

import (
	"fmt"
	"time"

	"go.dedis.ch/kyber/v3"

	"github.com/ac83ae/auti/benchmark/clolc/internal/blockchain/audchain"
	"github.com/ac83ae/auti/benchmark/clolc/internal/blockchain/localchain"
	"github.com/ac83ae/auti/benchmark/timecounter"
	"github.com/ac83ae/auti/internal/clolc/organization"
	"github.com/ac83ae/auti/internal/constants"
	"github.com/ac83ae/auti/internal/crypto"
)

func CEAccumulateCommitment(numOrganizations, iterations int) error {
	fmt.Println("[CLOLC-CE] Accumulate Commitment")
	fmt.Printf("Num org %d, Num iter: %d\n", numOrganizations, iterations)
	for i := 0; i < iterations; i++ {
		com, auditors, organizations := generateEntities(numOrganizations)
		_, err := com.InitializeEpoch(auditors, organizations)
		if err != nil {
			return err
		}
		dummyTXs := localchain.DummyHiddenTXWithCounterPartyID(organizations[1].ID, constants.MaxNumTXInEpoch)
		startTime := time.Now()
		if _, err = auditors[0].AccumulateCommitments(organizations[0].ID, dummyTXs); err != nil {
			return err
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func CEComputeB(numOrganizations, iterations int) error {
	fmt.Println("[CLOLC-CE] Compute B")
	fmt.Printf("Num org %d, Num iter: %d\n", numOrganizations, iterations)
	for i := 0; i < iterations; i++ {
		com, auditors, organizations := generateEntities(numOrganizations)
		_, err := com.InitializeEpoch(auditors, organizations)
		if err != nil {
			return err
		}
		randScalars1 := make([]kyber.Scalar, constants.MaxNumTXInEpoch)
		randScalars2 := make([]kyber.Scalar, constants.MaxNumTXInEpoch)
		for i := 0; i < constants.MaxNumTXInEpoch; i++ {
			randScalars1[i] = crypto.KyberSuite.Scalar().Pick(crypto.KyberSuite.RandomStream())
			randScalars2[i] = crypto.KyberSuite.Scalar().Pick(crypto.KyberSuite.RandomStream())
		}
		startTime := time.Now()
		if _, err = auditors[0].ComputeB(randScalars1, randScalars2); err != nil {
			return err
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func CEComputeC(numOrganizations, iterations int) error {
	fmt.Println("[CLOLC-CE] Compute C")
	fmt.Printf("Num org %d, Num iter: %d\n", numOrganizations, iterations)
	for i := 0; i < iterations; i++ {
		com, auditors, organizations := generateEntities(numOrganizations)
		_, err := com.InitializeEpoch(auditors, organizations)
		if err != nil {
			return err
		}
		randPoint1 := crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
		randPoint2 := crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
		startTime := time.Now()
		_ = auditors[0].ComputeC(randPoint1, randPoint2)
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func CEComputeD(numOrganizations, iterations int) error {
	fmt.Println("[CLOLC-CE] Compute D")
	fmt.Printf("Num org %d, Num iter: %d\n", numOrganizations, iterations)
	for i := 0; i < iterations; i++ {
		com, auditors, organizations := generateEntities(numOrganizations)
		_, err := com.InitializeEpoch(auditors, organizations)
		if err != nil {
			return err
		}
		randPoint1 := crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
		randPoint2 := crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
		startTime := time.Now()
		_ = auditors[0].ComputeD(randPoint1, randPoint2)
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func CEEncrypt(numOrganizations, iterations int) error {
	fmt.Println("[CLOLC-CE] Encrypt")
	fmt.Printf("Num org %d, Num iter: %d\n", numOrganizations, iterations)
	for i := 0; i < iterations; i++ {
		com, auditors, organizations := generateEntities(numOrganizations)
		_, err := com.InitializeEpoch(auditors, organizations)
		if err != nil {
			return err
		}
		counterPartyHashStr := organization.IDHashString(organizations[1].ID)
		_, publicKey, err := crypto.KeyGen()
		if err != nil {
			return err
		}
		randPoint1 := crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
		randPoint2 := crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
		randPoint3 := crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
		randPoint4 := crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
		startTime := time.Now()
		if _, err := auditors[0].EncryptConsistencyExamResult(
			organizations[0].ID, counterPartyHashStr, randPoint1, randPoint2, randPoint3, randPoint4, publicKey,
		); err != nil {
			return err
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func CEAudSubmitTX(numTotalTXs, iterations int) error {
	fmt.Println("[CLOLC-CE] Submit TX")
	fmt.Printf("Num total TXs %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		if _, err := audchain.SubmitTX(numTotalTXs); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func CEAudReadTX(numTotalTXs, iterations int) error {
	fmt.Println("[CLOLC-CE] Read TX")
	fmt.Printf("Num total TXs %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		if err := audchain.ReadTX(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func CEAudReadAllTXs(numTotalTXs, iterations int) error {
	fmt.Println("[CLOLC-CE] Read all TXs")
	fmt.Printf("Num total TXs %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		if err := audchain.ReadAllTXsByPage(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func CEAudPrepareTX(numTotalTXs int) error {
	fmt.Println("[CLOLC-CE] Prepare aud transaction")
	fmt.Printf("Num TX: %d\n", numTotalTXs)
	txIDs, err := audchain.SubmitTX(numTotalTXs)
	if err != nil {
		return err
	}
	if err = audchain.SaveTXIDs(txIDs); err != nil {
		return err
	}
	fmt.Println()
	return nil
}

func CEDecrypt(iterations int) error {
	fmt.Println("[CLOLC-CE] Decrypt")
	fmt.Printf("Num iter: %d\n", iterations)
	com, auditors, organizations := generateEntities(2)
	_, err := com.InitializeEpoch(auditors, organizations)
	if err != nil {
		return err
	}
	for i := 0; i < iterations; i++ {
		dummyTX, err := audchain.DummyOnChainTransaction()
		if err != nil {
			return err
		}
		orgIDHashStr := organization.IDHashString(organizations[0].ID)
		startTime := time.Now()
		if _, _, err := auditors[0].DecryptResAndB(orgIDHashStr, dummyTX); err != nil {
			return err
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func CECheck(iterations int) error {
	fmt.Println("[CLOLC-CE] Check")
	fmt.Printf("Num iter: %d\n", iterations)
	for i := 0; i < iterations; i++ {
		com, auditors, organizations := generateEntities(2)
		_, err := com.InitializeEpoch(auditors, organizations)
		if err != nil {
			return err
		}
		randPoints := make([]kyber.Point, 4)
		for i := 0; i < 4; i++ {
			randPoints[i] = crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
		}
		startTime := time.Now()
		_ = auditors[0].CheckResultConsistency(
			randPoints[0], randPoints[1], randPoints[2], randPoints[3],
		)
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}
