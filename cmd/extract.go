package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var extractCmd = &cobra.Command{
	Use:   "extract <vault> <filename> <output-dir>",
	Short: "Decrypt and extract a file from the vault",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		v, err := loadAndUnlock(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		}
		if err := v.Extract(args[1], args[2]); err != nil {
			fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		}
	},
}

var extractAllCmd = &cobra.Command{
	Use:   "extract-all <vault> <output-dir>",
	Short: "Decrypt and extract all files from the vault",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		v, err := loadAndUnlock(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		}
		if err := v.ExtractAll(args[1]); err != nil {
			fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
	rootCmd.AddCommand(extractAllCmd)
	
}
