package saferio

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func assertFilePathError(t *testing.T, filename string) {
	t.Helper()
	_, err := filePath("/tmp/data", filename)
	if err == nil {
		t.Errorf("filePath(%q) did not return an error", filename)
	}
}

func TestFilePath(t *testing.T) {
	t.Run("joins dir and filename", func(t *testing.T) {
		got, err := filePath("/tmp/data", "state.json")
		if err != nil {
			t.Fatalf("filePath returned an error: %v", err)
		}
		want := filepath.Join("/tmp/data", "state.json")
		if got != want {
			t.Errorf("filePath(%q, %q) = %q, want %q", "/tmp/data", "state.json", got, want)
		}
	})

	t.Run("rejects empty filename", func(t *testing.T) {
		assertFilePathError(t, "")
	})

	t.Run("rejects path traversal", func(t *testing.T) {
		cases := []string{"../etc/passwd", "subdir/state.json", "a/../b.json", "./file.json", "/", "./", "../", ".", ".."}
		for _, filename := range cases {
			assertFilePathError(t, filename)
		}
	})

	t.Run("rejects absolute paths", func(t *testing.T) {
		assertFilePathError(t, "/etc/passwd")
	})
}

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("existing file", func(t *testing.T) {
		path := filepath.Join(tmpDir, "exists.txt")
		if err := os.WriteFile(path, []byte("data"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		if !FileExists(tmpDir, "exists.txt") {
			t.Errorf("FileExists returned false for existing file")
		}
	})

	t.Run("non-existing file", func(t *testing.T) {
		if FileExists(tmpDir, "missing.txt") {
			t.Errorf("FileExists returned true for non-existing file")
		}
	})

	t.Run("ignores path traversal", func(t *testing.T) {
		if FileExists(tmpDir, "../missing.txt") {
			t.Errorf("FileExists should not escape dir")
		}
	})
}

func assertContentEquals(t *testing.T, got, want map[string]any) {
	t.Helper()
	if got["key"] != want["key"] || got["number"] != want["number"] {
		t.Errorf("unexpected JSON content: %v", got)
	}
}

func writeJSONToFile(t *testing.T, dir, filename string, data any) string {
	t.Helper()
	if err := WriteJSON(dir, filename, data); err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}
	return filepath.Join(dir, filename)
}

func readAndUnmarshal(t *testing.T, path string) map[string]any {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(content, &got); err != nil {
		t.Fatalf("written content is not valid JSON: %v", err)
	}
	return got
}

func TestWriteJSON(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("write valid JSON", func(t *testing.T) {
		data := map[string]any{"key": "value", "number": float64(42)}
		path := writeJSONToFile(t, tmpDir, "test.json", data)
		got := readAndUnmarshal(t, path)
		assertContentEquals(t, got, data)
	})

	t.Run("overwrite existing file", func(t *testing.T) {
		filename := "overwrite.json"
		path := filepath.Join(tmpDir, filename)
		if err := os.WriteFile(path, []byte("old"), 0644); err != nil {
			t.Fatalf("failed to create existing file: %v", err)
		}

		writeJSONToFile(t, tmpDir, filename, map[string]any{"new": true})
		content := readAndUnmarshal(t, path)
		if content["new"] != true {
			t.Errorf("file was not overwritten: %v", content)
		}
	})

	t.Run("invalid directory", func(t *testing.T) {
		err := WriteJSON("/nonexistent/directory", "test.json", map[string]any{})
		if err == nil {
			t.Fatal("expected an error for invalid directory, got nil")
		}
	})

	t.Run("invalid value", func(t *testing.T) {
		err := WriteJSON(tmpDir, "invalid.json", make(chan int))
		if err == nil {
			t.Fatal("expected an error for invalid JSON value, got nil")
		}
	})
}

func writeFile(t *testing.T, path string, content []byte) {
	t.Helper()
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
}

func assertJSONError(t *testing.T, err error, wantType string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected an error for %s, got nil", wantType)
	}
}

func TestReadJSON(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("read valid JSON", func(t *testing.T) {
		filename := "read.json"
		path := filepath.Join(tmpDir, filename)
		writeFile(t, path, []byte(`{"name":"test","count":3}`+"\n"))

		var got struct {
			Name  string `json:"name"`
			Count int    `json:"count"`
		}
		if err := ReadJSON(tmpDir, filename, &got); err != nil {
			t.Fatalf("ReadJSON failed: %v", err)
		}
		if got.Name != "test" || got.Count != 3 {
			t.Errorf("unexpected result: %+v", got)
		}
	})

	t.Run("file does not exist", func(t *testing.T) {
		var v map[string]any
		err := ReadJSON(tmpDir, "missing.json", &v)
		assertJSONError(t, err, "missing file")
	})

	t.Run("empty file", func(t *testing.T) {
		filename := "empty.json"
		path := filepath.Join(tmpDir, filename)
		writeFile(t, path, []byte{})

		var v map[string]any
		err := ReadJSON(tmpDir, filename, &v)
		assertJSONError(t, err, "empty file")
		if !errors.Is(err, io.EOF) {
			t.Errorf("expected io.EOF for empty file, got %T: %v", err, err)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		filename := "invalid.json"
		path := filepath.Join(tmpDir, filename)
		writeFile(t, path, []byte("not json"))

		var v map[string]any
		err := ReadJSON(tmpDir, filename, &v)
		assertJSONError(t, err, "invalid JSON")
		var syntaxErr *json.SyntaxError
		if !errors.As(err, &syntaxErr) {
			t.Errorf("expected a JSON syntax error, got %T: %v", err, err)
		}
	})
}

func TestWriteReadJSON(t *testing.T) {
	tmpDir := t.TempDir()
	filename := "roundtrip.json"

	want := map[string]any{"message": "hello", "items": []any{"a", "b"}}
	if err := WriteJSON(tmpDir, filename, want); err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	var got map[string]any
	if err := ReadJSON(tmpDir, filename, &got); err != nil {
		t.Fatalf("ReadJSON failed: %v", err)
	}

	if got["message"] != "hello" {
		t.Errorf("unexpected message: %v", got["message"])
	}
}
