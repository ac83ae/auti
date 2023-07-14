package task

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	mt "github.com/txaty/go-merkletree"
	"go.dedis.ch/kyber/v3"

	"github.com/ac83ae/auti/benchmark/timecounter"
	"github.com/ac83ae/auti/internal/closc/auditor"
	"github.com/ac83ae/auti/internal/closc/committee"
	"github.com/ac83ae/auti/internal/crypto"
)

func RVVerifyMerkleBatchProof(numTXs, iterations int) error {
	fmt.Println("[CLOSC-RV] Verify merkle batch proof")
	fmt.Printf("Num TXs: %d, Num iter: %d\n", numTXs, iterations)
	numTotalTXs := 1 << mergeTreeDepth
	dummyBlocks, dummyProofs, dummyRoot, err := genDummyDataBlockAndProof(mergeTreeDepth)
	if err != nil {
		return err
	}
	aud := auditor.New("aud", nil)
	com := committee.New("com", nil)
	for i := 0; i < iterations; i++ {
		indexes := randIndexes(numTXs, numTotalTXs)
		selectedBlocks := make([]mt.DataBlock, numTXs)
		selectedProofs := make([]*mt.Proof, numTXs)
		for j := 0; j < numTXs; j++ {
			selectedBlocks[j] = dummyBlocks[indexes[j]]
			selectedProofs[j] = dummyProofs[indexes[j]]
		}
		mergedProofByte, err := aud.MergeProof(selectedBlocks, selectedProofs)
		if err != nil {
			return err
		}
		startTime := time.Now()
		mergedProof, err := crypto.MerkleBatchProofUnmarshal(mergedProofByte)
		if err != nil {
			return err
		}
		res, err := com.VerifyMerkleBatchProof(selectedBlocks, mergedProof, dummyRoot)
		if err != nil {
			return err
		}
		if res != 1 {
			return fmt.Errorf("merkle batch proof verification failed")
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func RVSummarizeMerkleBatchProofVerificationResults(numResults, iterations int) error {
	fmt.Println("[CLOSC-RV] Summarize merkle batch proof verification results")
	fmt.Printf("Num results: %d, Num iter: %d\n", numResults, iterations)
	results := make([]uint, numResults)
	for i := 0; i < iterations; i++ {
		for j := 0; j < numResults; j++ {
			results[j] = uint(rand.Int() % 2)
		}
		com := committee.New("com", nil)
		startTime := time.Now()
		com.SummarizeMerkleBatchProofVerificationResults(results)
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func genDummyPoints(num int) []kyber.Point {
	results := make([]kyber.Point, num)
	var wg sync.WaitGroup
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(idx, step int) {
			defer wg.Done()
			for j := idx; j < num; j += step {
				results[j] = crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
			}
		}(i, numCPU)
	}
	wg.Wait()
	return results
}

func RVVerifyCommitments(numCommitments, iterations int) error {
	fmt.Println("[CLOSC-RV] Verify commitment")
	fmt.Printf("Num commitments: %d, Num iter: %d\n", numCommitments, iterations)
	for i := 0; i < iterations; i++ {
		commitments := genDummyPoints(numCommitments)
		com := committee.New("com", nil)
		startTime := time.Now()
		com.VerifyCommitment(commitments)
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}
