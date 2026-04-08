package vault

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/njten/gofenc/crypto"
)

// Remove deletes a file from the vault by index (1-based), original name, or encrypted filename.
func (v *Vault) Remove(identifier string) error {
	if v.IsVaultLocked() {
		return errors.New("vault is locked — run: gofenc unlock <vault>")
	}

	if v.MasterKey == nil {
		return errors.New("master key not loaded — unlock the vault first")
	}

	index := -1

	if n, err := strconv.Atoi(identifier); err == nil {
		if n < 1 || n > len(v.Config.Files) {
			return fmt.Errorf("invalid index %d — vault contains %d file(s)", n, len(v.Config.Files))
		}
		index = n - 1
	} else {
		for i, f := range v.Config.Files {
			decrypted, err := crypto.DecryptFilename(f.OriginalName, v.MasterKey)
			if err == nil && decrypted == identifier {
				index = i
				break
			}
			if f.EncryptedName == identifier {
				index = i
				break
			}
		}
	}

	if index == -1 {
		return fmt.Errorf("file not found: %s", identifier)
	}

	encFilePath := filepath.Join(v.FilesDir(), v.Config.Files[index].EncryptedName)
	if err := os.Remove(encFilePath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to delete encrypted file: %w", err)
	}

	removed := v.Config.Files[index].EncryptedName
	v.Config.Files = append(v.Config.Files[:index], v.Config.Files[index+1:]...)

	if err := v.Save(); err != nil {
		return fmt.Errorf("failed to save vault.json: %w", err)
	}

	fmt.Printf("Removed: %s\n", removed)
	return nil
}