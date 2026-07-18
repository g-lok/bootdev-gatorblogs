package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	URL      string `json:"db_url"`
	UserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configFilePath := fmt.Sprintf("%s/%s", homeDir, configFileName)

	return configFilePath, nil
}

func Read() (*Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

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

func write(cfg Config) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	err = os.WriteFile(configFilePath, data, 0o666)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (c *Config) SetUser(userName string) error {
	c.UserName = userName

	err := write(*c)
	if err != nil {
		return err
	}

	return nil
}
