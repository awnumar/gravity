package coffer

import (
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

// Exists checks if an entry exists and returns true or false.
func Exists(identifier []byte) bool {
	exists := false

	Coffer.View(func(tx *bolt.Tx) error {
		// Grab the bucket that we'll be using.
		bucket := tx.Bucket([]byte("coffer"))

		// Attempt to locate the entry.
		ct := bucket.Get(identifier)
		if ct != nil {
			exists = true
		}

		return nil
	})

	return exists
}

// Save saves a secret to the database.
func Save(identifier, ciphertext []byte) {
	Coffer.Update(func(tx *bolt.Tx) error {
		// Grab the bucket that we'll be using.
		bucket := tx.Bucket([]byte("coffer"))

		// Save the identifier:ciphertext pair to the coffer.
		bucket.Put(identifier, ciphertext)

		return nil
	})
}

// Retrieve retrieves a secret from the database.
func Retrieve(identifier []byte) []byte {
	var ciphertext []byte

	// Attempt to retrieve the ciphertext from the database.
	Coffer.View(func(tx *bolt.Tx) error {
		// Grab the bucket that we'll be using.
		bucket := tx.Bucket([]byte("coffer"))

		// Attempt to locate and grab the data from the coffer.
		ciphertext = append(ciphertext, bucket.Get(identifier)...)

		return nil
	})

	return ciphertext
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
