package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"

	"golang.org/x/crypto/chacha20poly1305"
)

// WrapKey encrypts the master key using the wrapping key and returns the ciphertext and nonce.
func WrapKey(masterKey, wrappingKey []byte, useAES bool) (ciphertext, nonce []byte, err error) {
	aead, err := newAEAD(wrappingKey, useAES)
	if err != nil {
		return nil, nil, err
	}

	nonce, err = GenerateRandom(aead.NonceSize())
	if err != nil {
		return nil, nil, err
	}

	ciphertext = aead.Seal(nil, nonce, masterKey, nil)
	return ciphertext, nonce, nil
}

// UnwrapKey decrypts the master key using the wrapping key.
func UnwrapKey(ciphertext, nonce, wrappingKey []byte, useAES bool) ([]byte, error) {
	aead, err := newAEAD(wrappingKey, useAES)
	if err != nil {
		return nil, err
	}

	masterKey, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("unable to unwrap key — bad password or damaged vault")
	}

	return masterKey, nil
}

// newAEAD creates an AEAD cipher using AES-256-GCM or ChaCha20-Poly1305.
func newAEAD(key []byte, useAES bool) (cipher.AEAD, error) {
	if useAES {
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}
		return cipher.NewGCM(block)
	}
	return chacha20poly1305.New(key)
}