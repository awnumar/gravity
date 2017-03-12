package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/crypto/blake2b"

	"github.com/libeclipse/pocket/coffer"
	"github.com/libeclipse/pocket/crypto"
	"github.com/libeclipse/pocket/input"
	"github.com/libeclipse/pocket/memory"
)

var (
	// The default cost factor for key deriviation.
	scryptCost = map[string]int{"N": 18, "r": 16, "p": 1}

	// Store the container ID globally.
	masterPassword []byte
)

func main() {
	// Setup the secret store.
	err := coffer.Setup()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer coffer.Close()

	// CleanupMemory in case of Ctrl+C
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		memory.Cleanup()
		os.Exit(0)
	}()

	// Launch CLI.
	err = cli()
	if err != nil {
		fmt.Println(err)
	}

	// Zero out and unlock any protected memory.
	memory.Cleanup()
}

func cli() error {
	var err error

	banner := `
                                  ▄▄
                                  ██                    ██
    ██▄███▄    ▄████▄    ▄█████▄  ██ ▄██▀    ▄████▄   ███████
    ██▀  ▀██  ██▀  ▀██  ██▀    ▀  ██▄██     ██▄▄▄▄██    ██
    ██    ██  ██    ██  ██        ██▀██▄    ██▀▀▀▀▀▀    ██
    ███▄▄██▀  ▀██▄▄██▀  ▀██▄▄▄▄█  ██  ▀█▄   ▀██▄▄▄▄█    ██▄▄▄
    ██ ▀▀▀      ▀▀▀▀      ▀▀▀▀▀   ▀▀   ▀▀▀    ▀▀▀▀▀      ▀▀▀▀
    ██
                        The guardian of super-secret things.
`
	fmt.Println(banner)

	masterPassword, err = input.GetMasterPassword()
	if err != nil {
		return err
	}
	fmt.Println("") // For formatting.

	help := `:: add       - Store some new data in the database.
:: get       - Retrieve some data from the database.
:: remove    - Remove some previously stored data.
:: decoys    - Add a variable number of decoys.
:: passwd    - Change the session's master password.
:: exit      - Exit the program.`

	for {
		cmd := strings.ToLower(strings.TrimSpace(string(input.Input("$ "))))

		switch cmd {
		case "passwd":
			masterPassword, err = input.GetMasterPassword()
			if err != nil {
				return err
			}
		case "add":
			err = add()
			if err != nil {
				return err
			}
		case "get":
			err = get()
			if err != nil {
				return err
			}
		case "remove":
			remove()
		case "decoys":
			decoys()
		case "exit":
			return nil
		default:
			fmt.Println(help)
		}
	}
}

func add() error {
	// Prompt the user for the identifier.
	identifier := input.Input("Enter a string to identify this data: ")

	// Derive the secure values for this "branch".
	masterKey, rootIdentifier := crypto.DeriveSecureValues(masterPassword, identifier, scryptCost)

	// Prompt user for the plaintext data.
	data := input.Input("Enter the data you wish to store: ")

	var padded []byte
	var err error
	for i := 0; i < len(data); i += 1024 {
		if i+1024 > len(data) {
			// Remaining data <= 1024.
			padded, err = crypto.Pad(data[len(data)-(len(data)%1024):], 1025)
		} else {
			// Split into chunks of 1024 bytes and pad.
			padded, err = crypto.Pad(data[i:i+1024], 1025)
		}
		if err != nil {
			return err
		}

		// Derive ID, encrypt and save to the database.
		err = coffer.Save(crypto.DeriveIdentifierN(rootIdentifier, i/1024), crypto.Encrypt(padded, masterKey))
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	fmt.Println(":: Saved that for you.")

	return nil
}

func get() error {
	// Prompt the user for the identifier.
	identifier := input.Input("Enter the string that identifies this data: ")

	// Derive the secure values for this "branch".
	masterKey, rootIdentifier := crypto.DeriveSecureValues(masterPassword, identifier, scryptCost)

	// Grab all the pieces.
	var plaintext []byte
	for n := 0; true; n++ {
		// Derive derived_identifier[n]
		ct, exists := coffer.Retrieve(crypto.DeriveIdentifierN(rootIdentifier, n))
		if exists != nil {
			// This one doesn't exist.
			break
		}

		// Decrypt this slice.
		pt := crypto.Decrypt(ct, masterKey)

		// Unpad this slice.
		unpadded, e := crypto.Unpad(pt)
		if e != nil {
			return e
		}

		// Append this slice of plaintext to the rest of it.
		plaintext = append(plaintext, unpadded...)
	}

	if len(plaintext) == 0 {
		fmt.Println("! There is nothing stored here")
		return nil
	}

	fmt.Printf(`
-----BEGIN DATA-----
%s
-----END DATA-----
`, plaintext)

	return nil
}

func remove() {
	// Prompt the user for the identifier.
	identifier := input.Input("Enter the string that identifies this data: ")

	// Derive the secure values for this "branch".
	_, rootIdentifier := crypto.DeriveSecureValues(masterPassword, identifier, scryptCost)

	// Delete all the pieces.
	count := 0
	for n := 0; true; n++ {
		// Get the DeriveIdentifierN for this n.
		derivedIdentifierN := crypto.DeriveIdentifierN(rootIdentifier, n)

		// Check if it exists.
		_, exists := coffer.Retrieve(derivedIdentifierN)
		if exists != nil {
			// This one doesn't exist.
			break
		}

		// It exists. Remove it.
		coffer.Delete(derivedIdentifierN)
		count++
	}

	if count != 0 {
		fmt.Println(":: Successfully removed data.")
	} else {
		fmt.Println("! There is nothing here to remove")
	}
}

func decoys() {
	var numberOfDecoys int
	var err error

	// Get the number of decoys to add as an int.
	for {
		numberOfDecoys, err = strconv.Atoi(string(input.Input("How many decoys do you want to add? ")))
		if err == nil {
			break
		}
	}

	count := 0
	for i := 0; i < numberOfDecoys; i++ {
		// Get some random bytes.
		randomBytes := crypto.GenerateRandomBytes(96)

		// Allocate 32 bytes as the key.
		var key [32]byte
		masterPassword := randomBytes[0:32]
		copy(key[:], masterPassword)

		// Allocate 32 bytes as the identifier.
		identifier := randomBytes[32:64]
		hashedIdentifier := blake2b.Sum256(identifier)

		// Allocate 32 bytes as the plaintext.
		plaintext := randomBytes[64:96]

		// Save to the database.
		coffer.Save(hashedIdentifier[:], crypto.Encrypt(plaintext, &key))

		// Increment counter.
		count++
		fmt.Printf("\rAdded %d/%d (%d%%)", count, numberOfDecoys, int(math.Floor(float64(count)/float64(numberOfDecoys)*100)))
	}
	fmt.Println("") // Formatting
}
