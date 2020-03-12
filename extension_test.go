package crx3

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtension_Pack(t *testing.T) {
	err := Extension("").Pack(nil)
	assert.Equal(t, ErrPathNotFound, err)
	err = Extension.Pack("/some/path", nil)
	assert.Error(t, err)
	err = Extension("").PackTo("", nil)
	assert.Equal(t, ErrPathNotFound, err)

	err = Extension("./testdata/pack/extension").Pack(nil)
	assert.Nil(t, err)
	assert.True(t, Extension("./testdata/pack/extension.crx").IsCRX3())
	assert.Nil(t, os.Remove("./testdata/pack/extension.crx"))
	assert.Nil(t, os.Remove("./testdata/pack/extension.crx.pem"))

	pk, err := NewPrivateKey()
	assert.Nil(t, err)
	customFilename := filepath.Join(os.TempDir(), "go.crx")
	err = Extension("./testdata/pack/extension").PackTo(customFilename, pk)
	assert.True(t, Extension(customFilename).IsCRX3())
	assert.Nil(t, os.Remove(customFilename))
}

func TestExtension_Unpack(t *testing.T) {
	err := Extension("").Unpack()
	assert.Equal(t, ErrPathNotFound, err)
	err = Extension("./testdata/unpack/extension.crx").Unpack()
	assert.Nil(t, err)
	assert.True(t, Extension("./testdata/unpack/extension").IsDir())
	assert.Nil(t, os.RemoveAll("./testdata/unpack/extension"))
}

func TestExtension_ToBase64(t *testing.T) {
	_, err := Extension("").ToBase64()
	assert.Equal(t, ErrPathNotFound, err)
	b, err := Extension("./testdata/unpack/extension.crx").ToBase64()
	assert.Nil(t, err)
	assert.NotEmpty(t, b)
}

func TestExtension_String(t *testing.T) {
	str := "./testdata/unpack/extension.crx"
	ext := Extension(str)
	assert.Equal(t, str, ext.String())
}

func TestExtension_Unzip(t *testing.T) {
	err := Extension("").Unzip()
	assert.Equal(t, ErrPathNotFound, err)
	err = Extension("./testdata/pack/extension.zip").Unzip()
	assert.Nil(t, err)
	assert.True(t, Extension("./testdata/pack/extension(1)").IsDir())
	assert.Nil(t, os.RemoveAll("./testdata/pack/extension(1)"))
}

func TestExtension_Zip(t *testing.T) {
	err := Extension("").Zip()
	assert.Equal(t, ErrPathNotFound, err)
	err = Extension("./testdata/unpack/extension.crx").Unpack()
	assert.Nil(t, err)
	assert.True(t, Extension("./testdata/unpack/extension").IsDir())
	err = Extension("./testdata/unpack/extension").Zip()
	assert.Nil(t, err)
	assert.True(t, Extension("./testdata/unpack/extension.zip").IsZip())
	assert.Nil(t, os.RemoveAll("./testdata/unpack/extension"))
	assert.Nil(t, os.RemoveAll("./testdata/unpack/extension.zip"))
}

func TestExtension_IsCRX3(t *testing.T) {
	ok := Extension("./testdata/unpack/extension.crx").IsCRX3()
	assert.True(t, ok)
	ok = Extension("./testdata/pack/extension.zip").IsCRX3()
	assert.False(t, ok)
}

func TestExtension_IsDir(t *testing.T) {
	ok := Extension("./testdata/pack/extension").IsDir()
	assert.True(t, ok)
	ok = Extension("./testdata/unpack/extension.crx").IsDir()
	assert.False(t, ok)
}

func TestExtension_IsZip(t *testing.T) {
	ok := Extension("./testdata/pack/extension.zip").IsZip()
	assert.True(t, ok)
	ok = Extension("./testdata/unpack/extension.crx").IsZip()
	assert.False(t, ok)
}
