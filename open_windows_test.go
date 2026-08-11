//go:build windows
// +build windows

package saferio

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenRD(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("open existing file for reading", func(t *testing.T) {
		path := filepath.Join(tmpDir, "readable.txt")
		if err := os.WriteFile(path, []byte("hello"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

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
		if err == nil {
			t.Fatal("expected an error for non-existing file, got nil")
		}
	})
}

func TestOpenAD(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("create and append to file", func(t *testing.T) {
		path := filepath.Join(tmpDir, "appendable.txt")

		f, err := OpenAD(path)
		if err != nil {
			t.Fatalf("OpenAD failed: %v", err)
		}
		if _, err := f.WriteString("first"); err != nil {
			t.Fatalf("WriteString failed: %v", err)
		}
		if err := f.Close(); err != nil {
			t.Fatalf("Close failed: %v", err)
		}

		f, err = OpenAD(path)
		if err != nil {
			t.Fatalf("OpenAD failed on second open: %v", err)
		}
		if _, err := f.WriteString("second"); err != nil {
			t.Fatalf("WriteString failed: %v", err)
		}
		if err := f.Close(); err != nil {
			t.Fatalf("Close failed: %v", err)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("ReadFile failed: %v", err)
		}
		if string(content) != "firstsecond" {
			t.Errorf("unexpected content: got %q, want %q", string(content), "firstsecond")
		}
	})
}
