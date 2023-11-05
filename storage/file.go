package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileBackend struct {
	storagePath string
}

func NewFileBackend(path string) (*FileBackend, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(absPath); err != nil {
		return nil, fmt.Errorf("path %s already exists", absPath)
	}

	if err := os.MkdirAll(absPath, 0700); err != nil {
		return nil, err
	}

	return &FileBackend{storagePath: absPath}, nil
}

func (f *FileBackend) Put(key string, value []byte) error {
	return nil
}

func (f *FileBackend) Get(key string) ([]byte, error) {
	return nil, nil
}

func (f *FileBackend) Delete(key string) error {
	return nil
}
