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
	CurrentVersion = 1
)

type CipherType string

const (
	CipherAESGCM   CipherType = "aes-256-gcm"
	CipherChaCha20 CipherType = "chacha20-poly1305"
)

type AuthType string

const (
	AuthPassword AuthType = "password"
	AuthKeyFile  AuthType = "keyfile"
)

type KDFParams struct {
	Time    uint32 `json:"time"`
	Memory  uint32 `json:"memory"`
	Threads uint8  `json:"threads"`
	Salt    string `json:"salt"`
}

type Config struct {
	Version            int         `json:"version"`
	Cipher             CipherType  `json:"cipher"`
	KDF                string      `json:"kdf"`
	KDFParams          KDFParams   `json:"kdf_params"`
	EncryptFileNames   bool        `json:"encrypt_filenames"`
	Auth               AuthType    `json:"auth"`
	EncryptedMasterKey string      `json:"encrypted_master_key"`
	MasterKeyNonce     string      `json:"master_key_nonce"`
	Files              []FileEntry `json:"files"`
}

type Vault struct {
	Path      string
	Config    Config
	MasterKey []byte
}

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
		return nil, errors.New("vault.json is not compatible")
	}

	return &Vault{Path: absPath, Config: cfg}, nil
}

func (v *Vault) Save() error {
	data, err := json.MarshalIndent(v.Config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(v.Path, VaultFileName), data, 0600)
}

func (v *Vault) FilesDir() string {
	return filepath.Join(v.Path, FilesDirName)
}

func (v *Vault) isLocked() bool {
	return v.MasterKey == nil
}

func (v *Vault) UseAES() bool {
	return v.Config.Cipher == CipherAESGCM
}

func Base64Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Base64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
