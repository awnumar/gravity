package main

import (
	"flag"
	"fmt"
	"os"

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
	flag.IntVar(&mode, "m", 0, "specify mode. 0 => retrieve secret (default); 1 => store secret; 2 => forget secret\n")

	// Specify the string that identifies the secret that we'll retrieve/store.
	id := flag.String("i", "", "identifies the entry that we'll store/retrieve; in case of your master password leaking, a strong id here may still protect your data\n")

	// Parse the flags.
	flag.Parse()

	// Verify that an identifier was specified.
	if len(*id) != 0 {
		identifier = []byte(*id)
	} else {
		fmt.Println("[!] id not specified; use -h for help")
		os.Exit(1)
	}

	// Run setup.
	auxiliary.Setup()

	// Prompt user for masterPassword.
	fmt.Print("[-] master password: ")
	masterPassword = auxiliary.GetPass()

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
	fmt.Println("[+] retrieving secret...")
	secret := secretData[string(identifier)]
	fmt.Println("[+] secret:", secret)
}

func store() {
	fmt.Println("[+] storing another secret...")
	secret := auxiliary.Input("[-] secret: ")
	// TODO: check if secret exists before overwriting.
	secretData[string(identifier)] = secret
	auxiliary.SaveSecrets(secretData)
}

func forget() {
	fmt.Println("[+] ready to forget secret")
}
