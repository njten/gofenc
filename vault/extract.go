package vault

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/njten/gofenc/crypto"
)

func (v *Vault) Extract(originalName, outputDir string) error  {
	if v.IsLocked() {
		return errors.New("vault is locked")
	}
	
	var entry *FileEntry
	for i, f := range v.Config.Files {
		if f.OriginalName == originalName {
			entry = &v.Config.Files[i]
			break
		}
	}
	
	if entry == nil {
		return fmt.Errorf("file not found in vault: %s", originalName)
	}
	
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	encFilePath := filepath.Join(v.FilesDir(), entry.EncryptedName)
	outFilePath := filepath.Join(outputDir, entry.OriginalName)
	
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
		fmt.Println("No files in the vault")
		return nil
	}
	
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	errCount := 0
	for _, f := range v.Config.Files {
		encFilePath := filepath.Join(v.FilesDir(), f.EncryptedName)
		outFilePath := filepath.Join(outputDir, f.OriginalName)
		
		if err := crypto.Decrypt(encFilePath, outFilePath, v.MasterKey, v.UseAES()); err != nil {
			fmt.Printf("Failed to extract %s - %v\n", f.OriginalName, err)
			errCount++
			continue
		}
		fmt.Printf("Extracted: %s\n", f.OriginalName)
	}
	if errCount > 0 {
		return fmt.Errorf("%d file(s) failed to extract", errCount)
	}
	return nil
}
