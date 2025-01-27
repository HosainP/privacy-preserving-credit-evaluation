package encryption

import (
	"testing"
)

// TestNewEncryptor tests the creation of a new Encryptor instance.
func TestNewEncryptor(t *testing.T) {
	encryptor := NewEncryptor()
	if encryptor == nil {
		t.Fatal("failed to create Encryptor")
	}
}

// TestEncryptPu tests public key encryption and decryption.
func TestEncryptPu(t *testing.T) {
	encryptor := NewEncryptor()

	// Test data
	data := int64(31000123)

	// Encrypt
	cipherText, err := encryptor.EncryptPu(data)
	if err != nil {
		t.Fatalf("failed to encrypt data: %v", err)
	}

	// Decrypt
	decryptedData, err := encryptor.Decrypt(cipherText)
	if err != nil {
		t.Fatalf("failed to decrypt data: %v", err)
	}

	// Verify
	if decryptedData != data {
		t.Errorf("decrypted data does not match original data: got %d, want %d", decryptedData, data)
	}
}

// TestEncryptPr tests private key encryption and decryption.
func TestEncryptPr(t *testing.T) {
	encryptor := NewEncryptor()

	// Test data
	data := int64(30678543)

	// Encrypt
	cipherText, err := encryptor.EncryptPr(data)
	if err != nil {
		t.Fatalf("failed to encrypt data: %v", err)
	}

	// Decrypt
	decryptedData, err := encryptor.Decrypt(cipherText)
	if err != nil {
		t.Fatalf("failed to decrypt data: %v", err)
	}

	// Verify
	if decryptedData != data {
		t.Errorf("decrypted data does not match original data: got %d, want %d", decryptedData, data)
	}
}

// TestDecryptWithKey tests decryption with a specific private key.
func TestDecryptWithKey(t *testing.T) {
	encryptor := NewEncryptor()

	data := int64(32787654)

	cipherText, err := encryptor.EncryptPu(data)
	if err != nil {
		t.Fatalf("failed to encrypt data: %v", err)
	}

	decryptedData, err := encryptor.DecryptWithKey(cipherText, encryptor.privateKey)
	if err != nil {
		t.Fatalf("failed to decrypt data with key: %v", err)
	}

	if decryptedData != data {
		t.Errorf("decrypted data does not match original data: got %d, want %d", decryptedData, data)
	}
}

// TestEncryptDecryptEdgeCases tests edge cases for encryption and decryption.
func TestEncryptDecryptEdgeCases(t *testing.T) {
	encryptor := NewEncryptor()

	tests := []struct {
		name string
		data int64
	}{
		{"Zero", 0},
		{"Negative Number", -123456},
		{"Large Number", 32900000}, // Close to PlaintextModulus / 2
		{"Small Number", 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt
			cipherText, err := encryptor.EncryptPu(tt.data)
			if err != nil {
				t.Fatalf("failed to encrypt data: %v", err)
			}

			// Decrypt
			decryptedData, err := encryptor.Decrypt(cipherText)
			if err != nil {
				t.Fatalf("failed to decrypt data: %v", err)
			}

			// Verify
			if decryptedData != tt.data {
				t.Errorf("decrypted data does not match original data: got %d, want %d", decryptedData, tt.data)
			}
		})
	}
}

// TestEncryptInvalidData tests encryption with invalid data.
//func TestEncryptInvalidData(t *testing.T) {
//	encryptor := NewEncryptor()
//
//	// Test data exceeding PlaintextModulus
//	data := int64(1000000008) // Exceeds PlaintextModulus (0x3ee0001 = 65,536,001)
//
//	// Encrypt
//_, err := encryptor.EncryptPu(data)
//if err == nil {
//	t.Error("expected error for data exceeding PlaintextModulus, got nil")
//}
//}

// TestDecryptInvalidCipherText tests decryption with invalid ciphertext.
//func TestDecryptInvalidCipherText(t *testing.T) {
//	encryptor := NewEncryptor()
//
//	// Invalid ciphertext
//invalidCipherText := &rlwe.Ciphertext{}
//
//	// Decrypt
//_, err := encryptor.Decrypt(invalidCipherText)
//if err == nil {
//	t.Error("expected error for invalid ciphertext, got nil")
//}
//}

