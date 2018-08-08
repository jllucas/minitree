package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Dummy tests for testing coverage.
func TestHistoryMemoryStore(t *testing.T) {
	tests := []struct {
		testname string
		key      Position
		value    Hash
	}{
		{"Store 0", NewPosition(0, 0), Hash{0x00}},
		{"Store 1", NewPosition(0, 1), Hash{0x01}},
	}

	store := NewHistoryMemoryStore()
	for _, test := range tests {
		store.Put(test.key, test.value)
		result := store.Get(test.key)
		assert.Equalf(t, test.value, result, "Values don't match in test: %s", test.testname)
	}
}

func TestHyperMemoryStore(t *testing.T) {
	tests := []struct {
		testname string
		key      [256]byte
		value    int
	}{
		{"Store 0", [256]byte{0x00}, 0},
		{"Store 1", [256]byte{0x01}, 1},
	}

	store := NewHyperMemoryStore()
	for _, test := range tests {
		store.Put(test.key, test.value)
		result := store.Get(test.key)
		assert.Equalf(t, test.value, result, "Values don't match in test: %s", test.testname)
	}
}
