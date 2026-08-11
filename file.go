package saferio

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func filePath(dir, filename string) (string, error) {
	if filename == "" {
		return "", fmt.Errorf("filename cannot be empty")
	}
	basename := filepath.Base(filename)
	if filename != basename || basename == "." || basename == ".." || basename == "/" {
		return "", fmt.Errorf("invalid filename: %s", filename)
	}
	return filepath.Join(dir, basename), nil
}

func FileExists(dir, filename string) bool {
	path, err := filePath(dir, filename)
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}

func WriteJSON(dir, filename string, v any) error {
	destPath, err := filePath(dir, filename)
	if err != nil {
		return err
	}

	newFile, err := os.CreateTemp(dir, "saferio_temp_")
	if err != nil {
		return err
	}
	defer func() {
		errRemove := os.Remove(newFile.Name())
		if errRemove != nil && !os.IsNotExist(errRemove) {
			fmt.Fprintf(os.Stderr, "Failed to remove temporary file: %s, error: %v", newFile.Name(), errRemove)
		}
	}()

	je := json.NewEncoder(newFile)
	err = je.Encode(v)
	if err != nil {
		_ = newFile.Close()
		return err
	}

	err = newFile.Close()
	if err != nil {
		return err
	}

	return replaceFile(newFile.Name(), destPath)
}

func ReadJSON(dir, filename string, v any) error {
	filePath, err := filePath(dir, filename)
	if err != nil {
		return err
	}

	file, err := OpenRD(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(v)
	if err != nil {
		return err
	}
	return nil
}
