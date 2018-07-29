package history

import (
	"crypto/sha256"
	"fmt"
	"math"

	"github.com/kr/pretty"
)

type hash []byte

type historyTree struct {
	version int
	content map[position]hash
}

/* "NewTree:Create new history tree." */
func NewTree() historyTree {
	t := historyTree{
		-1,
		make(map[position]hash),
	}
	return t
}

func (t historyTree) Prettyfy() {
	fmt.Printf("%# v \n", pretty.Formatter(t))
}

func computeHashLeaf(event []byte) hash {
	sha := sha256.New()
	left := []byte{'0'}
	full := append(left, event...)
	sha.Write(full)
	return sha.Sum(nil)
}

func computeHashInterior(event []byte, right []byte) hash {
	sha := sha256.New()
	left := []byte{'1'}
	partial := append(left, event...)
	full := append(partial, right...)
	sha.Write(full)
	return sha.Sum(nil)
}

func (t historyTree) getHash(pos position) hash {
	return t.content[pos]
}

func computeDepth(version int) int {
	return int(math.Ceil(math.Log2(float64(version + 1))))
}

func computeNodeValue(pos position) float64 {
	return float64(pos.index) + math.Exp2(float64(pos.layer-1))
}

func (t *historyTree) updatePostOrder(depth int, rootNode position) {
	var eventHash hash
	currentNode := rootNode

	if depth > 0 {
		currentNodeValue := computeNodeValue(currentNode)
		leftNode := getLeftNode(currentNode)
		rightNode := getRightNode(currentNode)

		if float64(t.version) >= currentNodeValue {
			t.updatePostOrder(depth-1, rightNode)
		} else {
			t.updatePostOrder(depth-1, leftNode)
		}

		hashLeft := t.getHash(leftNode)
		hashRight := t.getHash(rightNode)
		eventHash = computeHashInterior(hashLeft, hashRight)
		t.content[currentNode] = eventHash
		//		fmt.Println(currentNode)
	}
}

// Add event to the history tree.
func (t *historyTree) Add(event []byte) hash {
	t.version++
	// Add event hash
	pos := position{t.version, 0}
	eventHash := computeHashLeaf(event)
	t.content[pos] = eventHash
	//	fmt.Printf("%x\n", eventHash)
	//	fmt.Println(pos)

	// Compute root node and update tree in post-order
	depth := computeDepth(t.version)
	root := position{0, depth}
	t.updatePostOrder(depth, root)
	return t.content[root]
}

func (t historyTree) navigatePostOrder(depth int, currentNode, leaf position, proofTree historyTree) {

	if depth > 0 {
		currentNodeValue := computeNodeValue(currentNode)
		leftNode := getLeftNode(currentNode)
		rightNode := getRightNode(currentNode)

		if leaf.index >= int(currentNodeValue) {
			t.navigatePostOrder(depth-1, rightNode, leaf, proofTree)
			proofTree.content[leftNode] = t.getHash(leftNode)
		} else {
			t.navigatePostOrder(depth-1, leftNode, leaf, proofTree)
			if rightNode.index > proofTree.version {
				proofTree.content[rightNode] = hash{}
			} else {
				proofTree.content[rightNode] = t.getHash(rightNode)
			}
		}
	} else {
		proofTree.content[currentNode] = t.getHash(currentNode)
	}
}

func (t historyTree) GenerateMembershipProof(index int, commitment hash, version int) historyTree {
	if (version < index) || (version > t.version) {
		return historyTree{-1, map[position]hash{}}
	}

	proofTree := historyTree{
		version: version,
		content: make(map[position]hash),
	}

	depth := computeDepth(version)
	root := position{0, depth}
	leaf := position{index, 0}
	t.navigatePostOrder(depth, root, leaf, proofTree)

	return proofTree
}

func getCommonRoot(indexI, indexJ int, root position) position {
	rootValue := int(computeNodeValue(root))
	if (indexI >= rootValue) && (indexJ >= rootValue) {
		return getCommonRoot(indexI, indexJ, getRightNode(root))
	} else if (indexI < rootValue) && (indexJ < rootValue) {
		return getCommonRoot(indexI, indexJ, getLeftNode(root))
	} else {
		return root
	}
}

func (t *historyTree) cleanProof(indexI, indexJ int) {
	depth := computeDepth(indexJ)
	commonRoot := getCommonRoot(indexI, indexJ, position{0, depth})
	// Delete unnecessary nodes.
	if commonRoot.layer > 1 {
		leftNode := getLeftNode(commonRoot)
		rightNode := getRightNode(commonRoot)
		delete(t.content, leftNode)
		delete(t.content, rightNode)
	}
}

func (t historyTree) GenerateIncrementalProof(indexI, indexJ int) historyTree {
	if (indexJ < indexI) || (indexJ > t.version) {
		return historyTree{-1, map[position]hash{}}
	}

	mProofI := t.GenerateMembershipProof(indexI, hash{}, indexJ)
	mProofJ := t.GenerateMembershipProof(indexJ, hash{}, indexJ)

	for pos, hash := range mProofJ.content {
		mProofI.content[pos] = hash
	}
	mProofI.cleanProof(indexI, indexJ)
	return mProofI
}
