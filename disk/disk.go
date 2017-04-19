package disk

import (
	"errors"
	"fmt"
	"os"
)

// GetFileInfo gets the info for a file and returns it.
func GetFileInfo(path string) (os.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("! File does not exist")
		}
		return nil, err
	}

	return info, nil
}

// OpenFileRead opens a file path for reading and returns the file object.
func OpenFileRead(path string) (*os.File, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsPermission(err) {
			return nil, fmt.Errorf("! Insufficient permissions to open %s", path)
		}
		return nil, err
	}

	return f, nil
}

// OpenFileAppend opens a file path for appending and returns the file object. It
// creates the file if it does not exist. It returns an error if the file already exists.
func OpenFileAppend(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if os.IsExist(err) {
			return nil, fmt.Errorf("! %s already exists; cannot overwrite", path)
		} else if os.IsPermission(err) {
			return nil, fmt.Errorf("! Insufficient permissions to open %s", path)
		}
		return nil, err
	}

	return f, nil
}
