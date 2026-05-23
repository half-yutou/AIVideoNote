package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

const fallbackKeyHex = "766964656f6e6f74652d64656661756c742d6b65792d706c656173652d6368616e6765"

var key []byte

func EnsureKey(envKey string) error {
	if envKey != "" {
		return SetKey(envKey)
	}
	return SetKey(fallbackKeyHex)
}

func SetKey(hexKey string) error {
	k, err := hex.DecodeString(hexKey)
	if err != nil {
		return fmt.Errorf("ENCRYPTION_KEY must be a 64-character hex string: %w", err)
	}
	if len(k) != 32 {
		return fmt.Errorf("ENCRYPTION_KEY must be 32 bytes (64 hex characters), got %d bytes", len(k))
	}
	key = k
	return nil
}

func Encrypt(plaintext string) (string, error) {
	if key == nil {
		return "", fmt.Errorf("encryption key not set")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return hex.EncodeToString(ciphertext), nil
}

func Decrypt(encoded string) (string, error) {
	if key == nil {
		return "", fmt.Errorf("encryption key not set")
	}

	data, err := hex.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("decrypt: hex decode failed: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("decrypt: ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
