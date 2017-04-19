package crypto

import "golang.org/x/crypto/blake2b"

// GenDecoy generates and returns a single decoy.
func GenDecoy() (id, ct []byte) {
	// Get some random bytes.
	randomBytes := GenerateRandomBytes(64)

	// Allocate 32 bytes as the key.
	var key [32]byte
	masterPassword := randomBytes[0:32]
	copy(key[:], masterPassword)

	// Allocate 32 bytes as the identifier.
	identifier := randomBytes[32:64]
	hashedIdentifier := blake2b.Sum256(identifier)

	// Allocate 32 bytes as the plaintext.
	plaintext := make([]byte, 4096)

	// Encrypt/derive the final values.
	id = hashedIdentifier[:]
	ct = Encrypt(plaintext, &key)

	// Return the decoy to the caller.
	return
}
