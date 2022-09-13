package crx3

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	filename := "./testdata/unpack/extension.crx"
	id, err := ID(filename)
	assert.Nil(t, err)
	assert.Equal(t, "dgmchnekcpklnjppdmmjlgpmpohmpmgp", id)
}

func TestIDNegative(t *testing.T) {
	filename := filepath.Join(os.TempDir(), "extension.id.crx")
	buf := new(bytes.Buffer)
	buf.WriteString("Cr24")
	// version 4
	_ = binary.Write(buf, binary.LittleEndian, uint32(4))
	err := os.WriteFile(filename, buf.Bytes(), os.ModePerm)
	assert.Nil(t, err)
	id, err := ID(filename)
	assert.Error(t, err)
	assert.Empty(t, id)
	assert.Nil(t, os.Remove(filename))
}

func TestIDNegative_UnmarshalHeader(t *testing.T) {
	filename := filepath.Join(os.TempDir(), "extension.id.crx")
	buf := new(bytes.Buffer)
	buf.WriteString("Cr24")
	_ = binary.Write(buf, binary.LittleEndian, uint32(3))
	_ = binary.Write(buf, binary.LittleEndian, uint32(256))
	tmp := make([]byte, 512)
	for i := 0; i < 512; i++ {
		tmp[i] = byte(1)
	}
	buf.Write(tmp)
	err := os.WriteFile(filename, buf.Bytes(), os.ModePerm)
	assert.Nil(t, err)
	id, err := ID(filename)
	assert.Error(t, err)
	assert.Empty(t, id)
	assert.Nil(t, os.Remove(filename))
}

func TestIDNegative_UnmarshalSignedData(t *testing.T) {
	filename := filepath.Join(os.TempDir(), "extension.id.crx")
	buf := new(bytes.Buffer)
	buf.WriteString("Cr24")
	_ = binary.Write(buf, binary.LittleEndian, uint32(3))
	mockdata := []byte(`some data section`)
	header, err := makeHeader(mockdata, mockdata, mockdata)
	assert.Nil(t, err)
	_ = binary.Write(buf, binary.LittleEndian, uint32(len(header)))
	buf.Write(header)
	err = os.WriteFile(filename, buf.Bytes(), os.ModePerm)
	assert.Nil(t, err)
	id, err := ID(filename)
	assert.Error(t, err)
	assert.Empty(t, id)
	assert.Nil(t, os.Remove(filename))
}
