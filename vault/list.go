package vault

import (
	"errors"
	"fmt"
)

func (v *Vault) List() error {
	if len(v.Config.Files) == 0 {
		fmt.Println("No files in the vault")
		return nil
	}

	if v.isLocked() {
		return errors.New("vault is locked")
	}
	fmt.Println("%-40s %s\n", "original name", "encrypted name")
	fmt.Println("----------------------------------------------------------")
	for _, f := range v.Config.Files {
		fmt.Printf("%-40s %s\n", f.OriginalName, f.EncryptedName)
	}
	return nil
}
