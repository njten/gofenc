package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list <vault>",
	Short: "Lists all files in the vault",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		v, err := loadAndUnlock(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		}
		
		if err := v.List(); err != nil {
			fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
