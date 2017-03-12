package input

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/libeclipse/pocket/memory"

	"golang.org/x/crypto/ssh/terminal"
)

// GetPass takes a password from the user while doing all of the verifying stuff.
func GetPass() ([]byte, error) {
	// Prompt user for password.
	masterPassword, err := _secureInput("- Master password: ")
	if err != nil {
		return nil, err
	}

	// Prompt for password confirmation.
	confirmPassword, err := _secureInput("- Confirm password: ")
	if err != nil {
		return nil, err
	}

	// Check if password matches confirmation.
	if !bytes.Equal(masterPassword, confirmPassword) {
		return nil, errors.New("! Passwords do not match")
	}

	return masterPassword, nil
}

// Input reads from stdin while echoing back.
func Input(prompt string) []byte {
	// Output prompt.
	fmt.Print(prompt)

	// Create scanner and get input.
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	// Get the input out as bytes.
	data := scanner.Bytes()

	// Everything went well. Return the data.
	return data
}

// Get input without echoing and return a byte slice.
func _secureInput(prompt string) ([]byte, error) {
	// Output prompt.
	fmt.Print(prompt)

	// Get input without echoing back.
	input, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	memory.Protect(input)

	// Output a newline for formatting.
	fmt.Println()

	// Return password.
	return input, nil
}
