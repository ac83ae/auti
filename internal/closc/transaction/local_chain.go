package transaction

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	mt "github.com/txaty/go-merkletree"

	"github.com/ac83ae/auti/internal/crypto"
)

type LocalCommitmentPlain struct {
	Commitment []byte
}

func NewLocalCommitmentPlain(commitment []byte) *LocalCommitmentPlain {
	return &LocalCommitmentPlain{
		Commitment: commitment,
	}
}

func (l *LocalCommitmentPlain) Serialize() ([]byte, error) {
	return l.Commitment, nil
}

func (l *LocalCommitmentPlain) ToOnChain() *LocalCommitmentOnChain {
	return NewLocalCommitmentOnChain(hex.EncodeToString(l.Commitment))
}

type LocalCommitmentOnChain struct {
	Commitment string `json:"commitment"`
}

func NewLocalCommitmentOnChain(commitment string) *LocalCommitmentOnChain {
	return &LocalCommitmentOnChain{
		Commitment: commitment,
	}
}

func (l *LocalCommitmentOnChain) ToPlain() (*LocalCommitmentPlain, error) {
	commitment, err := hex.DecodeString(l.Commitment)
	if err != nil {
		return nil, err
	}
	return NewLocalCommitmentPlain(commitment), nil
}

func (l *LocalCommitmentOnChain) KeyVal() (string, []byte, error) {
	txJSON, err := json.Marshal(l)
	if err != nil {
		return "", nil, err
	}
	sha256Func := sha256.New()
	sha256Func.Write(txJSON)
	return hex.EncodeToString(sha256Func.Sum(nil)), txJSON, nil
}

type LocalPlain struct {
	Commitment  []byte
	MerkleRoot  []byte
	MerkleProof []byte
}

func NewLocalPlain(commitment, merkleRoot, merkleProof []byte) *LocalPlain {
	return &LocalPlain{
		Commitment:  commitment,
		MerkleRoot:  merkleRoot,
		MerkleProof: merkleProof,
	}
}

func NewLocalPlainFromProof(commitment, merkleRoot []byte, merkleProof *mt.Proof) (*LocalPlain, error) {
	merkleProofJSON, err := crypto.MerkleProofMarshal(merkleProof)
	if err != nil {
		return nil, err
	}
	return NewLocalPlain(commitment, merkleRoot, merkleProofJSON), nil
}

func (l *LocalPlain) Serialize() ([]byte, error) {
	return l.Commitment, nil
}

func (l *LocalPlain) ToOnChain() *LocalOnChain {
	return NewLocalOnChain(
		hex.EncodeToString(l.Commitment),
		hex.EncodeToString(l.MerkleRoot),
		hex.EncodeToString(l.MerkleProof),
	)
}

type LocalOnChain struct {
	Commitment  string `json:"commitment"`
	MerkleRoot  string `json:"merkle_root"`
	MerkleProof string `json:"merkle_proof"`
}

func NewLocalOnChain(commitment, merkleRoot, merkleProof string) *LocalOnChain {
	return &LocalOnChain{
		Commitment:  commitment,
		MerkleRoot:  merkleRoot,
		MerkleProof: merkleProof,
	}
}

func (l *LocalOnChain) KeyVal() (string, []byte, error) {
	txJSON, err := json.Marshal(l)
	if err != nil {
		return "", nil, err
	}
	sha256Func := sha256.New()
	sha256Func.Write(txJSON)
	return hex.EncodeToString(sha256Func.Sum(nil)), txJSON, nil
}

func (l *LocalOnChain) ToPlain() (*LocalPlain, error) {
	commitment, err := hex.DecodeString(l.Commitment)
	if err != nil {
		return nil, err
	}
	merkleRoot, err := hex.DecodeString(l.MerkleRoot)
	if err != nil {
		return nil, err
	}
	merkleProof, err := hex.DecodeString(l.MerkleProof)
	if err != nil {
		return nil, err
	}
	return NewLocalPlain(commitment, merkleRoot, merkleProof), nil
}
