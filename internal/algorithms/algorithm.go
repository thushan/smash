package algorithms

import (
	"hash"
	fnvh "hash/fnv"

	cxHash "github.com/cespare/xxhash"
)

type Algorithm int

const (
	xxhash Algorithm = iota
	fnv128
	fnv128a
)

// HashAlgorithms Used by CLI for validating --algorithm flag
var HashAlgorithms = map[int][]string{
	0: {"xxhash"},
	1: {"fnv128"},
	2: {"fnv128a"},
}

// New Instantiates a new representation of the Hash Algorithm.
func (a Algorithm) New() hash.Hash {
	switch a {
	case xxhash:
		return cxHash.New()
	case fnv128:
		return fnvh.New128()
	case fnv128a:
		return fnvh.New128a()
	}
	return fnvh.New128a()
}

// Index Returns the index for the Hash Algorithm
func (a Algorithm) Index() int {
	return int(a)
}

// Index Returns the human-readable representation of the Hash Algorithm
func (a Algorithm) String() string {
	return HashAlgorithms[a.Index()][0]
}
