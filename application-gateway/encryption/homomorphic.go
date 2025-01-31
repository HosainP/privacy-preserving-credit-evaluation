package encryption

import (
	"fmt"
	"github.com/tuneinsight/lattigo/v4/ckks"
	"github.com/tuneinsight/lattigo/v4/rlwe"
)

// CKKSHelper encapsulates CKKS parameters, keys, and operations
type CKKSHelper struct {
	Params       ckks.Parameters
	Encoder      ckks.Encoder
	EncryptorPu  rlwe.Encryptor
	EncryptorPr  rlwe.Encryptor
	Decryptor    rlwe.Decryptor
	Evaluator    ckks.Evaluator
	Relinearizer *rlwe.RelinearizationKey
	LogSlots     int
	Scale        rlwe.Scale
	secretKey    *rlwe.SecretKey
}

// NewCKKSHelper initializes a CKKSHelper instance
func NewCKKSHelper() *CKKSHelper {
	// Create CKKS parameters (using default parameter set PN12QP109)
	params, err := ckks.NewParametersFromLiteral(ckks.PN14QP438)
	if err != nil {
		panic(err)
	}

	// Generate keys
	kgen := ckks.NewKeyGenerator(params)
	sk, pk := kgen.GenKeyPair()
	rlk := kgen.GenRelinearizationKey(sk, 1)

	// Return the helper instance
	return &CKKSHelper{
		Params:       params,
		Encoder:      ckks.NewEncoder(params),
		EncryptorPu:  ckks.NewEncryptor(params, pk),
		EncryptorPr:  ckks.NewEncryptor(params, sk),
		Decryptor:    ckks.NewDecryptor(params, sk),
		Evaluator:    ckks.NewEvaluator(params, rlwe.EvaluationKey{Rlk: rlk}),
		Relinearizer: rlk,
		LogSlots:     params.LogSlots(),
		Scale:        params.DefaultScale(),
		secretKey:    sk,
	}
}

// EncryptPu Encrypt encodes and encrypts a value into a ciphertext
func (c *CKKSHelper) EncryptPu(value float64) *rlwe.Ciphertext {
	// Encode the value into a plaintext
	values := []complex128{complex(value, 0)} // Single value in slot
	plaintext := c.Encoder.EncodeNew(values, c.Params.MaxLevel(), c.Scale, c.LogSlots)

	// Encrypt the plaintext
	return c.EncryptorPu.EncryptNew(plaintext)
}

// EncryptPr Encrypt encodes and encrypts a value into a ciphertext
func (c *CKKSHelper) EncryptPr(value float64) *rlwe.Ciphertext {
	// Encode the value into a plaintext
	values := []complex128{complex(value, 0)} // Single value in slot
	plaintext := c.Encoder.EncodeNew(values, c.Params.MaxLevel(), c.Scale, c.LogSlots)

	// Encrypt the plaintext
	return c.EncryptorPr.EncryptNew(plaintext)
}

// Decrypt decrypts a ciphertext using the default secret key and decodes it
func (c *CKKSHelper) Decrypt(ciphertext *rlwe.Ciphertext) float64 {
	// Decrypt the ciphertext into a plaintext
	plaintext := c.Decryptor.DecryptNew(ciphertext)

	// Decode the plaintext to retrieve the value
	decoded := c.Encoder.Decode(plaintext, c.LogSlots)

	// Return the real part of the first slot
	return real(decoded[0])
}

// DecryptWithKey decrypts a ciphertext using a provided secret key and decodes it
func (c *CKKSHelper) DecryptWithKey(ciphertext *rlwe.Ciphertext, secretKey *rlwe.SecretKey) float64 {
	// Create a decryptor with the provided secret key
	decryptor := ckks.NewDecryptor(c.Params, secretKey)

	// Decrypt the ciphertext into a plaintext
	plaintext := decryptor.DecryptNew(ciphertext)

	// Decode the plaintext to retrieve the value
	decoded := c.Encoder.Decode(plaintext, c.LogSlots)

	// Return the real part of the first slot
	return real(decoded[0])
}

// Add adds two ciphertexts and returns the result
func (c *CKKSHelper) Add(ct1, ct2 *rlwe.Ciphertext) *rlwe.Ciphertext {
	// Create a new ciphertext to store the result
	result := ckks.NewCiphertext(c.Params, 1, ct1.Level())

	// Perform the addition
	c.Evaluator.Add(ct1, ct2, result)

	return result
}

// AddWithPlain adds a ciphertext with a plaintext and returns the result
func (c *CKKSHelper) AddWithPlain(ct *rlwe.Ciphertext, value float64) *rlwe.Ciphertext {
	// Encode the value into a plaintext
	values := []complex128{complex(value, 0)} // Single value in slot
	plaintext := c.Encoder.EncodeNew(values, ct.Level(), c.Scale, c.LogSlots)

	// Create a new ciphertext to store the result
	result := ckks.NewCiphertext(c.Params, 1, ct.Level())

	// Perform the addition
	c.Evaluator.Add(ct, plaintext, result)

	// Rescale the result to maintain the correct scale
	c.Evaluator.Rescale(result, c.Scale, result)

	return result
}

