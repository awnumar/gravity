package main

import (
	"flag"
	"fmt"

	"github.com/libeclipse/pocket/auxiliary"
	"github.com/libeclipse/pocket/crypto"
)

var (
	// This stores the mode that we're running in.
	// 0 => Retrieve secret.
	// 1 => Store secret.
	// 2 => Forget secret.
	mode int

	key        []byte
	identifier string

	secretData = make(map[string]interface{})
)

func main() {
	// Command line flag to determine mode at runtime.
	flag.IntVar(&mode, "m", 0, "specify mode: 0 => retrieve (default); 1 => store; 2 => forget")

	// Parse the flags.
	flag.Parse()

	// Verify mode is valid.
	if mode > 2 || mode < 0 {
		fmt.Println("[!] invalid mode; use -h for help")
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
	case 0:
		retrieve()
	case 1:
		store()
	case 2:
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

func store() {
	// Prompt the user for the secret that we'll store.
	secret := auxiliary.Input("[-] secret: ")

	// Check if there's a secret there already so we don't overwrite it.
	if secretData[identifier] == nil {
		// Store and save the id/secret pair.
		secretData[identifier] = crypto.Encrypt(secret, key)
		auxiliary.SaveSecrets(secretData)
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
	}
}
