package crypto

import (
	"encoding/base64"
	"reflect"
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

func TestPad(t *testing.T) {
	text := []byte("yellow submarine") // 16 bytes

	// Test when padTo < len(text)
	padded, err := Pad(text, 15)
	if err == nil {
		t.Error("Expected en error since inputs are invalid.")
	}

	// Test when padTo == len(text)
	padded, err = Pad(text, 16)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if len(padded) != len(text) {
		t.Error("expected length of padded=length of input; got len(padded)=", len(padded))
	}

	// Test when padTo > len(text)
	padded, err = Pad(text, 32)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if len(padded) != 32 {
		t.Error("expected length of padded=32; got", len(padded))
	}

	// Check if all padding bytes are correct.
	padLen := int(padded[len(padded)-1])
	for i := len(padded) - 1; i >= len(padded)-padLen; i-- {
		if padded[i] != padded[len(padded)-1] {
			t.Error("padding is invalid")
		}
	}
}

func TestUnpad(t *testing.T) {
	text := []byte("yellow submarine")

	// Test when len(text) == padTo
	padded, err := Pad(text, 16)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	unpadded := Unpad(padded)
	if !reflect.DeepEqual(padded, unpadded) {
		t.Error("padded16 should equal unpadded")
	}
	if !reflect.DeepEqual(unpadded, text) {
		t.Error("Unpad didn't work; got", unpadded)
	}

	// Test when len(text) < padTo
	padded, err = Pad(text, 32)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	unpadded = Unpad(padded)
	if !reflect.DeepEqual(unpadded, text) {
		t.Error("Unpad didn't work; got", unpadded)
	}
}
