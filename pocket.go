package main

import (
	"fmt"
	"os"

	"github.com/libeclipse/pocket/auxiliary"
	"github.com/libeclipse/pocket/coffer"
	"github.com/libeclipse/pocket/crypto"
)

var (
	scryptCost = map[string]int{"N": 18, "r": 16, "p": 1}
)

func main() {
	// Parse command line flags.
	mode, sc, err := auxiliary.ParseArgs(os.Args)
	if err != nil && err.Error() != "help" {
		fmt.Println(err)
		os.Exit(1)
	}

	if sc != nil {
		scryptCost = sc
	}

	// Setup the secret store.
	err = coffer.Setup()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer coffer.Close()

	// Launch appropriate function for run-mode.
	switch mode {
	case "add":
		err = add()
	case "get":
		err = retrieve()
	case "forget":
		err = forget()
	}

	// Output any errors that were returned.
	if err != nil {
		fmt.Println(err)
	}

	// Zero out and unlock any protected memory.
	crypto.CleanupMemory()
}

func add() error {
	// Prompt for masterPassword and identifier.
	masterPassword, identifier, err := auxiliary.GetPassAndID()
	if err != nil {
		return err
	}

	// Derive rootKey and rootIdentifier.
	masterKey, rootIdentifier := crypto.DeriveSecureValues(masterPassword, identifier, scryptCost)

	// Prompt user for the plaintext data.
	data, err := auxiliary.Input("[-] Data: ")
	if err != nil {
		return err
	}

	var padded []byte
	for i := 0; i < len(data); i += 1024 {
		if i+1024 > len(data) {
			// Remaining data <= 1024.
			padded, err = crypto.Pad(data[len(data)-(len(data)%1024):len(data)], 1025)
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
	// Prompt for masterPassword and identifier.
	masterPassword, identifier, err := auxiliary.GetPassAndID()
	if err != nil {
		return err
	}

	// Derive rootKey and rootIdentifier.
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

	fmt.Println("[+] Data:", string(plaintext))

	return nil
}

func forget() error {
	return nil
}
