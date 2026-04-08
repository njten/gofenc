package vault

import (
	"errors"
	"fmt"

	"github.com/njten/gofenc/crypto"
)

// Unlock derives the wrapping key from the user's secret, decrypts the master key,
// and removes the .locked file to enable vault operations.
func (v *Vault) Unlock(secret string) error {
	salt, err := Base64Decode(v.Config.KDFParams.Salt)
	if err != nil {
		return errors.New("failed to decode salt from vault.json")
	}

	wrappingKey, err := crypto.DeriveKey(secret, v.Config.Auth == AuthKeyFile, crypto.Argon2Params{
		Time:    v.Config.KDFParams.Time,
		Memory:  v.Config.KDFParams.Memory,
		Threads: v.Config.KDFParams.Threads,
		Salt:    salt,
	})
	if err != nil {
		return err
	}

	encMasterKey, err := Base64Decode(v.Config.EncryptedMasterKey)
	if err != nil {
		return errors.New("failed to decode master key from vault.json")
	}

	nonce, err := Base64Decode(v.Config.MasterKeyNonce)
	if err != nil {
		return errors.New("failed to decode nonce from vault.json")
	}

	masterKey, err := crypto.UnwrapKey(encMasterKey, nonce, wrappingKey, v.UseAES())
	if err != nil {
		return errors.New("bad password or damaged vault")
	}

	v.MasterKey = masterKey

	if err := v.setLocked(false); err != nil {
		return fmt.Errorf("failed to unlock vault: %w", err)
	}

	return nil
}