package main

import (
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
)

const (
	nonceSize      = 24
	encryptionIter = 2 << 16
)

func encryptData(plaintext []byte, password string, salt []byte) ([]byte, error) {
	if len(salt) == 0 {
		salt = make([]byte, 16)
		if _, err := rand.Read(salt); err != nil {
			return nil, err
		}
	}

	key, err := deriveKey(password, salt)
	if err != nil {
		return nil, err
	}

	var nonce [nonceSize]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, err
	}

	ciphertext := secretbox.Seal(nil, plaintext, &nonce, key)

	result := make([]byte, len(salt)+len(nonce)+len(ciphertext))
	copy(result, salt)
	copy(result[len(salt):], nonce[:])
	copy(result[len(salt)+len(nonce):], ciphertext)

	return result, nil
}

func decryptData(ciphertext []byte, password string) ([]byte, error) {
	if len(ciphertext) < 16+nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	salt := ciphertext[:16]
	nonce := ciphertext[16 : 16+nonceSize]
	encrypted := ciphertext[16+nonceSize:]

	key, err := deriveKey(password, salt)
	if err != nil {
		return nil, err
	}

	var nonceArray [nonceSize]byte
	copy(nonceArray[:], nonce)

	plaintext, ok := secretbox.Open(nil, encrypted, &nonceArray, key)
	if !ok {
		return nil, fmt.Errorf("decryption failed: invalid password or corrupted data")
	}

	return plaintext, nil
}

func deriveKey(password string, salt []byte) (*[32]byte, error) {
	key, err := scrypt.Key([]byte(password), salt, encryptionIter, 8, 1, 32)
	if err != nil {
		return nil, fmt.Errorf("key derivation failed: %w", err)
	}
	var keyArray [32]byte
	copy(keyArray[:], key)
	return &keyArray, nil
}


func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := make([]byte, padding)
	for i := 0; i < padding; i++ {
		padText[i] = byte(padding)
	}
	return append(data, padText...)
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("invalid padding")
	}
	padding := int(data[len(data)-1])
	if padding > len(data) || padding == 0 {
		return nil, fmt.Errorf("invalid padding")
	}
	for i := 0; i < padding; i++ {
		if data[len(data)-1-i] != byte(padding) {
			return nil, fmt.Errorf("invalid padding")
		}
	}
	return data[:len(data)-padding], nil
}
