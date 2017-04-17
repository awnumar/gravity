package meta

import (
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/libeclipse/dissident/coffer"
	"github.com/libeclipse/dissident/crypto"
	"github.com/libeclipse/dissident/memory"
)

var (
	metadata *gabs.Container
)

// Create creates a blank json object to hold the metadata and returns a pointer to it.
func Create() *gabs.Container {
	metadata = gabs.New()
	return metadata
}

// Set adds a value to metadata at the supplied path.
func Set(value interface{}, path string) {
	metadata.SetP(value, path)
}

// ExportBytes returns the JSON object in bytes.
func ExportBytes() []byte {
	return []byte(metadata.String())
}

// Save saves the metadata to the database.
func Save(rootIdentifier []byte, masterKey *[32]byte) {
	// Grab the metadata as bytes.
	data := ExportBytes()

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
