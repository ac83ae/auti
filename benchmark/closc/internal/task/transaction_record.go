package task

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	mt "github.com/txaty/go-merkletree"

	"github.com/ac83ae/auti/benchmark/closc/internal/blockchain/localchain"
	"github.com/ac83ae/auti/benchmark/closc/internal/blockchain/localchaincommit"
	"github.com/ac83ae/auti/benchmark/closc/internal/blockchain/localchainsc"
	"github.com/ac83ae/auti/benchmark/closc/internal/blockchain/orgchain"
	"github.com/ac83ae/auti/benchmark/timecounter"
	"github.com/ac83ae/auti/internal/constants"
	"github.com/ac83ae/auti/internal/crypto"
)

const hashByteLen = 32

var numCPU = runtime.NumCPU()

type randInput struct {
	amount       int64
	timestamp    int64
	receiverHash []byte
	counter      uint64
}

func generateRandInputs(num int) []randInput {
	results := make([]randInput, num)
	var wg sync.WaitGroup
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(idx, step int) {
			defer wg.Done()
			for j := idx; j < num; j += step {
				randBytes := make([]byte, constants.SecurityParameterBytes)
				_, err := crand.Read(randBytes)
				if err != nil {
					panic(err)
				}
				results[j] = randInput{
					amount:       rand.Int63(),
					timestamp:    time.Now().UnixNano(),
					receiverHash: randBytes,
					counter:      uint64(j),
				}
			}
		}(i, numCPU)
	}
	wg.Wait()
	return results
}

func TRCommitment(num, iterations int) error {
	fmt.Println("[CLOSC-TR] Commitment")
	fmt.Printf("Num: %d, Num iter: %d\n", num, iterations)
	for i := 0; i < iterations; i++ {
		randInputs := generateRandInputs(num)
		startTime := time.Now()
		for j := 0; j < num; j++ {
			if _, _, err := crypto.PedersonCommitWithHash(
				randInputs[j].amount,
				randInputs[j].timestamp,
				randInputs[j].receiverHash,
				randInputs[j].counter,
			); err != nil {
				return err
			}
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

type dummyDataBlock struct {
	data []byte
}

func (d dummyDataBlock) Serialize() ([]byte, error) {
	return d.data, nil
}

func generateDataBlocks(num int) []mt.DataBlock {
	results := make([]mt.DataBlock, num)
	var wg sync.WaitGroup
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(idx, step int) {
			defer wg.Done()
			for j := idx; j < num; j += step {
				randBytes := make([]byte, hashByteLen)
				_, err := crand.Read(randBytes)
				if err != nil {
					panic(err)
				}
				results[j] = dummyDataBlock{data: randBytes}
			}
		}(i, numCPU)
	}
	wg.Wait()
	return results
}

func TRMerkleProofGen(depth, iterations int) error {
	fmt.Println("[CLOSC-TR] Merkle proof generation")
	fmt.Printf("Depth: %d, Num iter: %d\n", depth, iterations)
	numDataBlock := 1 << depth
	dataBlocks := generateDataBlocks(numDataBlock)
	for i := 0; i < iterations; i++ {
		startTime := time.Now()
		_, _, err := crypto.GenerateMerkleProofs(dataBlocks)
		if err != nil {
			return err
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func TRLocalSubmitTX(numTXs, iterations int) error {
	fmt.Println("[CLOSC-TR] Local submit transaction")
	fmt.Printf("Num TXs: %d, Num iter: %d\n", numTXs, iterations)
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
	fmt.Println("[CLOSC-TR] Local prepare transaction")
	fmt.Printf("Num TXs: %d\n", numTotalTXs)
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
	fmt.Println("[CLOSC-TR] Local read transaction")
	fmt.Printf("Num TXs: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		if err := localchain.ReadTX(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TRLocalReadAllTXs(numTotals, iterations int) error {
	fmt.Println("[CLOSC-TR] Local read all transactions")
	fmt.Printf("Num TXs: %d, Num iter: %d\n", numTotals, iterations)
	for i := 0; i < iterations; i++ {
		if err := localchain.ReadAllTXsByPage(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TRLocalCommitmentSubmitTX(numTXs, iterations int) error {
	fmt.Println("[CLOSC-TR] Local submit commitment transaction")
	fmt.Printf("Num TXs: %d, Num iter: %d\n", numTXs, iterations)
	for i := 0; i < iterations; i++ {
		_, err := localchaincommit.SubmitTX(numTXs)
		if err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TRLocalCommitmentPrepareTX(numTotalTXs int) error {
	fmt.Println("[CLOSC-TR] Local prepare commitment transaction")
	fmt.Printf("Num TXs: %d\n", numTotalTXs)
	txIDs, err := localchaincommit.SubmitTX(numTotalTXs)
	if err != nil {
		return err
	}
	if err = localchaincommit.SaveTXIDs(txIDs); err != nil {
		return err
	}
	fmt.Println()
	return nil
}

func TRLocalCommitmentReadTX(numTotalTXs, iterations int) error {
	fmt.Println("[CLOSC-TR] Local read commitment transaction")
	fmt.Printf("Num TXs: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		if err := localchaincommit.ReadTX(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TRLocalCommitmentReadAllTXs(numTotals, iterations int) error {
	fmt.Println("[CLOSC-TR] Local read all commitment transactions")
	fmt.Printf("Num TXs: %d, Num iter: %d\n", numTotals, iterations)
	for i := 0; i < iterations; i++ {
		if err := localchaincommit.ReadAllTXsByPage(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TROrgSubmitTX(numTXs, iterations int) error {
	fmt.Println("[CLOSC-TR] Org submit transaction")
	fmt.Printf("Num TXs: %d, Num iter: %d\n", numTXs, iterations)
	for i := 0; i < iterations; i++ {
		_, err := orgchain.SubmitTX(numTXs)
		if err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TROrgPrepareTX(numTotalTXs int) error {
	fmt.Println("[CLOSC-TR] Org prepare transaction")
	fmt.Printf("Num TXs: %d\n", numTotalTXs)
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

func TROrgReadTX(numTotalTXs, iterations int) error {
	fmt.Println("[CLOSC-TR] Org read transaction")
	fmt.Printf("Num TXs: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		if err := orgchain.ReadTX(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TROrgReadAllTXs(numTotals, iterations int) error {
	fmt.Println("[CLOSC-TR] Org read all transactions")
	fmt.Printf("Num TXs: %d, Num iter: %d\n", numTotals, iterations)
	for i := 0; i < iterations; i++ {
		if err := orgchain.ReadAllTXsByPage(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TRLocalSC(numSC, iterations int) error {
	fmt.Println("[CLOSC-TR] Local smart contract")
	fmt.Printf("Num SC: %d, Num iter: %d\n", numSC, iterations)
	for i := 0; i < iterations; i++ {
		for j := 0; j < numSC; j++ {
			if _, err := localchainsc.SubmitTX(j); err != nil {
				return err
			}
		}
	}
	fmt.Println()
	return nil
}
