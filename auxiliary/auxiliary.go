package auxiliary

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// ParseArgs parses the command line arguments and returns the
// user's configuration options for the caller to then use.
func ParseArgs(args []string) (string, error) {
	helpMessage := fmt.Sprintf(`Usage: %s [arguments]

    help [mode] - Print detailed help regarding a mode.

    add 	- Add a new secure secret to storage.
    get 	- Retrieve a previously stored secret.
    forget 	- Remove a previously stored secret.`, args[0])

	if len(args) < 2 {
		fmt.Println(helpMessage)
		return "", errors.New("help")
	}

	switch args[1] {
	case "help":
		if len(args) < 3 {
			return "", errors.New("[!] The help command requires an argument")
		}

		switch args[2] {
		case "add":
			fmt.Printf(`Usage: %s add

    This mode is used for adding new secrets to the store.

    You'll be prompted to enter a password and an identifier. Both of those things
    together are used to derive the encryption key that protects your secrets, so
    a strong password is recommended. For the identifier, you should aim to use a
    phrase like 'l33t encryption key for them thingz init' instead of something like
    'encryption key' which could easily be guessed. There's also nothing stopping you
    from using random values for both fields, assuming that you can remember them.

    Speaking of not stopping you from doing things, you're also free to use different
    passwords for different entries. Aside from increasing security, this also has the
    side effect of allowing deniable encryption. Simply add a few legit-looking secrets
    with a decoy key and if you're ever forced to disclose your keys, just give up the
    decoys. The program adds its own decoys so you can claim that the other encrypted
    entries are just that: decoys.

    It should be noted that pocket will not ask you for a password confirmation, so
    make sure to try and retrieve the secrets you store, just to check that you entered
    everything correctly. This way you won't be sorry later when you can't decrypt
    what you added.
`, args[0])
		case "get":
			fmt.Printf(`Usage: %s get

    This mode is used for retrieving secrets from the store.

    You'll be prompted to enter a password and an identifier. The program will then
    derive the secure identifier and encryption key from both of these pieces of
    information, find and decrypt the secret, and then output it.
`, args[0])
		case "forget":
			fmt.Printf(`Usage: %s forget

    This mode is used for removing secrets from the store.

    You'll just need to enter the identifier for the entry and the program
    will derive the secure identifier, locate the entry, and remove it from
    the store.

    You won't be asked for a confirmation, so when you run forget, make sure
    that you mean it.
`, args[0])
		default:
			return "", errors.New("[!] Invalid argument to help")
		}
		return "", errors.New("help")
	case "add", "get", "forget":
		return args[1], nil
	default:
		return "", errors.New("[!] Invalid argument")
	}
}

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
	if _, err = os.Stat("./.pocket/secrets"); err == nil {
		// Apparently we have.
		return nil
	}

	// Create a directory to store our stuff in.
	err = os.Mkdir("./.pocket", 0700)
	if err != nil && !os.IsExist(err) {
		log.Fatalln(err)
	}

	// Create an empty storage file for the secrets.
	err = ioutil.WriteFile("./.pocket/secrets", []byte(""), 0700)
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
	ioutil.WriteFile("./.pocket/secrets", []byte(jsonFormattedSecrets), 0700)
}

// RetrieveSecrets retrieves the secrets from the disk.
func RetrieveSecrets() map[string]interface{} {
	// Read the raw JSON from the disk.
	jsonFormattedSecrets, err := ioutil.ReadFile("./.pocket/secrets")
	if err != nil {
		log.Fatalln(err)
	}

	// Convert JSON to map[string]interface{} type.
	secrets := make(map[string]interface{})
	json.Unmarshal(jsonFormattedSecrets, &secrets)

	// Return the secrets.
	return secrets
}
