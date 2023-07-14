package crypto

import (
	"math/rand"
	"testing"

	mt "github.com/txaty/go-merkletree"
)

// randIndexes generates the random indexes without duplication
func randIndexes(numIdx, max int) []int {
	if numIdx > max {
		numIdx = max
	}
	// Generate a pool of indexes
	pool := make([]int, max)
	for i := 0; i < max; i++ {
		pool[i] = i
	}
	// Shuffle the pool using Fisher-Yates algorithm
	for i := max - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		pool[i], pool[j] = pool[j], pool[i]
	}
	// Take the first numIdx elements from the shuffled pool
	return pool[:numIdx]
}

func TestNewMerkleBatchProof(t *testing.T) {
	tests := []struct {
		name              string
		numTotalBlocks    int
		numSelectedBlocks int
		wantErr           bool
	}{
		{
			name:              "100_10",
			numTotalBlocks:    100,
			numSelectedBlocks: 10,
			wantErr:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dummyBlocks := dummyDataBlocks(tt.numTotalBlocks)
			proofs, root, err := GenerateMerkleProofs(dummyBlocks)
			if err != nil {
				t.Errorf("GenerateMerkleProofs() error = %v", err)
				return
			}
			randIdxList := randIndexes(tt.numSelectedBlocks, tt.numTotalBlocks)
			selectedBlocks := make([]mt.DataBlock, tt.numSelectedBlocks)
			selectedProofs := make([]*mt.Proof, tt.numSelectedBlocks)
			for idx, randIdx := range randIdxList {
				selectedBlocks[idx] = dummyBlocks[randIdx]
				selectedProofs[idx] = proofs[randIdx]
			}
			mergedProof, err := NewMerkleBatchProof(selectedBlocks, selectedProofs)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMerkleBatchProof() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			ok, err := VerifyMerkleBatchProof(selectedBlocks, mergedProof, root)
			if err != nil {
				t.Errorf("VerifyMerkleBatchProof() error = %v", err)
				return
			}
			if !ok {
				t.Errorf("VerifyMerkleBatchProof() = %v, want %v", ok, true)
			}
		})
	}
}
