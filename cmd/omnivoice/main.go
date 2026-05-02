// Package main provides the entry point for the omnivoice CLI.
package main

import (
	"os"

	"github.com/plexusone/omnivoice/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
