package history

import (
	"strconv"
	"testing"

	"github.com/jllucas/minitree/common"
	"github.com/stretchr/testify/assert"
)

// Best practices from https://medium.com/@sebdah/go-best-practices-testing-3448165a0e18
func TestAdd(t *testing.T) {

	tests := map[string]struct {
		event            []byte
		expectedRootHash common.Hash
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
			Store: common.Store{
				{Index: 0, Layer: 1}: {},
				{Index: 2, Layer: 0}: {},
				{Index: 3, Layer: 0}: {},
			},
		},
		},
		"Index 2 in version 4": {2, 4, historyTree{version: 4,
			Store: common.Store{
				{Index: 0, Layer: 1}: {},
				{Index: 2, Layer: 0}: {},
				{Index: 3, Layer: 0}: {},
				{Index: 4, Layer: 2}: {},
			},
		},
		}, // Index > tree version.
		"Index greater than version": {100, 2, historyTree{version: -1,
			Store: common.NewStore(),
		},
		}, // Version > tree version.
		"Version greater than tree version": {2, 15, historyTree{version: -1,
			Store: common.NewStore(),
		},
		},
	}

	tree := NewTree()
	// Add events.
	for i := 0; i <= 7; i++ {
		tree.Add([]byte(strconv.Itoa((i))))
	}

	// Successful test cases.
	var commitment common.Hash // Unused
	for testname, test := range testsOK {
		membProof := tree.GenerateMembershipProof(test.index, commitment, test.version)
		assert.Equalf(t, test.expectedTree.version, membProof.version, "Versions don't match in test: %s", testname)
		assert.Equalf(t, len(test.expectedTree.Store), len(membProof.Store), "Number of nodes don't match in test: %s", testname)
		for node := range membProof.Store {
			assert.Contains(t, test.expectedTree.Store, node, "Node not found in expected tree in test: %s.", testname)
		}
	}

}

func TestGenerateIncrementalProof(t *testing.T) {

	testsOK := map[string]struct {
		indexI, indexJ int
		expectedTree   historyTree
	}{
		"Inc. proof between 2 and 3": {2, 3, historyTree{version: 3,
			Store: common.Store{
				{Index: 0, Layer: 1}: {},
				{Index: 2, Layer: 0}: {},
				{Index: 3, Layer: 0}: {},
			},
		},
		},
		"Inc. proof between 2 and 4": {2, 4, historyTree{version: 4,
			Store: common.Store{
				{Index: 0, Layer: 1}: {},
				{Index: 2, Layer: 0}: {},
				{Index: 3, Layer: 0}: {},
				{Index: 4, Layer: 0}: {},
				{Index: 5, Layer: 0}: {},
				{Index: 6, Layer: 1}: {},
			},
		},
		},
		"Inc. proof between 3 and 7": {3, 7, historyTree{version: 7,
			Store: common.Store{
				{Index: 0, Layer: 1}: {},
				{Index: 2, Layer: 0}: {},
				{Index: 3, Layer: 0}: {},
				{Index: 4, Layer: 1}: {},
				{Index: 6, Layer: 0}: {},
				{Index: 7, Layer: 0}: {},
			},
		},
		},
		"Inc. proof when i > j": {100, 2, historyTree{version: -1,
			Store: common.NewStore(),
		},
		},
		"Inc. proof when j > tree version": {2, 15, historyTree{version: -1,
			Store: common.NewStore(),
		},
		},
	}

	testsKO := map[string]struct {
		indexI, indexJ int
		expectedTree   historyTree
	}{
		"Wrong inc. proof between 2 and 4": {2, 4, historyTree{version: 1000,
			Store: common.Store{
				{Index: 100, Layer: 0}: {},
				{Index: 200, Layer: 0}: {},
				{Index: 300, Layer: 0}: {},
				{Index: 400, Layer: 0}: {},
				{Index: 500, Layer: 1}: {},
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
		assert.Equalf(t, len(test.expectedTree.Store), len(incrProof.Store), "Number of nodes don't match in test: %s", testname)
		for node := range incrProof.Store {
			assert.Contains(t, test.expectedTree.Store, node, "Node not found in expected tree in test: %s", testname)
		}
	}

	// Failure test cases.
	for _, test := range testsKO {
		incrProof := tree.GenerateIncrementalProof(test.indexI, test.indexJ)
		assert.NotEqualf(t, test.expectedTree.version, incrProof.version, "Versions match, but it shouldn't.")
		assert.NotEqualf(t, len(test.expectedTree.Store), len(incrProof.Store), "Number of nodes match, but it shouldn't.")
		for node := range incrProof.Store {
			assert.NotContains(t, test.expectedTree.Store, node, "Node found, but it was not expected.")
		}
	}
}
