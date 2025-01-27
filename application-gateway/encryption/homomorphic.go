package encryption

import (
	credit_evaluation "credit-evaluation/application-gateway/credit-evaluation"
	"fmt"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/schemes/bgv"
)

var bgvParams = bgv.ParametersLiteral{
	LogN:             14,
	LogQ:             []int{56, 55, 55, 54, 54, 54},
	LogP:             []int{55, 55},
	PlaintextModulus: 0x3ee0001,
}

type Encryptor struct {
	params      bgv.Parameters
	encoder     *bgv.Encoder
	decryptor   *rlwe.Decryptor
	encryptorPu *rlwe.Encryptor
	encryptorPr *rlwe.Encryptor
	evaluator   *bgv.Evaluator
	privateKey  *rlwe.SecretKey
}

func NewEncryptor() *Encryptor {
	params, err := bgv.NewParametersFromLiteral(bgvParams)
	if err != nil {
		panic(err)
	}

	encoder := bgv.NewEncoder(params)
	keyGen := rlwe.NewKeyGenerator(params)
	privateKey, publicKey := keyGen.GenKeyPairNew()

	// relinearization key for after multiplication
	rlk := keyGen.GenRelinearizationKeyNew(privateKey)
	evkSet := rlwe.NewMemEvaluationKeySet(rlk)
	evkSet.RelinearizationKey = rlk

	decryptor := rlwe.NewDecryptor(params, privateKey)
	encryptorPu := rlwe.NewEncryptor(params, publicKey)
	encryptorPr := rlwe.NewEncryptor(params, privateKey)

	evaluator := bgv.NewEvaluator(params, evkSet)
	return &Encryptor{
		params:      params,
		encoder:     encoder,
		decryptor:   decryptor,
		encryptorPu: encryptorPu,
		encryptorPr: encryptorPr,
		evaluator:   evaluator,
		privateKey:  privateKey,
	}
}

func (e *Encryptor) EncryptPu(data int64) (*rlwe.Ciphertext, error) {
	slice := []int64{data}
	plainText := bgv.NewPlaintext(e.params, e.params.MaxLevel())
	if err := e.encoder.Encode(slice, plainText); err != nil {
		fmt.Printf("failed to encode data: %v", err)
		return nil, err
	}

	cipherText, err := e.encryptorPu.EncryptNew(plainText)
	if err != nil {
		fmt.Printf("failed to encrypt data: %v", err)
		return nil, err
	}

	return cipherText, nil
}

func (e *Encryptor) EncryptPr(data int64) (*rlwe.Ciphertext, error) {
	slice := []int64{data}
	plainText := bgv.NewPlaintext(e.params, e.params.MaxLevel())
	if err := e.encoder.Encode(slice, plainText); err != nil {
		fmt.Printf("failed to encode data: %v", err)
		return nil, err
	}

	cipherText, err := e.encryptorPr.EncryptNew(plainText)
	if err != nil {
		fmt.Printf("failed to encrypt data: %v", err)
		return nil, err
	}

	return cipherText, nil
}

func (e *Encryptor) Decrypt(cipherText *rlwe.Ciphertext) (int64, error) {
	plainText := e.decryptor.DecryptNew(cipherText)

	result := make([]int64, e.params.N())
	err := e.encoder.Decode(plainText, result)
	if err != nil {
		return -1, err
	}

	return result[0], nil
}

func (e *Encryptor) DecryptWithKey(cipherText *rlwe.Ciphertext, privateKey *rlwe.SecretKey) (int64, error) {
	decryptor := rlwe.NewDecryptor(e.params, privateKey)
	plainText := decryptor.DecryptNew(cipherText)

	result := make([]int64, e.params.N())
	err := e.encoder.Decode(plainText, result)
	if err != nil {
		return -1, err
	}

	return result[0], nil
}

