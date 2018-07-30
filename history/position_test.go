package history

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLeftNode(t *testing.T) {
	tests := map[string]struct {
		node         position
		expectedNode position
	}{
		"Left of (0,1)": {position{0, 1}, position{0, 0}},
		"Left of (4,2)": {position{4, 2}, position{4, 1}},
		"Left of (0,3)": {position{0, 3}, position{0, 2}},
		"Left of (8,3)": {position{8, 3}, position{8, 2}},
	}

	for testname, test := range tests {
		leftNode := getLeftNode(test.node)
		assert.Equalf(t, test.expectedNode, leftNode, "Nodes don't match in test: %s", testname)
	}

}

func TestGetRigthNode(t *testing.T) {
	tests := map[string]struct {
		node         position
		expectedNode position
	}{
		"Right of (0,1)": {position{0, 1}, position{1, 0}},
		"Right of (4,2)": {position{4, 2}, position{6, 1}},
		"Right of (0,3)": {position{0, 3}, position{4, 2}},
		"Right of (0,4)": {position{0, 4}, position{8, 3}},
	}

	for testname, test := range tests {
		leftNode := getRightNode(test.node)
		assert.Equalf(t, test.expectedNode, leftNode, "Nodes don't match in test: %s", testname)
	}

}
