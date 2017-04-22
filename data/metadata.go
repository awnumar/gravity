package data

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

// NewMetadataObj creates a blank json object to hold the metadata
// and sets the globally accessible variable to it.
func NewMetadataObj() {
	metaObj = gabs.New()
}

// ResetMetaDataObj is an alias for NewMetadataObj(). It is used to reset
// the JSON object, removing all stored data.
func ResetMetaDataObj() {
	NewMetadataObj()
}

// SetMetaValue adds a value to metadata at the supplied path.
func SetMetaValue(value interface{}, path string) {
	metaObj.SetP(value, path)
}

// GetMetaValue gets a value from the metaObj at a path and returns it.
func GetMetaValue(path string) interface{} {
	value := metaObj.Path(path).Data()
	return value
}

// GetMetaLength retrieves the length of this data and returns it.
func GetMetaLength(path string) int64 {
	value := metaObj.Path(path).Data()
	if value == nil {
		fmt.Println("! No length field found; was importing interrupted?")
		memory.SafeExit(1)
	}

	return int64(value.(float64))
}

// ExportMetaBytes returns the JSON object in bytes.
func ExportMetaBytes() []byte {
	return []byte(metaObj.String())
}

// SaveMetadata saves the metadata to the database.
func SaveMetadata(rootIdentifier []byte, masterKey *[32]byte) {
	// Grab the metadata as bytes.
	data := ExportMetaBytes()

	var chunk []byte
	for i := 0; i < len(data); i += 4095 {
		if i+4095 > len(data) {
			// Remaining data <= 4095.
			chunk = data[len(data)-(len(data)%4095):]
		} else {
			// Split into chunks of 4095 bytes and pad.
			chunk = data[i : i+4095]
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

// RetrieveMetadata gets the metadata from the database and returns
func RetrieveMetadata(rootIdentifier []byte, masterKey *[32]byte) {
	// Declare variable to hold all of this metadata.
	var data []byte

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
		data = append(data, unpadded...)
	}

	if len(data) == 0 {
		// No data.
		return
	}

	// Set the global metadata JSON object to this data.
	metadataObj, err := gabs.ParseJSON(data)
	if err != nil {
		fmt.Println(err)
		memory.SafeExit(1)
	}

	// That went well. Set the global var to that object.
	metaObj = metadataObj
}

// RemoveMetadata deletes all the metadata related to an entry.
func RemoveMetadata(rootIdentifier []byte) {
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