// Add performs homomorphic addition on two ciphertexts.
func (e *Encryptor) Add(cipherText1, cipherText2 *rlwe.Ciphertext) (*rlwe.Ciphertext, error) {
	result, err := e.evaluator.AddNew(cipherText1, cipherText2)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Multiply performs homomorphic multiplication on two ciphertexts.
func (e *Encryptor) Multiply(cipherText1, cipherText2 *rlwe.Ciphertext) (*rlwe.Ciphertext, error) {
	result, err := e.evaluator.MulNew(cipherText1, cipherText2)
	if err != nil {
		return nil, err
	}
	err = e.evaluator.Relinearize(result, result)
	if err != nil {
		return nil, err
	} // Relinearize to reduce ciphertext size

	err = e.evaluator.Rescale(result, result)
	if err != nil {
		return nil, err
	} // Rescale to maintain precision

	return result, nil
}

// MultiplyWithPlain multiplies a ciphertext with a plaintext number.
func (e *Encryptor) MultiplyWithPlain(cipherText *rlwe.Ciphertext, plainNumber int64) (*rlwe.Ciphertext, error) {
	// Encode the plaintext number into a BGV plaintext
	plainText := bgv.NewPlaintext(e.params, e.params.MaxLevel())
	slice := []int64{plainNumber}
	if err := e.encoder.Encode(slice, plainText); err != nil {
		fmt.Printf("failed to encode plaintext number: %v\n", err)
		return nil, err
	}

	// Perform homomorphic multiplication
	result, err := e.evaluator.MulNew(cipherText, plainText)
	if err != nil {
		fmt.Printf("failed to multiply ciphertext with plaintext: %v\n", err)
		return nil, err
	}

	// Rescale to maintain precision (if needed by your parameters)
	err = e.evaluator.Rescale(result, result)
	if err != nil {
		fmt.Printf("failed to rescale result: %v\n", err)
		return nil, err
	}

	return result, nil
}

// DivideByPlain performs homomorphic division of a ciphertext by a plaintext number.
func (e *Encryptor) DivideByPlain(cipherText *rlwe.Ciphertext, plainNumber int64) (*rlwe.Ciphertext, error) {
	// Encode the plaintext number into a BGV plaintext
	plainText := bgv.NewPlaintext(e.params, e.params.MaxLevel())
	slice := []int64{plainNumber}
	if err := e.encoder.Encode(slice, plainText); err != nil {
		fmt.Printf("failed to encode plaintext number: %v\n", err)
		return nil, err
	}

	// Compute the multiplicative inverse of the plaintext number
	// This is required because homomorphic encryption typically supports multiplication and not division.
	// In the case of plaintexts, we multiply by the inverse of the plaintext number.
	// For example, dividing by 2 is equivalent to multiplying by 1/2.

	// Let's assume plainNumber is non-zero. The inverse of plainNumber in homomorphic encryption is calculated here.
	inversePlainText := bgv.NewPlaintext(e.params, e.params.MaxLevel())
	e.encoder.Encode([]int64{1 / plainNumber}, inversePlainText)

	// Perform homomorphic multiplication with the inverse of the plaintext number.
	result, err := e.evaluator.MulNew(cipherText, inversePlainText)
	if err != nil {
		fmt.Printf("failed to divide ciphertext by plaintext: %v\n", err)
		return nil, err
	}

	// Rescale to maintain precision (if needed by your parameters)
	err = e.evaluator.Rescale(result, result)
	if err != nil {
		fmt.Printf("failed to rescale result: %v\n", err)
		return nil, err
	}

	return result, nil
}

// Credit evaluation functions

// HomomorphicCreditEvaluation performs the credit evaluation using homomorphic encryption.
func (e *Encryptor) HomomorphicCreditEvaluation(age, salary, creditScore, dti *rlwe.Ciphertext) (*rlwe.Ciphertext, error) {
	// Step 1: Check if the preselection criteria are satisfied (age > MinAge and salary > MinSalary)
	isEligible, err := e.HomomorphicSatisfyPreselection(age, salary)
	if err != nil {
		return nil, fmt.Errorf("failed to check preselection: %v", err)
	}

	// If not eligible, return an encrypted -1
	if !isEligible {
		// Encrypt -1 to return as the result
		return e.EncryptPu(int64(-1))
	}

	// Step 2: Calculate the homomorphic score if preselection is passed
	encryptedScore, err := e.HomomorphicCalcScore(creditScore, dti)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate homomorphic score: %v", err)
	}

	// Step 3: Return the encrypted score
	return encryptedScore, nil
}

