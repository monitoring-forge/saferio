package saferio

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func FileExists(dir, filename string) bool {
	_, err := os.Stat(filepath.Join(dir, filename))
	return err == nil
}

func WriteJSON(dir, filename string, v any) error {
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

	destination := filepath.Join(dir, filename)
	if err := os.Rename(newFile.Name(), destination); err != nil {
		if removeErr := os.Remove(destination); removeErr != nil {
			return err
		}
		return os.Rename(newFile.Name(), destination)
	}
	return nil
}

func ReadJSON(dir, filename string, v any) error {
	filename = filepath.Join(dir, filename)

	file, err := OpenRD(filename)
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
