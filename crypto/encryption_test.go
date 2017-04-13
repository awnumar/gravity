package crypto

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func TestDecrypt(t *testing.T) {
	keySlice, _ := base64.StdEncoding.DecodeString("JNut6eJfb6ySwOac7FHe3bsSU75FpL/o776VD+oYWxk=")
	ciphertext, _ := base64.StdEncoding.DecodeString("5yiWqYEPgy9CbwMlJVxm3ge4h97X7Ptmvz6M3XLE2fLWpCo3F+VdcvU+Vrw=")

	// Correct key
	var key [32]byte
	copy(key[:], keySlice)
	plaintext := Decrypt(ciphertext, &key)
	if !bytes.Equal(plaintext, []byte("test")) {
		t.Error("Expected plaintext to be `test`; got", plaintext)
	}

	// Incorrect key
	var incorrectKey [32]byte
	copy(incorrectKey[:], []byte("yellow submarine"))
	plaintext = Decrypt(ciphertext, &incorrectKey)
	if plaintext != nil {
		t.Error("Expected plaintext to be nil; got", plaintext)
	}
}

func TestEncryptionCycle(t *testing.T) {
	plaintext := []byte("this is a test plaintext")

	var key [32]byte
	copy(key[:], []byte("yellow submarine"))

	ciphertext := Encrypt(plaintext, &key)
	decrypted := Decrypt(ciphertext, &key)

	if !bytes.Equal(decrypted, plaintext) {
		t.Error("Decrypted != Plaintext; decrypted =", string(decrypted))
	}
}
