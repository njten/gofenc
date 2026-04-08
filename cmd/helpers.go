package cmd

import (
	"errors"
	"fmt"
	"strings"
	"syscall"

	"github.com/njten/gofenc/vault"
	"golang.org/x/term"
)

// readSecret reads a password from stdin or returns a keyfile path depending on the auth type.
func readSecret(authType vault.AuthType) (string, error) {
	if authType == vault.AuthKeyFile {
		fmt.Print("keyfile path: ")
		var path string
		fmt.Scanln(&path)
		path = strings.TrimSpace(path)
		if path == "" {
			return "", errors.New("keyfile path cannot be empty")
		}
		return path, nil
	}

	fmt.Print("password: ")
	pw, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}
	if len(pw) == 0 {
		return "", errors.New("password cannot be empty")
	}
	return string(pw), nil
}

// loadAndUnlock loads a vault and unlocks it using the user's secret.
// For keyfile auth, the keyfile path is read from stdin.
func loadAndUnlock(vaultPath string) (*vault.Vault, error) {
	v, err := vault.Load(vaultPath)
	if err != nil {
		return nil, err
	}

	secret, err := readSecret(v.Config.Auth)
	if err != nil {
		return nil, err
	}

	if err := v.Unlock(secret); err != nil {
		return nil, err
	}

	return v, nil
}