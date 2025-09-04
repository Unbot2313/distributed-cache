package hash

import (
	"fmt"
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


func TestHashCollision(t *testing.T) {
    hasher := NewXXH3Hasher()
    hashes := make(map[uint32]string)
    collisions := 0
    
    // Test 10,000 keys
    for i := 0; i < 10000; i++ {
        key := fmt.Sprintf("key:%d", i)
        hash := hasher.Hash(key)
        
        if existing, exists := hashes[hash]; exists {
            collisions++
            t.Logf("Colision: %s y %s tienen como hash %d", key, existing, hash)
        } else {
            hashes[hash] = key
        }
    }
    
    // el xxh3 puede tener colisiones igual que todas las funciones hash,
	// sin embargo al ser no criptografica imagino que puede tener mas posibles colisiones
	// sobre todo si el output siendo de uint64 lo convierto en uint32
    if collisions >= 2 {
        t.Errorf("Demasidads colisiones: %d", collisions)
    }
}

