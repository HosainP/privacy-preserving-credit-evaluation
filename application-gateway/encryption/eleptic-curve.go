package encryption

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

type Signer struct {
	privateKey ecdsa.PrivateKey
}

func NewSigner(privateKey *ecdsa.PrivateKey) *Signer {
	return &Signer{
		privateKey: *privateKey,
	}
}

func GenKey() (*ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func (s *Signer) Sign(data []byte) ([]byte, error) {
	return ecdsa.SignASN1(rand.Reader, &s.privateKey, data)
}

func (s *Signer) Verify(pub *ecdsa.PublicKey, hash, sig []byte) bool {
	return ecdsa.VerifyASN1(pub, hash, sig)
}
