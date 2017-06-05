package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/0xAwn/dissident/coffer"
	"github.com/0xAwn/dissident/crypto"
	"github.com/0xAwn/dissident/data"
	"github.com/0xAwn/dissident/stdin"
	"github.com/0xAwn/memguard"
	"github.com/cheggaaa/pb"
)

var (
	// The default cost factor for key deriviation.
	scryptCost = map[string]int{"N": 18, "r": 16, "p": 1}

	// Store the container ID globally.
	masterPassword *memguard.LockedBuffer
)

func main() {
	// Setup the secret store.
	err := coffer.Setup()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer coffer.Close()

	// Cleanup memory when exiting.
	memguard.CatchInterrupt(func() {})
	defer memguard.DestroyAll()

	// Launch CLI.
	err = cli()
	if err != nil {
		fmt.Println(err)
	}
}

func cli() error {
	var err error

	help := `import [path] - Import a new file to the database.
export [path] - Retrieve data from the database and export to a file.
peak          - Grab data from the database and print it to the screen.
remove        - Remove some previously stored data from the database.
decoys        - Add a variable amount of random decoy data.
exit          - Exit the program.`

	masterPassword, err = stdin.GetMasterPassword()
	if err != nil {
		return err
	}
	fmt.Println("") // For formatting.

	for {
		cmd := strings.Split(strings.TrimSpace(stdin.Standard("$ ")), " ")

		switch cmd[0] {
		case "import":
			if len(cmd) < 2 {
				fmt.Println("! Missing argument: path")
			} else {
				importFromDisk(cmd[1])
			}
		case "export":
			if len(cmd) < 2 {
				fmt.Println("! Missing argument: path")
			} else {
				exportToDisk(cmd[1])
			}
		case "peak":
			peak()
		case "remove":
			remove()
		case "decoys":
			decoys()
		case "exit":
			return nil
		default:
			fmt.Println(help)
		}
	}
}

func importFromDisk(path string) {
	// Handle the file.
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("! %s does not exist\n", path)
		} else {
			fmt.Println(err)
		}
		return
	}

	if info.IsDir() {
		fmt.Println("! We can't handle directories yet")
		return
	}

	// Prompt the user for the identifier.
	identifier := stdin.Secure("- Secure identifier: ")

	// Derive the secure values for this "branch".
	fmt.Println("+ Generating root key...")
	masterKey, rootIdentifier := crypto.DeriveSecureValues(masterPassword, identifier, scryptCost)

	// Check if it exists already.
	derivedIdentifierN := crypto.DeriveIdentifierN(rootIdentifier, 0)
	if coffer.Exists(derivedIdentifierN) {
		fmt.Println("! Cannot overwrite existing entry")
		return
	}

	// Add the metadata to coffer.
	fmt.Println("+ Adding metadata...")
	data.MetaSetLength(info.Size(), rootIdentifier, masterKey)

	// Import this entry from disk.
	data.ImportData(path, info.Size(), rootIdentifier, masterKey)

	// Output status message.
	fmt.Println("+ Imported successfully.")
}

func exportToDisk(path string) {
	// Prompt the user for the identifier.
	identifier := stdin.Secure("- Secure identifier: ")

	// Derive the secure values for this "branch".
	fmt.Println("+ Generating root key...")
	masterKey, rootIdentifier := crypto.DeriveSecureValues(masterPassword, identifier, scryptCost)

	// Check if this entry exists.
	derivedIdentifierN := crypto.DeriveIdentifierN(rootIdentifier, 0)
	if !coffer.Exists(derivedIdentifierN) {
		fmt.Println("! This entry does not exist")
		return
	}

	// Export the entry.
	data.ExportData(path, rootIdentifier, masterKey)
}

func peak() {
	// Prompt the user for the identifier.
	identifier := stdin.Secure("- Secure identifier: ")

	// Derive the secure values for this "branch".
	fmt.Println("+ Generating root key...")
	masterKey, rootIdentifier := crypto.DeriveSecureValues(masterPassword, identifier, scryptCost)

	// Check if this entry exists.
	derivedIdentifierN := crypto.DeriveIdentifierN(rootIdentifier, 0)
	if !coffer.Exists(derivedIdentifierN) {
		fmt.Println("! This entry does not exist")
		return
	}

	// It exists, proceed to get data.
	data.ViewData(rootIdentifier, masterKey)
}

func remove() {
	// Prompt the user for the identifier.
	identifier := stdin.Secure("- Secure identifier: ")

	// Derive the secure values for this "branch".
	fmt.Println("+ Generating root key...")
	masterKey, rootIdentifier := crypto.DeriveSecureValues(masterPassword, identifier, scryptCost)

	// Check if this entry exists.
	derivedIdentifierN := crypto.DeriveIdentifierN(rootIdentifier, 0)
	if !coffer.Exists(derivedIdentifierN) {
		fmt.Println("! There is nothing here to remove")
		return
	}

	// Remove the data.
	data.RemoveData(rootIdentifier, masterKey)
}

func decoys() {
	var numberOfDecoys int
	var err error

	// Print some help information.
	fmt.Println(`
:: For deniable encryption, use this feature in conjunction with some fake data manually-added
   under a different master-password. Then if you are ever forced to hand over your keys,
   simply give up the fake data and claim that the rest of the entries in the database are decoys.

:: You do not necessarily have to make use of this feature. Rather, simply the fact that
   it exists allows you to claim that some or all of the entries in the database are decoys.
`)

	// Get the number of decoys to add as an int.
	for {
		numberOfDecoys, err = strconv.Atoi(stdin.Standard("How many decoys do you want to add? "))
		if err == nil {
			break
		}
		fmt.Println("! Input must be an integer")
	}

	// Create and configure the progress bar object.
	bar := pb.New64(int64(numberOfDecoys)).Prefix("+ Adding ")
	bar.ShowSpeed = true
	bar.SetUnits(pb.U_NO)
	bar.Start()

	for i := 0; i < numberOfDecoys; i++ {
		// Generate the decoy.
		identifier, ciphertext := crypto.GenDecoy()

		// Save to the database.
		coffer.Save(identifier, ciphertext)

		// Increment progress bar.
		bar.Increment()
	}
	bar.FinishPrint(fmt.Sprintf("+ Added %d decoys.", numberOfDecoys))
}
