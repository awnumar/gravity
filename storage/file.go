package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"lukechampine.com/frand"
)

// TODO: add locking to synchronise concurrent access to the same file
// https://stackoverflow.com/a/64612611

const (
	defaultFileSizeBytes = 8 * MB
)

type FileBackend struct {
	storagePath   string
	fileSizeBytes int
}

func NewFileBackend(path string) (*FileBackend, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(absPath); err == nil {
		return nil, fmt.Errorf("path %s already exists", absPath)
	}

	if err := os.MkdirAll(absPath, 0700); err != nil {
		return nil, err
	}

	return &FileBackend{
		storagePath:   absPath,
		fileSizeBytes: defaultFileSizeBytes,
	}, nil
}

func (b *FileBackend) Keys() ([]string, error) {
	entries, err := os.ReadDir(b.storagePath)
	if err != nil {
		return nil, err
	}

	var keys []string
	for _, entry := range entries {
		keys = append(keys, entry.Name())
	}

	return keys, nil
}

func (b *FileBackend) Put(key string, value []byte) error {
	f, err := os.Create(filepath.Join(b.storagePath, key))
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(value); err != nil {
		return err
	}

	return nil
}

func (b *FileBackend) Get(key string) ([]byte, error) {
	f, err := os.Open(filepath.Join(b.storagePath, key))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	value, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (b *FileBackend) Delete(key string) error {
	if !OverwriteToDelete {
		return os.Remove(filepath.Join(b.storagePath, key))
	}

	f, err := os.OpenFile(filepath.Join(b.storagePath, key), os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Seek(0, 0); err != nil {
		return err
	}

	if _, err := f.Write(frand.Bytes(b.fileSizeBytes)); err != nil {
		return err
	}

	return nil
}
