package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
)

func generateRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalln(err)
	}
	return b
}

// Encrypt takes a plaintext and a 32 byte key, encrypts the plaintext with
// said key using xSalsa20 with a Poly1305 MAC, and returns the ciphertext.
func Encrypt(plaintext string, key []byte) string {
	// Generate a random nonce.
	nonceSlice := generateRandomBytes(24)

	// Store it in an array.
	var nonce [24]byte
	copy(nonce[:], nonceSlice)

	// Store the symmetric key in an array.
	var secretKey [32]byte
	copy(secretKey[:], key)

	// Encrypt the plaintext.
	ciphertext := secretbox.Seal(nonce[:], []byte(plaintext), &nonce, &secretKey)

	// Return the base64 encoded ciphertext.
	return base64.StdEncoding.EncodeToString(ciphertext)
}

// Decrypt takes a ciphertext and a 32 byte key, decrypts the ciphertext with
// said key, and then returns the plaintext. If the key is incorrect, decryption
// fails and the program terminates with exit code 1.
func Decrypt(base64EncodedCiphertext string, key []byte) string {
	// Decode base64 encoded ciphertext into bytes.
	ciphertext, err := base64.StdEncoding.DecodeString(base64EncodedCiphertext)
	if err != nil {
		log.Fatalln(err)
	}

	// Grab the nonce from the ciphertext and store it in an array.
	var nonce [24]byte
	copy(nonce[:], ciphertext[:24])

	// Store the symmetric key in an array.
	var secretKey [32]byte
	copy(secretKey[:], key)

	// Decrypt the ciphertext and store the result.
	plaintext, okay := secretbox.Open([]byte{}, ciphertext[24:], &nonce, &secretKey)
	if !okay {
		fmt.Println("[!] decryption failed")
		os.Exit(1)
	}

	// Return the resulting plaintext.
	return string(plaintext)
}

// DeriveKey derives a 32 byte encryption key from a password and identifier.
func DeriveKey(password, identifier []byte) []byte {
	derivedKey, _ := scrypt.Key(password, identifier, 1<<20, 8, 1, 32)
	return derivedKey
}

// DeriveID hashes the identifier using Scrypt and returns a base64 encoded string.
func DeriveID(identifier []byte) string {
	dk, _ := scrypt.Key(identifier, []byte(""), 1<<18, 8, 1, 32)
	return base64.StdEncoding.EncodeToString(dk)
}
