package app

import (
	"os"

	"github.com/jandedobbeleer/aliae/src/cli"
)

// Execute is the application composition root entrypoint.
// It centralizes startup wiring while keeping main.go minimal.
func Execute(version string) {
	root := cli.NewRootCommand(version)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
