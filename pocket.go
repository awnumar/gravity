package main

import (
	"fmt"
	"log"
	"os"

	"github.com/libeclipse/pocket/auxiliary"
	"github.com/libeclipse/pocket/crypto"
)

var (
	scryptCost = map[string]int{"N": 18, "r": 8, "p": 1}
	secretData = make(map[string]interface{})
)

func main() {
	// Parse command line flags.
	mode, sc, err := auxiliary.ParseArgs(os.Args)
	if err != nil {
		if err.Error() == "help" {
			os.Exit(0)
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if sc != nil {
		scryptCost = sc
		fmt.Println(scryptCost)
	}

	// Run setup.
	auxiliary.Setup()

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
}

func retrieve() {
	// Get values from the user.
	values := auxiliary.GetInputs([]string{"password", "identifier"})

	// Derive and store identifier.
	fmt.Println("[+] Deriving secure identifier...")
	identifier := crypto.DeriveID([]byte(values[1]), scryptCost)

	// Derive and store encryption key.
	fmt.Println("[+] Deriving encryption key...")
	key := crypto.DeriveKey([]byte(values[0]), []byte(values[1]), scryptCost)

	secret := secretData[identifier]
	if secret != nil {
		secret, err := crypto.Unpad(crypto.Decrypt(secret.(string), key))
		if err != nil {
			// This should never happen.
			fmt.Println("[!] Invalid padding on decrypted secret.")
		}
		fmt.Println("[+] Secret:", string(secret))
	} else {
		fmt.Println("[+] There's nothing to see here.")
	}
}

func add() {
	// Get values from the user.
	values := auxiliary.GetInputs([]string{"password", "identifier", "secret"})

	// Derive and store identifier.
	fmt.Println("[+] Deriving secure identifier...")
	identifier := crypto.DeriveID([]byte(values[1]), scryptCost)

	// Derive and store encryption key.
	fmt.Println("[+] Deriving encryption key...")
	key := crypto.DeriveKey([]byte(values[0]), []byte(values[1]), scryptCost)

	// Check if there's a secret there already so we don't overwrite it.
	if secretData[identifier] == nil {
		// Store and save the id/secret pair.
		paddedSecret, err := crypto.Pad([]byte(values[2]), 1025)
		if err != nil {
			log.Fatalln(err)
		}
		secretData[identifier] = crypto.Encrypt(paddedSecret, key)
		auxiliary.SaveSecrets(secretData)

		fmt.Println("[+] Okay, I'll remember that.")
	} else {
		// Warn that there is already data here.
		fmt.Println("[!] Cannot overwrite existing entry.")
	}
}

func forget() {
	// Get values from the user.
	values := auxiliary.GetInputs([]string{"identifier"})

	// Derive and store identifier.
	fmt.Println("[+] Deriving secure identifier...")
	identifier := crypto.DeriveID([]byte(values[0]), scryptCost)

	// Check if there's actually something there.
	if secretData[identifier] != nil {
		// Delete the entry.
		delete(secretData, string(identifier))
		auxiliary.SaveSecrets(secretData)

		fmt.Println("[+] It is forgotten.")
	} else {
		fmt.Println("[+] Nothing to do.")
	}
}
