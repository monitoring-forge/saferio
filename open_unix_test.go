//go:build !windows
// +build !windows

package saferio

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTestFile(t *testing.T, path string, content []byte) {
	t.Helper()
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
}

func createSymlink(t *testing.T, target, link string) {
	t.Helper()
	createTestFile(t, target, []byte("secret"))
	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}
}

func TestOpenRD(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("open existing file for reading", func(t *testing.T) {
		path := filepath.Join(tmpDir, "readable.txt")
		createTestFile(t, path, []byte("hello"))

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

	t.Run("refuse to open symlink", func(t *testing.T) {
		target := filepath.Join(tmpDir, "target.txt")
		link := filepath.Join(tmpDir, "link.txt")
		createSymlink(t, target, link)

		_, err := OpenRD(link)
		require.Error(t, err, "expected an error opening symlink with O_NOFOLLOW, got nil")
	})
}

func writeAndClose(t *testing.T, f *os.File, s string) {
	t.Helper()
	if _, err := f.WriteString(s); err != nil {
		t.Fatalf("WriteString failed: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}

func TestOpenAD(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("create and append to file", func(t *testing.T) {
		path := filepath.Join(tmpDir, "appendable.txt")

		f, err := OpenAD(path)
		require.NoError(t, err, "OpenAD failed")
		writeAndClose(t, f, "first")

		f, err = OpenAD(path)
		require.NoError(t, err, "OpenAD failed on second open")
		writeAndClose(t, f, "second")

		content, err := os.ReadFile(path)
		require.NoError(t, err, "ReadFile failed")
		require.Equal(t, "firstsecond", string(content), "unexpected content")
	})

	t.Run("refuse to open symlink", func(t *testing.T) {
		target := filepath.Join(tmpDir, "target.txt")
		link := filepath.Join(tmpDir, "link.txt")
		createSymlink(t, target, link)

		_, err := OpenAD(link)
		require.Error(t, err, "expected an error opening symlink with O_NOFOLLOW, got nil")
	})
}
