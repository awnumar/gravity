package auxiliary

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// Setup sets up the environment.
func Setup() error {
	// Get the current user.
	user, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}

	// Change the working directory to the user's home.
	err = os.Chdir(user.HomeDir)
	if err != nil {
		log.Fatalln(err)
	}

	// Check if we've done this before.
	if _, err = os.Stat("./.envelope/secrets"); err == nil {
		// Apparently we have.
		return nil
	}

	// Create a directory to store our stuff in.
	err = os.Mkdir("./.envelope", 0700)
	if err != nil && !os.IsExist(err) {
		log.Fatalln(err)
	}

	// Create an empty storage file for the secrets.
	err = ioutil.WriteFile("./.envelope/secrets", []byte(""), 0700)
	if err != nil {
		log.Fatalln(err)
	}

	return nil
}

// GetPass prompts for input without echo.
func GetPass(prompt string) []byte {
	// Output prompt.
	fmt.Print(prompt)

	// Get input without echoing back.
	input, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatalln(err)
	}

	// Output a newline for formatting.
	fmt.Println()

	// Return password.
	return input
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
		log.Fatalln(err)
	}

	// Write the JSON to the disk.
	data := base64.StdEncoding.EncodeToString([]byte(jsonFormattedSecrets))
	ioutil.WriteFile("./.envelope/secrets", []byte(data), 0700)
}

// RetrieveSecrets retrieves the secrets from the disk.
func RetrieveSecrets() map[string]interface{} {
	// Read the raw JSON from the disk.
	jsonFormattedSecrets, err := ioutil.ReadFile("./.envelope/secrets")
	if err != nil {
		log.Fatalln(err)
	}

	// Convert the secrets from base64.
	jsonFormattedSecrets, err = base64.StdEncoding.DecodeString(string(jsonFormattedSecrets))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Convert JSON to map[string]interface{} type.
	secrets := make(map[string]interface{})
	json.Unmarshal(jsonFormattedSecrets, &secrets)

	// Return the secrets.
	return secrets
}
