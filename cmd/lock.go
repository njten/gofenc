package cmd

import (
	"fmt"
	"os"

	"github.com/njten/gofenc/vault"
	"github.com/spf13/cobra"
)

// lockCmd creates a .locked file in the vault directory, disabling add, remove and extract.
var lockCmd = &cobra.Command{
	Use:   "lock <vault>",
	Short: "Lock the vault — disables add and remove until unlocked",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		v, err := vault.Load(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}

		if err := v.Lock(); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(lockCmd)
}