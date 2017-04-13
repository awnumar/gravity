package crypto

import (
	"encoding/binary"

	"github.com/libeclipse/tranquil/memory"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/scrypt"
)

// DeriveSecureValues derives and returns a masterKey and rootIdentifier.
func DeriveSecureValues(masterPassword, identifier []byte, costFactor map[string]int) (*[32]byte, []byte) {
	// Concatenate the inputs.
	concatenatedValues := append(masterPassword, identifier...)
	memory.Protect(concatenatedValues)

	// Allocate and protect memory for the output of the hash function.
	rootKeySlice := make([]byte, 64)
	memory.Protect(rootKeySlice)

	// Allocate and protect memory for the 32 byte array that we'll return.
	var masterKey [32]byte
	memory.Protect(masterKey[:])

	// Derive rootKey.
	rootKeySlice, _ = scrypt.Key(concatenatedValues, []byte(""), 1<<uint(costFactor["N"]), costFactor["r"], costFactor["p"], 64)

	// Copy to the 32 byte array.
	copy(masterKey[:], rootKeySlice[0:32])

	// Slice and return respective values.
	return &masterKey, rootKeySlice[32:64]
}

// DeriveIdentifierN derives a value for derivedIdentifier for a value of `n`.
func DeriveIdentifierN(rootIdentifier []byte, n int) []byte {
	// Convert n to a byte slice.
	byteN := make([]byte, 4)
	binary.LittleEndian.PutUint32(byteN, uint32(n))

	// Derive derivedIdentifier.
	derivedIdentifier := blake2b.Sum256(append(rootIdentifier, byteN...))

	// Return as slice instead of array.
	return derivedIdentifier[:]
}
