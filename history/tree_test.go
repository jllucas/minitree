package history

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Best practices from https://medium.com/@sebdah/go-best-practices-testing-3448165a0e18
func TestAdd(t *testing.T) {

	tests := map[string]struct {
		event            []byte
		expectedRootHash hash
	}{
		"Add 0": {[]byte(strconv.Itoa((0))), []byte{0xf1, 0x53, 0x43, 0x92, 0x27, 0x9b, 0xdd, 0xbf, 0x9d, 0x43, 0xdd, 0xe8, 0x70, 0x1c, 0xb5, 0xbe, 0x14, 0xb8, 0x2f, 0x76, 0xec, 0x66, 0x7, 0xbf, 0x8d, 0x6a, 0xd5, 0x57, 0xf6, 0xf, 0x30, 0x4e}},
		"Add 1": {[]byte(strconv.Itoa((1))), []byte{0xfd, 0x9a, 0xf3, 0xd2, 0x89, 0xe1, 0xc7, 0xb2, 0x2c, 0x64, 0x8a, 0x4f, 0x31, 0x63, 0x81, 0xe7, 0xc, 0x63, 0xf8, 0x29, 0x4d, 0x3d, 0xc5, 0x7d, 0xbe, 0xc8, 0x2d, 0x32, 0x92, 0x40, 0x91, 0x5}},
	}

	tree := NewTree()
	for testname, test := range tests {
		rootHash := tree.Add(test.event)
		assert.Equalf(t, test.expectedRootHash, rootHash, "Hashes don't match in test: %s", testname)
	}

}

func TestGenerateMembershipProof(t *testing.T) {
	testsOK := map[string]struct {
		index, version int
		expectedTree   historyTree
	}{
		"Index 2 in version 2": {2, 2, historyTree{version: 2,
			content: map[position]hash{
				{index: 0, layer: 1}: {},
				{index: 2, layer: 0}: {},
				{index: 3, layer: 0}: {},
			},
		},
		},
		"Index 2 in version 4": {2, 4, historyTree{version: 4,
			content: map[position]hash{
				{index: 0, layer: 1}: {},
				{index: 2, layer: 0}: {},
				{index: 3, layer: 0}: {},
				{index: 4, layer: 2}: {},
			},
		},
		}, // Index > tree version.
		"Index greater than version": {100, 2, historyTree{version: -1,
			content: map[position]hash{},
		},
		}, // Version > tree version.
		"Version greater than tree version": {2, 15, historyTree{version: -1,
			content: map[position]hash{},
		},
		},
	}

	tree := NewTree()
	// Add events.
	for i := 0; i <= 7; i++ {
		tree.Add([]byte(strconv.Itoa((i))))
	}

	// Successful test cases.
	var commitment hash // Unused
	for testname, test := range testsOK {
		membProof := tree.GenerateMembershipProof(test.index, commitment, test.version)
		assert.Equalf(t, test.expectedTree.version, membProof.version, "Versions don't match in test: %s", testname)
		assert.Equalf(t, len(test.expectedTree.content), len(membProof.content), "Number of nodes don't match in test: %s", testname)
		for node := range membProof.content {
			assert.Contains(t, test.expectedTree.content, node, "Node not found in expected tree in test: %s.", testname)
		}
	}

}

func TestGenerateIncrementalProof(t *testing.T) {

	testsOK := map[string]struct {
		indexI, indexJ int
		expectedTree   historyTree
	}{
		"Inc. proof between 2 and 3": {2, 3, historyTree{version: 3,
			content: map[position]hash{
				{index: 0, layer: 1}: {},
				{index: 2, layer: 0}: {},
				{index: 3, layer: 0}: {},
			},
		},
		},
		"Inc. proof between 2 and 4": {2, 4, historyTree{version: 4,
			content: map[position]hash{
				{index: 0, layer: 1}: {},
				{index: 2, layer: 0}: {},
				{index: 3, layer: 0}: {},
				{index: 4, layer: 0}: {},
				{index: 5, layer: 0}: {},
				{index: 6, layer: 1}: {},
			},
		},
		},
		"Inc. proof between 3 and 7": {3, 7, historyTree{version: 7,
			content: map[position]hash{
				{index: 0, layer: 1}: {},
				{index: 2, layer: 0}: {},
				{index: 3, layer: 0}: {},
				{index: 4, layer: 1}: {},
				{index: 6, layer: 0}: {},
				{index: 7, layer: 0}: {},
			},
		},
		},
		"Inc. proof when i > j": {100, 2, historyTree{version: -1,
			content: map[position]hash{},
		},
		},
		"Inc. proof when j > tree version": {2, 15, historyTree{version: -1,
			content: map[position]hash{},
		},
		},
	}

	testsKO := map[string]struct {
		indexI, indexJ int
		expectedTree   historyTree
	}{
		"Wrong inc. proof between 2 and 4": {2, 4, historyTree{version: 1000,
			content: map[position]hash{
				{index: 100, layer: 0}: {},
				{index: 200, layer: 0}: {},
				{index: 300, layer: 0}: {},
				{index: 400, layer: 0}: {},
				{index: 500, layer: 1}: {},
			},
		},
		},
	}

	tree := NewTree()
	// Add events.
	for i := 0; i <= 7; i++ {
		tree.Add([]byte(strconv.Itoa((i))))
	}

	// Successful test cases.
	for testname, test := range testsOK {
		incrProof := tree.GenerateIncrementalProof(test.indexI, test.indexJ)
		assert.Equalf(t, test.expectedTree.version, incrProof.version, "Versions don't match in test: %s", testname)
		assert.Equalf(t, len(test.expectedTree.content), len(incrProof.content), "Number of nodes don't match in test: %s", testname)
		for node := range incrProof.content {
			assert.Contains(t, test.expectedTree.content, node, "Node not found in expected tree in test: %s", testname)
		}
	}

	// Failure test cases.
	for _, test := range testsKO {
		incrProof := tree.GenerateIncrementalProof(test.indexI, test.indexJ)
		assert.NotEqualf(t, test.expectedTree.version, incrProof.version, "Versions match, but it shouldn't.")
		assert.NotEqualf(t, len(test.expectedTree.content), len(incrProof.content), "Number of nodes match, but it shouldn't.")
		for node := range incrProof.content {
			assert.NotContains(t, test.expectedTree.content, node, "Node found, but it was not expected.")
		}
	}
}
