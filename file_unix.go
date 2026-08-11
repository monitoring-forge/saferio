//go:build !windows
// +build !windows

package saferio

import "os"

func replaceFile(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}
