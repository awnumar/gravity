package stdin

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"syscall"

	"github.com/libeclipse/memguard"

	"golang.org/x/crypto/ssh/terminal"
)

// GetMasterPassword takes the masterPassword from the user while doing all of the verifying stuff.
func GetMasterPassword() (*memguard.LockedBuffer, error) {
	// Prompt user for password and confirmation.
	masterPassword := Secure("- Master password: ")
	confirmPassword := Secure("- Confirm password: ")

	// Check if password matches confirmation.
	if !bytes.Equal(masterPassword.Buffer, confirmPassword.Buffer) {
		fmt.Println("! Passwords do not match")
		return GetMasterPassword()
	}

	// We no longer need this.
	confirmPassword.Destroy()

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
func Secure(prompt string) *memguard.LockedBuffer {
	// Output prompt.
	fmt.Print(prompt)

	// Get input without echoing back.
	rawinput, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println(err)
		memguard.SafeExit(1)
	}

	// Secure the input value.
	input, err := memguard.NewFromBytes(rawinput)
	if err != nil {
		fmt.Println(err)
		memguard.SafeExit(1)
	}

	// Output a newline for formatting.
	fmt.Println()

	// Return password.
	return input
}
