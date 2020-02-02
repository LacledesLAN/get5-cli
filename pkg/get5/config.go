package get5

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Name       string `json:"name"`
	ConfigPath string `json:"viperPath"`
	Port       string `json:"port"`
}

func GetConfig(configPath string) (Config, error) {

	var c Config

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return c, fmt.Errorf("could not find the config file at path %s", configPath)
	}

	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
		return c, err
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&c)
	if err != nil {
		log.Fatalf("Error decoding config file: %v", err)
	}

	return c, nil
}
