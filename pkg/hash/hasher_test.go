package hash

import (
	"testing"
)

func TestHashDeterministic(t *testing.T) {
    hasher := NewXXH3Hasher()
    key := "test:deterministic:key"
    
    hash1 := hasher.Hash(key)
    hash2 := hasher.Hash(key)
    
    if hash1 != hash2 {
        t.Errorf("El hash deberia ser determinista: %d != %d", hash1, hash2)
    }
}
