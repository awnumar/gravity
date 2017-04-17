package meta

import (
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/libeclipse/dissident/coffer"
	"github.com/libeclipse/dissident/crypto"
	"github.com/libeclipse/dissident/memory"
)

var (
	metaObj *gabs.Container
)

// New creates a blank json object to hold the metadata
// and sets the globally accessable variable to it.
func New() {
	metaObj = gabs.New()
}

// Reset is an alias for New(). It is used to reset
// the JSON object, removing all stored data.
func Reset() {
	New()
}

// Set adds a value to metadata at the supplied path.
func Set(value interface{}, path string) {
	metaObj.SetP(value, path)
}

// Get gets a value from the metaObj at a path and returns it.
func Get(path string) interface{} {
	value := metaObj.Path(path).Data()
	return value
}

// GetLength retrieves the length of this data and returns it.
func GetLength(path string) int64 {
	value := metaObj.Path(path).Data()
	if value == nil {
		fmt.Println("! No length field found")
		memory.SafeExit(1)
	}

	return int64(value.(float64))
}

// ExportBytes returns the JSON object in bytes.
func ExportBytes() []byte {
	return []byte(metaObj.String())
}

// Save saves the metadata to the database.
func Save(rootIdentifier []byte, masterKey *[32]byte) {
	// Grab the metadata as bytes.
	metadata := ExportBytes()

	var chunk []byte
	for i := 0; i < len(metadata); i += 4095 {
		if i+4095 > len(metadata) {
			// Remaining data <= 4095.
			chunk = metadata[len(metadata)-(len(metadata)%4095):]
		} else {
			// Split into chunks of 4095 bytes and pad.
			chunk = metadata[i : i+4095]
		}

		// Pad the chunk to standard size.
		padded, err := crypto.Pad(chunk, 4096)
		if err != nil {
			fmt.Println(err)
			memory.SafeExit(1)
		}

		// Save it to the database.
		coffer.Save(crypto.DeriveMetaIdentifierN(rootIdentifier, -i-1), crypto.Encrypt(padded, masterKey))
	}
}

// Retrieve gets the metadata from the database and returns
func Retrieve(rootIdentifier []byte, masterKey *[32]byte) {
	// Declare variable to hold all of this metadata.
	var metadata []byte

	for n := -1; true; n-- {
		ct := coffer.Retrieve(crypto.DeriveMetaIdentifierN(rootIdentifier, n))
		if ct == nil {
			// This one doesn't exist. //EOF
			break
		}

		// Decrypt this slice.
		pt, err := crypto.Decrypt(ct, masterKey)
		if err != nil {
			fmt.Println(err)
			memory.SafeExit(1)
		}

		// Unpad this slice.
		unpadded, e := crypto.Unpad(pt)
		if e != nil {
			fmt.Println(e)
			memory.SafeExit(1)
		}

		// Append this chunk to the metadata.
		metadata = append(metadata, unpadded...)
	}

	if len(metadata) == 0 {
		// No data.
		return
	}

	// Set the global metadata JSON object to this data.
	metadataObj, err := gabs.ParseJSON(metadata)
	if err != nil {
		fmt.Println(err)
		memory.SafeExit(1)
	}

	// That went well. Set the global var to that object.
	metaObj = metadataObj
}

// Remove deletes all the metadata related to an entry.
func Remove(rootIdentifier []byte) {
	for n := -1; true; n-- {
		// Get the DeriveIdentifierN for this n.
		derivedMetaIdentifierN := crypto.DeriveMetaIdentifierN(rootIdentifier, n)

		// Check if it exists.
		if coffer.Exists(derivedMetaIdentifierN) {
			coffer.Delete(derivedMetaIdentifierN)
		} else {
			break
		}
	}
}
