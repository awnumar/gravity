package crypto

import (
	"errors"
	"fmt"
)

// Pad implements byte padding.
func Pad(text []byte, padTo int) ([]byte, error) {
	// Check if input is even valid.
	if len(text) > padTo-1 {
		return nil, fmt.Errorf("! Length of data must not exceed %d bytes", padTo-1)
	}

	// Create a new slice to store the padded data since we don't want to mess with the original.
	padded := make([]byte, padTo)

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
		} else if text[i] == 1 {
			text = text[:len(text)-1]
			break
		} else {
			return nil, errors.New("! Invalid padding")
		}
	}

	// Copy to its own slice so we're not referencing useless data.
	unpadded := make([]byte, len(text))
	copy(unpadded, text)

	// That simple. We're done.
	return unpadded, nil
}
