package stdin

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"syscall"

	"github.com/libeclipse/dissident/memory"

	"golang.org/x/crypto/ssh/terminal"
)

// GetMasterPassword takes the masterPassword from the user while doing all of the verifying stuff.
func GetMasterPassword() ([]byte, error) {
	// Prompt user for password and confirmation.
	masterPassword := Secure("- Master password: ")
	confirmPassword := Secure("- Confirm password: ")

	// Check if password matches confirmation.
	if !bytes.Equal(masterPassword, confirmPassword) {
		fmt.Println("! Passwords do not match")
		return GetMasterPassword()
	}

	return masterPassword, nil
}

// Standard reads from stdin while echoing back.
func Standard(prompt string) string {
	// Output prompt.
	fmt.Print(prompt)

	// Declare scanner on stdin.
	scanner := bufio.NewScanner(os.Stdin)

	// Read bytes.
	scanner.Scan()

	// Everything went well. Return the data.
	return scanner.Text()
}

// Secure gets input without echoing and returns a byte slice.
func Secure(prompt string) []byte {
	// Output prompt.
	fmt.Print(prompt)

	// Get input without echoing back.
	input, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println(err)
		memory.SafeExit(1)
	}
	memory.Protect(input)

	// Output a newline for formatting.
	fmt.Println()

	// Return password.
	return input
}
