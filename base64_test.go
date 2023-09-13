package crx3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64(t *testing.T) {
	b, err := Base64("./testdata/dodyDol.crx")
	assert.Nil(t, err)
	assert.NotEmpty(t, b)

	b, err = Base64("./testdata/bobbyMol.zip")
	assert.Error(t, err)
	assert.Nil(t, b)

	extension, err := openCrxFile("./testdata/dodyDol.crx")
	assert.Nil(t, err)
	err = extension.Close()
	assert.Nil(t, err)
	b, err = encodeExtensionToBase64Str(extension)
	assert.Error(t, err)
	assert.Nil(t, b)
}
