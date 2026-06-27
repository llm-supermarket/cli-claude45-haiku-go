package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
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

	key := deriveKey(password, salt)

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

	key := deriveKey(password, salt)

	var nonceArray [nonceSize]byte
	copy(nonceArray[:], nonce)

	plaintext, ok := secretbox.Open(nil, encrypted, &nonceArray, key)
	if !ok {
		return nil, fmt.Errorf("decryption failed: invalid password or corrupted data")
	}

	return plaintext, nil
}

func deriveKey(password string, salt []byte) *[32]byte {
	key, _ := scrypt.Key([]byte(password), salt, encryptionIter, 8, 1, 32)
	var keyArray [32]byte
	copy(keyArray[:], key)
	return &keyArray
}

func encryptFilename(filename string, password string, salt []byte, encoding string) (string, error) {
	if len(salt) == 0 {
		salt = make([]byte, 16)
		for i := range salt {
			salt[i] = 0
		}
	}

	key := deriveKey(password, salt)
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	filenameBytes := []byte(filename)
	padded := pkcs7Pad(filenameBytes, aes.BlockSize)

	iv := make([]byte, aes.BlockSize)
	for i := range iv {
		iv[i] = 0
	}

	ciphertext := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded)

	switch encoding {
	case "base64":
		return base64.StdEncoding.EncodeToString(ciphertext), nil
	case "base32":
		return base32.StdEncoding.EncodeToString(ciphertext), nil
	default:
		return base64.StdEncoding.EncodeToString(ciphertext), nil
	}
}

func decryptFilename(encoded string, password string, encoding string) (string, error) {
	var ciphertext []byte
	var err error

	switch encoding {
	case "base64":
		ciphertext, err = base64.StdEncoding.DecodeString(encoded)
	case "base32":
		ciphertext, err = base32.StdEncoding.DecodeString(encoded)
	default:
		ciphertext, err = base64.StdEncoding.DecodeString(encoded)
	}

	if err != nil {
		return "", err
	}

	salt := make([]byte, 16)
	for i := range salt {
		salt[i] = 0
	}

	key := deriveKey(password, salt)
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	iv := make([]byte, aes.BlockSize)
	for i := range iv {
		iv[i] = 0
	}

	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	plaintext, err = pkcs7Unpad(plaintext)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
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
