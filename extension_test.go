package crx3

import (
	"crypto/rsa"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtension_ID(t *testing.T) {
	id, err := Extension("./testdata/dodyDol.crx").ID()
	assert.Nil(t, err)
	assert.Equal(t, "kpkcennohgffjdgaelocingbmkjnpjgc", id)
}

func TestExtension_IsEmpty(t *testing.T) {
	require.True(t, Extension("").IsEmpty())
	require.False(t, Extension("path/to/ext").IsEmpty())
}

func TestExtension_IsDir(t *testing.T) {
	ok := Extension("./testdata/extension").IsDir()
	assert.True(t, ok)
	ok = Extension("./testdata/dodyDol.crx").IsDir()
	assert.False(t, ok)
}

func TestExtension_IsZip(t *testing.T) {
	ok := Extension("./testdata/bobbyMol.zip").IsZip()
	assert.True(t, ok)
	ok = Extension("./testdata/dodyDol.crx").IsZip()
	assert.False(t, ok)
}

func TestExtension_IsCRX3(t *testing.T) {
	ok := Extension("./testdata/dodyDol.crx").IsCRX3()
	assert.True(t, ok)
	ok = Extension("./testdata/bobbyMol.zip").IsCRX3()
	assert.False(t, ok)
}

func TestExtension_Zip(t *testing.T) {
	basePath, err := os.MkdirTemp("", "ziptest")
	require.NoError(t, err)
	defer os.RemoveAll(basePath)

	src := "./testdata/bobbyMol.zip"
	dst := filepath.Join(basePath, "bobbyMol.zip")
	_, err = CopyFile(src, dst)
	require.NoError(t, err)
	require.NoError(t, UnzipTo(basePath, dst))
	os.Remove(dst)

	tests := []struct {
		name    string
		e       Extension
		assert  func()
		wantErr bool
	}{
		{
			name:    "should return error when extension is empty",
			e:       Extension(""),
			wantErr: true,
		},
		{
			name:    "should return error when dir does not exists",
			e:       Extension("path/not/exists"),
			wantErr: true,
		},
		{
			name: "should not return error when all params are valid",
			e:    Extension(filepath.Join(basePath, "bobbyMol")),
			assert: func() {
				expected := dst
				require.True(t, isZip(expected))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.Zip(); (err != nil) != tt.wantErr {
				t.Errorf("Extension.Zip() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assert != nil {
				tt.assert()
			}
		})
	}
}

func TestExtension_Unzip(t *testing.T) {
	basePath, err := os.MkdirTemp("", "unziptest")
	require.NoError(t, err)
	defer os.RemoveAll(basePath)

	src := "./testdata/bobbyMol.zip"
	dst := filepath.Join(basePath, "bobbyMol.zip")
	_, err = CopyFile(src, dst)
	require.NoError(t, err)

	tests := []struct {
		name    string
		e       Extension
		arrange func()
		assert  func()
		wantErr bool
	}{
		{
			name:    "should return error when extension is empty",
			e:       Extension(""),
			wantErr: true,
		},
		{
			name:    "should return error when zip does not exists",
			e:       Extension("file/not/found.zip"),
			wantErr: true,
		},
		{
			name: "should not return error when dir exists",
			e:    Extension(dst),
			arrange: func() {
				os.Mkdir(filepath.Join(basePath, "bobbyMol"), 0755)
			},
			assert: func() {
				expected := filepath.Join(basePath, "bobbyMol(1)")
				require.True(t, isDir(expected))
			},
		},
		{
			name: "should not return error when all params are valid",
			e:    Extension(dst),
			arrange: func() {
				os.RemoveAll(filepath.Join(basePath, "bobbyMol(1)"))
				os.RemoveAll(filepath.Join(basePath, "bobbyMol"))
			},
			assert: func() {
				expected := filepath.Join(basePath, "bobbyMol")
				require.True(t, isDir(expected))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.arrange != nil {
				tt.arrange()
			}
			if err := tt.e.Unzip(); (err != nil) != tt.wantErr {
				t.Errorf("Extension.Unzip() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assert != nil {
				tt.assert()
			}
		})
	}
}

func TestExtension_ToBase64(t *testing.T) {
	_, err := Extension("").Base64()
	assert.True(t, errors.Is(err, ErrPathNotFound))
	b, err := Extension("./testdata/dodyDol.crx").Base64()
	assert.Nil(t, err)
	assert.NotEmpty(t, b)
}

func TestExtension_Unpack(t *testing.T) {
	basePath, err := os.MkdirTemp("", "unpacktest")
	require.NoError(t, err)
	defer os.RemoveAll(basePath)

	src := "./testdata/dodyDol.crx"
	dst := filepath.Join(basePath, "dodyDol.crx")
	_, err = CopyFile(src, dst)
	require.NoError(t, err)

	tests := []struct {
		name    string
		e       Extension
		assert  func()
		wantErr bool
	}{
		{
			name:    "should return error when extension is empty",
			e:       Extension(""),
			wantErr: true,
		},
		{
			name: "should not return error when all params is valid",
			e:    Extension(dst),
			assert: func() {
				expected := filepath.Join(basePath, "dodyDol")
				require.True(t, isDir(expected))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.Unpack(); (err != nil) != tt.wantErr {
				t.Errorf("Extension.Unpack() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assert != nil {
				tt.assert()
			}
		})
	}
}

func TestExtension_PackTo(t *testing.T) {
	basePath, err := os.MkdirTemp("", "packtotest")
	require.NoError(t, err)
	defer os.RemoveAll(basePath)

	src := "./testdata/bobbyMol.zip"
	dst := filepath.Join(basePath, "bobbyMol.zip")
	_, err = CopyFile(src, dst)
	require.NoError(t, err)
	require.NoError(t, UnzipTo(basePath, dst))
	os.Remove(dst)

	type args struct {
		dst string
		pk  *rsa.PrivateKey
	}
	tests := []struct {
		name    string
		e       Extension
		assert  func()
		args    args
		wantErr bool
	}{
		{
			name:    "should return error when extension is empty",
			e:       Extension(""),
			wantErr: true,
		},
		{
			name: "should not return when all params are valid",
			e:    Extension(filepath.Join(basePath, "bobbyMol")),
			args: args{
				dst: filepath.Join(basePath, "bobbyMol.crx"),
			},
			assert: func() {
				require.True(t, fileExists(filepath.Join(basePath, "bobbyMol.crx")))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.PackTo(tt.args.dst, tt.args.pk); (err != nil) != tt.wantErr {
				t.Errorf("Extension.PackTo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assert != nil {
				tt.assert()
			}
		})
	}
}

func TestExtension_Pack(t *testing.T) {
	basePath, err := os.MkdirTemp("", "packtest")
	require.NoError(t, err)
	defer os.RemoveAll(basePath)

	src := "./testdata/bobbyMol.zip"
	dst := filepath.Join(basePath, "bobbyMol.zip")
	_, err = CopyFile(src, dst)
	require.NoError(t, err)
	require.NoError(t, UnzipTo(basePath, dst))
	os.Remove(dst)

	tests := []struct {
		name    string
		e       Extension
		assert  func()
		wantErr bool
	}{
		{
			name:    "should return error when extension is empty",
			e:       Extension(""),
			wantErr: true,
		},
		{
			name: "should not return when all params are valid",
			e:    Extension(filepath.Join(basePath, "bobbyMol")),
			assert: func() {
				require.True(t, fileExists(filepath.Join(basePath, "bobbyMol.crx")))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.Pack(nil); (err != nil) != tt.wantErr {
				t.Errorf("Extension.Pack() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assert != nil {
				tt.assert()
			}
		})
	}
}
