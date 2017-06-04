package crypto

import (
	"encoding/binary"
	"fmt"
	"runtime/debug"
	"unsafe"

	"github.com/libeclipse/memguard"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/scrypt"
)

// DeriveSecureValues derives and returns a masterKey and rootIdentifier.
func DeriveSecureValues(masterPassword, identifier *memguard.LockedBuffer, costFactor map[string]int) (*[32]byte, []byte) {
	// Allocate and protect memory for the concatenated values, and append the values to it.
	concatenatedValues, err := memguard.New(len(masterPassword)+len(identifier), false)
	if err != nil {
		fmt.Println(err)
		memguard.SafeExit(1)
	}
	concatenatedValues.Copy(masterPassword)
	concatenatedValues.CopyAt(identifier, len(masterPassword))

	// Derive the rootKey and then protect it.
	rootKeySlice, _ := scrypt.Key(
		concatenatedValues.Buffer, // Input data.
		[]byte(""),                // Salt.
		1<<uint(costFactor["N"]),  // Scrypt parameter N.
		costFactor["r"],           // Scrypt parameter r.
		costFactor["p"],           // Scrypt parameter p.
		64)                        // Output hash length.
	rootKey, _ := memguard.NewFromBytes(rootKeySlice, false)

	// Force the Go GC to do its job.
	debug.FreeOSMemory()

	// Get a pointer to the masterKey as an array.
	masterKeyArrayPtr := (*[32]byte)(unsafe.Pointer(&rootKey.Buffer[0]))

	// Slice and return respective values.
	return masterKeyArrayPtr, rootKey.Buffer[32:64]
}

// DeriveIdentifierN derives a value for derivedIdentifier for a value of `n`.
func DeriveIdentifierN(rootIdentifier []byte, n uint64) []byte {
	// Convert n to a byte slice.
	byteN := make([]byte, 8)
	binary.LittleEndian.PutUint64(byteN, n)

	// Append the uint64 to the root identifier.
	hashArg, _ := memguard.New(40, false)
	hashArg.Copy(rootIdentifier)
	copy(hashArg.Buffer[32:40], byteN)

	// Derive derivedIdentifier.
	derivedIdentifier := blake2b.Sum256(hashArg.Buffer)

	// Return as slice instead of array.
	return derivedIdentifier[:]
}

// DeriveMetaIdentifierN does the same as DeriveIdentifierN but uses signed integers instead of
// unsigned 64 bit unsigned. The intended purpose is for storing metadata and header information.
func DeriveMetaIdentifierN(rootIdentifier []byte, n int) []byte {
	// Convert n to a byte slice.
	byteN := make([]byte, 10)
	binary.PutVarint(byteN, int64(n))

	// Append the uint64 to the root identifier.
	hashArg, _ := memguard.New(42, false)
	hashArg.Copy(rootIdentifier)
	hashArg.CopyAt(byteN, 32)

	// Derive derivedIdentifier.
	derivedIdentifier := blake2b.Sum256(hashArg.Buffer)

	// Return as slice instead of array.
	return derivedIdentifier[:]
}
