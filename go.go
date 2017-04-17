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

	"golang.org/x/crypto/blake2b"

	"github.com/libeclipse/dissident/coffer"
	"github.com/libeclipse/dissident/crypto"
	"github.com/libeclipse/dissident/memory"
	"github.com/libeclipse/dissident/metadata"
	"github.com/libeclipse/dissident/stdin"
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
			err = remove()
			if err != nil {
				return err
			}
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

	f, err := os.Open(path)
	if err != nil {
		if os.IsPermission(err) {
			fmt.Printf("! Insufficient permissions to open %s\n", path)
		} else {
			fmt.Println(err)
		}
		return
	}
	defer f.Close()

	var chunkIndex uint64
	var totalImportedBytes int64
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
		totalImportedBytes += int64(b)

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

		// Output progress.
		fmt.Printf("\r+ Imported %d%% of %d bytes...", int(math.Floor(float64(totalImportedBytes)/float64(info.Size())*100)), info.Size())
	}

	// Add the metadata to coffer.
	fmt.Println("+ Adding metadata...")
	metadata.New()
	metadata.Set(totalImportedBytes, "length")
	metadata.Save(rootIdentifier, masterKey)
	metadata.Reset()

	fmt.Println("\n+ Imported successfully.")
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
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if os.IsExist(err) {
			fmt.Printf("! %s already exists; cannot overwrite\n", path)
		} else if os.IsPermission(err) {
			fmt.Printf("! Insufficient permissions to open %s\n", path)
		} else {
			fmt.Println(err)
		}
		return
	}
	defer f.Close()

	// Get the metadata first.
	metadata.New()
	metadata.Retrieve(rootIdentifier, masterKey)
	lenData := metadata.GetLength("length")
	metadata.Reset()

	// Grab the data.
	totalExportedBytes := 0
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
		totalExportedBytes += len(unpadded)
		memory.Wipe(pt)

		// Write and wipe data.
		f.Write(unpadded)
		memory.Wipe(unpadded)

		// Output progress.
		fmt.Printf("\r+ Exported %d%% of %d bytes...", int(math.Floor(float64(totalExportedBytes)/float64(lenData)*100)), lenData)
	}

	fmt.Printf("\n+ Saved to %s\n", path)
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
	fmt.Println("\n-----BEGIN PLAINTEXT-----")

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
		memory.Wipe(pt)

		// Write and wipe data.
		fmt.Print(string(unpadded))
		memory.Wipe(unpadded)
	}

	fmt.Println("-----END PLAINTEXT-----")
}

func remove() error {
	// Prompt the user for the identifier.
	identifier := stdin.Secure("- Secure identifier: ")

	// Derive the secure values for this "branch".
	fmt.Println("+ Generating root key...")
	_, rootIdentifier := crypto.DeriveSecureValues(masterPassword, identifier, scryptCost)

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
	}

	if count != 0 {
		fmt.Println("+ Successfully removed data.")
	} else {
		fmt.Println("! There is nothing here to remove")
	}

	return nil
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

	count := 0
	for i := 0; i < numberOfDecoys; i++ {
		// Get some random bytes.
		randomBytes := crypto.GenerateRandomBytes(64)

		// Allocate 32 bytes as the key.
		var key [32]byte
		masterPassword := randomBytes[0:32]
		copy(key[:], masterPassword)

		// Allocate 32 bytes as the identifier.
		identifier := randomBytes[32:64]
		hashedIdentifier := blake2b.Sum256(identifier)

		// Allocate 32 bytes as the plaintext.
		plaintext := make([]byte, 4096)

		// Save to the database.
		coffer.Save(hashedIdentifier[:], crypto.Encrypt(plaintext, &key))

		// Increment counter and display progress.
		count++
		fmt.Printf("\r+ Added %d%% of %d decoys...", int(math.Floor(float64(count)/float64(numberOfDecoys)*100)), numberOfDecoys)
	}
	fmt.Println("") // For formatting.
}
