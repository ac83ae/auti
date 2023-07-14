package crypto

import (
	"bytes"
	"container/heap"
	"encoding/json"
	"errors"

	mt "github.com/txaty/go-merkletree"
)

type ProofNode struct {
	Coordinate [2]int
	Data       []byte
}

// MerkleBatchProof is the batched proof for a set of data blocks.
type MerkleBatchProof struct {
	Nodes   []ProofNode `json:"nodes"`
	Indexes []int       `json:"indexes"`
}

func MerkleBatchProofMarshal(proof *MerkleBatchProof) ([]byte, error) {
	return json.Marshal(proof)
}

func MerkleBatchProofUnmarshal(data []byte) (*MerkleBatchProof, error) {
	var proof MerkleBatchProof
	if err := json.Unmarshal(data, &proof); err != nil {
		return nil, err
	}
	return &proof, nil
}

type proofNodePQ []ProofNode

func (p *proofNodePQ) Len() int {
	return len(*p)
}

func (p *proofNodePQ) Less(i, j int) bool {
	if (*p)[i].Coordinate[0] == (*p)[j].Coordinate[0] {
		return (*p)[i].Coordinate[1] < (*p)[j].Coordinate[1]
	}
	return (*p)[i].Coordinate[0] < (*p)[j].Coordinate[0]
}

func (p *proofNodePQ) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

func (p *proofNodePQ) Push(x interface{}) {
	*p = append(*p, x.(ProofNode))
}

func (p *proofNodePQ) Pop() interface{} {
	oldPQ := *p
	n := len(oldPQ)
	x := oldPQ[n-1]
	*p = oldPQ[0 : n-1]
	return x
}

// NewMerkleBatchProof creates a batched proof for the given data blocks and their individual proofs.
func NewMerkleBatchProof(dataBlocks []mt.DataBlock, proofs []*mt.Proof) (*MerkleBatchProof, error) {
	if len(dataBlocks) != len(proofs) {
		return nil, errors.New("number of data blocks and proofs must be equal")
	}
	if len(dataBlocks) == 0 {
		return nil, errors.New("no data blocks or proofs given")
	}
	treeDepth := len(proofs[0].Siblings)
	maxNodeIdx := uint32(1<<uint(treeDepth)) - 1
	necessaryProofCoordinateMap := make(map[[2]int]bool)
	nodeBuffer := make([][][]byte, treeDepth+1)
	maxNumBlocks := 1 << uint(treeDepth)
	for i := 0; i < treeDepth+1; i++ {
		nodeBuffer[i] = make([][]byte, maxNumBlocks)
		maxNumBlocks = (maxNumBlocks + 1) / 2
	}
	batchProofs := new(MerkleBatchProof)
	batchProofs.Indexes = make([]int, len(dataBlocks))
	for i := 0; i < len(dataBlocks); i++ {
		blockBytes, err := dataBlocks[i].Serialize()
		if err != nil {
			return nil, err
		}
		blockHash, _ := hashFunc(blockBytes)
		nodeIdx := maxNodeIdx - proofs[i].Path
		batchProofs.Indexes[i] = int(nodeIdx)
		nodeBuffer[0][nodeIdx] = blockHash
		if _, ok := necessaryProofCoordinateMap[[2]int{0, int(nodeIdx)}]; ok {
			delete(necessaryProofCoordinateMap, [2]int{0, int(nodeIdx)})
		}
		pathNode := make([]byte, len(blockHash))
		copy(pathNode, blockHash)
		path := proofs[i].Path
		for j := 0; j < treeDepth; j++ {
			direction := path&1 == 1
			path >>= 1
			proofIdx := nodeIdx - 1
			if direction {
				proofIdx = nodeIdx + 1
			}
			if nodeBuffer[j][proofIdx] == nil {
				necessaryProofCoordinateMap[[2]int{j, int(proofIdx)}] = true
				nodeBuffer[j][proofIdx] = proofs[i].Siblings[j]
			}
			nodeIdx >>= 1
			if nodeBuffer[j+1][nodeIdx] == nil {
				if direction {
					pathNode, _ = hashFunc(append(pathNode, proofs[i].Siblings[j]...))
				} else {
					pathNode, _ = hashFunc(append(proofs[i].Siblings[j], pathNode...))
				}
				nodeBuffer[j+1][nodeIdx] = pathNode
			} else {
				pathNode = make([]byte, len(nodeBuffer[j+1][nodeIdx]))
				copy(pathNode, nodeBuffer[j+1][nodeIdx])
				if _, ok := necessaryProofCoordinateMap[[2]int{j + 1, int(nodeIdx)}]; ok {
					delete(necessaryProofCoordinateMap, [2]int{j + 1, int(nodeIdx)})
				} else {
					break
				}
			}
		}
	}

	qp := new(proofNodePQ)
	for coordinate := range necessaryProofCoordinateMap {
		heap.Push(qp, ProofNode{
			Coordinate: coordinate,
			Data:       nodeBuffer[coordinate[0]][coordinate[1]],
		})
	}
	batchProofs.Nodes = make([]ProofNode, 0, len(*qp))
	for len(*qp) > 0 {
		batchProofs.Nodes = append(batchProofs.Nodes, heap.Pop(qp).(ProofNode))
	}
	return batchProofs, nil
}

// VerifyMerkleBatchProof verifies the batched proof against the given root hash.
func VerifyMerkleBatchProof(dataBlocks []mt.DataBlock, batchProof *MerkleBatchProof, root []byte) (bool, error) {
	if len(dataBlocks) != len(batchProof.Indexes) {
		return false, errors.New("number of data blocks and proof data  must be equal")
	}
	if len(dataBlocks) == 0 {
		return false, errors.New("no data blocks given")
	}
	pq := new(proofNodePQ)
	for _, proofNode := range batchProof.Nodes {
		heap.Push(pq, proofNode)
	}
	for idx, block := range dataBlocks {
		blockBytes, err := block.Serialize()
		if err != nil {
			return false, err
		}
		//blockHash, _ := hashFunc(blockBytes)
		heap.Push(pq, ProofNode{
			Coordinate: [2]int{0, batchProof.Indexes[idx]},
			Data:       blockBytes,
		})
	}
	for len(*pq) > 1 {
		node1 := heap.Pop(pq).(ProofNode)
		node2 := heap.Pop(pq).(ProofNode)
		if node1.Coordinate[0] != node2.Coordinate[0] {
			return false, errors.New("invalid proof")
		}
		hash, err := hashFunc(append(node1.Data, node2.Data...))
		if err != nil {
			return false, err
		}
		heap.Push(pq, ProofNode{
			Coordinate: [2]int{node1.Coordinate[0] + 1, node1.Coordinate[1] / 2},
			Data:       hash,
		})
	}
	if len(*pq) == 0 {
		return false, errors.New("invalid proof")
	}
	node := heap.Pop(pq).(ProofNode)
	return bytes.Equal(node.Data, root), nil
}
