package vault

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/njten/gofenc/crypto"
)

// InitOptions holds the configuration options for creating a new vault.
type InitOptions struct {
	Cipher CipherType
	Auth   AuthType
	Secret string
}

// Init creates a new vault at the given path with the provided options.
// If auth is keyfile, a keyfile is automatically generated next to the vault.
func Init(vaultPath string, opts InitOptions) error {
	if err := validateNewVaultPath(vaultPath); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Join(vaultPath, FilesDirName), 0700); err != nil {
		return fmt.Errorf("the vault structure cannot be created: %w", err)
	}

	// generate keyfile automatically if auth is keyfile
	if opts.Auth == AuthKeyFile {
		keyfilePath := vaultPath + ".key"
		keyfileData, err := crypto.GenerateRandom(32)
		if err != nil {
			return fmt.Errorf("failed to generate keyfile: %w", err)
		}
		if err := os.WriteFile(keyfilePath, keyfileData, 0600); err != nil {
			return fmt.Errorf("failed to write keyfile: %w", err)
		}
		opts.Secret = keyfilePath
		fmt.Printf("Keyfile generated: %s — keep it safe!\n", keyfilePath)
	}

	salt, err := crypto.GenerateRandom(32)
	if err != nil {
		return err
	}

	kdfParams := crypto.DefaultArgon2Params(salt)
	wrappingKey, err := crypto.DeriveKey(opts.Secret, opts.Auth == AuthKeyFile, kdfParams)
	if err != nil {
		return err
	}

	masterKey, err := crypto.GenerateRandom(32)
	if err != nil {
		return err
	}

	encMasterKey, nonce, err := crypto.WrapKey(masterKey, wrappingKey, opts.Cipher == CipherAESGCM)
	if err != nil {
		return err
	}

	v := &Vault{
		Path: vaultPath,
		Config: Config{
			Version: CurrentVersion,
			Cipher:  opts.Cipher,
			KDF:     "argon2id",
			KDFParams: KDFParams{
				Time:    kdfParams.Time,
				Memory:  kdfParams.Memory,
				Threads: kdfParams.Threads,
				Salt:    Base64Encode(salt),
			},
			Auth:               opts.Auth,
			EncryptedMasterKey: Base64Encode(encMasterKey),
			MasterKeyNonce:     Base64Encode(nonce),
		},
	}

	if err := v.Save(); err != nil {
		return fmt.Errorf("the vault cannot be saved: %w", err)
	}

	fmt.Printf("Your new vault is at %s\n", vaultPath)
	fmt.Printf("Algorithm: %s\n", opts.Cipher)
	fmt.Printf("Auth: %s\n", opts.Auth)
	return nil
}

// validateNewVaultPath checks that the given path is suitable for a new vault.
func validateNewVaultPath(path string) error {
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errors.New("path is not a directory")
	}
	entries, _ := os.ReadDir(path)
	for _, e := range entries {
		if e.Name() == VaultFileName {
			return errors.New("vault.json already exists")
		}
	}
	return nil
}