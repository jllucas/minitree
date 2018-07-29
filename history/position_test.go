package history

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLeftNode(t *testing.T) {
	tests := []struct {
		node         position
		expectedNode position
	}{
		{position{0, 1}, position{0, 0}},
		{position{4, 2}, position{4, 1}},
		{position{0, 3}, position{0, 2}},
		{position{8, 3}, position{8, 2}},
	}

	for _, test := range tests {
		leftNode := getLeftNode(test.node)
		assert.Equalf(t, test.expectedNode, leftNode, "Nodes don't match")
	}

}

func TestGetRigthNode(t *testing.T) {
	tests := []struct {
		node         position
		expectedNode position
	}{
		{position{0, 1}, position{1, 0}},
		{position{4, 2}, position{6, 1}},
		{position{0, 3}, position{4, 2}},
		{position{0, 4}, position{8, 3}},
	}

	for _, test := range tests {
		leftNode := getRightNode(test.node)
		assert.Equalf(t, test.expectedNode, leftNode, "Nodes don't match")
	}

}
