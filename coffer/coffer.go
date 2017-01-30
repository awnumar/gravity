package coffer

import (
	"errors"
	"log"
	"os"
	"os/user"
	"reflect"

	"github.com/boltdb/bolt"
)

var (
	// Coffer is a pointer to the database object.
	Coffer *bolt.DB
)

// Setup sets up the environment.
func Setup() {
	// Ascertain the path to the secret store.
	user, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}

	// Check if we've done this before.
	if _, err = os.Stat(user.HomeDir + "/.pocket"); err != nil {
		// Apparently we haven't.

		// Create a directory to store our stuff in.
		err = os.Mkdir(user.HomeDir+"/.pocket", 0700)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// Open the database file.
	db, err := bolt.Open(user.HomeDir+"/.pocket/coffer.bolt", 0700, nil)
	if err != nil {
		log.Fatal(err)
	}
	Coffer = db

	// Create the bucket to guarantee it exists.
	Coffer.Update(func(tx *bolt.Tx) error {
		_, _ = tx.CreateBucketIfNotExists([]byte("coffer"))
		return nil
	})
}

// SaveSecret saves the secrets to the disk.
func SaveSecret(identifier, ciphertext []byte) error {
	return Coffer.Update(func(tx *bolt.Tx) error {
		// Grab the bucket that we'll be using.
		bucket := tx.Bucket([]byte("coffer"))

		// Check if this identifier already exists.
		value := bucket.Get(identifier)
		if value != nil {
			// It does; abort and return an error.
			return errors.New("[!] Cannot overwrite existing entry")
		}

		// Save the identifier:ciphertext pair to the coffer.
		bucket.Put(identifier, ciphertext)

		return nil
	})
}

// RetrieveSecret retrieves the secrets from the disk.
func RetrieveSecret(identifier []byte) ([]byte, error) {
	// Allocate space to hold the ciphertext.
	var ciphertext []byte

	// Attempt to retrieve the ciphertext from the database.
	if err := Coffer.View(func(tx *bolt.Tx) error {
		// Grab the bucket.
		bucket := tx.Bucket([]byte("secrets"))

		// Iterate over all the keys.
		c := bucket.Cursor()
		for id, ct := c.First(); id != nil; id, ct = c.Next() {
			if reflect.DeepEqual(id, identifier) {
				ciphertext = append(ciphertext, ct...)
				return nil
			}
		}

		// We didn't find that key; return an error.
		return errors.New("[!] Nothing to see here")
	}); err != nil {
		return nil, err
	}

	return ciphertext, nil
}
