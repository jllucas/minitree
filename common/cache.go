package common

type Cache interface {
	Put()
	Get()
	Delete()
	Exist()
}

type MemoryCache map[Position]Hash

func NewCache() MemoryCache {
	return make(map[Position]Hash)
}

func (s MemoryCache) Get(pos Position) Hash {
	return s[pos]
}

func (s MemoryCache) Put(pos Position, event Hash) {
	s[pos] = event
}
