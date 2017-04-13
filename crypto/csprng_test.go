package crypto

import "testing"

func TestGenerateRandomBytes(t *testing.T) {
	randomBytes := GenerateRandomBytes(32)
	if len(randomBytes) != 32 {
		t.Error("Expected length to be 32; got", len(randomBytes))
	}
}
