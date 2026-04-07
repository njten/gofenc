package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

func EncryptFilename(name string, masterKey []byte) (string, error) {
	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return "", err
	}
	
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	
	nonce, err := GenerateRandom(aead.NonceSize())
	if err != nil {
		return "", err
	}
	
	encrypted := aead.Seal(nonce, nonce, []byte(name), nil)
	return base64.URLEncoding.EncodeToString(encrypted), nil
}

func DecryptFilename(encrypted string, masterKey []byte) (string, error) {
	data, err := base64.URLEncoding.DecodeString(encrypted)
	if err != nil {
		return "", errors.New("failed to decode encrypted filename")
	}
	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return "", err
	}
	
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	
	nonceSize := aead.NonceSize()
	if len(encrypted) < nonceSize {
		return "", errors.New("encrypted filename is too short")
	}
	
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("failed to decrypt filename - wrong key or corrupted data")
	}
	
	return string(plaintext), nil
}