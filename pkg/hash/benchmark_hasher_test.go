package hash

import "testing"

// Benchmark de cuanto tarda en hacer hash de las keys
func BenchmarkHash(b *testing.B) {
    hasher := NewXXH3Hasher()
    key := "benchmark:test:key:1234567890"
    
    for range b.N {
        hasher.Hash(key)
    }
}