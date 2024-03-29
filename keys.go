package crx3

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

// NewPrivateKey returns a new RSA private key with a bit size of 2048.
func NewPrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}

// SavePrivateKey saves the provided 'key' private key to the specified 'filename'.
// If 'key' is nil, it generates a new private key and saves it to the file.
func SavePrivateKey(filename string, key *rsa.PrivateKey) error {
	if key == nil {
		key, _ = NewPrivateKey()
	}
	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	bytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return err
	}
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: bytes,
	}
	_, err = fd.Write(pem.EncodeToMemory(block))
	return err
}

// LoadPrivateKey loads the RSA private key from the specified 'filename' into memory.
// It returns the loaded private key or an error if the key cannot be loaded.
func LoadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(buf)
	if block == nil {
		return nil, ErrPrivateKeyNotFound
	}
	r, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return r.(*rsa.PrivateKey), nil
}
