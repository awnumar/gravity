package crypto

import (
	"crypto/rand"
	"errors"
	"fmt"

	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/sys/unix"
)

var (
	// Monitor if mlock() worked or not.
	mlock = true
)

// ProtectMemory calls mlock() and prevents sensitive information from being written to SWAP.
func ProtectMemory(data []byte) {
	err := unix.Mlock(data)
	if err != nil && mlock {
		// It failed once, probably won't work next time. Supress further warnings.
		mlock = false
		fmt.Printf("[!] Warning: Could not mlock() sensitive data; it might get written to SWAP [%s]", err)
	}
}

// Encrypt takes a plaintext and a 32 byte key, encrypts the plaintext with
// said key using xSalsa20 with a Poly1305 MAC, and returns the ciphertext.
func Encrypt(plaintext []byte, key *[32]byte) ([]byte, error) {
	// Generate a random nonce.
	nonceSlice, err := generateRandomBytes(24)
	if err != nil {
		return nil, err
	}

	// Store it in an array.
	var nonce [24]byte
	copy(nonce[:], nonceSlice)

	// Encrypt and return the plaintext.
	return secretbox.Seal(nonce[:], plaintext, &nonce, key), nil
}

// Decrypt takes a ciphertext and a 32 byte key, decrypts the ciphertext with
// said key, and then returns the plaintext.
func Decrypt(ciphertext []byte, key *[32]byte) ([]byte, error) {
	// Grab the nonce from the ciphertext and store it in an array.
	var nonce [24]byte
	copy(nonce[:], ciphertext[:24])

	// Decrypt the ciphertext and store the result.
	plaintext, okay := secretbox.Open([]byte{}, ciphertext[24:], &nonce, key)
	if !okay {
		// This shouldn't happen.
		return nil, errors.New("[!] Decryption of data failed")
	}

	// Protect the plaintext.
	ProtectMemory(plaintext)

	// Return the resulting plaintext.
	return plaintext, nil
}

// DeriveKey derives a 32 byte encryption key from a password and identifier.
func DeriveKey(password, identifier []byte, cost map[string]int) *[32]byte {
	// Allocate and protect  memory for the output of the hash function.
	derivedKeySlice := make([]byte, 32)
	ProtectMemory(derivedKeySlice)

	// Allocate and protect memory for the 32 byte array that we'll return.
	var derivedKey [32]byte
	ProtectMemory(derivedKey[:])

	// Derive the key and store in the memory we allocated above.
	derivedKeySlice, _ = scrypt.Key(password, identifier, 1<<uint(cost["N"]), cost["r"], cost["p"], 32)

	// Copy to the 32 byte array.
	copy(derivedKey[:], derivedKeySlice)

	// Return a pointer.
	return &derivedKey
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
		return nil, fmt.Errorf("[!] Length of data must not exceed %d bytes", padTo-1)
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
	original := make([]byte, len(text))
	ProtectMemory(original)
	copy(original, text)

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

	// Copy to its own slice so we're not referencing useless data.
	unpadded := make([]byte, len(text))
	ProtectMemory(unpadded)
	copy(unpadded, text)

	// That simple.  We're done.
	return unpadded, nil
}

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
