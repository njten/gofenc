package cmd

import (
	"fmt"
	"os"

	"github.com/njten/gofenc/vault"
	"github.com/spf13/cobra"
)

// unlockCmd verifies the user's secret, loads the master key and removes the .locked file.
var unlockCmd = &cobra.Command{
	Use:   "unlock <vault>",
	Short: "Unlock the vault — enables add and remove",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		v, err := vault.Load(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}

		secret, err := readSecret(v.Config.Auth)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}

		if err := v.Unlock(secret); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Println("vault unlocked")
	},
}

func init() {
	rootCmd.AddCommand(unlockCmd)
}