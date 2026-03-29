package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <vault> <file>",
	Short: "Add and encrypt a file into the vault",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		v, err := loadAndUnlock(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr,"error", err)
			os.Exit(1)
		}
		
		if err := v.Add(args[1]); err != nil {
			fmt.Fprintln(os.Stderr,"error", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
