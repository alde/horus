package decryptor

// The Decryptor interface
type Decryptor interface {
	Decrypt(secret []byte) ([]byte, error)
}
