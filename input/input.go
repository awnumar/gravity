package input

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"syscall"

	"github.com/libeclipse/pocket/memory"

	"golang.org/x/crypto/ssh/terminal"
)

// GetMasterPassword takes the masterPassword from the user while doing all of the verifying stuff.
func GetMasterPassword() ([]byte, error) {
	// Prompt user for password.
	masterPassword, err := SecureInput("- Master password: ")
	if err != nil {
		return nil, err
	}

	// Prompt for password confirmation.
	confirmPassword, err := SecureInput("- Confirm password: ")
	if err != nil {
		return nil, err
	}

	// Check if password matches confirmation.
	if !bytes.Equal(masterPassword, confirmPassword) {
		fmt.Println("! Passwords do not match")
		return GetMasterPassword()
	}

	return masterPassword, nil
}

// Input reads from stdin while echoing back.
func Input(prompt string) string {
	// Output prompt.
	fmt.Print(prompt)

	// Declare scanner on stdin.
	scanner := bufio.NewScanner(os.Stdin)

	// Read bytes.
	scanner.Scan()

	// Everything went well. Return the data.
	return scanner.Text()
}

// SecureInput gets input without echoing and returns a byte slice.
func SecureInput(prompt string) ([]byte, error) {
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
