package crypto

import (
	"encoding/base64"
	"reflect"
	"testing"
)

var (
	scryptCost = map[string]int{"N": 18, "r": 8, "p": 1}
)

func TestGenerateRandomBytes(t *testing.T) {
	randomBytes := generateRandomBytes(32)
	if len(randomBytes) != 32 {
		t.Error("Expected length to be 32; got", len(randomBytes))
	}
}

func TestDecrypt(t *testing.T) {
	key, _ := base64.StdEncoding.DecodeString("JNut6eJfb6ySwOac7FHe3bsSU75FpL/o776VD+oYWxk=")
	ciphertext, _ := base64.StdEncoding.DecodeString("5yiWqYEPgy9CbwMlJVxm3ge4h97X7Ptmvz6M3XLE2fLWpCo3F+VdcvU+Vrw=")

	// Correct key
	plaintext := string(Decrypt(ciphertext, key))
	if plaintext != "test" {
		t.Error("Expected plaintext to be `test`; got", plaintext)
	}
}

func TestDeriveKey(t *testing.T) {
	derivedKey := base64.StdEncoding.EncodeToString(DeriveKey([]byte("password"), []byte("identifier"), scryptCost))
	if derivedKey != "rjbQVprXRtR4z3ZYGxfcBIYLj3exf/ftMVpdsc6YKGo=" {
		t.Error("Expected `rjbQVprXRtR4z3ZYGxfcBIYLj3exf/ftMVpdsc6YKGo=`; got", derivedKey)
	}
}

func TestDeriveID(t *testing.T) {
	derivedKey := base64.StdEncoding.EncodeToString(DeriveID([]byte("identifier"), scryptCost))
	if derivedKey != "HRd9/hpzbvfCEnhfNTIMPnGHOhTFEZSoVrdcBOrQT7w=" {
		t.Error("Expected `HRd9/hpzbvfCEnhfNTIMPnGHOhTFEZSoVrdcBOrQT7w=`; got", derivedKey)
	}
}

func TestPad(t *testing.T) {
	text := []byte("yellow submarine") // 16 bytes

	// Test when padTo < len(text)
	padded, err := Pad(text, 15)
	if err == nil {
		t.Error("Expected an error since inputs are invalid; padded:", padded)
	}

	// Test when padTo == len(text)
	padded, err = Pad(text, 16)
	if err == nil {
		t.Error("Expected an error since inputs are invalid; padded:", padded)
	}

	// Test when padTo-1 = len(text)
	padded, err = Pad(text, 17)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if len(padded) != 17 {
		t.Error("expected length of padded=32; got", len(padded))
	}

	// Test when padTo > len(text)
	padded, err = Pad(text, 32)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if len(padded) != 32 {
		t.Error("expected length of padded=32; got", len(padded))
	}

	// Test when padTo >> len(text)
	padded, err = Pad(text, 4096)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if len(padded) != 4096 {
		t.Error("expected length of padded=32; got", len(padded))
	}
}

func TestUnpad(t *testing.T) {
	text := []byte("yellow submarine") // 16 bytes

	// Test when len(text) == padTo-1
	padded, _ := Pad(text, 17)
	unpadded, err := Unpad(padded)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if !reflect.DeepEqual(unpadded, text) {
		t.Error("Unpad didn't work; got", unpadded)
	}

	// Test when len(text) < padTo
	padded, err = Pad(text, 32)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	unpadded, err = Unpad(padded)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if !reflect.DeepEqual(unpadded, text) {
		t.Error("Unpad didn't work; got", unpadded)
	}

	// Test when len(text) << padTo
	padded, err = Pad(text, 4096)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	unpadded, err = Unpad(padded)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if !reflect.DeepEqual(unpadded, text) {
		t.Error("Unpad didn't work; got", unpadded)
	}

	// Test invalid padding.
	unpadded, err = Unpad(text)
	if err == nil {
		t.Error("Expected an error since inputs are invalid; unpadded:", unpadded)
	}
	if !reflect.DeepEqual(unpadded, text) {
		t.Error("Unpadded != text with invalid input; unpadded:", unpadded)
	}
}
