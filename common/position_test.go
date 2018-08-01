package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLeftNode(t *testing.T) {
	tests := map[string]struct {
		node         Position
		expectedNode Position
	}{
		"Left of (0,1)": {NewPosition(0, 1), NewPosition(0, 0)},
		"Left of (4,2)": {NewPosition(4, 2), NewPosition(4, 1)},
		"Left of (0,3)": {NewPosition(0, 3), NewPosition(0, 2)},
		"Left of (8,3)": {NewPosition(8, 3), NewPosition(8, 2)},
	}

	for testname, test := range tests {
		leftNode := test.node.GetLeftNode()
		assert.Equalf(t, test.expectedNode, leftNode, "Nodes don't match in test: %s", testname)
	}

}

func TestGetRigthNode(t *testing.T) {
	tests := map[string]struct {
		node         Position
		expectedNode Position
	}{
		"Right of (0,1)": {NewPosition(0, 1), NewPosition(1, 0)},
		"Right of (4,2)": {NewPosition(4, 2), NewPosition(6, 1)},
		"Right of (0,3)": {NewPosition(0, 3), NewPosition(4, 2)},
		"Right of (0,4)": {NewPosition(0, 4), NewPosition(8, 3)},
	}

	for testname, test := range tests {
		leftNode := test.node.GetRightNode()
		assert.Equalf(t, test.expectedNode, leftNode, "Nodes don't match in test: %s", testname)
	}

}

func TestComputeNodeValue(t *testing.T) {
	tests := map[string]struct {
		node          Position
		expectedValue float64
	}{
		"Value of (0,1)": {NewPosition(0, 1), 1},
		"Value of (4,2)": {NewPosition(4, 2), 6},
		"Value of (0,3)": {NewPosition(0, 3), 4},
		"Value of (0,4)": {NewPosition(0, 4), 8},
	}

	for testname, test := range tests {
		value := test.node.ComputeNodeValue()
		assert.Equalf(t, test.expectedValue, value, "Values don't match in test: %s", testname)
	}
}
