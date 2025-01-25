package encryption

import (
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
}

func NewEncryptor() *Encryptor {
	params, err := bgv.NewParametersFromLiteral(bgvParams)
	if err != nil {
		panic(err)
	}

	encoder := bgv.NewEncoder(params)
	keyGen := rlwe.NewKeyGenerator(params)
	privateKey, publicKey := keyGen.GenKeyPairNew()
	decryptor := rlwe.NewDecryptor(params, privateKey)
	encryptorPu := rlwe.NewEncryptor(params, publicKey)
	encryptorPr := rlwe.NewEncryptor(params, privateKey)

	evaluator := bgv.NewEvaluator(params, nil)
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
		return nil, err
	}

	cipherText, err := e.encryptorPu.EncryptNew(plainText)
	if err != nil {
		return nil, err
	}

	return cipherText, nil
}

func (e *Encryptor) EncryptPr(data int64) (*rlwe.Ciphertext, error) {
	slice := []int64{data}
	plainText := bgv.NewPlaintext(e.params, e.params.MaxLevel())
	if err := e.encoder.Encode(slice, plainText); err != nil {
		return nil, err
	}

	cipherText, err := e.encryptorPr.EncryptNew(plainText)
	if err != nil {
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

func (e *Encryptor) DecryptWithKey(text string) (string, error) {
	panic("implement me")
}
