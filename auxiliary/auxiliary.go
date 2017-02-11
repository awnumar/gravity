package auxiliary

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/libeclipse/pocket/crypto"

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

// Input takes input from the user.
func Input(prompt string) ([]byte, error) {
	// Output prompt.
	fmt.Print(prompt)

	// Create scanner and get input.
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	// Check the length of the data.
	data := scanner.Bytes()
	crypto.ProtectMemory(data)
	if len(data) == 0 {
		return nil, errors.New("[!] Length of input must be non-zero")
	}

	// Everything went well. Return the data.
	return data, nil
}

// GetPass prompts for input without echoing back.
func GetPass() ([]byte, error) {
	// Prompt for password.
	password, err := _getPass("[-] Password: ")
	if err != nil {
		return nil, err
	}

	// Check if length of password is zero.
	if len(password) == 0 {
		return nil, errors.New("[!] Length of password must be non-zero")
	}

	// Prompt for password confirmation.
	confirmPassword, err := _getPass("[-] Confirm password: ")
	if err != nil {
		return nil, err
	}

	// Check if password matches confirmation.
	if !bytes.Equal(password, confirmPassword) {
		return nil, errors.New("[!] Passwords do not match")
	}

	// Return the password.
	return password, nil
}

func _getPass(prompt string) ([]byte, error) {
	// Output prompt.
	fmt.Print(prompt)

	// Get input without echoing back.
	input, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	crypto.ProtectMemory(input)

	// Output a newline for formatting.
	fmt.Println()

	// Return password.
	return input, nil
}
