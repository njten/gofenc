package vault

import (
	"fmt"
	"os"
)

// Delete permanently removes the vault directory and all its contents from disk.
func Delete(vaultPath string) error {
	if _, err := Load(vaultPath); err != nil {
		return fmt.Errorf("not a valid vault: %w", err)
	}

	if err := os.RemoveAll(vaultPath); err != nil {
		return fmt.Errorf("failed to delete vault: %w", err)
	}

	fmt.Printf("Vault deleted: %s\n", vaultPath)
	return nil
}