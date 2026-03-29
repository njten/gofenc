package cmd

import (
	"fmt"
	"os"

	"github.com/njten/gofenc/vault"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init <path>",
	Short: "Initialize a new vault",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cipher, _ := cmd.Flags().GetString("cipher")
		auth, _ := cmd.Flags().GetString("auth")
		encryptNames, _ := cmd.Flags().GetBool("encrypt-names")
		
		var cipherType vault.CipherType
		switch cipher {
		case "aes-gcm":
			cipherType = vault.CipherAESGCM
		case "chacha20":
			cipherType = vault.CipherChaCha20
		default:
			fmt.Fprintln(os.Stderr, "unknown cipher — use aes-gcm or chacha20")
			os.Exit(1)
		}
		
		var authType vault.AuthType
		switch auth {
		case "password":
			authType = vault.AuthPassword
		case "keyfile":
			authType = vault.AuthKeyFile
		default:
			fmt.Fprintln(os.Stderr, "unknown auth — use password or keyfile")
			os.Exit(1)
				
		}
		
		secret, err := readSecret(authType)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		
		opts := vault.InitOptions{
			Cipher: 			cipherType,
			EncryptFileNames: 	encryptNames,
			Auth:   			authType,
			Secret: 			secret,
			
		}
		
		if err := vault.Init(args[0], opts); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().String("cipher", "aes-gcm", "encryption algorithm: aes-gcm or chacha20")
	initCmd.Flags().String("auth", "password", "authentication method: password or keyfile")
	initCmd.Flags().Bool("encrypt-names", false, "encrypt filenames inside the vault")
}