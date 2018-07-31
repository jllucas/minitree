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
		"Left of (0,1)": {Position{0, 1}, Position{0, 0}},
		"Left of (4,2)": {Position{4, 2}, Position{4, 1}},
		"Left of (0,3)": {Position{0, 3}, Position{0, 2}},
		"Left of (8,3)": {Position{8, 3}, Position{8, 2}},
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
		"Right of (0,1)": {Position{0, 1}, Position{1, 0}},
		"Right of (4,2)": {Position{4, 2}, Position{6, 1}},
		"Right of (0,3)": {Position{0, 3}, Position{4, 2}},
		"Right of (0,4)": {Position{0, 4}, Position{8, 3}},
	}

	for testname, test := range tests {
		leftNode := test.node.GetRightNode()
		assert.Equalf(t, test.expectedNode, leftNode, "Nodes don't match in test: %s", testname)
	}

}
