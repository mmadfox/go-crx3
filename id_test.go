package crx3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	filename := "./testdata/unpack/extension.crx"
	id, err := ID(filename)
	assert.Nil(t, err)
	assert.Equal(t, "dgmchnekcpklnjppdmmjlgpmpohmpmgp", id)
}
