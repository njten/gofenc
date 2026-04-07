package cmd

import (
	"bufio"
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
			fmt.Println("Are you sure you want to permanently delete vault '%s'? [y/N]: ", args[0])
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(strings.ToLower(input))
			if input != "y\n" {
				fmt.Println("Aborting")
				return
			}
		}
		if err := vault.Delete(args[0]); err != nil {
			fmt.Fprintln(os.Stderr, "error", err)
				os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
		deleteCmd.Flags().Bool("force", false, "skip confirmation prompt")
}