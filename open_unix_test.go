//go:build !windows
// +build !windows

package saferio

import (
	"os"
	"path/filepath"
	"testing"
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

func assertError(t *testing.T, err error, message string) {
	t.Helper()
	if err == nil {
		t.Fatal(message)
	}
}

func TestOpenRD(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("open existing file for reading", func(t *testing.T) {
		path := filepath.Join(tmpDir, "readable.txt")
		createTestFile(t, path, []byte("hello"))

		f, err := OpenRD(path)
		if err != nil {
			t.Fatalf("OpenRD failed: %v", err)
		}
		defer f.Close()

		buf := make([]byte, 5)
		n, err := f.Read(buf)
		if err != nil {
			t.Fatalf("Read failed: %v", err)
		}
		if string(buf[:n]) != "hello" {
			t.Errorf("unexpected content: got %q, want %q", string(buf[:n]), "hello")
		}
	})

	t.Run("open non-existing file", func(t *testing.T) {
		path := filepath.Join(tmpDir, "not_exist.txt")
		_, err := OpenRD(path)
		assertError(t, err, "expected an error for non-existing file, got nil")
	})

	t.Run("refuse to open symlink", func(t *testing.T) {
		target := filepath.Join(tmpDir, "target.txt")
		link := filepath.Join(tmpDir, "link.txt")
		createSymlink(t, target, link)

		_, err := OpenRD(link)
		assertError(t, err, "expected an error opening symlink with O_NOFOLLOW, got nil")
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
		if err != nil {
			t.Fatalf("OpenAD failed: %v", err)
		}
		writeAndClose(t, f, "first")

		f, err = OpenAD(path)
		if err != nil {
			t.Fatalf("OpenAD failed on second open: %v", err)
		}
		writeAndClose(t, f, "second")

		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("ReadFile failed: %v", err)
		}
		if string(content) != "firstsecond" {
			t.Errorf("unexpected content: got %q, want %q", string(content), "firstsecond")
		}
	})

	t.Run("refuse to open symlink", func(t *testing.T) {
		target := filepath.Join(tmpDir, "target.txt")
		link := filepath.Join(tmpDir, "link.txt")
		createSymlink(t, target, link)

		_, err := OpenAD(link)
		assertError(t, err, "expected an error opening symlink with O_NOFOLLOW, got nil")
	})
}
