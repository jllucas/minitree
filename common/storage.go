package common

type storage interface {
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
