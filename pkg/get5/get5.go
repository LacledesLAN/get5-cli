package get5

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// OpenSchemaFile loads a get 5 configuration json from a filesystem path.
func OpenSchemaFile(path string, cfg *Match) error {
	path = strings.TrimSpace(path)
	if len(path) == 0 {
		return errors.New("cannot load a file using an empty or whitespace-only path")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("couldn't determine absolute path of %q: %w", path, err)
	}

	if stats, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("path %q does not exist: %w", absPath, err)
	} else if stats.IsDir() {
		return fmt.Errorf("path %q is a directory and not a file", absPath)
	}

	jsonBytes, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("couldn't read get5 config file %q: %w", absPath, err)
	}

	if !json.Valid(jsonBytes) {
		return fmt.Errorf("input json file format was invalid.")
	}

	if err = json.Unmarshal(jsonBytes, &cfg); err != nil {
		return fmt.Errorf("couldn't unmarshal get5 configuration file %q: %w", absPath, err)
	}

	return nil
}

// SaveSchemaFile saves a get 5 configuration to a json file
func SaveSchemaFile(c Match, path string) error {
	if err := sanitizeMatch(&c); err != nil {
		return fmt.Errorf("generated schema file was invalid: %w", err)
	}

	path = strings.TrimSpace(path)
	if len(path) == 0 {
		return errors.New("cannot save a file using an empty or whitespace-only path")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("couldn't determine absolute path of %q: %w", path, err)
	}

	fileBytes, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return fmt.Errorf("unable to encode get5 configuration as JSON: %w", err)
	}

	fileBytes = append(fileBytes, []byte("\n")...)

	if !json.Valid(fileBytes) {
		return fmt.Errorf("generated json file format was invalid.")
	}

	fh, err := os.OpenFile(absPath, os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return fmt.Errorf("couldn't open file %q for writing: %w", absPath, err)
	}
	defer fh.Close()

	if err := fh.Truncate(0); err != nil {
		return fmt.Errorf("couldn't clear contents of destination file, before writing")
	}

	if _, err := fh.Write(fileBytes); err != nil {
		return fmt.Errorf("couldn't write get5 configuration to file %q : %w", absPath, err)
	}

	return nil
}
