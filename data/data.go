package data

import (
	"fmt"
	"io"
	"math"
	"os"

	"github.com/cheggaaa/pb"
	"github.com/libeclipse/dissident/coffer"
	"github.com/libeclipse/dissident/crypto"
	"github.com/libeclipse/dissident/memory"
)

// ImportData reads a file from the disk and imports it.
func ImportData(path string, fileSize int64, rootIdentifier []byte, masterKey *[32]byte) {
	// Open the file.
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

	// Start the progress bar.
	bar := pb.New64(fileSize).Prefix("+ Importing ")
	bar.ShowSpeed = true
	bar.SetUnits(pb.U_BYTES)
	bar.Start()

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
}

// ExportData exports data from coffer to the disk.
func ExportData(path string, rootIdentifier []byte, masterKey *[32]byte) {
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
	lenData := MetaGetLength("length", rootIdentifier, masterKey)

	// Start the progress bar object.
	bar := pb.New64(lenData).Prefix("+ Exporting ")
	bar.ShowSpeed = true
	bar.SetUnits(pb.U_BYTES)
	bar.Start()

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

// ViewData grabs the data from coffer and writes it to stdout.
func ViewData(rootIdentifier []byte, masterKey *[32]byte) {
	// Get the metadata first.
	lenData := MetaGetLength("length", rootIdentifier, masterKey)

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

// RemoveData removes data from coffer.
func RemoveData(rootIdentifier []byte, masterKey *[32]byte) {
	// Get the metadata first.
	lenData := MetaGetLength("length", rootIdentifier, masterKey)

	// Start the progress bar.
	bar := pb.New64(int64(math.Ceil(float64(lenData) / 4096))).Prefix("+ Removing ")
	bar.ShowCounters = false
	bar.SetUnits(pb.U_NO)
	bar.Start()

	// Remove all metadata.
	MetaRemoveData(rootIdentifier)

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
