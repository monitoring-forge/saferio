//go:build !windows
// +build !windows

package saferio

import (
	"os"
	"syscall"
)

// OpenRD opens a file for reading with O_NOFOLLOW flag to prevent symlink attacks.
// This function is specific to Unix-like systems where syscall.O_NOFOLLOW is available.
func OpenRD(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_RDONLY|syscall.O_NOFOLLOW, 0)
}

// osOpenAD opens a file for appending with O_NOFOLLOW flag to prevent symlink attacks.
// This function is specific to Unix-like systems where syscall.O_NOFOLLOW is available.
func osOpenAD(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND|syscall.O_NOFOLLOW, 0644)
}
