package cmd

import (
	"fmt"

	"github.com/njten/gofenc/crypto"
	"github.com/njten/gofenc/vault"
	"github.com/spf13/cobra"
)

// lockCmd creates a .locked file in the vault directory, disabling add, remove and extract.
var lockCmd = &cobra.Command{
	Use:   "lock <vault>",
	Short: "Lock the vault",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		v, err := vault.Load(args[0])
		if err != nil {
			return err
		}

		secret, err := readSecret(v.Config.Auth)
		if err != nil {
			return err
		}

		salt, err := vault.Base64Decode(v.Config.KDFParams.Salt)
		if err != nil {
			return err
		}

		wrappingKey, err := crypto.DeriveKey(secret, v.Config.Auth == vault.AuthKeyFile, crypto.Argon2Params{
			Time:    v.Config.KDFParams.Time,
			Memory:  v.Config.KDFParams.Memory,
			Threads: v.Config.KDFParams.Threads,
			Salt:    salt,
		})
		if err != nil {
			return err
		}

		encMasterKey, err := vault.Base64Decode(v.Config.EncryptedMasterKey)
		if err != nil {
			return err
		}
		nonce, err := vault.Base64Decode(v.Config.MasterKeyNonce)
		if err != nil {
			return err
		}

		masterKey, err := crypto.UnwrapKey(encMasterKey, nonce, wrappingKey, v.UseAES())
		if err != nil {
			return fmt.Errorf("bad password or damaged vault")
		}

		return v.Lock(masterKey, wrappingKey)
	},
}