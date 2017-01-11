// This is a draft script that is being used to develop this feature and will subsequently be removed.

package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

// Pad implements ANSI X.923-type padding.
func Pad(text []byte, padTo int) ([]byte, error) {
	// Check if input is even valid.
	if len(text) > padTo {
		return nil, errors.New("pad: input length greater than padTo length")
	}

	// Add NULL bytes, leaving four at the end.
	padLen := padTo - len(text)
	for c := 1; c <= padLen-4; c++ {
		text = append(text, byte(0))
	}

	// Add the number of padding bytes as the last four bytes.
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, int32(padLen))
	if err != nil {
		return nil, err
	}
	text = append(text, buf.Bytes()...)

	// Return padded byte slice.
	return text, nil
}

// Unpad reverses ANSI X.923-type padding.
func Unpad(text []byte) ([]byte, error) {
	// Get the supposed length of the padding.
	padLenBytes := text[len(text)-4 : len(text)]

	// Convert length from bytes to int32.
	var padLen int32
	buf := bytes.NewReader(padLenBytes)
	err := binary.Read(buf, binary.BigEndian, &padLen)
	if err != nil {
		return nil, err
	}

	// Return everything except the padding.
	return text[:int32(len(text))-padLen], nil
}

func main() {
	text := "yellow submarine"
	fmt.Println([]byte(text))

	padded, err := Pad([]byte(text), 32)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(padded)

	unpadded, err := Unpad(padded)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(unpadded)
}
