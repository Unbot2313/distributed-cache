package hash

import (
	"strconv"

	"github.com/zeebo/xxh3"
)

type Hasher interface {
    Hash(key int) uint32
    HashWithSeed(key int, seed uint64) uint32
}

// struct para interface
type XXH3Hasher struct{}

// Constructor
func NewXXH3Hasher() *XXH3Hasher {
    return &XXH3Hasher{}
}

// Hash key
func (h *XXH3Hasher) Hash(key int) uint32 {
    keyStr := strconv.Itoa(key)
    return uint32(xxh3.HashString(keyStr))
}

// Hash con seed para los nodos virtuales
func (h *XXH3Hasher) HashWithSeed(key int, seed uint64) uint32 {
    keyStr := strconv.Itoa(key)
    return uint32(xxh3.HashStringSeed(keyStr, seed))
}