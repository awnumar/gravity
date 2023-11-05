package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/awnumar/memguard"
	"github.com/docker/go-units"
)

var args = os.Args

func main() {
	cleanup := func() {
		// Sync and close disk-backed database.
		closeDB()

		// Purge sensitive information from memory.
		memguard.Purge()
		fmt.Println("[i] Memory successfully purged... exiting.")
	}
	memguard.CatchSignal(func(_ os.Signal) {
		cleanup() // Call cleanup on catching a signal.
	})
	defer cleanup() // Run cleanup after returning.

	// Parse command line arguments.
	if args[1] == "seal" {
		if len(args) != 3 {
			goto help
		}

		// Get files.
		files, err := Files(args[2])
		if err != nil {
			outputError(err)
			return
		}

		// Output state.
		totalSize := func() (total int64) {
			for file := range files {
				total += files[file].Size
			}
			return
		}()
		fmt.Printf("[i] Encrypting %d files from \"%s\" (%s)\n", len(files), args[2], units.BytesSize(float64(totalSize)))

		// Read key from standard input directly into secure buffer.
		key, err := input("[?] Enter master key: ")
		if err != nil {
			outputError(err)
			return
		}

		// Derive root key from user key.
		fmt.Println("[i] Processing key...")
		pocket := GetPocket(key)

		// Initialise identifier.
		id, idMemory, err := pocket.Identifier()
		if err != nil {
			outputError(err)
			return
		}

		// Get local copy of key
		key, err = pocket.Key.Open()
		if err != nil {
			outputError(err)
			return
		}

		// Process each file
		var buffer [4096]byte
		for file, fileInfo := range files {
			fmt.Printf("[+] Sealing %s (%s)\n", fileInfo.Path, units.BytesSize(float64(fileInfo.Size)))

			// Handle metadata
			metadata, err := json.Marshal(fileInfo)
			if err != nil {
				outputError(err)
				return
			}
			for i := 0; i < len(metadata); i += 4095 {
				var size int
				if i+4095 > len(metadata) {
					size = copy(buffer[:], metadata[len(metadata)-(len(metadata)%4095):])
				} else {
					size = copy(buffer[:], metadata[i:i+4095])
				}
				buffer[size] = 1 // padding

				ct, _ := Encrypt(buffer[:], key.Bytes())
				if err := Put(id.Derive(idMemory, uint64(file), uint64(2*i/4095+1)), ct); err != nil {
					outputError(err)
					return
				}
				memguard.WipeBytes(buffer[:])
			}

			// Handle file contents
			f, err := os.Open(fileInfo.Path)
			if err != nil {
				outputError(err)
				return
			}
			for c := uint64(0); ; c += 2 {
				n, _ := io.ReadFull(f, buffer[:4095])
				if n == 0 {
					break
				}
				buffer[n] = 1 // padding
				ct, _ := Encrypt(buffer[:], key.Bytes())
				if err := Put(id.Derive(idMemory, uint64(file), c), ct); err != nil {
					outputError(err)
					return
				}
				memguard.WipeBytes(buffer[:])
			}

		}

		return
	} else if args[1] == "open" {
		if len(args) != 3 {
			goto help
		}

		if err := os.Mkdir(args[2], os.ModeDir|os.ModePerm); err != nil {
			outputError(err)
			return
		}
		if err := os.Chdir(args[2]); err != nil {
			outputError(err)
			return
		}

		// Read key from standard input directly into secure buffer.
		key := input("[?] Enter master key: ")

		// Derive root key from user key.
		fmt.Println("[i] Processing key...")
		pocket := GetPocket(key)

		// Initialise identifier.
		id, idMemory, err := pocket.Identifier()
		if err != nil {
			outputError(err)
			return
		}

		// Get local copy of key
		key, err = pocket.Key.Open()
		if err != nil {
			outputError(err)
			return
		}

		// Extract data
		var buffer [4096]byte
		for i := uint64(0); ; i++ { // for every file...
			// Handle metadata
			var metadata []byte
			for j := uint64(1); ; j += 2 {
				chunk, err := Get(id.Derive(idMemory, i, j))
				if err != nil {
					// eof
					break
				}
				n, err := Decrypt(chunk, key.Bytes(), buffer[:])
				if err != nil {
					outputError(err)
					return
				}
				if n != 4096 {
					outputError(errors.New("error invalid plaintext size"))
					return
				}
				for k := len(buffer) - 1; ; k-- {
					if buffer[k] == 0 {
						continue
					} else if buffer[k] == 1 {
						metadata = append(metadata, buffer[:k]...)
						break
					} else {
						outputError(errors.New("error invalid padding"))
						return
					}
				}
				memguard.WipeBytes(buffer[:])
			}
			if len(metadata) == 0 {
				//eof
				break
			}
			var md FileInfo
			if err := json.Unmarshal(metadata, &md); err != nil {
				outputError(err)
				return
			}

			fmt.Printf("[+] Extracting %s (%s)\n", md.Path, units.BytesSize(float64(md.Size)))

			// Create output directory if needed
			if dir := filepath.Dir(md.Path); dir != "" {
				if err := os.MkdirAll(dir, os.ModeDir|os.ModePerm); err != nil {
					outputError(err)
					return
				}
			}
			// Create output file.
			file, err := os.Create(md.Path)
			if err != nil {
				outputError(err)
				return
			}

			// Handle contents
			for j := uint64(0); ; j += 2 {
				chunk, err := Get(id.Derive(idMemory, i, j))
				if err != nil {
					// eof
					break
				}
				n, err := Decrypt(chunk, key.Bytes(), buffer[:])
				if err != nil {
					outputError(err)
					return
				}
				if n != 4096 {
					outputError(errors.New("error invalid plaintext size"))
					return
				}
				for k := len(buffer) - 1; ; k-- {
					if buffer[k] == 0 {
						continue
					} else if buffer[k] == 1 {
						if _, err := file.Write(buffer[:k]); err != nil {
							outputError(err)
							return
						}
						break
					} else {
						outputError(errors.New("error invalid padding"))
						return
					}
				}
			}
			file.Close()
		}
		return
	} else if args[1] == "wipe" {
		if len(args) != 2 {
			goto help
		}
		// todo
	}

help:
	help()
}

func help() {
	fmt.Printf(`Usage: %s {command}

Commands:
	help			print help information
	seal {path}		encrypt and store data at given path
	open {path}		decrypt and extract data and write to given path
	wipe			removes all data associated with an entry from the database
	`, args[0])
}

func outputError(err error) {
	fmt.Fprintln(os.Stderr, "[!]", err)
}
