package crypto

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/libeclipse/memguard"
)

func TestDecrypt(t *testing.T) {
	keySlice, _ := base64.StdEncoding.DecodeString("JNut6eJfb6ySwOac7FHe3bsSU75FpL/o776VD+oYWxk=")
	ciphertext, _ := base64.StdEncoding.DecodeString("5yiWqYEPgy9CbwMlJVxm3ge4h97X7Ptmvz6M3XLE2fLWpCo3F+VdcvU+Vrw=")

	key, _ := memguard.NewFromBytes(keySlice, false)

	// Correct key
	plaintext, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Error("Unexpected err:", err)
	}
	if !bytes.Equal(plaintext, []byte("test")) {
		t.Error("Expected plaintext to be `test`; got", plaintext)
	}

	// Incorrect key
	key.Copy([]byte("lel"))
	plaintext, err = Decrypt(ciphertext, key)
	if err == nil {
		t.Error("Expected error; got nil")
	}
	if plaintext != nil {
		t.Error("Expected plaintext to be nil; got", plaintext)
	}
}

func TestEncryptionCycle(t *testing.T) {
	plaintext := []byte("this is a test plaintext")

	key, _ := memguard.New(32, false)

	ciphertext := Encrypt(plaintext, key)
	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Error("Unexpected err:", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Error("Decrypted != Plaintext; decrypted =", string(decrypted))
	}
}
