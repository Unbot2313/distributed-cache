package hash

import "github.com/zeebo/xxh3"

type Hasher interface {
    Hash(key string) uint32
    HashWithSeed(key string, seed uint64) uint32
}

// struct para interface
type XXH3Hasher struct{}

// Constructor
func NewXXH3Hasher() *XXH3Hasher {
    return &XXH3Hasher{}
}

// Hash key
func (h *XXH3Hasher) Hash(key string) uint32 {
    return uint32(xxh3.HashString(key))
}

// Hash con seed para los nodos virtuales
func (h *XXH3Hasher) HashWithSeed(key string, seed uint64) uint32 {
    return uint32(xxh3.HashStringSeed(key, seed))
}