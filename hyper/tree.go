package hyper

import "github.com/jllucas/minitree/common"

type HyperTree struct {
	DefaultHashes []common.Hash
	Cache         common.HyperStore
	Store         common.HyperStore
	Hasher        common.Hasher
}

func NewHyperTree(depth uint8, cache, store common.HyperStore, hasher common.Hasher) HyperTree {
	htree := HyperTree{
		make([]common.Hash, depth),
		cache,
		store,
		hasher,
	}
	htree.generateDefaultHashes(depth)
	return htree
}

func (t *HyperTree) generateDefaultHashes(depth uint8) {
	hasher := common.NewSha256Hasher()
	t.DefaultHashes[0] = hasher.DoHash([]byte{0x00})
	for i := 1; i < int(depth); i++ {
		t.DefaultHashes[i] = hasher.DoHash(t.DefaultHashes[i-1], t.DefaultHashes[i-1])
	}
}

func (t HyperTree) GetDefaultHash(depth uint8) common.Hash {
	return t.DefaultHashes[depth]
}

func (t *HyperTree) Add() {

}

func (t HyperTree) GenerateMembershipProof() {

}
