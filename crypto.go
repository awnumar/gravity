package main

import (
	"errors"
	"unsafe"

	"golang.org/x/crypto/nacl/secretbox"

	"github.com/awnumar/memguard"
)

// Overhead is the size by which the ciphertext exceeds the plaintext.
const Overhead int = secretbox.Overhead + 24 // auth + nonce

// ErrInvalidKeyLength is returned when attempting to encrypt or decrypt with a key that is not exactly 32 bytes in size.
var ErrInvalidKeyLength = errors.New("<gravity::core::ErrInvalidKeyLength> key must be exactly 32 bytes")

// ErrBufferTooSmall is returned when the decryption function, Open, is given an output buffer that is too small to hold the plaintext. In practice the plaintext will be Overhead bytes smaller than the ciphertext returned by the encryption function, Seal.
var ErrBufferTooSmall = errors.New("<gravity::core::ErrBufferTooSmall> the given buffer is too small to hold the plaintext")

// ErrDecryptionFailed is returned when the attempted decryption fails. This can occur if the given key is incorrect or if the ciphertext is invalid.
var ErrDecryptionFailed = errors.New("<gravity::core::ErrDecryptionFailed> decryption failed")

// Encrypt takes a plaintext message and a 32 byte key and returns an authenticated ciphertext.
func Encrypt(plaintext, key []byte) ([]byte, error) {
	// Check the length of the key is correct.
	if len(key) != 32 {
		return nil, ErrInvalidKeyLength
	}

	// Get a reference to the key's underlying array without making a copy.
	k := (*[32]byte)(unsafe.Pointer(&key[0]))

	// Allocate space for and generate a nonce value.
	var nonce [24]byte
	memguard.ScrambleBytes(nonce[:])

	// Encrypt m and return the result.
	return secretbox.Seal(nonce[:], plaintext, &nonce, k), nil
}

/*
Decrypt decrypts a given ciphertext with a given 32 byte key and writes the result to the start of a given buffer.

The buffer must be large enough to contain the decrypted data. This is in practice Overhead bytes less than the length of the ciphertext returned by the Seal function above. This value is the size of the nonce plus the size of the Poly1305 authenticator.

The size of the decrypted data is returned.
*/
func Decrypt(ciphertext, key []byte, output []byte) (int, error) {
	// Check the length of the key is correct.
	if len(key) != 32 {
		return 0, ErrInvalidKeyLength
	}

	// Check the capacity of the given output buffer.
	if cap(output) < (len(ciphertext) - Overhead) {
		return 0, ErrBufferTooSmall
	}

	// Get a reference to the key's underlying array without making a copy.
	k := (*[32]byte)(unsafe.Pointer(&key[0]))

	// Retrieve and store the nonce value.
	var nonce [24]byte
	copy(nonce[:], ciphertext[:24])

	// Decrypt and return the result.
	m, ok := secretbox.Open(nil, ciphertext[24:], &nonce, k)
	if ok { // Decryption successful.
		copy(output[:cap(output)], m) // Move plaintext to given output buffer.
		memguard.WipeBytes(m)         // Wipe source buffer.
		return len(m), nil            // Return length of decrypted plaintext.
	}

	// Decryption unsuccessful. Either the key was wrong or the authentication failed.
	return 0, ErrDecryptionFailed
}
