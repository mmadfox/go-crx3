package crx3

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

var defaultKeySize = 4096

// SetDefaultKeySize sets the global default key size for RSA key generation.
// It accepts key sizes of 2048, 3072, or 4096 bits. If a size outside these
// specified options is passed, the function panics to prevent the use of
// an unsupported key size, which could lead to security vulnerabilities.
// This strict enforcement helps ensure that only strong, widely accepted
// key sizes are used throughout the application.
//
// Usage:
//
//	SetDefaultKeySize(2048)  // sets the default key size to 2048 bits
//	SetDefaultKeySize(4096)  // sets the default key size to 4096 bits
//
// Valid key sizes are:
//   - 2048 bits
//   - 3072 bits
//   - 4096 bits
//
// Panics:
//
//	This function will panic if any key size other than the above mentioned
//	valid sizes is attempted to be set. This is a deliberate design choice
//	to catch incorrect key size settings during the development phase.
func SetDefaultKeySize(size int) {
	switch size {
	case 2048, 3072, 4096:
		defaultKeySize = size
	default:
		panic("invalid key size: only 2048, 3072, or 4096 bits are allowed")
	}
}

// NewPrivateKey returns a new RSA private key with a bit size of 4096.
func NewPrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, defaultKeySize)
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
