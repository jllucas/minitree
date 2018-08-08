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

	db, _ := NewBadger()
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

	db, _ := NewBadger()
	defer db.Close()
	for _, test := range tests {
		db.Insert(test.key, test.value)
		result, err := db.Get(test.key)
		assert.Equalf(t, test.value, result, "Values don't match in test: %s", test.testname)
		assert.Equalf(t, test.err, err, "Errors don't match in test: %s", test.testname)
	}
}

func TestGetRange(t *testing.T) {
	tests := []struct {
		testname      string
		start         []byte
		end           []byte
		expectedRange KVRange
	}{
		{"Get range 0-5", []byte(strconv.Itoa(0)), []byte(strconv.Itoa(5)), KVRange{}},
		{"Get range 100-200", []byte(strconv.Itoa(100)), []byte(strconv.Itoa(200)), KVRange{}},
	}

	db, _ := NewBadger()
	defer db.Close()
	for i := 0; i < 1e3; i++ {
		db.Insert([]byte(strconv.Itoa(i)), []byte(strconv.Itoa(i)))
	}

	for _, test := range tests {
		result, _ := db.GetRange(test.start, test.end)
		assert.Equalf(t, test.expectedRange, result, "Ranges don't match in test: %s", test.testname)
	}
}
