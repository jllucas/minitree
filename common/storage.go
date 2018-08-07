package common

type Storage interface {
	Put()
	Get()
}

type Store map[Position]Hash

func NewStore() Store {
	return make(map[Position]Hash)
}

func (s Store) Get(pos Position) Hash {
	return s[pos]
}

func (s Store) Put(pos Position, event Hash) {
	s[pos] = event
}

type HyperStore map[[256]byte]int

func NewHyperStore() HyperStore {
	return make(map[[256]byte]int)
}

func (s HyperStore) Get(key Hash) int {
	var toArray [256]byte
	copy(toArray[:], key)
	return s[toArray]
}

func (s HyperStore) Put(key Hash, value int) {
	var toArray [256]byte
	copy(toArray[:], key)
	s[toArray] = value
}