// TestAddition tests homomorphic addition on encrypted data.
func TestAddition(t *testing.T) {
	encryptor := NewEncryptor()

	// Test data
	data1 := int64(3456789)
	data2 := int64(17654321)

	// Encrypt data
	cipherText1, err := encryptor.EncryptPu(data1)
	if err != nil {
		t.Fatalf("failed to encrypt data1: %v", err)
	}

	cipherText2, err := encryptor.EncryptPu(data2)
	if err != nil {
		t.Fatalf("failed to encrypt data2: %v", err)
	}

	// Perform homomorphic addition
	sumCipherText, err := encryptor.Add(cipherText1, cipherText2)
	if err != nil {
		t.Fatalf("failed to add ciphertexts: %v", err)
	}

	// Decrypt the result
	sum, err := encryptor.Decrypt(sumCipherText)
	if err != nil {
		t.Fatalf("failed to decrypt sum: %v", err)
	}

	// Verify the result
	expectedSum := data1 + data2
	if sum != expectedSum {
		t.Errorf("addition result is incorrect: got %d, want %d", sum, expectedSum)
	}
}

// TestMultiplication tests homomorphic multiplication on encrypted data.
func TestMultiplication(t *testing.T) {
	encryptor := NewEncryptor()

	// Test data
	data1 := int64(12345)
	data2 := int64(789)

	// Encrypt data
	cipherText1, err := encryptor.EncryptPu(data1)
	if err != nil {
		t.Fatalf("failed to encrypt data1: %v", err)
	}

	cipherText2, err := encryptor.EncryptPu(data2)
	if err != nil {
		t.Fatalf("failed to encrypt data2: %v", err)
	}

	// Perform homomorphic multiplication
	productCipherText, err := encryptor.Multiply(cipherText1, cipherText2)
	if err != nil {
		t.Fatalf("failed to multiply ciphertexts: %v", err)
	}

	// Decrypt the result
	product, err := encryptor.Decrypt(productCipherText)
	if err != nil {
		t.Fatalf("failed to decrypt product: %v", err)
	}

	// Verify the result
	expectedProduct := data1 * data2
	if product != expectedProduct {
		t.Errorf("multiplication result is incorrect: got %d, want %d", product, expectedProduct)
	}
}

// TestMultiplyWithPlain tests homomorphic multiplication of ciphertext with plaintext data.
func TestMultiplyWithPlain(t *testing.T) {
	encryptor := NewEncryptor()

	// Test data
	data := int64(12345)
	plainNumber := int64(789)

	// Encrypt data
	cipherText, err := encryptor.EncryptPu(data)
	if err != nil {
		t.Fatalf("failed to encrypt data: %v", err)
	}

	// Perform homomorphic multiplication with plaintext
	resultCipherText, err := encryptor.MultiplyWithPlain(cipherText, plainNumber)
	if err != nil {
		t.Fatalf("failed to multiply ciphertext with plaintext: %v", err)
	}

	// Decrypt the result
	result, err := encryptor.Decrypt(resultCipherText)
	if err != nil {
		t.Fatalf("failed to decrypt result: %v", err)
	}

	// Verify the result
	expectedResult := data * plainNumber
	if result != expectedResult {
		t.Errorf("multiplication with plaintext result is incorrect: got %d, want %d", result, expectedResult)
	}
}

// TestAdditionAndMultiplication tests a combination of addition and multiplication.
func TestAdditionAndMultiplication(t *testing.T) {
	encryptor := NewEncryptor()

	// Test data
	data1 := int64(123)
	data2 := int64(456)
	data3 := int64(789)

	// Encrypt data
	cipherText1, err := encryptor.EncryptPu(data1)
	if err != nil {
		t.Fatalf("failed to encrypt data1: %v", err)
	}

	cipherText2, err := encryptor.EncryptPu(data2)
	if err != nil {
		t.Fatalf("failed to encrypt data2: %v", err)
	}

	cipherText3, err := encryptor.EncryptPu(data3)
	if err != nil {
		t.Fatalf("failed to encrypt data3: %v", err)
	}

	// Perform homomorphic addition: data1 + data2
	sumCipherText, err := encryptor.Add(cipherText1, cipherText2)
	if err != nil {
		t.Fatalf("failed to add ciphertexts: %v", err)
	}

	// Perform homomorphic multiplication: (data1 + data2) * data3
	productCipherText, err := encryptor.Multiply(sumCipherText, cipherText3)
	if err != nil {
		t.Fatalf("failed to multiply ciphertexts: %v", err)
	}

	// Decrypt the result
	result, err := encryptor.Decrypt(productCipherText)
	if err != nil {
		t.Fatalf("failed to decrypt result: %v", err)
	}

	// Verify the result
	expectedResult := (data1 + data2) * data3
	if result != expectedResult {
		t.Errorf("combined operation result is incorrect: got %d, want %d", result, expectedResult)
	}
}

