package main

import (
	"fmt"
	"os"

	"github.com/g-lok/bootdev-gatorblogs/cmd"
)

const userName = "g"

// func getConfig() config.Config {
// 	currentCfg, err := config.Read()
// 	if err != nil {
// 		errMsg := fmt.Errorf("failed to read .gatorconfig.json: %w", err)
// 		fmt.Fprintf(os.Stderr, "error: %v\n", errMsg)
// 		os.Exit(1)
// 	}
//
// 	return *currentCfg
// }

func main() {
	// currentCfg := getConfig()
	//
	// err := currentCfg.SetUser(userName)
	// if err != nil {
	// 	errMsg := fmt.Errorf("failed to set UserName: %w", err)
	// 	fmt.Fprintf(os.Stderr, "error: %v\n", errMsg)
	// }
	//
	// updatedConfig := getConfig()
	//
	// cfg, err := json.MarshalIndent(updatedConfig, "", "  ")
	// if err != nil {
	// 	errMsg := fmt.Errorf("failed to marshall .gatorconfig: %w", err)
	// 	fmt.Fprintf(os.Stderr, "error: %v\n", errMsg)
	// 	os.Exit(1)
	// }
	//
	// fmt.Println(string(cfg))
	// args := os.Args
	// if len(args) <= 2 {
	// 	fmt.Println("error: cli requires at least 2 argument keywords")
	// }
	err := cmd.Root()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
