package crx3

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	filename := "./testdata/dodyDol.crx"
	id, err := ID(filename)
	assert.Nil(t, err)
	assert.Equal(t, "kpkcennohgffjdgaelocingbmkjnpjgc", id)
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

func TestIDFromPubKey(t *testing.T) {
	var pubKey1 = []byte(`MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAj/u/XDdjlDyw7gHEtaaasZ9GdG8WOKAyJzXd8HFrDtz2Jcuy7er7MtWvHgNDA0bwpznbI5YdZeV4UfCEsA4SrA5b3MnWTHwA1bgbiDM+L9rrqvcadcKuOlTeN48Q0ijmhHlNFbTzvT9W0zw/GKv8LgXAHggxtmHQ/Z9PP2QNF5O8rUHHSL4AJ6hNcEKSBVSmbbjeVm4gSXDuED5r0nwxvRtupDxGYp8IZpP5KlExqNu1nbkPc+igCTIB6XsqijagzxewUHCdovmkb2JNtskx/PMIEv+TvWIx2BzqGp71gSh/dV7SJ3rClvWd2xj8dtxG8FfAWDTIIi0qZXWn2QhizQIDAQAB`)
	var pubKey2 = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAj/u/XDdjlDyw7gHEtaaa
sZ9GdG8WOKAyJzXd8HFrDtz2Jcuy7er7MtWvHgNDA0bwpznbI5YdZeV4UfCEsA4S
rA5b3MnWTHwA1bgbiDM+L9rrqvcadcKuOlTeN48Q0ijmhHlNFbTzvT9W0zw/GKv8
LgXAHggxtmHQ/Z9PP2QNF5O8rUHHSL4AJ6hNcEKSBVSmbbjeVm4gSXDuED5r0nwx
vRtupDxGYp8IZpP5KlExqNu1nbkPc+igCTIB6XsqijagzxewUHCdovmkb2JNtskx
/PMIEv+TvWIx2BzqGp71gSh/dV7SJ3rClvWd2xj8dtxG8FfAWDTIIi0qZXWn2Qhi
zQIDAQAB
-----END PUBLIC KEY-----
`)
	expectedExtensionID := "lfoeajgcchlidpicbabpmckkejpckcfb"
	type args struct {
		pubKey []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "should return error when public key is empty",
			args:    args{pubKey: []byte{}},
			wantErr: true,
		},
		{
			name:    "should reutrn error when public key is invalid",
			args:    args{pubKey: []byte(strings.Repeat("a", 128))},
			wantErr: true,
		},
		{
			name:    "should return extension ID from public key",
			args:    args{pubKey: pubKey1},
			want:    expectedExtensionID,
			wantErr: false,
		},
		{
			name:    "should return error when base64 decoding fails",
			args:    args{pubKey: []byte(strings.Repeat("~", 128))},
			wantErr: true,
		},
		{
			name:    "should return extension ID from formatted public key",
			args:    args{pubKey: pubKey2},
			want:    expectedExtensionID,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IDFromPubKey(tt.args.pubKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("IDFromPubKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IDFromPubKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
