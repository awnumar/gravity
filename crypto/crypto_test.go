package crypto

import (
	"encoding/base64"
	"testing"
)

func TestGenerateRandomBytes(t *testing.T) {
	randomBytes := generateRandomBytes(32)
	if len(randomBytes) != 32 {
		t.Error("Expected length to be 32; got", len(randomBytes))
	}
}

func TestDecrypt(t *testing.T) {
	key, _ := base64.StdEncoding.DecodeString("JNut6eJfb6ySwOac7FHe3bsSU75FpL/o776VD+oYWxk=")
	ciphertext := "5yiWqYEPgy9CbwMlJVxm3ge4h97X7Ptmvz6M3XLE2fLWpCo3F+VdcvU+Vrw="
	plaintext := Decrypt(ciphertext, key)
	if plaintext != "test" {
		t.Error("Expected plaintext to be `test`; got", plaintext)
	}
}

func TestDeriveKey(t *testing.T) {
	derivedKey := base64.StdEncoding.EncodeToString(DeriveKey([]byte("password"), []byte("identifier")))
	if derivedKey != "rjbQVprXRtR4z3ZYGxfcBIYLj3exf/ftMVpdsc6YKGo=" {
		t.Error("Expected `rjbQVprXRtR4z3ZYGxfcBIYLj3exf/ftMVpdsc6YKGo=`; got", derivedKey)
	}
}

func TestDeriveID(t *testing.T) {
	derivedKey := DeriveID([]byte("identifier"))
	if derivedKey != "HRd9/hpzbvfCEnhfNTIMPnGHOhTFEZSoVrdcBOrQT7w=" {
		t.Error("Expected `HRd9/hpzbvfCEnhfNTIMPnGHOhTFEZSoVrdcBOrQT7w=`; got", derivedKey)
	}
}
