package vault

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func (v *Vault) Remove(originalName string) error {
	if v.IsLocked() {
		return errors.New("vault is locked")
	}

	index := -1
	for i, f := range v.Config.Files {
		if f.OriginalName == originalName {
			index = i
			break
		}
	}
	if index == -1 {
		return fmt.Errorf("file not found in vault: %s", originalName)
	}

	encFilePath := filepath.Join(v.FilesDir(), v.Config.Files[index].EncryptedName)
	if err := os.Remove(encFilePath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to remove encrypted file: %w", err)
	}

	removed := v.Config.Files[index].EncryptedName
	v.Config.Files = append(v.Config.Files[:index], v.Config.Files[index+1:]...)

	if err := v.Save(); err != nil {
		return fmt.Errorf("failed to save vault.json: %w", err)
	}
	fmt.Printf("Removed: %s (%s)\n", originalName, removed)
	return nil
}