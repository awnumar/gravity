package main

import (
	"fmt"
	"os"

	"github.com/libeclipse/pocket/auxiliary"
	"github.com/libeclipse/pocket/crypto"
)

var (
	mode string

	key        []byte
	identifier string

	secretData = make(map[string]interface{})
)

func main() {
	// Parse command line flags.
	if len(os.Args) < 2 {
		fmt.Println("[!] mode not specified; use `pocket help`")
		os.Exit(1)
	}

	argument := os.Args[1]
	if argument == "-h" || argument == "--help" || argument == "help" || argument == "-help" {
		fmt.Printf("Usage: %s [get|add|forget]\n", os.Args[0])
		os.Exit(2)
	} else {
		mode = argument
	}

	// Verify that mode is valid.
	if mode != "get" && mode != "add" && mode != "forget" {
		fmt.Println("[!] invalid mode; use `pocket help`")
		os.Exit(1)
	}

	// Run setup.
	auxiliary.Setup()

	// Prompt user for the password without echoing back.
	password := auxiliary.GetPass("[-] password: ")

	// Prompt user for identifier.
	id := []byte(auxiliary.Input("[-] identifier: "))

	// Derive and store encryption key.
	fmt.Println("[+] deriving encryption key...")
	key = crypto.DeriveKey(password, id)

	// Derive and store identifier.
	fmt.Println("[+] deriving secure identifier...")
	identifier = crypto.DeriveID(id)

	// Grab pre-saved secrets.
	secretData = auxiliary.RetrieveSecrets()

	// Launch appropriate function for run-mode.
	switch mode {
	case "get":
		retrieve()
	case "add":
		add()
	case "forget":
		forget()
	}

	// Clear sensitive data from memory. (Probably not secure, but good enough.)
	key = []byte("")
	identifier = ""
}

func retrieve() {
	secret := secretData[identifier]
	if secret != nil {
		fmt.Println("[+] secret:", crypto.Decrypt(secret.(string), key))
	} else {
		fmt.Println("[+] nothing to see here")
	}
}

func add() {
	// Prompt the user for the secret that we'll store.
	secret := auxiliary.Input("[-] secret: ")

	// Check if there's a secret there already so we don't overwrite it.
	if secretData[identifier] == nil {
		// Store and save the id/secret pair.
		secretData[identifier] = crypto.Encrypt([]byte(secret), key)
		auxiliary.SaveSecrets(secretData)

		fmt.Println("[+] ok, i'll remember that")
	} else {
		// Warn that there is already data here.
		fmt.Println("[!] cannot overwrite existing entry")
	}
}

func forget() {
	// Check if there's actually something there.
	if secretData[identifier] != nil {
		// Decryption here serves no cryptographic purpose. The reason for it is
		// so that deleting the entry through the application isn't trivial. Of
		// course the attacker could still simply just `rm -rf ~/.pocket/secrets`
		crypto.Decrypt(secretData[identifier].(string), key)

		// Delete the entry. This code will never be reached if the decryption failed.
		delete(secretData, string(identifier))
		auxiliary.SaveSecrets(secretData)

		fmt.Println("[+] it is forgotten")
	} else {
		fmt.Println("[+] nothing to do")
	}
}
