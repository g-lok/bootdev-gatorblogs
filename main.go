package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/g-lok/bootdev-gatorblogs/internal/config"
)

const userName = "g"

func main() {
	currentCfg, err := config.Read()
	if err != nil {
		err := fmt.Errorf("failed to read .gatorconfig.json: %w", err)
		fmt.Fprintf(os.Stderr, "error: %w", err)
	}

	err = currentCfg.SetUser(userName)
	if err != nil {
		err := fmt.Errorf("failed to set UserName: %w", err)
		fmt.Fprintf(os.Stderr, "error: %w", err)
	}

	updatedConfig, err := config.Read()
	if err != nil {
		err := fmt.Errorf("failed to read .gatorconfig.json: %w", err)
		fmt.Fprintf(os.Stderr, "error: %w", err)
	}

	cfg, err := json.MarshalIndent(updatedConfig, "", "  ")
	if err != nil {
		err := fmt.Errorf("failed to marshall .gatorconfig: %w", err)
		fmt.Fprintf(os.Stderr, "error: %w", err)
	}

	fmt.Println(string(cfg))
}
