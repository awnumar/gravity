package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
)

func generateRandomBytes(n int) []byte {
	// Create a byte slice (b) of size n to store the random bytes.
	b := make([]byte, n)

	// Read n bytes into b; throw an error if number of bytes read != n.
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalln(err)
	}

	// Return the CSPR bytes.
	return b
}

// Encrypt takes a plaintext and a 32 byte key, encrypts the plaintext with
// said key using xSalsa20 with a Poly1305 MAC, and returns the ciphertext.
func Encrypt(plaintext []byte, key []byte) string {
	// Generate a random nonce.
	nonceSlice := generateRandomBytes(24)

	// Store it in an array.
	var nonce [24]byte
	copy(nonce[:], nonceSlice)

	// Store the symmetric key in an array.
	var secretKey [32]byte
	copy(secretKey[:], key)

	// Encrypt the plaintext.
	ciphertext := secretbox.Seal(nonce[:], plaintext, &nonce, &secretKey)

	// Return the base64 encoded ciphertext.
	return base64.StdEncoding.EncodeToString(ciphertext)
}

// Decrypt takes a ciphertext and a 32 byte key, decrypts the ciphertext with
// said key, and then returns the plaintext. If the key is incorrect, decryption
// fails and the program terminates with exit code 1.
func Decrypt(base64EncodedCiphertext string, key []byte) []byte {
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
		fmt.Println("[!] Decryption failed")
		os.Exit(1)
	}

	// Return the resulting plaintext.
	return plaintext
}

// DeriveKey derives a 32 byte encryption key from a password and identifier.
func DeriveKey(password, identifier []byte, costFactor map[string]int) []byte {
	derivedKey, _ := scrypt.Key(password, identifier, 1<<uint(costFactor["N"]), costFactor["r"], costFactor["p"], 32)
	return derivedKey
}

// DeriveID hashes the identifier using Scrypt and returns a base64 encoded string.
func DeriveID(identifier []byte, costFactor map[string]int) string {
	derivedKey, _ := scrypt.Key(identifier, []byte(""), 1<<uint(costFactor["N"]), costFactor["r"], costFactor["p"], 32)
	return base64.StdEncoding.EncodeToString(derivedKey)
}

// Pad implements byte padding.
func Pad(text []byte, padTo int) ([]byte, error) {
	// Check if input is even valid.
	if len(text) > padTo-1 {
		return nil, errors.New("pad: text length must not exceed (padTo-1)")
	}

	// Add the compulsory byte of value `1`.
	text = append(text, byte(1))

	// Determine number of zeros to add.
	padLen := padTo - len(text)

	// Append the determined number of zeroes to the text.
	for n := 1; n <= padLen; n++ {
		text = append(text, byte(0))
	}

	// Return padded byte slice.
	return text, nil
}

// Unpad reverses byte padding.
func Unpad(text []byte) []byte {
	// Iterate over the text backwards,
	// removing the appropriate padding bytes.
	for i := len(text) - 1; i >= 0; i-- {
		if text[i] == 0 {
			text = text[:len(text)-1]
			continue
		}
		if text[i] == 1 {
			text = text[:len(text)-1]
			break
		}
	}

	// That simple.  We're done.
	return text
}
