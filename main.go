package main

import (
	"fmt"
	"os"

	"github.com/g-lok/bootdev-gatorblogs/cmd"

	_ "github.com/lib/pq"
)

func main() {
	err := cmd.Root()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
