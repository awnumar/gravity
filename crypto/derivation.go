package crypto

import (
	"encoding/binary"

	"github.com/libeclipse/dissident/memory"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/scrypt"
)

// DeriveSecureValues derives and returns a masterKey and rootIdentifier.
func DeriveSecureValues(masterPassword, identifier []byte, costFactor map[string]int) (*[32]byte, []byte) {
	// Allocate and protect memory for the concatenated values, and append the values to it.
	concatenatedValues := memory.MakeProtected(len(masterPassword) + len(identifier))
	copy(concatenatedValues[:len(masterPassword)], masterPassword)
	copy(concatenatedValues[len(masterPassword):], identifier)

	// Allocate and protect memory for the output of the hash function, and put the output into it.
	rootKeySlice := memory.MakeProtected(64)
	rootKeySlice, _ = scrypt.Key(
		concatenatedValues,       // Input data.
		[]byte(""),               // Salt.
		1<<uint(costFactor["N"]), // Scrypt parameter N.
		costFactor["r"],          // Scrypt parameter r.
		costFactor["p"],          // Scrypt parameter p.
		64)                       // Output hash length.

	// Allocate a protected array to hold the key, and copy the key into it.
	var masterKey [32]byte
	memory.Protect(masterKey[:])
	copy(masterKey[:], rootKeySlice[0:32])

	// Slice and return respective values.
	return &masterKey, rootKeySlice[32:64]
}

// DeriveIdentifierN derives a value for derivedIdentifier for a value of `n`.
func DeriveIdentifierN(rootIdentifier []byte, n uint64) []byte {
	// Convert n to a byte slice.
	byteN := make([]byte, 8)
	binary.LittleEndian.PutUint64(byteN, n)

	// Append the uint64 to the root identifier.
	hashArg := memory.MakeProtected(32)
	copy(hashArg, rootIdentifier)
	hashArg = append(hashArg, byteN...)

	// Derive derivedIdentifier.
	derivedIdentifier := blake2b.Sum256(hashArg)

	// Return as slice instead of array.
	return derivedIdentifier[:]
}
