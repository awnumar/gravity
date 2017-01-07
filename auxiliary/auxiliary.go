package auxiliary

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"syscall"

	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/ssh/terminal"
)

// Setup sets up the environment.
func Setup() error {
	// Get the current user.
	user, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Change the working directory to the user's home.
	err = os.Chdir(user.HomeDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Check if we've done this before.
	if _, err = os.Stat("./.envelope"); err == nil {
		// Apparently we have.
		return nil
	}

	// Create a directory to store our stuff in.
	err = os.Mkdir("./.envelope", 0700)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}

// GetPass prompts for input without echo.
func GetPass() []byte {
	input, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println()
	return input
}

// DeriveKey derives a 32 byte encryption key from a masterPassword and identifier.
func DeriveKey(masterPassword, identifier []byte, N uint) []byte {
	dk, _ := scrypt.Key(masterPassword, identifier, 1<<N, 8, 1, 32)
	return dk
}

// Input takes input from the user.
func Input(prompt string) string {
	// Output prompt.
	fmt.Print(prompt)

	// Create scanner and get input.
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	// Return the inputted data.
	return scanner.Text()
}

// SaveSecrets saves the secrets to the disk.
func SaveSecrets(secrets map[string]interface{}) {
	// Convert interface{} into raw JSON.
	jsonFormattedSecrets, err := json.Marshal(secrets)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Write the JSON to the disk.
	ioutil.WriteFile("./.envelope/secrets.enc", []byte(jsonFormattedSecrets), 0700)
}

// RetrieveSecrets retrieves the secrets from the disk.
func RetrieveSecrets() map[string]interface{} {
	// Read the raw JSON from the disk.
	jsonFormattedSecrets, err := ioutil.ReadFile("./.envelope/secrets.enc")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Convert the JSON into an interface{} type.
	secrets := make(map[string]interface{})
	err = json.Unmarshal(jsonFormattedSecrets, &secrets)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Return the secrets.
	return secrets
}
