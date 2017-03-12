package main

import (
	"fmt"
	"strings"

	"github.com/libeclipse/pocket/coffer"
	"github.com/libeclipse/pocket/crypto"
	"github.com/libeclipse/pocket/input"
	"github.com/libeclipse/pocket/memory"
)

var (
	// The default cost factor for key deriviation.
	scryptCost = map[string]int{"N": 18, "r": 16, "p": 1}

	// Store the current session's secure values.
	masterKey      *[32]byte
	rootIdentifier []byte
)

func main() {
	// Setup the secret store.
	err := coffer.Setup()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer coffer.Close()

	// Launch CLI.
	err = cli()
	if err != nil {
		fmt.Println(err)
	}

	// Zero out and unlock any protected memory.
	memory.Cleanup()
}

func cli() error {
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

	fmt.Println(`
:: You find yourself in a field containing an infinite number of trees,
each with an infinity of branches. You are told that there are secrets here
of the darkest kind, but no one has any idea of where they are.

:: It is said many a great traveller has let life slip through his fingers
like sand, searching this tantalising place.

:: I see in your eyes that you are determined to do the same. In that case,
I will need the name of the tree that you wish to view...
`)

	masterPassword, err := input.GetPass()
	if err != nil {
		return err
	}

	fmt.Println(`
:: Ah, a good choice. But I fear you do not have the time to search it all.

:: The chosen tree has an infinite number of branches. Which one would you
like to take a closer look at?
`)

	identifier := input.Input("- Identifier: ")

	fmt.Println("\n:: Climbing tree...")

	// Derive the secure values for this "branch".
	masterKey, rootIdentifier = crypto.DeriveSecureValues(masterPassword, identifier, scryptCost)

	// Check if there's something here.
	_, exists := coffer.Retrieve(crypto.DeriveIdentifierN(rootIdentifier, 0))
	if exists != nil {
		// Nope, nothing here.
		fmt.Println("\n:: It doesn't look like there's anything here.")

		conf := input.Input("\n- Would you like to hide something? ")
		if strings.Contains(strings.ToLower(string(conf)), "y") {
			return add()
		}
	} else {
		fmt.Println("\n:: Oh, you found something!")

		conf := string(input.Input("\n- Would you like to (V)iew or (R)emove this? "))
		if strings.ToLower(conf) == "v" {
			return retrieve()
		} else if strings.ToLower(conf) == "r" {
			return forget()
		}
	}

	return nil
}

func add() error {
	// Prompt user for the plaintext data.
	data := input.Input("\n- Enter the data that you wish to hide: ")

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
			return err
		}
	}

	fmt.Println("\n:: Ah, that is done. I doubt anyone will ever find it.")

	return nil
}

func retrieve() error {
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
		pt, e := crypto.Decrypt(ct, masterKey)
		if e != nil {
			return e
		}

		// Unpad this slice.
		unpadded, e := crypto.Unpad(pt)
		if e != nil {
			return e
		}

		// Append this slice of plaintext to the rest of it.
		plaintext = append(plaintext, unpadded...)
	}

	fmt.Println("\n:: Here it is:")

	fmt.Printf(`
-----BEGIN SECRET PLAINTEXT-----
%s
-----END SECRET PLAINTEXT-----
`, plaintext)

	return nil
}

func forget() error {
	// Delete all the pieces.
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
	}

	fmt.Println("\n:: You successfully destroyed whatever was stored here.")

	return nil
}
