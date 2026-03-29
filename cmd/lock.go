package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var lockCmd = &cobra.Command{
	Use:   "lock <vault> <input-dir>",
	Short: "Encrypt all files from a directory into the vault and remove the originals",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		v, err := loadAndUnlock(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		}
		
		if err := v.Lock(args[1]); err != nil {
			fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		}
	},	
}

func init() {
	rootCmd.AddCommand(lockCmd)
}
