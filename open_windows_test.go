//go:build windows
// +build windows

package saferio

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenRD(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("open existing file for reading", func(t *testing.T) {
		path := filepath.Join(tmpDir, "readable.txt")
		if err := os.WriteFile(path, []byte("hello"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		f, err := OpenRD(path)
		require.NoError(t, err, "OpenRD failed")
		defer f.Close()

		buf := make([]byte, 5)
		n, err := f.Read(buf)
		require.NoError(t, err, "Read failed")
		require.Equal(t, "hello", string(buf[:n]), "unexpected content")
	})

	t.Run("open non-existing file", func(t *testing.T) {
		path := filepath.Join(tmpDir, "not_exist.txt")
		_, err := OpenRD(path)
		require.Error(t, err, "expected an error for non-existing file, got nil")
	})
}

func TestOpenAD(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("create and append to file", func(t *testing.T) {
		path := filepath.Join(tmpDir, "appendable.txt")

		f, err := OpenAD(path)
		require.NoError(t, err, "OpenAD failed on first open")
		_, err = f.WriteString("first")
		require.NoError(t, err, "WriteString failed")
		err = f.Close()
		require.NoError(t, err, "Close failed")

		f, err = OpenAD(path)
		require.NoError(t, err, "OpenAD failed on second open")
		_, err = f.WriteString("second")
		require.NoError(t, err, "WriteString failed")
		err = f.Close()
		require.NoError(t, err, "Close failed")

		content, err := os.ReadFile(path)
		require.NoError(t, err, "ReadFile failed")
		require.Equal(t, "firstsecond", string(content), "unexpected content")
	})
}
