package history

import (
	"bytes"
	"fmt"
)

// Recursive search will never reach depth 0.
func (t historyTree) computeHashPostOrder(depth int, rootNode position) hash {
	var leftHash, rightHash hash

	leftNode := getLeftNode(rootNode)
	if value, ok := t.content[leftNode]; ok {
		leftHash = value
	} else {
		leftHash = t.computeHashPostOrder(depth-1, leftNode)
	}

	rightNode := getRightNode(rootNode)
	if value, ok := t.content[rightNode]; ok {
		rightHash = value
	} else {
		rightHash = t.computeHashPostOrder(depth-1, rightNode)
	}

	return computeHashInterior(leftHash, rightHash)
}

// Not implemented: comparison between eventHash parameter and Xi.
func (t historyTree) VerifyMembershipProof(index int, rootHash hash, eventHash hash) bool {
	depth := computeDepth(t.version)
	root := position{0, depth}
	computedRootHash := t.computeHashPostOrder(depth, root)
	fmt.Printf("\nComputed root hash: %x", computedRootHash)
	fmt.Printf("\nRoot hash: %x", rootHash)
	return bytes.Equal(computedRootHash, rootHash)
}

// indexI, indexJ are useless if we can guess both index from their commitments.
func (t historyTree) VerifyIncrementalProof(commitmentI, commitmentJ hash, indexI, indexJ int) bool {
	prunedTree := t.GenerateMembershipProof(indexI, commitmentI, indexI)
	prunedTree.Prettyfy()
	Iverified := prunedTree.VerifyMembershipProof(indexI, commitmentI, hash{})
	Jverified := t.VerifyMembershipProof(indexJ, commitmentJ, hash{})
	return (Iverified && Jverified)
}
