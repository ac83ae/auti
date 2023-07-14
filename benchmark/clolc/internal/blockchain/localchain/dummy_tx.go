package localchain

import (
	crand "crypto/rand"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"go.dedis.ch/kyber/v3"

	"github.com/ac83ae/auti/internal/clolc/organization"
	"github.com/ac83ae/auti/internal/clolc/transaction"
)

var numCPUs = runtime.NumCPU()

func DummyOnChainTransactions(numTXs int) []*transaction.LocalOnChain {
	results := make([]*transaction.LocalOnChain, numTXs)
	var wg sync.WaitGroup
	for i := 0; i < numCPUs; i++ {
		wg.Add(1)
		go func(idx, step int) {
			defer wg.Done()
			for j := idx; j < numTXs; j += step {
				dummyTX, err := DummyOnChainTransaction()
				if err != nil {
					panic(err)
				}
				results[j] = dummyTX
			}
		}(i, numCPUs)
	}
	wg.Wait()
	return results
}

func DummyOnChainTransaction() (*transaction.LocalOnChain, error) {
	dummyCounterPartyBytes := make([]byte, 32)
	_, err := crand.Read(dummyCounterPartyBytes)
	if err != nil {
		return nil, err
	}
	plainTX, err := DummyPlainTransaction()
	if err != nil {
		return nil, err
	}
	hiddenTX, _, _, err := plainTX.Hide()
	if err != nil {
		return nil, err
	}
	return hiddenTX.ToOnChain(), nil
}

func DummyPlainTransactions(numTXs int) []*transaction.LocalPlain {
	results := make([]*transaction.LocalPlain, numTXs)
	var wg sync.WaitGroup
	for i := 0; i < numCPUs; i++ {
		wg.Add(1)
		go func(idx, step int) {
			defer wg.Done()
			var err error
			for j := idx; j < numTXs; j += step {
				results[j], err = DummyPlainTransaction()
				if err != nil {
					panic(err)
				}
			}
		}(i, numCPUs)
	}
	wg.Wait()
	return results
}

func DummyPlainTransaction() (*transaction.LocalPlain, error) {
	currTimestamp := time.Now().UnixNano()
	randAmount := rand.Float64()
	dummyCounterPartyBytes := make([]byte, 32)
	_, err := crand.Read(dummyCounterPartyBytes)
	if err != nil {
		return nil, err
	}
	plainTX := transaction.NewLocalPlain(
		string(dummyCounterPartyBytes),
		randAmount,
		currTimestamp,
	)
	return plainTX, nil
}

func DummyHiddenTXCommitments(numTXs int) []kyber.Point {
	results := make([]kyber.Point, numTXs)
	var wg sync.WaitGroup
	for i := 0; i < numCPUs; i++ {
		wg.Add(1)
		go func(idx, step int) {
			defer wg.Done()
			for j := idx; j < numTXs; j += step {
				tx, err := DummyPlainTransaction()
				if err != nil {
					panic(err)
				}
				_, commitment, _, err := tx.Hide()
				if err != nil {
					panic(err)
				}
				results[j] = commitment
			}
		}(i, numCPUs)
	}
	wg.Wait()
	return results
}

func DummyHiddenTXWithCounterPartyID(counterPartyID organization.TypeID, numTXs int) []*transaction.LocalHidden {
	results := make([]*transaction.LocalHidden, numTXs)
	var wg sync.WaitGroup
	for i := 0; i < numCPUs; i++ {
		wg.Add(1)
		go func(idx, step int) {
			defer wg.Done()
			for j := idx; j < numTXs; j += step {
				tx, err := DummyPlainTransaction()
				tx.CounterParty = string(counterPartyID)
				if err != nil {
					panic(err)
				}
				hiddenTX, _, _, err := tx.Hide()
				if err != nil {
					panic(err)
				}
				results[j] = hiddenTX
			}
		}(i, numCPUs)
	}
	wg.Wait()
	return results
}
