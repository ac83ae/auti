package audchain

import (
	crand "crypto/rand"
	"runtime"
	"sync"

	"go.dedis.ch/kyber/v3/group/edwards25519"

	"github.com/ac83ae/auti/internal/closc/transaction"
)

var (
	numCPUs     = runtime.NumCPU()
	kyberSuite  = edwards25519.NewBlakeSHA256Ed25519()
	hashByteLen = 32
)

func DummyOnChainTransactions(numTXs int) []*transaction.AudOnChain {
	results := make([]*transaction.AudOnChain, numTXs)
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

func DummyOnChainTransaction() (*transaction.AudOnChain, error) {
	randIDBytes := make([]byte, 32)
	_, err := crand.Read(randIDBytes)
	if err != nil {
		return nil, err
	}
	randPoint := kyberSuite.Point().Pick(kyberSuite.RandomStream())
	randHashBytes := make([]byte, hashByteLen)
	_, err = crand.Read(randHashBytes)
	if err != nil {
		return nil, err
	}
	tx, err := transaction.NewAudPlainFromPoint(randPoint, randHashBytes)
	if err != nil {
		return nil, err
	}
	return tx.ToOnChain(), nil
}
