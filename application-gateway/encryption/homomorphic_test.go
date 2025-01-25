package encryption

import (
	"fmt"
	"testing"
)

func TestEncryptor(t *testing.T) {
	encryptor := NewEncryptor()

	t.Run("Encrypt and Decrypt with Public Key", func(t *testing.T) {
		data := int64(32 * 1000 * 1000)

		cipherText, err := encryptor.EncryptPu(data)
		if err != nil {
			t.Fatalf("failed to encrypt with public key: %v", err)
		}

		decryptedData, err := encryptor.Decrypt(cipherText)
		if err != nil {
			t.Fatalf("failed to decrypt with public key: %v", err)
		}

		if decryptedData != data {
			t.Errorf("decrypted data mismatch: got %v, want %v", decryptedData, data)
		}

		fmt.Printf("decrypted data: %v", decryptedData)
	})

	t.Run("Encrypt and Decrypt with Private Key", func(t *testing.T) {
		data := int64(20000)

		cipherText, err := encryptor.EncryptPr(data)
		if err != nil {
			t.Fatalf("failed to encrypt with private key: %v", err)
		}

		decryptedData, err := encryptor.Decrypt(cipherText)
		if err != nil {
			t.Fatalf("failed to decrypt with private key: %v", err)
		}

		if decryptedData != data {
			t.Errorf("decrypted data mismatch: got %v, want %v", decryptedData, data)
		}
	})

	//t.Run("Invalid Ciphertext", func(t *testing.T) {
	//	invalidCipher := &rlwe.Ciphertext{
	//		Value: nil,
	//	}
	//	_, err := encryptor.Decrypt(invalidCipher)
	//	if err == nil {
	//		t.Error("expected an error for invalid ciphertext, got none")
	//	}
	//})
}
