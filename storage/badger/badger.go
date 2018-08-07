package badger

import (
	"github.com/dgraph-io/badger"
)

type Badger struct {
	*badger.DB
}

func NewBadgerStore() (*Badger, error) {
	opts := badger.DefaultOptions
	opts.Dir = "/tmp/badger"
	opts.ValueDir = "/tmp/badger"
	db, err := badger.Open(opts)
	return &Badger{db}, err
}

func (db Badger) Get(key []byte) ([]byte, error) {
	output := make([][]byte, 1)

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		val, err := item.Value()
		if err != nil {
			return err
		}
		output[0] = make([]byte, len(val))
		copy(output[0], val)
		return nil
	})

	return output[0], err
}

func (db Badger) Insert(key, value []byte) error {
	err := db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
	return err
}
