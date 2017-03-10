package main

import (
	"errors"
	"fmt"

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
	// THIS IS TEMPORARY UNTIL THE CLI IS DONE.
	m, _ := input.Input("[-] Mode: ")
	mode := string(m)

	// Setup the secret store.
	err := coffer.Setup()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer coffer.Close()

	// Get the secure inputs from the user.
	masterPassword, err := input.GetPass()
	if err != nil {
		fmt.Println(err)
		return
	}
	identifier, err := input.Input("[-] Identifier: ")
	if err != nil {
		fmt.Println(err)
		return
	}
	memory.Protect(identifier)

	// Derive the secure values for this session.
	masterKey, rootIdentifier = crypto.DeriveSecureValues(masterPassword, identifier, scryptCost)

	// Launch appropriate function for run-mode.
	switch mode {
	case "add":
		err = add()
	case "get":
		err = retrieve()
	case "forget":
		err = forget()
	default:
		err = errors.New("[!] Unknown mode")
	}

	// Output any errors that were returned.
	if err != nil {
		fmt.Println(err)
	}

	// Zero out and unlock any protected memory.
	memory.Cleanup()
}

func add() error {
	// Prompt user for the plaintext data.
	data, err := input.Input("[-] Data: ")
	if err != nil {
		return err
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

		// Derive ID, encrypt and save to the database.
		err := coffer.Save(crypto.DeriveIdentifierN(rootIdentifier, i/1024), crypto.Encrypt(padded, masterKey))
		if err != nil {
			return err
		}
	}

	fmt.Println("[+] Saved")

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

	if len(plaintext) == 0 {
		return errors.New("[!] Nothing to see here")
	}

	fmt.Println("[+] Data:", string(plaintext))

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

	fmt.Println("[+] It is forgotten.")

	return nil
}
