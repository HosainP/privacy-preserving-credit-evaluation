package models

import (
	"credit-evaluation/application-gateway/encryption"
	"time"
)

type Document struct {
	ID      string `json:"ID"`
	OrgID   string `json:"OrgID"`
	OwnerID string `json:"OwnerID"`

	Title string            `json:"Title"`
	Time  time.Time         `json:"Time"`
	Data  map[string]string `json:"Data"`

	OrgSignature   string `json:"OrgSignature"`
	OwnerSignature string `json:"OwnerSignature"`
}

func (d *Document) Decrypt(privateKey string) (map[string]string, error) {
	var decryptedData map[string]string
	for key, value := range d.Data {
		decryptedData[key], _ = encryption.Decrypt([]byte(value), privateKey)
	}
	return decryptedData, nil
}
