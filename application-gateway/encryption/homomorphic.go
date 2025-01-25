package encryption

import (
	"fmt"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/schemes/bgv"
)

var bgvParams = bgv.ParametersLiteral{
	LogN:             14,
	LogQ:             []int{56, 55, 55, 54, 54, 54},
	LogP:             []int{55, 55},
	PlaintextModulus: 0x3ee0001, // 1,000,000,007 in decimal (-+500,000,000)
}

type Encryptor struct {
	params      bgv.Parameters
	encoder     *bgv.Encoder
	decryptor   *rlwe.Decryptor
	encryptorPu *rlwe.Encryptor
	encryptorPr *rlwe.Encryptor
	evaluator   *bgv.Evaluator
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
