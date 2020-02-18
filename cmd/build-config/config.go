package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Paths used by get5-cli
type Paths struct {
	SRCDS  string `json:"srcds"`
	Input  string `json:"input"`
	Output string `json:"output"`
}

// Config for get5-cli
type Config struct {
	Paths Paths `json:"paths"`
}

// LoadConfig loads a get5-cli configuration from filesystem path
func LoadConfig(p string, c *Config) error {
	p = strings.TrimSpace(p)
	if len(p) == 0 {
		return errors.New("cannot load a file using an empty or whitespace-only path")
	}

	absPath, err := filepath.Abs(p)
	if err != nil {
		return fmt.Errorf("couldn't determine absolute path of %q: %w", p, err)
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

	if err = json.Unmarshal(jsonBytes, &c); err != nil {
		return fmt.Errorf("couldn't unmarshal get5 configuration file %q: %w", absPath, err)
	}

	return nil
}
