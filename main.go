package main

import (
	"os"

	"github.com/Radical-Egg/steamcmd-healthcheck/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
