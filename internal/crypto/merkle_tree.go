package crypto

import (
	"crypto/sha256"
	"encoding/json"

	mt "github.com/txaty/go-merkletree"
)

var (
	sha256Digest     = sha256.New()
	merkleTreeConfig = &mt.Config{
		HashFunc:           hashFunc,
		DisableLeafHashing: true,
	}
)

func hashFunc(data []byte) ([]byte, error) {
	defer sha256Digest.Reset()
	sha256Digest.Write(data)
	return sha256Digest.Sum(nil), nil
}

func GenerateMerkleProofs(dataBlocks []mt.DataBlock) ([]*mt.Proof, []byte, error) {
	tree, err := mt.New(merkleTreeConfig, dataBlocks)
	if err != nil {
		return nil, nil, err
	}
	return tree.Proofs, tree.Root, nil
}

func VerifyMerkleProof(block mt.DataBlock, proof *mt.Proof, root []byte) (bool, error) {
	return mt.Verify(block, proof, root, merkleTreeConfig)
}

type MerkleProof struct {
	Siblings [][]byte `json:"siblings"`
	Path     uint32   `json:"path"`
}

func MerkleProofMarshal(proof *mt.Proof) ([]byte, error) {
	return json.Marshal(&MerkleProof{
		Siblings: proof.Siblings,
		Path:     proof.Path,
	})
}

func MerkleProofUnmarshal(data []byte) (*mt.Proof, error) {
	var proof MerkleProof
	if err := json.Unmarshal(data, &proof); err != nil {
		return nil, err
	}
	return &mt.Proof{
		Siblings: proof.Siblings,
		Path:     proof.Path,
	}, nil
}
