package crx3

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
)

// NewPrivateKey returns a new private key.
func NewPrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}

// SavePrivateKey saves private key to file.
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

// LoadPrivateKey loads the private key from a file into memory.
func LoadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	buf, err := ioutil.ReadFile(filename)
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
