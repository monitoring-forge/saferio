//go:build windows
// +build windows

package saferio

import "os"

// OpenRD opens a file for reading with appropriate flags for Windows.
// Since Windows doesn't support O_NOFOLLOW, we simply open the file with O_RDONLY.
func OpenRD(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_RDONLY, 0)
}

// osOpenAD opens a file for appending with appropriate flags for Windows.
// Since Windows doesn't support O_NOFOLLOW, we simply open the file with O_WRONLY|O_CREATE|O_APPEND.
func osOpenAD(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
}
