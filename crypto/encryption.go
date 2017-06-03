package crypto

import (
	"errors"
	"unsafe"

	"github.com/libeclipse/memguard"
	"golang.org/x/crypto/nacl/secretbox"
)

// Encrypt takes a plaintext and a 32 byte key, encrypts the plaintext with
// said key using xSalsa20 with a Poly1305 MAC, and returns the ciphertext.
func Encrypt(plaintext []byte, key *memguard.LockedBuffer) []byte {
	// Generate a random nonce.
	nonceSlice := GenerateRandomBytes(24)

	// Store it in an array.
	var nonce [24]byte
	copy(nonce[:], nonceSlice)

	// Get the key as an array.
	keyArrayPtr := (*[32]byte)(unsafe.Pointer(&key.Buffer[0]))

	// Encrypt and return the plaintext.
	return secretbox.Seal(nonce[:], plaintext, &nonce, keyArrayPtr)
}

// Decrypt takes a ciphertext and a 32 byte key, decrypts the ciphertext with
// said key, and then returns the plaintext.
func Decrypt(ciphertext []byte, key *memguard.LockedBuffer) ([]byte, error) {
	// Grab the nonce from the ciphertext and store it in an array.
	var nonce [24]byte
	copy(nonce[:], ciphertext[:24])

	// Get the key as an array.
	keyArrayPtr := (*[32]byte)(unsafe.Pointer(&key.Buffer[0]))

	// Decrypt the ciphertext and store the result.
	plaintext, okay := secretbox.Open([]byte{}, ciphertext[24:], &nonce, keyArrayPtr)
	if !okay {
		return nil, errors.New("! Decryption failed; data is likely corrupted")
	}

	// Return the resulting plaintext.
	return plaintext, nil
}
