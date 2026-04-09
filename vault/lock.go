package vault

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/njten/gofenc/crypto"
)

// Lock verifies access, re-wraps the master key and hides the files directory.
func (v *Vault) Lock(masterKey, wrappingKey []byte) error {
	if v.IsVaultLocked() {
		fmt.Println("vault is already locked")
		return nil
	}

	encMasterKey, nonce, err := crypto.WrapKey(masterKey, wrappingKey, v.UseAES())
	if err != nil {
		return fmt.Errorf("failed to re-wrap master key: %w", err)
	}

	v.Config.EncryptedMasterKey = Base64Encode(encMasterKey)
	v.Config.MasterKeyNonce = Base64Encode(nonce)

	if err := v.Save(); err != nil {
		return fmt.Errorf("failed to save vault: %w", err)
	}

	for i := range masterKey {
		masterKey[i] = 0
	}

	hiddenDir := filepath.Join(v.Path, HiddenFilesDirName)
	if err := os.Rename(v.FilesDir(), hiddenDir); err != nil {
		return fmt.Errorf("failed to hide files directory: %w", err)
	}

	if err := v.setLocked(true); err != nil {
		os.Rename(hiddenDir, v.FilesDir()) // rollback
		return fmt.Errorf("failed to lock vault: %w", err)
	}

	fmt.Println("vault locked")
	return nil
}