// Multiply multiplies two ciphertexts and returns the result
func (c *CKKSHelper) Multiply(ct1, ct2 *rlwe.Ciphertext) *rlwe.Ciphertext {
	// Create a new ciphertext to store the result
	result := ckks.NewCiphertext(c.Params, 1, ct1.Level())

	// Perform the multiplication
	c.Evaluator.MulRelin(ct1, ct2, result)

	// Rescale the result to maintain the correct scale
	c.Evaluator.Rescale(result, c.Scale, result)

	return result
}

// MultiplyPlain multiplies a ciphertext by a plaintext and returns the result
func (c *CKKSHelper) MultiplyPlain(ct *rlwe.Ciphertext, value float64) *rlwe.Ciphertext {
	// Encode the value into a plaintext
	values := []complex128{complex(value, 0)} // Single value in slot
	plaintext := c.Encoder.EncodeNew(values, ct.Level(), c.Scale, c.LogSlots)

	// Create a new ciphertext to store the result
	result := ckks.NewCiphertext(c.Params, 1, ct.Level())

	// Perform the multiplication
	c.Evaluator.Mul(ct, plaintext, result)

	// Rescale the result to maintain the correct scale
	c.Evaluator.Rescale(result, c.Scale, result)

	return result
}

// DivideByPlain divides a ciphertext by a plaintext value
func (c *CKKSHelper) DivideByPlain(ct *rlwe.Ciphertext, value float64) *rlwe.Ciphertext {
	// Compute the reciprocal of the plaintext value
	reciprocal := 1.0 / value

	// Multiply the ciphertext by the reciprocal
	return c.MultiplyPlain(ct, reciprocal)
}

// Divide divides two ciphertexts (ct1 / ct2)
func (c *CKKSHelper) Divide(ct1, ct2 *rlwe.Ciphertext) *rlwe.Ciphertext {
	// Decrypt ct2 to get the divisor value
	divisor := c.Decrypt(ct2)

	// Compute the reciprocal of the divisor
	reciprocal := 1.0 / divisor

	// Multiply ct1 by the reciprocal
	return c.MultiplyPlain(ct1, reciprocal)
}

// //////////////////////////////////// CREDIT EVALUATION ////////////////////////////////////////////

const MinSalary = 10 * 1000 * 1000
const MinAge = 18

const MaxCreditScore = 850
const MinCreditScore = 300
const MinDTI = 0.01

const W0 = 0.5
const W1 = 0.5

// CreditEvaluation evaluates credit eligibility using encrypted inputs
func (c *CKKSHelper) CreditEvaluation(helper *CKKSHelper, ageCiphertext *rlwe.Ciphertext, salaryCiphertext *rlwe.Ciphertext, creditScoreCiphertext *rlwe.Ciphertext, dtiCiphertext *rlwe.Ciphertext) (*rlwe.Ciphertext, error) {
	// Check preselection (age and salary)
	//preselectionResult := satisfyPreselection(helper, ageCiphertext, salaryCiphertext)

	// If preselection fails, return an invalid result
	//if helper.Decrypt(preselectionResult) == 0 {
	//	return helper.EncryptPu(-1), nil
	//}

	// Calculate the credit score
	scoreCiphertext := c.calcScore(helper, creditScoreCiphertext, dtiCiphertext)

	// Apply sigmoid to the score
	//resultCiphertext := sigmoid(helper, scoreCiphertext)

	return scoreCiphertext, nil
}

// satisfyPreselection checks if age and salary meet the minimum requirements
//func (c *CKKSHelper) satisfyPreselection(helper *CKKSHelper, ageCiphertext *rlwe.Ciphertext, salaryCiphertext *rlwe.Ciphertext) *rlwe.Ciphertext {
//	// Encrypt the minimum age and salary
//	minAgeCiphertext := helper.EncryptPu(float64(MinAge))
//	minSalaryCiphertext := helper.EncryptPu(float64(MinSalary))
//
//	// Compare age and salary with minimums
//	ageCheck := helper.Evaluator.SubNew(ageCiphertext, minAgeCiphertext)
//	salaryCheck := helper.Evaluator.SubNew(salaryCiphertext, minSalaryCiphertext)
//
//	// Check if both conditions are met (age > MinAge AND salary > MinSalary)
//	ageCheck = helper.Evaluator.IsGreaterThan(ageCheck, helper.EncryptPu(0))
//	salaryCheck = helper.Evaluator.IsGreaterThan(salaryCheck, helper.EncryptPu(0))
//
//	// Combine results using multiplication (logical AND)
//	result := helper.Evaluator.MulRelinNew(ageCheck, salaryCheck)
//	helper.Evaluator.Rescale(result, helper.Scale, result)
//
//	return result
//}

