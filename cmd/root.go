package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gofenc",
	Short: "A simple file encryption tool using vaults",
	Long:  "gofenc is a CLI tool for encrypting files and folders using a vault-based approach.",
}

// Execute runs the root command and exits on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}