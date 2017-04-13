package crypto

import (
	"crypto/rand"
	"fmt"

	"github.com/libeclipse/tranquil/memory"
)

// GenerateRandomBytes generates cryptographically secure random bytes.
func GenerateRandomBytes(n int) []byte {
	// Create a byte slice (b) of size n to store the random bytes.
	b := make([]byte, n)

	// Read n bytes into b; throw an error if number of bytes read != n.
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println(err)
		memory.SafeExit(1)
	}

	// Return the CSPR bytes.
	return b
}
