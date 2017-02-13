package coffer

import (
	"errors"
	"os"
	"os/user"

	"github.com/boltdb/bolt"
)

var (
	// Coffer is a pointer to the database object.
	Coffer *bolt.DB
)

// Setup sets up the environment.
func Setup() error {
	// Ascertain the path to the secret store.
	user, err := user.Current()
	if err != nil {
		return err
	}

	// Check if we've done this before.
	if _, err = os.Stat(user.HomeDir + "/.pocket"); err != nil {
		// Apparently we haven't.

		// Create a directory to store our stuff in.
		err = os.Mkdir(user.HomeDir+"/.pocket", 0700)
		if err != nil {
			return err
		}
	}

	// Open the database file.
	db, err := bolt.Open(user.HomeDir+"/.pocket/coffer.bolt", 0700, nil)
	if err != nil {
		return err
	}
	Coffer = db

	// Create the bucket to guarantee it exists.
	Coffer.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("coffer"))
		return nil
	})

	return nil
}

// Save saves a secret to the database.
func Save(identifier, ciphertext []byte) error {
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

// Retrieve retrieves a secret from the database.
func Retrieve(identifier []byte) ([]byte, error) {
	// Allocate space to hold the ciphertext.
	var ciphertext []byte

	// Attempt to retrieve the ciphertext from the database.
	if err := Coffer.View(func(tx *bolt.Tx) error {
		// Grab the bucket that we'll be using.
		bucket := tx.Bucket([]byte("coffer"))

		id := bucket.Get(identifier)
		if id == nil {
			// We didn't find that key; return an error.
			return errors.New("[!] Nothing to see here")
		}

		ciphertext = append(ciphertext, id...)

		return nil
	}); err != nil {
		return nil, err
	}

	return ciphertext, nil
}

// Delete deletes an entry from the database.
func Delete(identifier []byte) {
	Coffer.Update(func(tx *bolt.Tx) error {
		// Grab the bucket that we'll be using.
		bucket := tx.Bucket([]byte("coffer"))

		// Delete the entry.
		bucket.Delete(identifier)

		return nil
	})
}

// Close closes the database object.
func Close() {
	Coffer.Close()
}
