package crypto

import (
	"bytes"
	"encoding/base64"
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
	keySlice, _ := base64.StdEncoding.DecodeString("JNut6eJfb6ySwOac7FHe3bsSU75FpL/o776VD+oYWxk=")
	ciphertext, _ := base64.StdEncoding.DecodeString("5yiWqYEPgy9CbwMlJVxm3ge4h97X7Ptmvz6M3XLE2fLWpCo3F+VdcvU+Vrw=")

	// Correct key
	var key [32]byte
	copy(key[:], keySlice)
	plaintext, err := Decrypt(ciphertext, &key)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(plaintext, []byte("test")) {
		t.Error("Expected plaintext to be `test`; got", plaintext)
	}

	// Incorrect key
	var incorrectKey [32]byte
	copy(incorrectKey[:], []byte("yellow submarine"))
	plaintext, err = Decrypt(ciphertext, &incorrectKey)
	if err == nil {
		t.Error("Expected error; got nil")
	}
	if plaintext != nil {
		t.Error("Expected plaintext to be nil; got", plaintext)
	}
}

func TestEncryptionCycle(t *testing.T) {
	plaintext := []byte("this is a test plaintext")

	var key [32]byte
	copy(key[:], []byte("yellow submarine"))

	ciphertext := Encrypt(plaintext, &key)
	decrypted, err := Decrypt(ciphertext, &key)
	if err != nil {
		t.Error("Unexpected error:", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Error("Decrypted != Plaintext; decrypted =", string(decrypted))
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
	if !bytes.Equal(unpadded, text) {
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
	if !bytes.Equal(unpadded, text) {
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
	if !bytes.Equal(unpadded, text) {
		t.Error("Unpad didn't work; got", unpadded)
	}

	// Test invalid padding.
	unpadded, err = Unpad(text)
	if err == nil {
		t.Error("Expected an error since inputs are invalid; unpadded:", unpadded)
	}
	if unpadded != nil {
		t.Error("Expected unpadded to be nil; unpadded =", unpadded)
	}
}
