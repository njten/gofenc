package vault

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/njten/gofenc/crypto"
)

// ExtractByIndex decrypts a file by its 1-based index from the vault
func (v *Vault) ExtractByIndex(index int, outputDir string) error {
	if v.IsVaultLocked() {
		return errors.New("vault is locked — run: gofenc unlock <vault>")
	}

	if v.MasterKey == nil {
		return errors.New("master key not loaded — unlock the vault first")
	}

	if index < 1 || index > len(v.Config.Files) {
		return fmt.Errorf("invalid index %d — vault contains %d file(s)", index, len(v.Config.Files))
	}

	entry := v.Config.Files[index-1]

	// always decrypt the stored filename
	originalName, err := crypto.DecryptFilename(entry.OriginalName, v.MasterKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt filename: %w", err)
	}

	if err := os.MkdirAll(outputDir, 0700); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	encFilePath := filepath.Join(v.FilesDir(), entry.EncryptedName)
	outFilePath := filepath.Join(outputDir, originalName)

	if err := crypto.Decrypt(encFilePath, outFilePath, v.MasterKey, v.UseAES()); err != nil {
		return fmt.Errorf("decryption failed: %w", err)
	}

	fmt.Printf("Extracted: %s -> %s\n", originalName, outFilePath)
	return nil
}

// ExtractAll decrypts all files from the vault into the output directory
func (v *Vault) ExtractAll(outputDir string) error {
	if v.IsVaultLocked() {
		return errors.New("vault is locked — run: gofenc unlock <vault>")
	}

	if v.MasterKey == nil {
		return errors.New("master key not loaded — unlock the vault first")
	}

	if len(v.Config.Files) == 0 {
		fmt.Println("vault is empty — nothing to extract")
		return nil
	}

	if err := os.MkdirAll(outputDir, 0700); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	errCount := 0
	for _, f := range v.Config.Files {
		// always decrypt the stored filename
		originalName, err := crypto.DecryptFilename(f.OriginalName, v.MasterKey)
		if err != nil {
			fmt.Printf("Failed to decrypt filename: %v\n", err)
			errCount++
			continue
		}

		encFilePath := filepath.Join(v.FilesDir(), f.EncryptedName)
		outFilePath := filepath.Join(outputDir, originalName)

		if err := crypto.Decrypt(encFilePath, outFilePath, v.MasterKey, v.UseAES()); err != nil {
			fmt.Printf("Failed to extract %s — %v\n", originalName, err)
			errCount++
			continue
		}
		fmt.Printf("Extracted: %s\n", originalName)
	}

	if errCount > 0 {
		return fmt.Errorf("%d file(s) failed to extract", errCount)
	}
	return nil
}