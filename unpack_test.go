package crx3

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	basePath, err := os.MkdirTemp("", "unpacktest")
	require.NoError(t, err)
	defer os.RemoveAll(basePath)

	dst := filepath.Join(basePath, "dodyDol.crx")
	src := "./testdata/dodyDol.crx"
	_, err = CopyFile(src, dst)
	require.NoError(t, err)

	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		assert  func()
		wantErr bool
	}{
		{
			name:    "should return error when filename is empty",
			wantErr: true,
		},
		{
			name: "should return error when file is not in crx format",
			args: args{
				filename: "./testdata/bobbyMol.zip",
			},
			wantErr: true,
		},
		{
			name: "should not return error when all params are valid",
			args: args{
				filename: dst,
			},
			assert: func() {
				expected := filepath.Join(basePath, "dodyDol")
				require.True(t, isDir(expected))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Unpack(tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("Unpack() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assert != nil {
				tt.assert()
			}
		})
	}
}

func TestUnpackTo(t *testing.T) {
	basePath, err := os.MkdirTemp("", "unpacktotest")
	require.NoError(t, err)
	defer os.RemoveAll(basePath)

	type args struct {
		filename string
		dirname  string
	}
	tests := []struct {
		name    string
		args    args
		assert  func()
		wantErr bool
	}{
		{
			name: "should return error when dir does not exists and create failed",
			args: args{
				filename: "./testdata/dodyDol.crx",
				dirname:  "/not/not/not",
			},
			wantErr: true,
		},
		{
			name: "should not return error when all params are valid",
			args: args{
				filename: "./testdata/dodyDol.crx",
				dirname:  basePath,
			},
			assert: func() {
				expected := filepath.Join(basePath, "dodyDol")
				require.True(t, isDir(expected))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UnpackTo(tt.args.filename, tt.args.dirname); (err != nil) != tt.wantErr {
				t.Errorf("UnpackTo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assert != nil {
				tt.assert()
			}
		})
	}
}
