package badger

import (
	"strconv"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/stretchr/testify/assert"
)

// ADD ERROR COVERAGE
func TestInsert(t *testing.T) {
	tests := []struct {
		testname string
		key      []byte
		value    []byte
		err      error
	}{
		{"Insert 0", []byte(strconv.Itoa(0)), []byte(strconv.Itoa(0)), nil},
		{"Insert empty key", []byte{}, []byte(strconv.Itoa(1)), badger.ErrEmptyKey},
	}

	db, _ := NewBadgerStore()
	defer db.Close()
	for _, test := range tests {
		err := db.Insert(test.key, test.value)
		assert.Equal(t, test.err, err, "Insert error in test: %s", test.testname)
	}
}

// ADD ERROR COVERAGE
func TestGet(t *testing.T) {
	tests := []struct {
		testname string
		key      []byte
		value    []byte
		err      error
	}{
		{"Get 0", []byte(strconv.Itoa(0)), []byte(strconv.Itoa(0)), nil},
		{"Get 1", []byte(strconv.Itoa(1)), []byte(strconv.Itoa(1)), nil},
		{"Get 2", []byte(strconv.Itoa(2)), []byte(strconv.Itoa(2)), nil},
		{"Get 3", []byte(strconv.Itoa(3)), []byte(strconv.Itoa(3)), nil},
	}

	db, _ := NewBadgerStore()
	defer db.Close()
	for _, test := range tests {
		db.Insert(test.key, test.value)
		result, err := db.Get(test.key)
		assert.Equalf(t, test.value, result, "Values don't match in test: %s", test.testname)
		assert.Equalf(t, test.err, err, "Errors don't match in test: %s", test.testname)
	}
}

func TestGetRange(t *testing.T) {

}
