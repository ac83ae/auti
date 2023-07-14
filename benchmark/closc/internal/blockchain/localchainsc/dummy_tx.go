package localchainsc

import (
	crand "crypto/rand"
	"math/rand"

	mt "github.com/txaty/go-merkletree"

	"github.com/ac83ae/auti/internal/closc/transaction"
	"github.com/ac83ae/auti/internal/crypto"
)

const (
	hashByteLen = 32
	proofDepth  = 5
)

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
	return plainTX.ToOnChain(), nil
}

func DummyPlainTransaction() (*transaction.LocalPlain, error) {
	randScalar := crypto.KyberSuite.Scalar().Pick(crypto.KyberSuite.RandomStream())
	randPoint := crypto.KyberSuite.Point().Mul(randScalar, nil)
	dummyCommitment, err := randPoint.MarshalBinary()
	if err != nil {
		return nil, err
	}
	dummyRoot := make([]byte, hashByteLen)
	_, err = crand.Read(dummyRoot)
	if err != nil {
		return nil, err
	}
	siblings := make([][]byte, proofDepth)
	for i := 0; i < proofDepth; i++ {
		randBytes := make([]byte, hashByteLen)
		_, err = crand.Read(randBytes)
		if err != nil {
			return nil, err
		}
		siblings[i] = randBytes
	}
	dummyProof := &mt.Proof{
		Siblings: siblings,
		Path:     rand.Uint32(),
	}
	dummyProofBytes, err := crypto.MerkleProofMarshal(dummyProof)
	if err != nil {
		return nil, err
	}
	return transaction.NewLocalPlain(dummyCommitment, dummyRoot, dummyProofBytes), nil
}
