package main

import (
	"os"

	"github.com/gravitational/gamma/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
