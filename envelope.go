package main

import (
	"flag"
	"fmt"

	"github.com/libeclipse/envelope/auxiliary"
)

var (
	// This stores the mode that we're running in.
	// 0 => Retrieve secret.
	// 1 => Store secret.
	// 2 => Forget secret.
	mode int

	masterPassword []byte
	identifier     []byte

	secretData = make(map[string]interface{})
)

func main() {
	// Command line flag to determine mode at runtime.
	flag.IntVar(&mode, "m", 0, "specify mode: 0 => retrieve (default); 1 => store; 2 => forget")

	// Parse the flags.
	flag.Parse()

	// Run setup.
	auxiliary.Setup()

	// Prompt user for the master password.
	masterPassword = auxiliary.GetPass("[-] master password: ")

	// Prompt user for identifier.
	identifier = []byte(auxiliary.Input("[-] identifier: "))

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
}

func retrieve() {
	secret := secretData[string(identifier)]
	if secret != nil {
		fmt.Println("[+] secret:", secret)
	} else {
		fmt.Println("[+] nothing to see here")
	}
}

func store() {
	secret := auxiliary.Input("[-] secret: ")
	// TODO: check if secret exists before overwriting.
	secretData[string(identifier)] = secret
	auxiliary.SaveSecrets(secretData)
}

func forget() {
	fmt.Println("[+] ready to forget secret")
}
