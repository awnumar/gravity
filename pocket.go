package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/libeclipse/pocket/auxiliary"
	"github.com/libeclipse/pocket/crypto"
	"github.com/libeclipse/pocket/store"
	"github.com/libeclipse/pocket/store/coffer"
)

var (
	// options
	storeFlag  = flag.String("store", "coffer", "set which store to use")
	scryptCost = map[string]int{"N": 18, "r": 8, "p": 1}

	secretStore store.Store
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

	// Determine which store to use
	switch *storeFlag {
	default:
		fallthrough
	case "coffer":
		var err error
		secretStore, err = coffer.Setup()
		if err != nil {
			log.Fatalln(err)
		}
	}

	// Launch appropriate function for run-mode.
	switch mode {
	case "add":
		add()
	case "get":
		retrieve()
	case "forget":
		forget()
	}
}

func add() {
	// Get values from the user.
	values := auxiliary.GetInputs([]string{"password", "identifier", "data"})

	// Derive and store identifier.
	fmt.Println("[+] Deriving secure identifier...")
	identifier := crypto.DeriveID([]byte(values[1]), scryptCost)

	// Derive and store encryption key.
	fmt.Println("[+] Deriving encryption key...")
	key := crypto.DeriveKey([]byte(values[0]), []byte(values[1]), scryptCost)

	// Store and save the id/data pair.
	paddedData, err := crypto.Pad([]byte(values[2]), 1025)
	if err != nil {
		log.Fatalln(err)
	}

	// Encrypt the padded data.
	encryptedData := crypto.Encrypt(paddedData, key)

	// Save the identifier:data pair in the database.
	err = secretStore.Save(identifier, encryptedData)
	if err != nil {
		// Cannot overwrite existing entry.
		fmt.Println(err)
	} else {
		fmt.Println("[+] Okay, I'll remember that.")
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

	data, err := secretStore.Retrieve(identifier)
	if err != nil {
		// Entry not found.
		fmt.Println(err)
	} else {
		data, err = crypto.Unpad(crypto.Decrypt(data, key))
		if err != nil {
			// This should never happen.
			fmt.Println("[!] Invalid padding on decrypted data")
		} else {
			fmt.Println("[+] Data:", string(data))
		}
	}
}

func forget() {
	// Get values from the user.
	values := auxiliary.GetInputs([]string{"identifier"})

	// Derive and store identifier.
	fmt.Println("[+] Deriving secure identifier...")
	identifier := crypto.DeriveID([]byte(values[0]), scryptCost)

	// Delete the entry.
	secretStore.Delete(identifier)
	fmt.Println("[+] It is forgotten.")
}
