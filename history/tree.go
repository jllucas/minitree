package history

import (
	"crypto/sha256"
	"fmt"
	"math"

	"github.com/jllucas/minitree/common"
	"github.com/kr/pretty"
)

type HistoryTree struct {
	version int
	Store   common.HistoryMemoryStore
}

func NewHistoryTree() HistoryTree {
	return HistoryTree{-1, common.NewHistoryMemoryStore()}
}

func (t HistoryTree) Prettyfy() {
	fmt.Printf("%# v \n", pretty.Formatter(t))
}

func computeHashLeaf(event []byte) common.Hash {
	sha := sha256.New()
	left := []byte{'0'}
	full := append(left, event...)
	sha.Write(full)
	return sha.Sum(nil)
}

func computeHashInterior(event []byte, right []byte) common.Hash {
	sha := sha256.New()
	left := []byte{'1'}
	partial := append(left, event...)
	full := append(partial, right...)
	sha.Write(full)
	return sha.Sum(nil)
}

func (t HistoryTree) getHash(pos common.Position) common.Hash {
	return t.Store.Get(pos)
}

func computeDepth(version int) int {
	return int(math.Ceil(math.Log2(float64(version + 1))))
}

func (t *HistoryTree) updatePostOrder(depth int, rootNode common.Position) {
	currentNode := rootNode

	if depth > 0 {
		currentNodeValue := currentNode.ComputeNodeValue()
		leftNode := currentNode.GetLeftNode()
		rightNode := currentNode.GetRightNode()

		if float64(t.version) >= currentNodeValue {
			t.updatePostOrder(depth-1, rightNode)
		} else {
			t.updatePostOrder(depth-1, leftNode)
		}

		hashLeft := t.getHash(leftNode)
		hashRight := t.getHash(rightNode)
		t.Store[currentNode] = computeHashInterior(hashLeft, hashRight)
	}
}

// Add event to the history tree.
func (t *HistoryTree) Add(event []byte) common.Hash {
	t.version++
	// Add event hash
	pos := common.NewPosition(t.version, 0)
	t.Store[pos] = computeHashLeaf(event)

	// Compute root node and update tree in post-order
	depth := computeDepth(t.version)
	root := common.NewPosition(0, depth)
	t.updatePostOrder(depth, root)
	return t.Store.Get(root)
}

func (t HistoryTree) navigatePostOrder(depth int, currentNode, leaf common.Position, proofTree HistoryTree) {

	if depth > 0 {
		currentNodeValue := currentNode.ComputeNodeValue()
		leftNode := currentNode.GetLeftNode()
		rightNode := currentNode.GetRightNode()

		if leaf.Index >= int(currentNodeValue) {
			t.navigatePostOrder(depth-1, rightNode, leaf, proofTree)
			proofTree.Store[leftNode] = t.getHash(leftNode)
		} else {
			t.navigatePostOrder(depth-1, leftNode, leaf, proofTree)
			if rightNode.Index > proofTree.version {
				proofTree.Store[rightNode] = common.Hash{}
			} else {
				proofTree.Store[rightNode] = t.getHash(rightNode)
			}
		}
	} else {
		proofTree.Store[currentNode] = t.getHash(currentNode)
	}
}

func (t HistoryTree) GenerateMembershipProof(index int, commitment common.Hash, version int) HistoryTree {
	if (version < index) || (version > t.version) {
		return HistoryTree{-1, common.NewHistoryMemoryStore()}
	}

	proofTree := HistoryTree{
		version: version,
		Store:   common.NewHistoryMemoryStore(),
	}

	depth := computeDepth(version)
	root := common.NewPosition(0, depth)
	leaf := common.NewPosition(index, 0)
	t.navigatePostOrder(depth, root, leaf, proofTree)

	return proofTree
}

func getCommonRoot(indexI, indexJ int, root common.Position) common.Position {
	rootValue := int(root.ComputeNodeValue())
	if (indexI >= rootValue) && (indexJ >= rootValue) {
		return getCommonRoot(indexI, indexJ, root.GetRightNode())
	} else if (indexI < rootValue) && (indexJ < rootValue) {
		return getCommonRoot(indexI, indexJ, root.GetLeftNode())
	} else {
		return root
	}
}

func (t *HistoryTree) cleanProof(indexI, indexJ int) {
	depth := computeDepth(indexJ)
	commonRoot := getCommonRoot(indexI, indexJ, common.Position{0, depth})
	// Delete unnecessary nodes.
	if commonRoot.Layer > 1 {
		leftNode := commonRoot.GetLeftNode()
		rightNode := commonRoot.GetRightNode()
		delete(t.Store, leftNode)
		delete(t.Store, rightNode)
	}
}

func (t HistoryTree) GenerateIncrementalProof(indexI, indexJ int) HistoryTree {
	if (indexJ < indexI) || (indexJ > t.version) {
		return HistoryTree{-1, common.NewHistoryMemoryStore()}
	}

	mProofI := t.GenerateMembershipProof(indexI, common.Hash{}, indexJ)
	mProofJ := t.GenerateMembershipProof(indexJ, common.Hash{}, indexJ)

	for pos, hash := range mProofJ.Store {
		mProofI.Store[pos] = hash
	}
	mProofI.cleanProof(indexI, indexJ)
	return mProofI
}
