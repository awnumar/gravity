package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/cheggaaa/pb"
	"github.com/libeclipse/dissident/coffer"
	"github.com/libeclipse/dissident/crypto"
	"github.com/libeclipse/dissident/disk"
	"github.com/libeclipse/dissident/memory"
	"github.com/libeclipse/dissident/metadata"
	"github.com/libeclipse/dissident/stdin"
	"github.com/libeclipse/dissident/ui"
)

var (
	// The default cost factor for key deriviation.
	scryptCost = map[string]int{"N": 18, "r": 16, "p": 1}

	// Store the container ID globally.
	masterPassword []byte
)

func main() {
	// Setup the secret store.
	err := coffer.Setup()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer coffer.Close()

	// CleanupMemory in case of Ctrl+C
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		memory.SafeExit(0)
	}()

	// Launch CLI.
	err = cli()
	if err != nil {
		fmt.Println(err)
	}

	// Zero out and unlock any protected memory.
	memory.Cleanup()
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
	info, err := disk.GetFileInfo(path)
	if err != nil {
		fmt.Println(err)
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

	f, err := disk.OpenFileRead(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	// Add the metadata to coffer.
	fmt.Println("+ Adding metadata...")
	metadata.New()
	metadata.Set(info.Size(), "length")
	metadata.Save(rootIdentifier, masterKey)
	metadata.Reset()

	// Start the progress bar.
	bar := ui.StartBar(info.Size(), "+ Importing ", pb.U_BYTES, true, true)

	// Import the data.
	var chunkIndex uint64
	buffer := make([]byte, 4095)
	for {
		b, err := f.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return
		}
		bar.Add(b) // Increment the progress bar.

		data := make([]byte, b)
		copy(data, buffer[:b])

		// Pad data and wipe the buffer.
		data, err = crypto.Pad(data, 4096)
		if err != nil {
			fmt.Println(err)
			return
		}
		memory.Wipe(buffer)

		// Save it and wipe plaintext.
		coffer.Save(crypto.DeriveIdentifierN(rootIdentifier, chunkIndex), crypto.Encrypt(data, masterKey))
		memory.Wipe(data)

		// Increment counter.
		chunkIndex++
	}
	// We're done. End the progress bar.
	bar.Finish()

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

	// Atempt to open the file now.
	f, err := disk.OpenFileAppend(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	// Get the metadata first.
	metadata.New()
	metadata.Retrieve(rootIdentifier, masterKey)
	lenData := metadata.GetLength("length")
	metadata.Reset()

	// Start the progress bar object.
	bar := ui.StartBar(lenData, "+ Exporting ", pb.U_BYTES, true, true)

	// Grab the data.
	for n := new(uint64); true; *n++ {
		// Derive derived_identifier[n]
		ct := coffer.Retrieve(crypto.DeriveIdentifierN(rootIdentifier, *n))
		if ct == nil {
			// This one doesn't exist. //EOF
			break
		}

		// Decrypt this slice.
		pt, err := crypto.Decrypt(ct, masterKey)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Unpad this slice and wipe old one.
		unpadded, e := crypto.Unpad(pt)
		if e != nil {
			fmt.Println(e)
			return
		}
		bar.Add(len(unpadded)) // Increment the progress bar.
		memory.Wipe(pt)

		// Write and wipe data.
		f.Write(unpadded)
		memory.Wipe(unpadded)
	}
	// We're done. End the progress bar.
	bar.FinishPrint(fmt.Sprintf("+ Saved to %s", path))

	// Compare length in metadata to actual exported length.
	if bar.Get() != lenData {
		fmt.Println("! Data incomplete; database may be corrupt")
	}
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

	// It exists, proceed.

	// Get the metadata first.
	metadata.New()
	metadata.Retrieve(rootIdentifier, masterKey)
	lenData := metadata.GetLength("length")
	metadata.Reset()

	fmt.Println("\n-----BEGIN PLAINTEXT-----")

	var totalExportedBytes int64
	for n := new(uint64); true; *n++ {
		// Derive derived_identifier[n]
		ct := coffer.Retrieve(crypto.DeriveIdentifierN(rootIdentifier, *n))
		if ct == nil {
			// This one doesn't exist. //EOF
			break
		}

		// Decrypt this slice.
		pt, err := crypto.Decrypt(ct, masterKey)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Unpad this slice and wipe old one.
		unpadded, e := crypto.Unpad(pt)
		if e != nil {
			fmt.Println(e)
			return
		}
		totalExportedBytes += int64(len(unpadded))
		memory.Wipe(pt)

		// Write and wipe data.
		fmt.Print(string(unpadded))
		memory.Wipe(unpadded)
	}

	fmt.Println("-----END PLAINTEXT-----")

	// Compare length in metadata to actual exported length.
	if totalExportedBytes != lenData {
		fmt.Println("! Data incomplete; database may be corrupt")
	}
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

	// Get the metadata first.
	metadata.New()
	metadata.Retrieve(rootIdentifier, masterKey)
	lenData := metadata.GetLength("length")
	metadata.Reset()

	// Start the progress bar.
	bar := ui.StartBar(int64(math.Ceil(float64(lenData)/4096)), "+ Removing ", pb.U_NO, false, false)
	bar.ShowCounters = false
	bar.Start()

	// Remove all metadata.
	metadata.Remove(rootIdentifier)

	// Delete all the pieces.
	count := 0
	for n := new(uint64); true; *n++ {
		// Get the DeriveIdentifierN for this n.
		derivedIdentifierN := crypto.DeriveIdentifierN(rootIdentifier, *n)

		// Check if it exists.
		if coffer.Exists(derivedIdentifierN) {
			coffer.Delete(derivedIdentifierN)
			count++
		} else {
			break
		}

		// Increment progress bar.
		bar.Increment()
	}
	// We're done. End the progress bar.
	bar.FinishPrint("+ Successfully removed data.")
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
	bar := ui.StartBar(int64(numberOfDecoys), "+ Adding ", pb.U_NO, true, true)

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
