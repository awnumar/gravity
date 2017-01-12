package main

import (
	"fmt"
	"log"
	"os"

	"github.com/libeclipse/pocket/auxiliary"
	"github.com/libeclipse/pocket/crypto"
)

var (
	mode string

	identifier []byte

	secretData = make(map[string]interface{})
)

func main() {
	// Parse command line flags.
	mode, err := auxiliary.ParseArgs(os.Args)
	if err != nil {
		if err.Error() == "help" {
			os.Exit(0)
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// Run setup.
	auxiliary.Setup()

	// Prompt user for identifier.
	identifier = []byte(auxiliary.Input("[-] Identifier: "))
	if len(identifier) < 1 {
		fmt.Println("[!] Length of identifier must be non-zero.")
		os.Exit(1)
	}

	// Grab pre-saved secrets.
	secretData = auxiliary.RetrieveSecrets()

	// Launch appropriate function for run-mode.
	switch mode {
	case "get": // id pass
		// Prompt user for the password without echoing back.
		password := auxiliary.GetPass("[-] Password: ")
		if len(password) < 1 {
			fmt.Println("[!] Length of password must be non-zero.")
			os.Exit(1)
		}
		retrieve(password)
	case "add": // id pass secret
		// Prompt user for the password without echoing back.
		password := auxiliary.GetPass("[-] Password: ")
		if len(password) < 1 {
			fmt.Println("[!] Length of password must be non-zero.")
			os.Exit(1)
		}
		add(password)
	case "forget": // id
		forget()
	}

	// Clear sensitive data from memory.
	// (Probably not secure, but good enough.)
	identifier = []byte("")
}

func retrieve(password []byte) {
	// Derive and store identifier.
	fmt.Println("[+] Deriving secure identifier...")
	id := crypto.DeriveID(identifier)

	// Derive and store encryption key.
	fmt.Println("[+] Deriving encryption key...")
	key := crypto.DeriveKey(password, identifier)

	secret := secretData[id]
	if secret != nil {
		secret := crypto.Unpad(crypto.Decrypt(secret.(string), key))
		fmt.Println("[+] Secret:", string(secret))
	} else {
		fmt.Println("[+] There's nothing to see here.")
	}
}

func add(password []byte) {
	// Derive and store identifier.
	fmt.Println("[+] Deriving secure identifier...")
	id := crypto.DeriveID(identifier)

	// Derive and store encryption key.
	fmt.Println("[+] Deriving encryption key...")
	key := crypto.DeriveKey(password, identifier)

	// Prompt the user for the secret that we'll store.
	secret := auxiliary.Input("[-] Input secret: ")
	if len(secret) < 1 || len(secret) > 1024 {
		fmt.Println("[!] Length of secret must be between 1-1024 bytes.")
		os.Exit(1)
	}

	// Check if there's a secret there already so we don't overwrite it.
	if secretData[id] == nil {
		// Store and save the id/secret pair.
		paddedSecret, err := crypto.Pad([]byte(secret), 1025)
		if err != nil {
			log.Fatalln(err)
		}
		secretData[id] = crypto.Encrypt(paddedSecret, key)
		auxiliary.SaveSecrets(secretData)

		fmt.Println("[+] Okay, I'll remember that.")
	} else {
		// Warn that there is already data here.
		fmt.Println("[!] Cannot overwrite existing entry.")
	}
}

func forget() {
	// Derive and store identifier.
	fmt.Println("[+] Deriving secure identifier...")
	id := crypto.DeriveID(identifier)

	// Check if there's actually something there.
	if secretData[id] != nil {
		// Delete the entry. This code will never be reached if the decryption failed.
		delete(secretData, string(identifier))
		auxiliary.SaveSecrets(secretData)

		fmt.Println("[+] It is forgotten.")
	} else {
		fmt.Println("[+] Nothing to do.")
	}
}
