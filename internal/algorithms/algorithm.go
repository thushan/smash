package algorithms

import (
	md5h "crypto/md5"
	sha256h "crypto/sha256"
	sha512h "crypto/sha512"
	"hash"
	fnvh "hash/fnv"

	"github.com/spaolacci/murmur3"

	cxHash "github.com/cespare/xxhash"
)

type Algorithm int

const (
	Xxhash Algorithm = iota
	Fnv128
	Fnv128a
	Murmur3_128
	Murmur3_64
	Murmur3_32
	Md5
	Sha256
	Sha512
)

// HashAlgorithms Used by CLI for validating --algorithm flag
var HashAlgorithms = map[int][]string{
	0: {"xxhash"},
	1: {"fnv128"},
	2: {"fnv128a", "fnv"},
	3: {"murmur3-128", "murmur3"},
	4: {"murmur3-64"},
	5: {"murmur3-32"},
	6: {"md5"},
	7: {"sha-256", "sha256"},
	8: {"sha-512", "sha512"},
}

// New Instantiates a new representation of the Hash Algorithm.
func (a Algorithm) New() hash.Hash {
	switch a {
	case Xxhash:
		return cxHash.New()
	case Fnv128:
		return fnvh.New128()
	case Fnv128a:
		return fnvh.New128a()
	case Murmur3_32:
		return murmur3.New32()
	case Murmur3_64:
		return murmur3.New64()
	case Murmur3_128:
		return murmur3.New128()
	case Md5:
		return md5h.New()
	case Sha256:
		return sha256h.New()
	case Sha512:
		return sha512h.New()
	}
	return cxHash.New()
}

// Index Returns the index for the Hash Algorithm
func (a Algorithm) Index() int {
	return int(a)
}

// Index Returns the human-readable representation of the Hash Algorithm
func (a Algorithm) String() string {
	return HashAlgorithms[a.Index()][0]
}