// calcScore calculates the credit score using normalized credit score and DTI
func (c *CKKSHelper) calcScore(helper *CKKSHelper, creditScoreCiphertext *rlwe.Ciphertext, dtiCiphertext *rlwe.Ciphertext) *rlwe.Ciphertext {
	// Normalize credit score
	minCreditScoreCiphertext := helper.EncryptPu(MinCreditScore)
	maxCreditScoreCiphertext := helper.EncryptPu(MaxCreditScore)
	normalizedCreditScore := helper.Evaluator.SubNew(creditScoreCiphertext, minCreditScoreCiphertext)
	normalizedCreditScore = helper.Divide(normalizedCreditScore, helper.Evaluator.SubNew(maxCreditScoreCiphertext, minCreditScoreCiphertext))

	// Normalize DTI (ensure DTI is not too small)
	//minDTICiphertext := helper.EncryptPu(MinDTI)
	//dtiCiphertext = helper.Evaluator.MaxNew(dtiCiphertext, minDTICiphertext)

	// Calculate the score: (W0 * normalizedCreditScore) + (W1 * (1 / dti))
	w0Ciphertext := helper.EncryptPu(W0)
	w1Ciphertext := helper.EncryptPu(W1)

	term1 := helper.Evaluator.MulRelinNew(normalizedCreditScore, w0Ciphertext)
	fmt.Println("1", c.Decrypt(term1))
	helper.Evaluator.Rescale(term1, helper.Scale, term1)
	fmt.Println("2", c.Decrypt(term1))

	term2, _ := helper.Evaluator.InverseNew(dtiCiphertext, 5)
	term2 = helper.Evaluator.MulRelinNew(term2, w1Ciphertext)
	fmt.Println("3", c.Decrypt(term2))
	helper.Evaluator.Rescale(term2, helper.Scale, term2)
	fmt.Println("4", c.Decrypt(term2))

	scoreCiphertext := helper.Add(term1, term2)
	fmt.Println("5", c.Decrypt(scoreCiphertext))

	return scoreCiphertext
}

// sigmoid applies the sigmoid function to the input ciphertext
func (c *CKKSHelper) sigmoid(helper *CKKSHelper, xCiphertext *rlwe.Ciphertext) *rlwe.Ciphertext {
	// Compute sigmoid(x) = 1 / (1 + exp(-x))

	// Step 1: Negate x
	negXCiphertext := helper.Evaluator.NegNew(xCiphertext)

	// Step 2: Approximate exp(-x) using Taylor series:
	// exp(-x) ≈ 1 - x + x^2/2 - x^3/6 + x^4/24 - x^5/120

	x2 := helper.Multiply(negXCiphertext, negXCiphertext)
	fmt.Println(c.Decrypt(negXCiphertext))
	fmt.Println(c.Decrypt(x2))
	x3 := helper.Multiply(x2, negXCiphertext)
	fmt.Println(c.Decrypt(x3))
	//x4 := helper.Multiply(x3, negXCiphertext)
	//x5 := helper.Multiply(x4, negXCiphertext)

	term0 := helper.EncryptPu(1.0) // 1
	fmt.Println("term0", c.Decrypt(term0))

	term1 := helper.MultiplyPlain(negXCiphertext, 1.0) // -x
	fmt.Println("term1", c.Decrypt(term1))

	term2 := helper.MultiplyPlain(x2, 1.0/2.0) // x^2 / 2
	fmt.Println("term2", c.Decrypt(term2))

	term3 := helper.MultiplyPlain(x3, 1.0/6.0) // -x^3 / 6
	fmt.Println("term3", c.Decrypt(term3))

	//term4 := helper.MultiplyPlain(x4, 1.0/24.0) // x^4 / 24
	//fmt.Println("term4", c.Decrypt(term4))
	//
	//term5 := helper.MultiplyPlain(x5, -1.0/120.0) // -x^5 / 120
	//fmt.Println("term5", c.Decrypt(term5))

	// Sum up the terms for exp(-x)
	expNegXCiphertext := helper.Add(term0, term1)
	fmt.Println(c.Decrypt(expNegXCiphertext))
	expNegXCiphertext = helper.Add(expNegXCiphertext, term2)
	fmt.Println(c.Decrypt(expNegXCiphertext))
	expNegXCiphertext = helper.Add(expNegXCiphertext, term3)
	fmt.Println(c.Decrypt(expNegXCiphertext))
	//expNegXCiphertext = helper.Evaluator.AddNew(expNegXCiphertext, term4)
	//fmt.Println(c.Decrypt(expNegXCiphertext))
	//expNegXCiphertext = helper.Evaluator.AddNew(expNegXCiphertext, term5)
	//fmt.Println(c.Decrypt(expNegXCiphertext))

	// Step 3: Compute 1 + exp(-x)
	result := helper.AddWithPlain(expNegXCiphertext, 1)
	fmt.Println("result", c.Decrypt(result))
	// Step 4: Inverse of the denominator: 1 / (1 + exp(-x))
	//result, _ := helper.Evaluator.InverseNew(denominator, 5) // use 5 as the precision of the approximation

	return result
}
