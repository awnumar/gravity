package coffer

import (
	"os"
	"os/user"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	// Coffer is a pointer to the database object.
	Coffer *leveldb.DB
)

// Setup sets up the environment.
func Setup() error {
	// Ascertain the path to the secret store.
	user, err := user.Current()
	if err != nil {
		return err
	}

	// Check if we've done this before.
	if _, err = os.Stat(user.HomeDir + "/.dissident"); err != nil {
		// Apparently we haven't.

		// Create a directory to store our stuff in.
		err = os.Mkdir(user.HomeDir+"/.dissident", 0700)
		if err != nil {
			return err
		}
	}

	// Open the database file.
	Coffer, err = leveldb.OpenFile(user.HomeDir+"/.dissident/coffer", nil)
	if err != nil {
		return err
	}

	return nil
}

// Exists checks if an entry exists and returns true or false.
func Exists(identifier []byte) bool {
	_, err := Coffer.Get(identifier, nil)
	if err != nil {
		return false
	}

	return true
}

// Save saves a secret to the database.
func Save(identifier, ciphertext []byte) {
	Coffer.Put(identifier, ciphertext, nil)
}

// Retrieve retrieves a secret from the database.
func Retrieve(identifier []byte) []byte {
	data, _ := Coffer.Get(identifier, nil)

	return data
}

// Delete deletes an entry from the database.
func Delete(identifier []byte) {
	Coffer.Delete(identifier, nil)
}

// Close closes the database object.
func Close() {
	Coffer.Close()
}
