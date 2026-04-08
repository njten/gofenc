package cmd

import (
	"fmt"
	"os"

	"github.com/njten/gofenc/vault"
	"github.com/spf13/cobra"
)

// removeCmd deletes a file from the vault by index or filename.
var removeCmd = &cobra.Command{
	Use:   "remove <vault> <index|filename>",
	Short: "Remove a file from the vault",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		v, err := vault.Load(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}

		if v.IsVaultLocked() {
			fmt.Fprintln(os.Stderr, "error: vault is locked — run: gofenc unlock <vault>")
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

		if err := v.Remove(args[1]); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}