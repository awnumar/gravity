package crypto

import (
	"crypto/rand"
	"errors"
	"fmt"

	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
)

func generateRandomBytes(n int) ([]byte, error) {
	// Create a byte slice (b) of size n to store the random bytes.
	b := make([]byte, n)

	// Read n bytes into b; throw an error if number of bytes read != n.
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	// Return the CSPR bytes.
	return b, nil
}

// Encrypt takes a plaintext and a 32 byte key, encrypts the plaintext with
// said key using xSalsa20 with a Poly1305 MAC, and returns the ciphertext.
func Encrypt(plaintext []byte, key [32]byte) ([]byte, error) {
	// Generate a random nonce.
	nonceSlice, err := generateRandomBytes(24)
	if err != nil {
		return nil, err
	}

	// Store it in an array.
	var nonce [24]byte
	copy(nonce[:], nonceSlice)

	// Encrypt the plaintext.
	ciphertext := secretbox.Seal(nonce[:], plaintext, &nonce, &key)

	// Return the base64 encoded ciphertext.
	return ciphertext, nil
}

// Decrypt takes a ciphertext and a 32 byte key, decrypts the ciphertext with
// said key, and then returns the plaintext.
func Decrypt(ciphertext []byte, key [32]byte) ([]byte, error) {
	// Grab the nonce from the ciphertext and store it in an array.
	var nonce [24]byte
	copy(nonce[:], ciphertext[:24])

	// Decrypt the ciphertext and store the result.
	plaintext, okay := secretbox.Open([]byte{}, ciphertext[24:], &nonce, &key)
	if !okay {
		// This shouldn't happen.
		return nil, errors.New("[!] Decryption of data failed")
	}

	// Return the resulting plaintext.
	return plaintext, nil
}

// DeriveKey derives a 32 byte encryption key from a password and identifier.
func DeriveKey(password, identifier []byte, cost map[string]int) [32]byte {
	//LOCKTHIS
	derivedKeySlice, _ := scrypt.Key(password, identifier, 1<<uint(cost["N"]), cost["r"], cost["p"], 32)

	// Convert to fixed-size array.
	var derivedKey [32]byte //LOCKTHIS
	copy(derivedKey[:], derivedKeySlice)

	return derivedKey
}

// DeriveID hashes the identifier using Scrypt and returns a base64 encoded string.
func DeriveID(identifier []byte, cost map[string]int) []byte {
	derivedKey, _ := scrypt.Key(identifier, []byte(""), 1<<uint(cost["N"]), cost["r"], cost["p"], 32)
	return derivedKey
}

// Pad implements byte padding.
func Pad(text []byte, padTo int) ([]byte, error) {
	// Check if input is even valid.
	if len(text) > padTo-1 {
		return nil, errors.New(fmt.Sprint("[!] Length of data must not exceed ", padTo-1, " bytes"))
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
func Unpad(text []byte) ([]byte, error) {
	// Keep a copy of the original just in case.
	var original = text

	// Iterate over the text backwards,
	// removing the appropriate padding bytes.
	for i := len(text) - 1; i >= 0; i-- {
		if text[i] == 0 {
			text = text[:len(text)-1]
			continue
		} else if text[i] == 1 {
			text = text[:len(text)-1]
			break
		} else {
			return original, errors.New("unpad: invalid padding")
		}
	}

	// That simple.  We're done.
	return text, nil
}
