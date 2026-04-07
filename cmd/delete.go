package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/njten/gofenc/vault"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <vault>",
	Short: "Permanently delete a vault and all its contents",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")

		if !force {
			fmt.Printf("Are you sure you want to permanently delete vault '%s'? [y/N]: ", args[0])
			var input string
			fmt.Scan(&input)
			input = strings.TrimSpace(strings.ToLower(input))
			if input != "y" {
				fmt.Println("Aborted.")
				return
			}
		}

		if err := vault.Delete(args[0]); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
		deleteCmd.Flags().Bool("force", false, "skip confirmation prompt")
}