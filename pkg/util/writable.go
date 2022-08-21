package util

import (
	"io/ioutil"
	"os"
)

// IsWritableDir returns true if the given directory path is writable
func IsWritableDir(path string) (bool, error) {
	// Create a temprary directory
	dir, err := ioutil.TempDir(path, "")
	defer os.RemoveAll(dir)
	if err != nil {
		return false, nil
	}
	return true, nil
}
