package history

import (
	"bytes"

	"github.com/jllucas/minitree/common"
)

// Recursive search will never reach depth 0.
func (t HistoryTree) computeHashPostOrder(depth int, rootNode common.Position) common.Hash {
	var leftHash, rightHash common.Hash

	leftNode := rootNode.GetLeftNode()
	if value, ok := t.Store[leftNode]; ok {
		leftHash = value
	} else {
		leftHash = t.computeHashPostOrder(depth-1, leftNode)
	}

	rightNode := rootNode.GetRightNode()
	if value, ok := t.Store[rightNode]; ok {
		rightHash = value
	} else {
		rightHash = t.computeHashPostOrder(depth-1, rightNode)
	}

	return computeHashInterior(leftHash, rightHash)
}

// Not implemented: comparison between eventHash parameter and Xi.
func (t HistoryTree) VerifyMembershipProof(index int, rootHash common.Hash, eventHash common.Hash) bool {
	depth := computeDepth(t.version)
	root := common.NewPosition(0, depth)
	computedRootHash := t.computeHashPostOrder(depth, root)
	return bytes.Equal(computedRootHash, rootHash)
}

// indexI, indexJ are useless if we can guess both index from their commitments.
func (t HistoryTree) VerifyIncrementalProof(commitmentI, commitmentJ common.Hash, indexI, indexJ int) bool {
	prunedTree := t.GenerateMembershipProof(indexI, commitmentI, indexI)
	Iverified := prunedTree.VerifyMembershipProof(indexI, commitmentI, common.Hash{})
	Jverified := t.VerifyMembershipProof(indexJ, commitmentJ, common.Hash{})
	return (Iverified && Jverified)
}
