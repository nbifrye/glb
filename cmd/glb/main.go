package main

import (
	"fmt"
	"os"

	"github.com/nbifrye/glb/internal/cmdutils"
	"github.com/nbifrye/glb/internal/commands"
)

var version = "dev"

func main() {
	f := cmdutils.NewFactory()
	rootCmd := commands.NewRootCmd(f, version)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
