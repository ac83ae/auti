package task

import (
	"fmt"
	"time"

	"go.dedis.ch/kyber/v3/group/edwards25519"

	"github.com/ac83ae/auti/benchmark/clolc/internal/blockchain/localchain"
	"github.com/ac83ae/auti/benchmark/clolc/internal/blockchain/orgchain"
	"github.com/ac83ae/auti/benchmark/timecounter"
)

func TRLocalSubmitTX(numTXs, iterations int) error {
	fmt.Println("[CLOLC-TR] Local submit transaction")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTXs, iterations)
	for i := 0; i < iterations; i++ {
		_, err := localchain.SubmitTX(numTXs)
		if err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TRLocalPrepareTX(numTotalTXs int) error {
	fmt.Println("[CLOLC-TR] Prepare local transaction")
	fmt.Printf("Num TX: %d\n", numTotalTXs)
	txIDs, err := localchain.SubmitTX(numTotalTXs)
	if err != nil {
		return err
	}
	if err = localchain.SaveTXIDs(txIDs); err != nil {
		return err
	}
	fmt.Println()
	return nil
}

func TRLocalReadTX(numTotalTXs, iterations int) error {
	fmt.Println("[CLOLC-TR] Local read transaction")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		if err := localchain.ReadTX(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TRLocalReadAllTXs(numTotalTXs, iterations int) error {
	fmt.Println("[CLOLC-TR] Local read all transactions")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		if err := localchain.ReadAllTXsByPage(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TRCommitment(numTotalTXs, iterations int) error {
	fmt.Println("[CLOLC-TR] Commitment")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		dummyTXs := localchain.DummyPlainTransactions(numTotalTXs)
		startTime := time.Now()
		for _, tx := range dummyTXs {
			if _, _, _, err := tx.Hide(); err != nil {
				return err
			}
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func TRAccumulate(numTotalTXs, iterations int) error {
	fmt.Println("[CLOLC-TR] Accumulate")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		dummyCommitments := localchain.DummyHiddenTXCommitments(numTotalTXs)
		kyberSuite := edwards25519.NewBlakeSHA256Ed25519()
		accumulator := kyberSuite.Point().Null()
		startTime := time.Now()
		for _, commitment := range dummyCommitments {
			accumulator = accumulator.Add(accumulator, commitment)
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func TROrgPrepareTX(numTotalTXs int) error {
	fmt.Println("[CLOLC-TR] Prepare transaction")
	fmt.Printf("Num TX: %d\n", numTotalTXs)
	txIDs, err := orgchain.SubmitTX(numTotalTXs)
	if err != nil {
		return err
	}
	if err = orgchain.SaveTXIDs(txIDs); err != nil {
		return err
	}
	fmt.Println()
	return nil
}

func TROrgSubmitTX(numTXs, iterations int) error {
	fmt.Println("[CLOLC-TR] Submit transaction")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTXs, iterations)
	for i := 0; i < iterations; i++ {
		if _, err := orgchain.SubmitTX(numTXs); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TROrgReadTX(numTotalTXs, iterations int) error {
	fmt.Println("[CLOLC-TR] Read transaction")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		if err := orgchain.ReadTX(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TROrgReadAllTXs(numTotalTXs, iterations int) error {
	fmt.Println("[CLOLC-TR] Read all transactions")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		if err := orgchain.ReadAllTXsByPage(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}
