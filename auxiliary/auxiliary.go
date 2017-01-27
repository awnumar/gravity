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
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	// ErrHelp is for when the user just wanted help.
	ErrHelp = errors.New("help")
)

// ParseArgs parses the command line arguments and returns the
// user's configuration options for the caller to then use.
func ParseArgs(args []string) (string, map[string]int, error) {
	helpMessage := fmt.Sprintf(`Usage: %s [mode] [options]

    help [mode] - Print detailed help regarding a mode.

    add 	- Add a new secure secret to storage.
    get 	- Retrieve a previously stored secret.
    forget 	- Remove a previously stored secret.

    -c N,r,p	- Cost factors for scrypt. If you don't
                  understand this, leave it alone. Note
                  that, if set, you'll have to specify
                  these settings on every retrieval and
                  also every time you delete an entry.
                  The defaults are '18,8,1'.`, args[0])

	if len(args) < 2 {
		fmt.Println(helpMessage)
		return "", nil, ErrHelp
	}

	switch args[1] {
	case "help":
		if len(args) < 3 {
			fmt.Println(helpMessage)
			return "", nil, ErrHelp
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
			return "", nil, errors.New("[!] Invalid mode passed to help")
		}
		return "", nil, ErrHelp
	case "add", "get", "forget":
		if len(args) > 2 && args[2] == "-c" {
			if len(args) < 4 {
				return "", nil, errors.New("[!] Nothing passed to -c")
			}

			costFactorParams := strings.Split(args[3], ",")
			if len(costFactorParams) != 3 {
				return "", nil, errors.New("[!] Invalid number of arguments passed to -c")
			}

			N, err := strconv.ParseInt(costFactorParams[0], 10, 0)
			if err != nil {
				return "", nil, errors.New("[!] Arguments to -c must be integers")
			}
			r, err := strconv.ParseInt(costFactorParams[1], 10, 0)
			if err != nil {
				return "", nil, errors.New("[!] Arguments to -c must be integers")
			}
			p, err := strconv.ParseInt(costFactorParams[2], 10, 0)
			if err != nil {
				return "", nil, errors.New("[!] Arguments to -c must be integers")
			}

			if !(N > 1) {
				return "", nil, errors.New("[!] N must be more than 1")
			}

			scryptCost := map[string]int{"N": int(N), "r": int(r), "p": int(p)}

			return args[1], scryptCost, nil
		}
		return args[1], nil, nil
	default:
		return "", nil, errors.New("[!] Invalid option")
	}
}

// GetInputs takes a slice of required inputs,
// asks the user for them, and returns them in
// the same order.
func GetInputs(required []string) []string {
	for i, input := range required {
		switch input {
		case "password":
			password := GetPass("[-] Password: ")
			if len(password) < 1 {
				fmt.Println("[!] Length of password must be non-zero.")
				os.Exit(1)
			}
			required[i] = string(password)
		case "identifier":
			identifier := Input("[-] Identifier: ")
			if len(identifier) < 1 {
				fmt.Println("[!] Length of identifier must be non-zero.")
				os.Exit(1)
			}
			required[i] = identifier
		case "secret":
			secret := Input("[-] Secret: ")
			if len(secret) < 1 || len(secret) > 1024 {
				fmt.Println("[!] Length of secret must be between 1-1024 bytes.")
				os.Exit(1)
			}
			required[i] = secret
		}
	}
	return required
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
