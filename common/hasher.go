package common

import "crypto/sha256"

type Hash []byte

type Hasher interface {
	DoHash(events ...[]byte) Hash
}

/*
sha256 Hasher: current hashing method.
*/
type Sha256Hasher struct{}

func NewSha256Hasher() Sha256Hasher {
	return Sha256Hasher{}
}

func (h Sha256Hasher) DoHash(events ...[]byte) Hash {
	sha := sha256.New()

	for i := 0; i < len(events); i++ {
		sha.Write(events[i])
	}

	return sha.Sum(nil)
}

/*
XOR Hasher: needed for testing.
*/
type XORHasher struct{}

func NewXORHasher() XORHasher {
	return XORHasher{}
}

// From https://golang.org/src/crypto/cipher/xor.go
func (h XORHasher) DoHash(eventA, eventB []byte) Hash {
	n := len(eventA)
	if len(eventB) < n {
		n = len(eventB)
	}
	dst := make([]byte, n)
	for i := 0; i < n; i++ {
		dst[i] = eventA[i] ^ eventB[i]
	}
	return dst
}
