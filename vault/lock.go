package vault

import "fmt"

// Lock creates a .locked file in the vault directory, disabling add, remove and extract.
func (v *Vault) Lock() error {
	if v.IsVaultLocked() {
		fmt.Println("vault is already locked")
		return nil
	}

	if err := v.setLocked(true); err != nil {
		return fmt.Errorf("failed to lock vault: %w", err)
	}

	fmt.Println("vault locked")
	return nil
}