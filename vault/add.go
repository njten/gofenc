package vault

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/njten/gofenc/crypto"
)

// FileEntry maps the original encrypted filename to its UUID-based storage name.
type FileEntry struct {
	OriginalName  string `json:"original_name"`
	EncryptedName string `json:"encrypted_name"`
}

// Add encrypts a file or all files in a directory and stores them in the vault.
func (v *Vault) Add(filePath string) error {
	if v.IsVaultLocked() {
		return errors.New("vault is locked — run: gofenc unlock <vault>")
	}

	info, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file or directory not found: %s", filePath)
	}

	if info.IsDir() {
		return v.addDirectory(filePath)
	}

	return v.addFile(filePath)
}

// addDirectory encrypts all files in the given directory and adds them to the vault.
func (v *Vault) addDirectory(dirPath string) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println("directory is empty — nothing to add")
		return nil
	}

	errCount := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if err := v.addFile(filepath.Join(dirPath, entry.Name())); err != nil {
			fmt.Printf("✗ Failed to add %s: %v\n", entry.Name(), err)
			errCount++
		}
	}

	if errCount > 0 {
		return fmt.Errorf("%d file(s) failed to add", errCount)
	}
	return nil
}

// addFile encrypts a single file and stores it in the vault files/ directory.
func (v *Vault) addFile(filePath string) error {
	encFileName := uuid.New().String() + ".enc"
	encFilePath := filepath.Join(v.FilesDir(), encFileName)

	if err := crypto.Encrypt(filePath, encFilePath, v.MasterKey, v.UseAES()); err != nil {
		return fmt.Errorf("encryption failed: %w", err)
	}

	originalName := filepath.Base(filePath)
	storedName, err := crypto.EncryptFilename(originalName, v.MasterKey)
	if err != nil {
		os.Remove(encFilePath)
		return fmt.Errorf("failed to encrypt filename: %w", err)
	}

	v.Config.Files = append(v.Config.Files, FileEntry{
		OriginalName:  storedName,
		EncryptedName: encFileName,
	})

	if err := v.Save(); err != nil {
		os.Remove(encFilePath)
		return fmt.Errorf("failed to save vault.json: %w", err)
	}

	fmt.Printf("Added: %s -> %s\n", originalName, encFileName)
	return nil
}