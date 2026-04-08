package vault

import (
	"errors"
	"fmt"

	"github.com/njten/gofenc/crypto"
)

// List prints all files stored in the vault with their index and storage name.
func (v *Vault) List() error {
	if v.IsVaultLocked() {
		return errors.New("vault is locked — unlock it first")
	}

	if len(v.Config.Files) == 0 {
		fmt.Println("vault is empty")
		return nil
	}

	fmt.Printf("%-5s %-40s %s\n", "index", "filename", "stored as")
	fmt.Println("------------------------------------------------------------------")
	for i, f := range v.Config.Files {
		name, err := crypto.DecryptFilename(f.OriginalName, v.MasterKey)
		if err != nil {
			name = "[unreadable]"
		}
		fmt.Printf("%-5d %-40s %s\n", i+1, name, f.EncryptedName)
	}
	return nil
}