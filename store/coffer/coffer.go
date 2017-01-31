package coffer

import (
	"errors"
	"fmt"
	"os"
	"os/user"

	"github.com/boltdb/bolt"
)

// Coffer holds a reference to the BoltDB database, along with any possible
// metadata required.
type Coffer struct {
	bolt *bolt.DB
}

const (
	bucketName = "coffer"
)

var (
	// MasterBucket defines a bucket name to be used to store key/values in
	MasterBucket = []byte(bucketName)
)

// Setup sets up the environment.
func Setup() (*Coffer, error) {
	dir, err := poketDirectory()
	if err != nil {
		return nil, err
	}

	// Open the database file.
	db, err := bolt.Open(dir, 0700, nil)
	if err != nil {
		return nil, err
	}

	// Create our initial bucket
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(MasterBucket)
		return err
	})

	return &Coffer{
		bolt: db,
	}, nil
}

// Save saves a secret to the database.
func (coffer *Coffer) Save(identifier, ciphertext []byte) error {
	return coffer.bolt.Update(func(tx *bolt.Tx) error {
		// Grab the bucket that we'll be using.
		bucket := tx.Bucket(MasterBucket)

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
func (coffer *Coffer) Retrieve(identifier []byte) ([]byte, error) {
	var ciphertext []byte

	// Attempt to retrieve the ciphertext from the database.
	err := coffer.bolt.View(func(tx *bolt.Tx) error {
		// Grab the bucket that we'll be using.
		bucket := tx.Bucket(MasterBucket)

		// Find our key
		k, v := bucket.Cursor().Seek(identifier)
		if k == nil {
			return ErrNotfound
		}

		// Set the found value
		ciphertext = append(ciphertext, v...)
		return nil
	})

	return ciphertext, err
}

// Delete deletes an entry from the database.
func (coffer *Coffer) Delete(identifier []byte) {
	coffer.bolt.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(MasterBucket).Delete(identifier)
	})
}

// Close closes the database object.
func (coffer *Coffer) Close() error {
	return coffer.bolt.Close()
}

/**
 *	Helpers
 */

func poketDirectory() (string, error) {
	// Ascertain the path to the secret store.
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	// Check if we've done thiasdas before.
	if _, err := os.Stat(user.HomeDir + "/.pocket/"); err == nil {
		// Directory exists
		return fmt.Sprintf("%s/.pocket/", user.HomeDir), nil
	}

	// Create the directory
	if err := os.Mkdir(
		fmt.Sprintf("%s/.pocket", user.HomeDir),
		0700,
	); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/.pocket/", user.HomeDir), nil
}

// Errors
var (
	ErrNotfound = errors.New("key not found")
)
