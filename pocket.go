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
		memory.SafeExit(0)
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

	help := `
:: add       - Add new data to the database. The data's security relies on
               both the master password and the identifier that is supplied.

:: get       - Retrieve data from the database. The data can only be retrieved
               with both the correct password and identifier.

:: remove    - Remove some previously stored data. To locate the data to remove,
               both the correct password and identifier must be supplied.

:: decoys    - This feature lets you add a variable amount of random decoy data
               that is indistinguishable from real data. Note that this data cannot
               later be removed from the database since it cannot be differentiated.

:: exit      - Exit the program.
`

	masterPassword, err = input.GetMasterPassword()
	if err != nil {
		return err
	}
	fmt.Println("") // For formatting.

	for {
		cmd := strings.ToLower(strings.TrimSpace(string(input.Input("$ "))))

		switch cmd {
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
			err = remove()
			if err != nil {
				return err
			}
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
	identifier, err := input.SecureInput("- Identifier: ")
	if err != nil {
		return err
	}

	// Derive the secure values for this "branch".
	fmt.Println("+ Generating root key...")
	masterKey, rootIdentifier := crypto.DeriveSecureValues(masterPassword, identifier, scryptCost)

	// Prompt user for the plaintext data.
	data := input.Input("- Data: ")

	// Check if it exists already.
	derivedIdentifierN := crypto.DeriveIdentifierN(rootIdentifier, 0)
	if coffer.Exists(derivedIdentifierN) {
		fmt.Println("! Cannot overwrite existing entry")
		return nil
	}

	var padded []byte
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

		// Save it.
		coffer.Save(crypto.DeriveIdentifierN(rootIdentifier, i/1024), crypto.Encrypt(padded, masterKey))

		// Wipe the padded data.
		memory.Wipe(padded)
	}

	// Wipe the plaintext.
	memory.Wipe(data)

	fmt.Println("+ Saved that for you.")

	return nil
}

func get() error {
	// Prompt the user for the identifier.
	identifier, err := input.SecureInput("- Identifier: ")
	if err != nil {
		return err
	}

	// Derive the secure values for this "branch".
	fmt.Println("+ Generating root key...")
	masterKey, rootIdentifier := crypto.DeriveSecureValues(masterPassword, identifier, scryptCost)

	// Grab all the pieces.
	var plaintext []byte
	for n := 0; true; n++ {
		// Derive derived_identifier[n]
		ct, exists := coffer.Retrieve(crypto.DeriveIdentifierN(rootIdentifier, n))
		if exists == false {
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

		// Wipe the plaintext slice.
		memory.Wipe(unpadded)
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

	// Wipe the plaintext.
	memory.Wipe(plaintext)

	return nil
}

func remove() error {
	// Prompt the user for the identifier.
	identifier, err := input.SecureInput("- Identifier: ")
	if err != nil {
		return err
	}

	// Derive the secure values for this "branch".
	fmt.Println("+ Generating root key...")
	_, rootIdentifier := crypto.DeriveSecureValues(masterPassword, identifier, scryptCost)

	// Delete all the pieces.
	count := 0
	for n := 0; true; n++ {
		// Get the DeriveIdentifierN for this n.
		derivedIdentifierN := crypto.DeriveIdentifierN(rootIdentifier, n)

		// Check if it exists.
		if coffer.Exists(derivedIdentifierN) == false {
			break
		}

		// It exists. Remove it.
		coffer.Delete(derivedIdentifierN)
		count++
	}

	if count != 0 {
		fmt.Println("+ Successfully removed data.")
	} else {
		fmt.Println("! There is nothing here to remove")
	}

	return nil
}

func decoys() {
	var numberOfDecoys int
	var err error

	// Print some help information.
	fmt.Println(`
:: For deniable encryption, use this feature in conjunction with some fake data manually-added
   under a different master-password. Then if you are ever forced to hand over your keys,
   simply give up the fake data and claim that the rest of the entries in the database are decoys.

:: You do not necessarily have to make use of this feature. Rather, simply the fact that
   it exists allows you to claim that some or all of the entries in the database are decoys.
`)

	// Get the number of decoys to add as an int.
	for {
		numberOfDecoys, err = strconv.Atoi(string(input.Input("How many decoys do you want to add? ")))
		if err == nil {
			break
		}
		fmt.Println("! Input must be an integer")
	}

	count := 0
	for i := 0; i < numberOfDecoys; i++ {
		// Get some random bytes.
		randomBytes := crypto.GenerateRandomBytes(64)

		// Allocate 32 bytes as the key.
		var key [32]byte
		masterPassword := randomBytes[0:32]
		copy(key[:], masterPassword)

		// Allocate 32 bytes as the identifier.
		identifier := randomBytes[32:64]
		hashedIdentifier := blake2b.Sum256(identifier)

		// Allocate 32 bytes as the plaintext.
		plaintext := make([]byte, 1025)

		// Save to the database.
		coffer.Save(hashedIdentifier[:], crypto.Encrypt(plaintext, &key))

		// Increment counter.
		count++
		fmt.Printf("\r+ Added %d/%d (%d%%)", count, numberOfDecoys, int(math.Floor(float64(count)/float64(numberOfDecoys)*100)))
	}
	fmt.Println("") // For formatting.
}
