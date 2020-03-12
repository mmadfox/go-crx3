package crx3

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	pk, err := NewPrivateKey()
	assert.Nil(t, err)
	assert.NotZero(t, pk.Size())

	pem := filepath.Join(os.TempDir(), "key.pem")
	err = SavePrivateKey(pem, pk)
	assert.Nil(t, err)
	assert.True(t, fileExists(pem))

	pem2, err := LoadPrivateKey(pem)
	assert.Nil(t, err)
	assert.Equal(t, pem2.Size(), pk.Size())
	assert.Nil(t, os.Remove(pem))
}