// HomomorphicSatisfyPreselection checks if the age and salary satisfy the preselection conditions.
func (e *Encryptor) HomomorphicSatisfyPreselection(age, salary *rlwe.Ciphertext) (bool, error) {
	// Define constants for the preselection checks
	MinAge := credit_evaluation.MinAge
	MinSalary := credit_evaluation.MinSalary

	// Convert MinAge and MinSalary to plaintext
	minAgePlain := bgv.NewPlaintext(e.params, e.params.MaxLevel())
	e.encoder.Encode([]int64{int64(MinAge)}, minAgePlain)

	minSalaryPlain := bgv.NewPlaintext(e.params, e.params.MaxLevel())
	e.encoder.Encode([]int64{int64(MinSalary)}, minSalaryPlain)

	// Compare salary with MinSalary: salary > MinSalary
	salaryDiff, err := e.evaluator.SubNew(salary, minSalaryPlain)
	if err != nil {
		return false, fmt.Errorf("failed to subtract MinSalary from salary: %v", err)
	}

	// Compare age with MinAge: age > MinAge
	ageDiff, err := e.evaluator.SubNew(age, minAgePlain)
	if err != nil {
		return false, fmt.Errorf("failed to subtract MinAge from age: %v", err)
	}

	// Now we check if the results are greater than zero using homomorphic operations.
	// This step depends on how your homomorphic encryption library handles comparisons.
	// A common approach is to check whether the resulting ciphertext is non-negative.

	// In a homomorphic setting, we would ideally need a function to "extract" or "check" if the result is positive.
	// We simulate this by simply comparing the decrypted result of `ageDiff` and `salaryDiff` to see if both are positive.

	// Decrypting and checking (for testing purposes)
	ageResult, err := e.Decrypt(ageDiff)
	if err != nil {
		return false, fmt.Errorf("failed to decrypt ageDiff: %v", err)
	}

	salaryResult, err := e.Decrypt(salaryDiff)
	if err != nil {
		return false, fmt.Errorf("failed to decrypt salaryDiff: %v", err)
	}

	// Check if both conditions are satisfied
	if ageResult > 0 && salaryResult > 0 {
		return true, nil
	}

	return false, nil
}

// HomomorphicCalcScore calculates the score based on the encrypted credit score and DTI.
func (e *Encryptor) HomomorphicCalcScore(creditScore, dti *rlwe.Ciphertext) (*rlwe.Ciphertext, error) {
	panic("implement me")
}

// HomomorphicSigmoid approximates the sigmoid function homomorphically using a cubic polynomial: a_0 + a_1*x + a_2*x^3.
func (e *Encryptor) HomomorphicSigmoid(x *rlwe.Ciphertext) (*rlwe.Ciphertext, error) {
	// Coefficients for the cubic polynomial approximation of the sigmoid
	a0 := int64(0)      // Constant term
	a1 := int64(1)      // Linear term coefficient
	a2 := int64(-1)     // Cubic term coefficient (inverted for simplicity)
	scale := int64(100) // Scaling factor for better precision

	// Encode coefficients as plaintexts
	a0Plain := bgv.NewPlaintext(e.params, e.params.MaxLevel())
	a1Plain := bgv.NewPlaintext(e.params, e.params.MaxLevel())
	a2Plain := bgv.NewPlaintext(e.params, e.params.MaxLevel())

	e.encoder.Encode([]int64{a0 * scale}, a0Plain)
	e.encoder.Encode([]int64{a1 * scale}, a1Plain)
	e.encoder.Encode([]int64{a2 * scale}, a2Plain)

	// Scale input ciphertext
	scaledX, err := e.MultiplyWithPlain(x, scale)
	if err != nil {
		return nil, fmt.Errorf("failed to scale input: %v", err)
	}

	// Compute x^3
	x2, err := e.Multiply(scaledX, scaledX) // x^2
	if err != nil {
		return nil, fmt.Errorf("failed to compute x^2: %v", err)
	}

	x3, err := e.Multiply(x2, scaledX) // x^3
	if err != nil {
		return nil, fmt.Errorf("failed to compute x^3: %v", err)
	}

	// Compute a2 * x^3
	a2x3, err := e.evaluator.MulNew(x3, a2Plain)
	if err != nil {
		return nil, fmt.Errorf("failed to compute a2 * x^3: %v", err)
	}

	// Compute a1 * x
	a1x, err := e.evaluator.MulNew(scaledX, a1Plain)
	if err != nil {
		return nil, fmt.Errorf("failed to compute a1 * x: %v", err)
	}

	// Add a2 * x^3 + a1 * x
	polyResult, err := e.evaluator.AddNew(a1x, a2x3)
	if err != nil {
		return nil, fmt.Errorf("failed to add a1 * x and a2 * x^3: %v", err)
	}

	// Add a0 to the result
	finalResult, err := e.evaluator.AddNew(polyResult, a0Plain)
	if err != nil {
		return nil, fmt.Errorf("failed to add a0: %v", err)
	}

	// Rescale to match the original scale (optional based on parameters)
	err = e.evaluator.Rescale(finalResult, finalResult)
	if err != nil {
		return nil, fmt.Errorf("failed to rescale final result: %v", err)
	}

	return finalResult, nil
}
