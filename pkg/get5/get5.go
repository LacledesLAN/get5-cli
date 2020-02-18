package get5

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// FromFile loads a get 5 configuration json from a filesystem path.
func FromFile(path string, cfg *Config) error {
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

	jsonBytes, err := ioutil.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("couldn't read get5 config file %q: %w", absPath, err)
	}

	if err = json.Unmarshal(jsonBytes, &cfg); err != nil {
		return fmt.Errorf("couldn't unmarshal get5 configuration file %q: %w", absPath, err)
	}

	return nil
}

// SaveFile saves a get 5 configuration to a json file
func (c Config) SaveFile(path string) error {
	path = strings.TrimSpace(path)
	if len(path) == 0 {
		return errors.New("cannot save a file using an empty or whitespace-only path")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("couldn't determine absolute path of %q: %w", path, err)
	}

	fileBytes, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("unable to encode get5 configuration as JSON: %w", err)
	}

	fh, err := os.OpenFile(absPath, os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return fmt.Errorf("couldn't open file %q for writing: %w", absPath, err)
	}
	defer fh.Close()

	if _, err := fh.Write(fileBytes); err != nil {
		return fmt.Errorf("couldn't write get5 configuration to file %q : %w", absPath, err)
	}

	return nil
}
