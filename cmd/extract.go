package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/njten/gofenc/vault"
	"github.com/spf13/cobra"
)

// extractCmd decrypts a single file from the vault by its index.
var extractCmd = &cobra.Command{
	Use:   "extract <vault> <index> <output-dir>",
	Short: "Decrypt and extract a file from the vault by its index (see: gofenc list)",
	Args:  cobra.ExactArgs(3),
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

		index, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: index must be a number")
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

		if err := v.ExtractByIndex(index, args[2]); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	},
}

// extractAllCmd decrypts all files from the vault into the output directory.
var extractAllCmd = &cobra.Command{
	Use:   "extract-all <vault> <output-dir>",
	Short: "Decrypt and extract all files from the vault",
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

		if err := v.ExtractAll(args[1]); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
	rootCmd.AddCommand(extractAllCmd)
}