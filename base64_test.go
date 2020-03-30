package crx3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64(t *testing.T) {
	b, err := Base64("./testdata/unpack/extension.crx")
	assert.Nil(t, err)
	assert.NotEmpty(t, b)

	b, err = Base64("./testdata/pack/extension.zip")
	assert.Error(t, err)
	assert.Nil(t, b)

	extension, err := openCrxFile("./testdata/unpack/extension.crx")
	assert.Nil(t, err)
	err = extension.Close()
	assert.Nil(t, err)
	b, err = encodeExtensionToBase64Str(extension)
	assert.Error(t, err)
	assert.Nil(t, b)
}
