//go:build windows
// +build windows

package saferio

import (
	"fmt"
	"os"
)

// replaceFile replaces the old file with the new file on Windows.
// It first attempts to rename the old file to the new file. If that fails (e.g., if the new file already exists),
// it removes the new file and tries to rename again.
func replaceFile(oldPath, newPath string) error {
	if err := os.Rename(oldPath, newPath); err != nil {
		if removeErr := os.Remove(newPath); removeErr != nil {
		    return fmt.Errorf("rename %q to %q failed: %w; removing destination failed: %v", oldPath, newPath, err, removeErr)
		}
		return os.Rename(oldPath, newPath)
	}
	return nil
}