// TestDivideByPlain tests homomorphic division of ciphertext by a plaintext number.
func TestDivideByPlain(t *testing.T) {
	encryptor := NewEncryptor()

	// Test data
	data := int64(12345)
	plainNumber := int64(5)

	// Encrypt data
	cipherText, err := encryptor.EncryptPu(data)
	if err != nil {
		t.Fatalf("failed to encrypt data: %v", err)
	}

	// Perform homomorphic division by plaintext
	resultCipherText, err := encryptor.DivideByPlain(cipherText, plainNumber)
	if err != nil {
		t.Fatalf("failed to divide ciphertext by plaintext: %v", err)
	}

	// Decrypt the result
	result, err := encryptor.Decrypt(resultCipherText)
	if err != nil {
		t.Fatalf("failed to decrypt result: %v", err)
	}

	// Verify the result
	expectedResult := data / plainNumber
	if result != expectedResult {
		t.Errorf("division result is incorrect: got %d, want %d", result, expectedResult)
	}
}

////////////////// credit evaluation tests

// TestHomomorphicCreditEvaluation tests the homomorphic credit evaluation process.
func TestHomomorphicCreditEvaluation(t *testing.T) {
	encryptor := NewEncryptor()

	// Test data: Assume Ciphertext for age, salary, credit score, and DTI
	age, _ := encryptor.EncryptPu(int64(30))          // Example age
	salary, _ := encryptor.EncryptPu(int64(50000))    // Example salary
	creditScore, _ := encryptor.EncryptPu(int64(750)) // Example credit score
	dti, _ := encryptor.EncryptPu(int64(3))           // Example DTI //todo 0.3

	// Test case when preselection criteria are met
	encryptedResult, err := encryptor.HomomorphicCreditEvaluation(age, salary, creditScore, dti)
	if err != nil {
		t.Fatalf("failed to perform homomorphic credit evaluation: %v", err)
	}

	// Decrypt and verify the result
	result, err := encryptor.Decrypt(encryptedResult)
	if err != nil {
		t.Fatalf("failed to decrypt the result: %v", err)
	}

	// Assuming the expected score should be a positive value when eligible
	if result <= 0 {
		t.Errorf("unexpected evaluation result: got %v, expected positive", result)
	}

	// Test case when preselection criteria are not met (e.g., salary too low)
	age, _ = encryptor.EncryptPu(int64(30))
	salary, _ = encryptor.EncryptPu(int64(20000)) // Below threshold for salary
	encryptedResult, err = encryptor.HomomorphicCreditEvaluation(age, salary, creditScore, dti)
	if err != nil {
		t.Fatalf("failed to perform homomorphic credit evaluation: %v", err)
	}

	// Decrypt and verify the result
	result, err = encryptor.Decrypt(encryptedResult)
	if err != nil {
		t.Fatalf("failed to decrypt the result: %v", err)
	}

	// Ineligible, so result should be -1
	if result != -1 {
		t.Errorf("unexpected evaluation result: got %v, expected -1", result)
	}
}

// TestHomomorphicSatisfyPreselection tests the homomorphic preselection satisfaction.
func TestHomomorphicSatisfyPreselection(t *testing.T) {
	encryptor := NewEncryptor()

	// Test data
	age, _ := encryptor.EncryptPu(int64(35))                  // Example age
	salary, _ := encryptor.EncryptPu(int64(10*1000*1000 + 1)) // Example salary

	// Test when preselection criteria are satisfied
	isEligible, err := encryptor.HomomorphicSatisfyPreselection(age, salary)
	if err != nil {
		t.Fatalf("failed to check preselection: %v", err)
	}

	if !isEligible {
		t.Errorf("expected eligibility, but got not eligible")
	}

	// Test when preselection criteria are not satisfied (low salary)
	salary, _ = encryptor.EncryptPu(int64(15000)) // Below threshold for salary
	isEligible, err = encryptor.HomomorphicSatisfyPreselection(age, salary)
	if err != nil {
		t.Fatalf("failed to check preselection: %v", err)
	}

	if isEligible {
		t.Errorf("expected not eligible, but got eligible")
	}
}

// TestHomomorphicCalcScore tests the calculation of the homomorphic score.
func TestHomomorphicCalcScore(t *testing.T) {
	encryptor := NewEncryptor()

	// Test data
	creditScore, _ := encryptor.EncryptPu(int64(700)) // Example credit score
	dti, _ := encryptor.EncryptPu(int64(25))          // Example DTI //todo 0.25

	// Test the homomorphic score calculation
	encryptedScore, err := encryptor.HomomorphicCalcScore(creditScore, dti)
	if err != nil {
		t.Fatalf("failed to calculate homomorphic score: %v", err)
	}

	// Decrypt and verify the result
	score, err := encryptor.Decrypt(encryptedScore)
	if err != nil {
		t.Fatalf("failed to decrypt score: %v", err)
	}

	// Assuming the score should be a positive value after calculation
	if score <= 0 {
		t.Errorf("unexpected score: got %v, expected positive", score)
	}
}
