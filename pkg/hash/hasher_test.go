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

func TestHashWithSeed(t *testing.T) {
    hasher := NewXXH3Hasher()
    key := "test:seed:key"
    
    hash1 := hasher.HashWithSeed(key, 0)
    hash2 := hasher.HashWithSeed(key, 1)
    hash3 := hasher.HashWithSeed(key, 0) // misma seed
    
    // la seed deberia variar el hash producido
    if hash1 == hash2 {
        t.Errorf("Una diferente seed deberia provocar un hash diferente")
    }
    
    // misma key y seed deberia tener el mismo hash
    if hash1 != hash3 {
        t.Errorf("Hash con misma key y semilla deberia dar igual, no fue deterministico: %d != %d", hash1, hash3)
    }
}
