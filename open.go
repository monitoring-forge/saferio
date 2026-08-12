package saferio

import (
	"fmt"
	"os"
)

// OpenRD opens a file for reading with appropriate flags to prevent symlink attacks.
func OpenRD(filename string) (*os.File, error) {
	return osOpenRD(filename)
}

// OpenAD opens a file for appending with appropriate flags to prevent symlink attacks.
// It first checks if the file is a symlink using os.Lstat. If it is, an error is returned.
// If the file does not exist, it proceeds to open the file for appending.
func OpenAD(filename string) (*os.File, error) {
	if fi, err := os.Lstat(filename); err == nil {
		if fi.Mode()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("file is a symlink: %s", filename)
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	return osOpenAD(filename)
}
