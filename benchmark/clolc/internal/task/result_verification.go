package task

import (
	"fmt"
	"time"

	"github.com/ac83ae/auti/benchmark/clolc/internal/blockchain/audchain"
	"github.com/ac83ae/auti/benchmark/clolc/internal/blockchain/orgchain"
	"github.com/ac83ae/auti/benchmark/timecounter"
)

func RVVerifyOrgAndAudResult(numOrganizations, iterations int) error {
	fmt.Println("[CLOLC-RV] Verify org and aud result")
	fmt.Printf("Num org %d, Num iter: %d\n", numOrganizations, iterations)
	for i := 0; i < iterations; i++ {
		com, auditors, organizations := generateEntities(numOrganizations)
		_, err := com.InitializeEpoch(auditors, organizations)
		if err != nil {
			return err
		}
		dummyOrgChainTX, err := orgchain.DummyOnChainTransaction()
		if err != nil {
			return err
		}
		dummyAudChainTX, err := audchain.DummyOnChainTransaction()
		if err != nil {
			return err
		}
		startTime := time.Now()
		if _, err = com.VerifyOrgAndAudResult(
			organizations[0].ID,
			auditors[0].ID,
			dummyOrgChainTX,
			dummyAudChainTX,
		); err != nil {
			return err
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func RVVerifyAuditPairResult(numOrganizations, iterations int) error {
	fmt.Println("[CLOLC-RV] Verify audit pair result")
	fmt.Printf("Num org %d, Num iter: %d\n", numOrganizations, iterations)
	for i := 0; i < iterations; i++ {
		com, auditors, organizations := generateEntities(numOrganizations)
		_, err := com.InitializeEpoch(auditors, organizations)
		if err != nil {
			return err
		}
		dummyTX1, err := audchain.DummyOnChainTransaction()
		if err != nil {
			return err
		}
		dummyTX2, err := audchain.DummyOnChainTransaction()
		if err != nil {
			return err
		}
		startTime := time.Now()
		if _, err = com.VerifyAuditPairResult(
			organizations[0].ID,
			organizations[1].ID,
			auditors[0].ID,
			auditors[1].ID,
			dummyTX1,
			dummyTX2,
		); err != nil {
			return err
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}
