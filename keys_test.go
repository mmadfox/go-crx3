package crx3

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSavePrivateKey(t *testing.T) {
	filename := filepath.Join(os.TempDir(), "key.pem")
	err := SavePrivateKey(filename, nil)
	assert.Nil(t, err)
	assert.FileExists(t, filename)
	assert.Nil(t, os.Remove(filename))

	key, err := NewPrivateKey()
	assert.Nil(t, err)
	assert.Nil(t, key.Validate())
	err = SavePrivateKey(filename, key)
	assert.Nil(t, err)
	assert.FileExists(t, filename)
	assert.Nil(t, os.Remove(filename))
}

func TestSavePrivateKeyNegative(t *testing.T) {
	filename := "/path/to/path/key.pem"
	err := SavePrivateKey(filename, nil)
	assert.Error(t, err)
}

func TestLoadPrivateKey(t *testing.T) {
	filename := filepath.Join(os.TempDir(), "key.pem")
	err := SavePrivateKey(filename, nil)
	assert.Nil(t, err)
	assert.FileExists(t, filename)

	key, err := LoadPrivateKey(filename)
	assert.Nil(t, err)
	assert.Nil(t, key.Validate())
	assert.Nil(t, os.Remove(filename))
}

func TestLoadPrivateKeyNegative(t *testing.T) {
	key, err := LoadPrivateKey("/path/to/key.pem")
	assert.Error(t, err)
	assert.Nil(t, key)

	key, err = LoadPrivateKey("./testdata/pack/somefile.fg")
	assert.Error(t, err)
	assert.Nil(t, key)
}

func TestPrivateKeyToPEM(t *testing.T) {
	key, err := NewPrivateKey()
	assert.Nil(t, err)
	assert.Nil(t, key.Validate())
	pem := PrivateKeyToPEM(key)
	assert.NotNil(t, pem)
	assert.NotEmpty(t, pem)
}
