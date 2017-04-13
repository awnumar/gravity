package crypto

import "testing"

func TestGenerateRandomBytes(t *testing.T) {
	randomBytes := GenerateRandomBytes(32)
	if randomBytes == nil {
		t.Error("Returned bytes not random; got nil.")
	}
	if len(randomBytes) != 32 {
		t.Error("Expected length to be 32; got", len(randomBytes))
	}
}
