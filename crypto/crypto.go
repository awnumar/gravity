package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/libeclipse/pocket/crypto/memlock"

	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
)

var (
	// Count of how many goroutines there are.
	lockersCount int

	// Let the goroutines know we're exiting.
	isExiting = make(chan bool)

	// Used to wait for goroutines to finish before exiting.
	lockers sync.WaitGroup
)

// ProtectMemory prevents memory from being paged to disk, follows it
// around until program exit, then zeros it out and unlocks it.
func ProtectMemory(data []byte) {
	// Increment counters since we're starting another goroutine.
	lockersCount++ // Normal counter.
	lockers.Add(1) // WaitGroup counter.

	// Run as a goroutine so that callers don't have to be explicit.
	go func(b []byte) {
		// Monitor if we managed to lock b.
		lockSuccess := true

		// Prevent memory from being paged to disk.
		err := memlock.Lock(b)
		if err != nil {
			lockSuccess = false
			fmt.Printf("[!] Failed to lock %p; will still zero it out on exit. [Err: %s]\n", &b, err)
		}

		// Wait for the signal to let us know we're exiting.
		<-isExiting

		// Zero out the memory.
		for i := 0; i < len(b); i++ {
			b[i] = byte(0)
		}

		// If we managed to lock earlier, unlock.
		if lockSuccess {
			err := memlock.Unlock(b)
			if err != nil {
				fmt.Printf("[!] Failed to unlock %p [Err: %s]\n", &b, err)
			}
		}

		// We're done. Decrement WaitGroup counter.
		lockers.Done()
	}(data)
}

// CleanupMemory instructs the goroutines to cleanup the
// memory they've been watching and waits for them to finish.
func CleanupMemory() {
	for n := 0; n < lockersCount; n++ {
		isExiting <- true
	}
	lockers.Wait()
}

// Encrypt takes a plaintext and a 32 byte key, encrypts the plaintext with
// said key using xSalsa20 with a Poly1305 MAC, and returns the ciphertext.
func Encrypt(plaintext []byte, key *[32]byte) []byte {
	// Generate a random nonce.
	nonceSlice := generateRandomBytes(24)

	// Store it in an array.
	var nonce [24]byte
	copy(nonce[:], nonceSlice)

	// Encrypt and return the plaintext.
	return secretbox.Seal(nonce[:], plaintext, &nonce, key)
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

// DeriveSecureValues derives and returns a masterKey and rootIdentifier.
func DeriveSecureValues(masterPassword, identifier []byte, costFactor map[string]int) (*[32]byte, []byte) {
	// Concatenate the inputs.
	concatenatedValues := append(masterPassword, identifier...)
	ProtectMemory(concatenatedValues)

	// Allocate and protect memory for the output of the hash function.
	rootKeySlice := make([]byte, 64)
	ProtectMemory(rootKeySlice)

	// Allocate and protect memory for the 32 byte array that we'll return.
	var masterKey [32]byte
	ProtectMemory(masterKey[:])

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
	derivedIdentifier := sha256.Sum256(append(rootIdentifier, byteN...))

	// Return as slice instead of array.
	return derivedIdentifier[:]
}

// Pad implements byte padding.
func Pad(text []byte, padTo int) ([]byte, error) {
	// Check if input is even valid.
	if len(text) > padTo-1 {
		return nil, fmt.Errorf("[!] Length of data must not exceed %d bytes", padTo-1)
	}

	// Create a new slice to store the padded data since we don't want to mess with the original.
	padded := make([]byte, padTo)
	ProtectMemory(padded)

	// Copy text into new slice.
	copy(padded, text)

	// Add the compulsory byte of value `1`.
	padded[len(text)] = byte(1)

	// Return padded byte slice.
	return padded, nil
}

// Unpad reverses byte padding.
func Unpad(text []byte) ([]byte, error) {
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
			return nil, errors.New("unpad: invalid padding")
		}
	}

	// Copy to its own slice so we're not referencing useless data.
	unpadded := make([]byte, len(text))
	ProtectMemory(unpadded)
	copy(unpadded, text)

	// That simple. We're done.
	return unpadded, nil
}

func generateRandomBytes(n int) []byte {
	// Create a byte slice (b) of size n to store the random bytes.
	b := make([]byte, n)

	// Read n bytes into b; throw an error if number of bytes read != n.
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println(err)
		CleanupMemory()
		os.Exit(1)
	}

	// Return the CSPR bytes.
	return b
}
