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
	}

	// Run setup.
	auxiliary.Setup()

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

	secret, err := auxiliary.RetrieveSecret(identifier)
	if err != nil {
		fmt.Println(err)
	} else {
		secret, err = crypto.Unpad(crypto.Decrypt(secret, key))
		if err != nil {
			// This should never happen.
			fmt.Println("[!] Invalid padding on decrypted secret.")
		} else {
			fmt.Println("[+] Secret:", string(secret))
		}
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

	// Store and save the id/secret pair.
	paddedSecret, err := crypto.Pad([]byte(values[2]), 1025)
	if err != nil {
		log.Fatalln(err)
	}

	// Encrypt the padded secret.
	encryptedSecret := crypto.Encrypt(paddedSecret, key)

	// Save the identifier:secret pair in the database.
	err = auxiliary.SaveSecret(identifier, encryptedSecret)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("[+] Okay, I'll remember that.")
	}
}

func forget() {
	// Get values from the user.
	values := auxiliary.GetInputs([]string{"identifier"})

	// Derive and store identifier.
	fmt.Println("[+] Deriving secure identifier...")
	identifier := crypto.DeriveID([]byte(values[0]), scryptCost)
	_ = identifier
}
