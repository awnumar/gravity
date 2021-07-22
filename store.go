package main

import (
	"fmt"

	"github.com/awnumar/memguard"
	"git.mills.io/prologic/bitcask"
)

var database = func() *bitcask.Bitcask {
	db, err := bitcask.Open("store")
	if err != nil {
		memguard.SafePanic(err)
	}
	return db
}()

// Put puts a key value pair in the database
func Put(key, value []byte) error {
	return database.Put(key, value)
}

// Get gets a value for a key from the database
func Get(key []byte) ([]byte, error) {
	return database.Get(key)
}

func closeDB() {
	fmt.Println("[i] Compacting database...")
	database.Merge()
	fmt.Println("[i] Syncing data with disk...")
	database.Sync()
	fmt.Println("[i] Closing database...")
	database.Close()
}
