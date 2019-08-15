package main

import (
	"os"
	"path/filepath"
)

// FileInfo represents a single file.
type FileInfo struct {
	Path string // Relative path from directory root.
	Size int64  // Size of the file in bytes.
}

// Files walks a given path and returns a slice of the files within it.
func Files(path string) (files []FileInfo, err error) {
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Discard directories.
		if info.IsDir() {
			return nil
		}

		// Construct FileInfo object and append to slice.
		files = append(files, FileInfo{
			Path: path,
			Size: info.Size(),
		})

		return nil
	})

	return
}
