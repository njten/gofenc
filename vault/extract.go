package vault

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/njten/gofenc/crypto"
)

func (v *Vault) Extract(originalName, outputDir string) error {
	if v.IsLocked() {
		return errors.New("vault is locked — unlock it first")
	}

	var entry *FileEntry
	for i, f := range v.Config.Files {
		storedName := f.OriginalName
		if v.Config.EncryptFileNames {
			decrypted, err := crypto.DecryptFilename(f.OriginalName, v.MasterKey)
			if err != nil {
				continue
			}
			storedName = decrypted
		}

		if storedName == originalName {
			entry = &v.Config.Files[i]
			break
		}
	}

	if entry == nil {
		return fmt.Errorf("file not found in vault: %s", originalName)
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

func (v *Vault) ExtractAll(outputDir string) error {
	if v.IsLocked() {
		return errors.New("vault is locked")
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
		originalName := f.OriginalName
		if v.Config.EncryptFileNames {
			decrypted, err := crypto.DecryptFilename(f.OriginalName, v.MasterKey)
			if err != nil {
				fmt.Printf("Failed to decrypt filename: %v\n", err)
				errCount++
				continue
			}
			originalName = decrypted
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