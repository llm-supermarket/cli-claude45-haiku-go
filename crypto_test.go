package main

import (
	"bytes"
	"testing"
)

func TestEncryptDecryptSymmetric(t *testing.T) {
	plaintext := []byte("Hello, World! This is a test message.")
	password := "mypassword123"
	salt := make([]byte, 16)
	copy(salt, "testsalt12345678")

	encrypted, err := encryptData(plaintext, password, salt)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	decrypted, err := decryptData(encrypted, password)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("decrypted data doesn't match original: got %q, want %q", decrypted, plaintext)
	}
}

func TestEncryptDecryptWithoutSalt(t *testing.T) {
	plaintext := []byte("Test data without salt")
	password := "securepass"

	encrypted, err := encryptData(plaintext, password, nil)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	decrypted, err := decryptData(encrypted, password)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("decrypted data doesn't match original")
	}
}

func TestDecryptWithWrongPassword(t *testing.T) {
	plaintext := []byte("Secret message")
	password := "correctpassword"
	salt := make([]byte, 16)
	copy(salt, "testsalt12345678")

	encrypted, err := encryptData(plaintext, password, salt)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	_, err = decryptData(encrypted, "wrongpassword")
	if err == nil {
		t.Fatal("decryption with wrong password should fail")
	}
}

func TestEncryptDecryptBinaryData(t *testing.T) {
	plaintext := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD, 0x00, 0x00}
	password := "testpass"
	salt := make([]byte, 16)
	copy(salt, "binarytestsalt!")

	encrypted, err := encryptData(plaintext, password, salt)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	decrypted, err := decryptData(encrypted, password)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("binary data doesn't match")
	}
}

func TestEncryptDecryptLargeData(t *testing.T) {
	plaintext := make([]byte, 10*1024*1024)
	for i := range plaintext {
		plaintext[i] = byte(i % 256)
	}
	password := "largedata"
	salt := make([]byte, 16)
	copy(salt, "largetestsalt!!!!")

	encrypted, err := encryptData(plaintext, password, salt)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	decrypted, err := decryptData(encrypted, password)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("large data doesn't match")
	}
}


func TestPKCS7Padding(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		blockSize int
	}{
		{"empty", []byte{}, 16},
		{"exact_block", make([]byte, 16), 16},
		{"partial_block", make([]byte, 10), 16},
		{"multiple_blocks", make([]byte, 32), 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			padded := pkcs7Pad(tt.data, tt.blockSize)
			if len(padded)%tt.blockSize != 0 {
				t.Errorf("padded length not multiple of blockSize: %d", len(padded))
			}

			unpadded, err := pkcs7Unpad(padded)
			if err != nil {
				t.Fatalf("unpadding failed: %v", err)
			}

			if !bytes.Equal(tt.data, unpadded) {
				t.Errorf("unpadded data doesn't match original")
			}
		})
	}
}
