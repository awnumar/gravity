package output

import (
	"fmt"

	"github.com/libeclipse/tranquil/coffer"
	"github.com/libeclipse/tranquil/crypto"
	"github.com/libeclipse/tranquil/memory"
)

// PipeFullEntry returns a channel and sends through all the processed plaintext chunks.
func PipeFullEntry(rootIdentifier []byte, masterKey *[32]byte) chan []byte {
	chunksChan := make(chan []byte)

	go func(chan []byte) {
		for n := 0; true; n++ {
			// Derive derived_identifier[n]
			ct := coffer.Retrieve(crypto.DeriveIdentifierN(rootIdentifier, n))
			if ct == nil {
				// This one doesn't exist. //EOF
				chunksChan <- nil
				break
			}

			// Decrypt this slice.
			pt := crypto.Decrypt(ct, masterKey)

			// Unpad this slice and wipe old one.
			unpadded, e := crypto.Unpad(pt)
			if e != nil {
				fmt.Println(e)
				return
			}
			memory.Wipe(pt)

			// Send the processed plaintext chunk to the caller.
			chunksChan <- unpadded
		}
	}(chunksChan)

	return chunksChan
}
