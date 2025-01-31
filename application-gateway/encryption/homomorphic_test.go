package encryption

import (
	"math"
	"testing"
)

const precision = 1e-4

// TestEncryptDecryptPublic ensures that encryption with the public key and decryption works correctly
func TestEncryptDecryptPublic(t *testing.T) {
	helper := NewCKKSHelper()

	originalValue := 42.0
	ciphertext := helper.EncryptPu(originalValue)

	// Decrypt the ciphertext
	decryptedValue := helper.Decrypt(ciphertext)

	// Check the result
	if math.Abs(decryptedValue-originalValue) > precision {
		t.Errorf("Decrypt failed with public encryption. Got %f, expected %f", decryptedValue, originalValue)
	}
}

// TestEncryptDecryptPrivate ensures that encryption with the private key and decryption works correctly
func TestEncryptDecryptPrivate(t *testing.T) {
	helper := NewCKKSHelper()

	originalValue := 84.0
	ciphertext := helper.EncryptPr(originalValue)

	// Decrypt the ciphertext
	decryptedValue := helper.Decrypt(ciphertext)

	// Check the result
	if math.Abs(decryptedValue-originalValue) > precision {
		t.Errorf("Decrypt failed with private encryption. Got %f, expected %f", decryptedValue, originalValue)
	}
}

// TestDecryptWithKey ensures that decryption with a custom key works correctly
func TestDecryptWithKey(t *testing.T) {
	helper := NewCKKSHelper()

	originalValue := 123.45
	ciphertext := helper.EncryptPu(originalValue)

	// Decrypt the ciphertext using the default secret key
	decryptedValue := helper.DecryptWithKey(ciphertext, helper.secretKey)

	// Check the result
	if math.Abs(decryptedValue-originalValue) > precision {
		t.Errorf("DecryptWithKey failed. Got %f, expected %f", decryptedValue, originalValue)
	}
}

// TestEncryptDecryptWithHomomorphicOperation tests encryption, a simple operation, and decryption
func TestEncryptDecryptWithHomomorphicOperation(t *testing.T) {
	helper := NewCKKSHelper()

	originalValue := 12.34
	ciphertext := helper.EncryptPu(originalValue)

	// Perform a homomorphic operation (multiply by 2)
	helper.Evaluator.MultByConst(ciphertext, 2.0, ciphertext)

	// Decrypt the ciphertext
	decryptedValue := helper.Decrypt(ciphertext)

	// Check the result
	expectedValue := originalValue * 2
	if math.Abs(decryptedValue-expectedValue) > precision {
		t.Errorf("Decrypt failed after homomorphic operation. Got %f, expected %f", decryptedValue, expectedValue)
	}
}

// TestAdd tests homomorphic addition of two ciphertexts
func TestAdd(t *testing.T) {
	helper := NewCKKSHelper()

	value1 := 3.5
	value2 := 2.0

	// Encrypt the values
	ct1 := helper.EncryptPu(value1)
	ct2 := helper.EncryptPu(value2)

	// Perform homomorphic addition
	ctAdd := helper.Add(ct1, ct2)

	// Decrypt the result
	result := helper.Decrypt(ctAdd)

	// Check the result
	expected := value1 + value2
	if math.Abs(result-expected) > precision {
		t.Errorf("Add failed. Got %f, expected %f", result, expected)
	}
}

func TestAddWithPlain(t *testing.T) {
	helper := NewCKKSHelper()

	value1 := 3.5
	value2 := 2.0

	// Encrypt the first value
	ct1 := helper.EncryptPu(value1)

	// Perform homomorphic addition with a plaintext
	ctMulPlain := helper.AddWithPlain(ct1, value2)

	// Decrypt the result
	result := helper.Decrypt(ctMulPlain)

	// Check the result
	expected := value1 + value2
	if math.Abs(result-expected) > precision {
		t.Errorf("MultiplyPlain failed. Got %f, expected %f", result, expected)
	}
}

// TestMultiply tests homomorphic multiplication of two ciphertexts
func TestMultiply(t *testing.T) {
	helper := NewCKKSHelper()

	value1 := 3.5
	value2 := 2.0

	// Encrypt the values
	ct1 := helper.EncryptPu(value1)
	ct2 := helper.EncryptPu(value2)

	// Perform homomorphic multiplication
	ctMul := helper.Multiply(ct1, ct2)

	// Decrypt the result
	result := helper.Decrypt(ctMul)

	// Check the result
	expected := value1 * value2
	if math.Abs(result-expected) > precision {
		t.Errorf("Multiply failed. Got %f, expected %f", result, expected)
	}
}

// TestMultiplyPlain tests homomorphic multiplication of a ciphertext with a plaintext
func TestMultiplyPlain(t *testing.T) {
	helper := NewCKKSHelper()

	value1 := 3.5
	value2 := 2.0

	// Encrypt the first value
	ct1 := helper.EncryptPu(value1)

	// Perform homomorphic multiplication with a plaintext
	ctMulPlain := helper.MultiplyPlain(ct1, value2)

	// Decrypt the result
	result := helper.Decrypt(ctMulPlain)

	// Check the result
	expected := value1 * value2
	if math.Abs(result-expected) > precision {
		t.Errorf("MultiplyPlain failed. Got %f, expected %f", result, expected)
	}
}

// TestDivideByPlain tests homomorphic division of a ciphertext by a plaintext
func TestDivideByPlain(t *testing.T) {
	helper := NewCKKSHelper()

	value1 := 10.0
	value2 := 2.0

	// Encrypt the first value
	ct1 := helper.EncryptPu(value1)

	// Perform homomorphic division by a plaintext
	ctDivPlain := helper.DivideByPlain(ct1, value2)

	// Decrypt the result
	result := helper.Decrypt(ctDivPlain)

	// Check the result
	expected := value1 / value2
	if math.Abs(result-expected) > precision {
		t.Errorf("DivideByPlain failed. Got %f, expected %f", result, expected)
	}
}

// TestDivide tests homomorphic division of two ciphertexts
func TestDivide(t *testing.T) {
	helper := NewCKKSHelper()

	value1 := 10.0
	value2 := 2.0

	// Encrypt the values
	ct1 := helper.EncryptPu(value1)
	ct2 := helper.EncryptPu(value2)

	// Perform homomorphic division
	ctDiv := helper.Divide(ct1, ct2)

	// Decrypt the result
	result := helper.Decrypt(ctDiv)

	// Check the result
	expected := value1 / value2
	if math.Abs(result-expected) > precision {
		t.Errorf("Divide failed. Got %f, expected %f", result, expected)
	}
}

// //////////////////////////////////// CREDIT EVALUATION ////////////////////////////////////////////

// TestSigmoid tests the homomorphic sigmoid function
func TestSigmoid(t *testing.T) {
	helper := NewCKKSHelper()

	// Test input
	value := 2.5

	// Encrypt the input value
	ct := helper.EncryptPu(value)

	// Apply the homomorphic sigmoid function
	ctSigmoid := helper.sigmoid(helper, ct)

	// Decrypt the result
	result := helper.Decrypt(ctSigmoid)

	// Compute the expected result using the standard math package
	expected := 1.0 / (1.0 + math.Exp(-value))

	// Check if the result is within the acceptable precision
	if math.Abs(result-expected) > precision {
		t.Errorf("Sigmoid failed. Got %f, expected %f", result, expected)
	}
}
