package vault

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/njten/gofenc/crypto"
)

type FileEntry struct {
	OriginalName  string `json:"original_name"`
	EncryptedName string `json:"encrypted_name"`
}

func (v *Vault) Add(filePath string) error {
	if v.isLocked() {
		return errors.New("vault is locked")
	}
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file not found: %s", filePath)
	}

	encFileName := uuid.New().String() + ".enc"
	encFilePath := filepath.Join(v.FilesDir(), encFileName)

	if err := crypto.Encrypt(filePath, encFileName, v.MasterKey, v.UseAES()); err != nil {
		return fmt.Errorf("encryption failed: %w", err)
	}

	originalName := filepath.Base(filePath)
	v.Config.Files = append(v.Config.Files, FileEntry{
		OriginalName:  originalName,
		EncryptedName: encFileName,
	})

	if err := v.Save(); err != nil {
		os.Remove(encFilePath)
		return fmt.Errorf("failed to save vault.json: %w", err)
	}

	fmt.Printf("Added: %s → %s\\n", originalName, encFileName)
	return nil
}
