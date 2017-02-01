package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/libeclipse/pocket/auxiliary"
	"github.com/libeclipse/pocket/coffer"
	"github.com/libeclipse/pocket/crypto"
)

var (
	scryptCost = map[string]int{"N": 18, "r": 8, "p": 1}
)

func main() {
	// Parse command line flags.
	mode, sc, err := auxiliary.ParseArgs(os.Args)
	if err != nil {
		if err.Error() == "help" {
			return
		}
		fmt.Println(err)
		return
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

	if err != nil {
		fmt.Println(err)
	}
}

func add() error {
	// Prompt user for password.
	password, err := auxiliary.GetPass()
	if err != nil {
		return err
	}

	// Prompt user for the identifier.
	identifier, err := auxiliary.Input("[-] Identifier: ")
	if err != nil {
		return err
	}

	// Prompt user for the plaintext data.
	data, err := auxiliary.Input("[-] Data: ")
	if err != nil {
		return err
	}

	// Pad the data.
	paddedData, err := crypto.Pad(data, 1025)
	if err != nil {
		return err
	}

	// Derive and store encryption key.
	fmt.Println("[+] Deriving encryption key...")
	key := crypto.DeriveKey(password, identifier, scryptCost) //LOCKTHIS

	// Encrypt the padded data.
	encryptedData, err := crypto.Encrypt(paddedData, key)
	if err != nil {
		return err
	}

	// Derive and store secure identifier.
	fmt.Println("[+] Deriving secure identifier...")
	secureIdentifier := crypto.DeriveID(identifier, scryptCost)

	// Save the identifier:data pair in the database.
	err = coffer.Save(secureIdentifier, encryptedData)
	if err != nil {
		// Cannot overwrite existing entry.
		return err
	}

	fmt.Println("[+] Okay, I'll remember that.")

	return nil
}

func retrieve() error {
	// Prompt user for password.
	password, err := auxiliary.GetPass()
	if err != nil {
		return err
	}

	// Prompt user for the identifier.
	identifier, err := auxiliary.Input("[-] Identifier: ")
	if err != nil {
		return err
	}

	// Derive and store identifier.
	fmt.Println("[+] Deriving secure identifier...")
	secureIdentifier := crypto.DeriveID(identifier, scryptCost)

	data, err := coffer.Retrieve(secureIdentifier)
	if err != nil {
		// Entry not found.
		return err
	}

	// Derive and store encryption key.
	fmt.Println("[+] Deriving encryption key...")
	key := crypto.DeriveKey(password, identifier, scryptCost)

	// Decrypt the data.
	data, err = crypto.Decrypt(data, key) //LOCKTHIS
	if err != nil {
		return err
	}

	// Unpad the data.
	data, err = crypto.Unpad(data)
	if err != nil {
		// This should never happen.
		return errors.New("[!] Invalid padding on decrypted data")
	}

	fmt.Println("[+] Data:", string(data))

	return nil
}

func forget() error {
	// Prompt user for the identifier.
	identifier, err := auxiliary.Input("[-] Identifier: ")
	if err != nil {
		return err
	}

	// Derive and store identifier.
	fmt.Println("[+] Deriving secure identifier...")
	secureIdentifier := crypto.DeriveID(identifier, scryptCost)

	// Delete the entry.
	coffer.Delete(secureIdentifier)
	fmt.Println("[+] It is forgotten.")

	return nil
}
