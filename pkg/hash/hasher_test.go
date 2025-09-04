package hash

import (
	"testing"
)

// no se si dejar este
func TestHash(t *testing.T){
	Hasher := NewXXH3Hasher()

	// probar combinaciones
	testCases := []string{"server-1", "server-2", "server-3", "server-1"}
	hash1 := Hasher.Hash(testCases[0])
	hash2 := Hasher.Hash(testCases[1])
	hash3 := Hasher.Hash(testCases[2])
	hash4 := Hasher.Hash(testCases[3])
	if hash1 != hash4 {
		t.Errorf("El hash de %s y %s han tenido resultados diferentes", testCases[0], testCases[3])
	}
	if hash1 == hash2 {
		t.Errorf("El hash de %s y %s han sido iguales", testCases[0], testCases[1])
	}
	if hash1 == hash3 {
		t.Errorf("El hash de %s y %s han sido iguales", testCases[0], testCases[2])
	}
	if hash2 == hash3 {
		t.Errorf("El hash de %s y %s han sido iguales", testCases[1], testCases[2])
	}
}