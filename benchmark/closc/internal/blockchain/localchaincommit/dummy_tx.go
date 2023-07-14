package localchaincommit

import (
	crand "crypto/rand"
	"runtime"
	"sync"

	"github.com/ac83ae/auti/internal/closc/transaction"
	"github.com/ac83ae/auti/internal/crypto"
)

var numCPUs = runtime.NumCPU()

func DummyCommitmentOnChainTransactions(numTXs int) []*transaction.LocalCommitmentOnChain {
	results := make([]*transaction.LocalCommitmentOnChain, numTXs)
	var wg sync.WaitGroup
	for i := 0; i < numCPUs; i++ {
		wg.Add(1)
		go func(idx, step int) {
			defer wg.Done()
			for j := idx; j < numTXs; j += step {
				dummyTX, err := DummyCommitmentOnChainTransaction()
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

func DummyCommitmentOnChainTransaction() (*transaction.LocalCommitmentOnChain, error) {
	dummyCounterPartyBytes := make([]byte, 32)
	_, err := crand.Read(dummyCounterPartyBytes)
	if err != nil {
		return nil, err
	}
	plainTX, err := DummyCommitmentPlainTransaction()
	if err != nil {
		return nil, err
	}
	return plainTX.ToOnChain(), nil
}

func DummyCommitmentPlainTransactions(numTXs int) []*transaction.LocalCommitmentPlain {
	results := make([]*transaction.LocalCommitmentPlain, numTXs)
	var wg sync.WaitGroup
	for i := 0; i < numCPUs; i++ {
		wg.Add(1)
		go func(idx, step int) {
			defer wg.Done()
			var err error
			for j := idx; j < numTXs; j += step {
				results[j], err = DummyCommitmentPlainTransaction()
				if err != nil {
					panic(err)
				}
			}
		}(i, numCPUs)
	}
	wg.Wait()
	return results
}

func DummyCommitmentPlainTransaction() (*transaction.LocalCommitmentPlain, error) {
	randScalar := crypto.KyberSuite.Scalar().Pick(crypto.KyberSuite.RandomStream())
	randPoint := crypto.KyberSuite.Point().Mul(randScalar, nil)
	dummyCommitment, err := randPoint.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return transaction.NewLocalCommitmentPlain(dummyCommitment), nil
}
