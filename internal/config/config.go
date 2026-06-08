package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	URL      string `json:"db_url"`
	Username string `json:"current_user_name"`
}

func Read() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configFilePath := fmt.Sprintf("%s/.gatorconfig.json", homeDir)

	configFile, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	configByte, err := io.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(configByte, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
