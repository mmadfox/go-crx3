package crx3

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPack(t *testing.T) {
	pk, err := NewPrivateKey()
	assert.Nil(t, err)

	var have, want string

	// pack unpacked extension
	have = "./testdata/pack/extension"
	want = "./testdata/pack/extension.crx"
	err = Pack(have, pk)
	assert.Nil(t, err)
	assert.True(t, fileExists(want))
	assert.True(t, isCRC(want))
	err = os.Remove(want)
	assert.Nil(t, err)

	// pack zip extension
	have = "./testdata/pack/extension.zip"
	want = "./testdata/pack/extension.crx"
	err = Pack(have, pk)
	assert.Nil(t, err)
	assert.True(t, fileExists(want))
	assert.True(t, isCRC(want))
	err = os.Remove(want)
	assert.Nil(t, err)

	// pack without private key
	have = "./testdata/pack/extension.zip"
	want = "./testdata/pack/extension.crx"
	wantPem := "./testdata/pack/extension.pem"
	err = Pack(have, nil)
	assert.Nil(t, err)
	assert.True(t, fileExists(want))
	assert.True(t, isCRC(want))
	assert.True(t, fileExists(wantPem))
	err = os.Remove(want)
	assert.Nil(t, err)
	err = os.Remove(wantPem)

	// pack unsupported type
	have = "./testdata/pack/somefile.fg"
	err = Pack(have, nil)
	assert.Error(t, err)
	assert.Equal(t, ErrUnsupportedFileFormat, err)
}
