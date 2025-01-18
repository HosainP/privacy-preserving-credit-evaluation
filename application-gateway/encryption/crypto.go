package encryption

func Encrypt(data []byte, publicKey string) ([]byte, error) {
	return []byte(data), nil // todo
}

func Decrypt(data []byte, privateKey string) (string, error) {
	return string(data), nil // todo
}
