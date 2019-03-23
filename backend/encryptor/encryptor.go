package encryptor

// The Encryptor interface
type Encryptor interface {
	Encrypt(bytes []byte) ([]byte, error)
}
