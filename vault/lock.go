package vault

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/njten/gofenc/crypto"
)

func (v *Vault) Lock(inputDir string) error {
	if v.IsLocked() {
		return errors.New("vault is already locked")
	}
	
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return fmt.Errorf("cannot read input directory: %w", err)
	}
	
	if len(entries) == 0 {
		fmt.Println("No files to lock")
		return nil
	}
	
	errCount := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		filePath := filepath.Join(inputDir, entry.Name())
		
		alreadyAdded := false
		for _, f := range v.Config.Files {
			if f.OriginalName == entry.Name() {
				alreadyAdded = true
				break
			}
		}
		
		if alreadyAdded {
			fmt.Printf("Skipped (already in vault): %s\n", entry.Name())
			continue
		}
		
		encFileName := entry.Name()
		
		if err := crypto.Encrypt(filePath, encFileName, v.MasterKey, v.UseAES()); err != nil {
			fmt.Printf("Failed to encrypt %s - %v\n", entry.Name(), err)
				errCount++
				continue
		}
		
		v.Config.Files = append(v.Config.Files, FileEntry{
			OriginalName:  entry.Name(),
			EncryptedName: encFileName + ".enc",
		})
		
		if err := os.Remove(filePath); err != nil {
			fmt.Printf("Encrypted but failed to remove plaintext: %s\n", entry.Name())
		} else {
			fmt.Printf("Locked: %s\n", entry.Name())
		}
	}
	
	if err := v.Save(); err != nil {
		return fmt.Errorf("failed to save vault.json: %w", err)
	}
	
	if errCount > 0 {
		return fmt.Errorf("%d file(s) failed to lock", errCount)
	}
	
	return nil
}