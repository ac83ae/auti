package crypto

import (
	"crypto/rand"
	"testing"

	mt "github.com/txaty/go-merkletree"
)

type dummyDataBlock []byte

func (d dummyDataBlock) Serialize() ([]byte, error) {
	return d, nil
}

func dummyDataBlocks(numBlocks int) []mt.DataBlock {
	dataBlocks := make([]mt.DataBlock, numBlocks)
	for i := 0; i < numBlocks; i++ {
		randBytes := make([]byte, 32)
		if _, err := rand.Read(randBytes); err != nil {
			return nil
		}
		dataBlocks[i] = dummyDataBlock(randBytes)
	}
	return dataBlocks
}

func TestGenerateMerkleProofs(t *testing.T) {
	tests := []struct {
		name          string
		numdataBlocks int
		wantErr       bool
	}{
		{
			name:          "1",
			numdataBlocks: 1,
			wantErr:       true,
		},
		{
			name:          "4",
			numdataBlocks: 4,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dummyBlocks := dummyDataBlocks(tt.numdataBlocks)
			proofs, root, err := GenerateMerkleProofs(dummyBlocks)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateMerkleProofs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for idx, proof := range proofs {
				ok, err := VerifyMerkleProof(dummyBlocks[idx], proof, root)
				if err != nil {
					return
				}
				if !ok {
					t.Errorf("VerifyMerkleProof() = %v, want %v", ok, true)
				}
			}
		})
	}
}
