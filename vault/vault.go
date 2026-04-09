package vault

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	VaultFileName  = "vault.json"
	FilesDirName   = "files"
	HiddenFilesDirName = ".files"
	LockedVault    = ".locked"
	CurrentVersion = 1
)

// CipherType represents the supported encryption algorithms.
type CipherType string

const (
	CipherAESGCM   CipherType = "aes-256-gcm"
	CipherChaCha20 CipherType = "chacha20-poly1305"
)

// AuthType represents the supported authentication methods.
type AuthType string

const (
	AuthPassword AuthType = "password"
	AuthKeyFile  AuthType = "keyfile"
)

// KDFParams holds the Argon2id key derivation parameters stored in vault.json.
type KDFParams struct {
	Time    uint32 `json:"time"`
	Memory  uint32 `json:"memory"`
	Threads uint8  `json:"threads"`
	Salt    string `json:"salt"`
}

// Config represents the vault.json configuration file.
type Config struct {
	Version            int         `json:"version"`
	Cipher             CipherType  `json:"cipher"`
	KDF                string      `json:"kdf"`
	KDFParams          KDFParams   `json:"kdf_params"`
	Auth               AuthType    `json:"auth"`
	EncryptedMasterKey string      `json:"encrypted_master_key"`
	MasterKeyNonce     string      `json:"master_key_nonce"`
	Files              []FileEntry `json:"files"`
}

// Vault represents an open vault with its configuration and optionally loaded master key.
type Vault struct {
	Path      string
	Config    Config
	MasterKey []byte
}

// Load reads and parses vault.json from the given path and returns a Vault instance.
func Load(vaultPath string) (*Vault, error) {
	absPath, err := filepath.Abs(vaultPath)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filepath.Join(absPath, VaultFileName))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.New("vault.json not found")
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, errors.New("vault.json is damaged or corrupted")
	}

	if cfg.Version != CurrentVersion {
		return nil, errors.New("vault version is not compatible")
	}

	return &Vault{Path: absPath, Config: cfg}, nil
}

// Save writes the current vault configuration to vault.json.
func (v *Vault) Save() error {
	data, err := json.MarshalIndent(v.Config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(v.Path, VaultFileName), data, 0600)
}

// FilesDir returns the absolute path to the files/ directory inside the vault.
func (v *Vault) FilesDir() string {
	return filepath.Join(v.Path, FilesDirName)
}

// IsVaultLocked reports whether the vault is locked by checking for a .locked file.
func (v *Vault) IsVaultLocked() bool {
	_, err := os.Stat(filepath.Join(v.Path, LockedVault))
	return err == nil
}

// setLocked creates or removes the .locked file to change the vault lock state.
func (v *Vault) setLocked(locked bool) error {
	lockPath := filepath.Join(v.Path, LockedVault)
	if locked {
		f, err := os.Create(lockPath)
		if err != nil {
			return err
		}
		f.Close()
		return nil
	}
	err := os.Remove(lockPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// UseAES reports whether the vault is configured to use AES-256-GCM.
func (v *Vault) UseAES() bool {
	return v.Config.Cipher == CipherAESGCM
}

// Base64Encode encodes a byte slice to a base64 string.
func Base64Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// Base64Decode decodes a base64 string to a byte slice.
func Base64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}