package coffer

import (
	"errors"
	"log"
	"os"
	"os/user"
	"reflect"

	"github.com/boltdb/bolt"
)

// Setup sets up the environment.
func Setup() error {
	// Get the current user.
	user, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}

	// Change the working directory to the user's home.
	err = os.Chdir(user.HomeDir)
	if err != nil {
		log.Fatalln(err)
	}

	// Check if we've done this before.
	if _, err = os.Stat("./.pocket"); err != nil {
		// Apparently we have.

		// Create a directory to store our stuff in.
		err = os.Mkdir("./.pocket", 0700)
		if err != nil && !os.IsExist(err) {
			log.Fatalln(err)
		}
	}

	return nil
}

// SaveSecret saves the secrets to the disk.
func SaveSecret(identifier, ciphertext []byte) error {
	// Open the database.
	db, err := bolt.Open("./.pocket/secrets", 0700, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Update(func(tx *bolt.Tx) error {
		// Create bucket if it doesn't exist.
		bucket, _ := tx.CreateBucketIfNotExists([]byte("secrets"))

		// Check if this identifier already exists.
		key := bucket.Get(identifier)
		if key != nil {
			return errors.New("[!] Cannot overwrite existing entry")
		}

		// Save the identifier/ciphertext pair to the bucket.
		bucket.Put(identifier, ciphertext)

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// RetrieveSecret retrieves the secrets from the disk.
func RetrieveSecret(identifier []byte) ([]byte, error) {
	// Open the database.
	db, err := bolt.Open("./.pocket/secrets", 0700, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Allocate space to hold the ciphertext.
	var ciphertext []byte

	// Attempt to retrieve the ciphertext from the database.
	if err := db.View(func(tx *bolt.Tx) error {
		// Grab the bucket.
		bucket := tx.Bucket([]byte("secrets"))
		if bucket == nil {
			// It doesn't exist.
			return errors.New("[!] Nothing to see here")
		}

		// Iterate over all the keys.
		c := bucket.Cursor()
		for id, ct := c.First(); id != nil; id, ct = c.Next() {
			if reflect.DeepEqual(id, identifier) {
				ciphertext = append(ciphertext, ct...)
				return nil
			}
		}

		return errors.New("[!] Nothing to see here")
	}); err != nil {
		return nil, err
	}

	return ciphertext, nil
}
