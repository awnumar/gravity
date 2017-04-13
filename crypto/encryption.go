package crypto

import "golang.org/x/crypto/nacl/secretbox"

// Encrypt takes a plaintext and a 32 byte key, encrypts the plaintext with
// said key using xSalsa20 with a Poly1305 MAC, and returns the ciphertext.
func Encrypt(plaintext []byte, key *[32]byte) []byte {
	// Generate a random nonce.
	nonceSlice := GenerateRandomBytes(24)

	// Store it in an array.
	var nonce [24]byte
	copy(nonce[:], nonceSlice)

	// Encrypt and return the plaintext.
	return secretbox.Seal(nonce[:], plaintext, &nonce, key)
}

// Decrypt takes a ciphertext and a 32 byte key, decrypts the ciphertext with
// said key, and then returns the plaintext.
func Decrypt(ciphertext []byte, key *[32]byte) []byte {
	// Grab the nonce from the ciphertext and store it in an array.
	var nonce [24]byte
	copy(nonce[:], ciphertext[:24])

	// Decrypt the ciphertext and store the result.
	plaintext, _ := secretbox.Open([]byte{}, ciphertext[24:], &nonce, key)

	// Return the resulting plaintext.
	return plaintext
}
