package common

type HistoryStore interface {
	Put(pos Position, event Hash)
	Get(pos Position) Hash
}

type HyperStore interface {
	Put(key [256]byte, index int)
	Get(key [256]byte) int
}

type HistoryMemoryStore map[Position]Hash

func NewHistoryMemoryStore() HistoryMemoryStore {
	return make(map[Position]Hash)
}

func (s HistoryMemoryStore) Get(pos Position) Hash {
	return s[pos]
}

func (s HistoryMemoryStore) Put(pos Position, event Hash) {
	s[pos] = event
}

type HyperMemoryStore map[[256]byte]int

func NewHyperMemoryStore() HyperMemoryStore {
	return make(map[[256]byte]int)
}

func (s HyperMemoryStore) Get(key [256]byte) int {
	return s[key]
}

func (s HyperMemoryStore) Put(key [256]byte, index int) {
	s[key] = index
}

/* type BadgerStore badger.Badger

func NewBadgerStore() BadgerStore {
	db, _ := badger.NewBadger()
	return db
}

func (s BadgerStore) Get(key [256]byte) int {
	value, _ := s.Get(key)
	return 1
}

func (s BadgerStore) Put(key Hash, value int) error {
	return nil
} */
