package auxiliary

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
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
	helpMessage := fmt.Sprintf(`Usage: %s mode [-c N,r,p int]

:: Modes

  help          Display this usage message.

  add           Add a new entry to the store.
  get           Retrieve an existing entry.
  forget        Remove an existing entry.

:: Options

  -c N,r,p      Specify custom cost factors for scrypt. (default: 18,8,1)

Further help and usage information can be found in the README file or on the project page.`, args[0])

	if len(args) < 2 {
		fmt.Println(helpMessage)
		return "", nil, ErrHelp
	}

	switch args[1] {
	case "help":
		fmt.Println(helpMessage)
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
			if len(password) == 0 {
				fmt.Println("[!] Length of password must be non-zero")
				os.Exit(1)
			}
			confirmPassword := GetPass("[-] Confirm password: ")
			if !reflect.DeepEqual(password, confirmPassword) {
				fmt.Println("[!] Passwords do not match")
				os.Exit(1)
			}
			required[i] = string(password)
		case "identifier":
			identifier := Input("[-] Identifier: ")
			if len(identifier) < 1 {
				fmt.Println("[!] Length of identifier must be non-zero")
				os.Exit(1)
			}
			required[i] = identifier
		case "data":
			data := Input("[-] Data: ")
			if len(data) < 1 || len(data) > 1024 {
				fmt.Println("[!] Length of data must be between 1-1024 bytes")
				os.Exit(1)
			}
			required[i] = data
		}
	}
	return required
}

// GetPass prompts for input without echoing back.
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
