package crx3

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func readZip(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

func TestZip(t *testing.T) {
	type args struct {
		dirname string
	}
	tests := []struct {
		name    string
		args    args
		assert  func(dst *bytes.Buffer)
		wantErr error
	}{
		{
			name: "should return error when file path not exists",
			args: args{
				dirname: "/some/path/to/extension",
			},
			wantErr: ErrPathNotFound,
		},
		{
			name: "should not return error when all params are valid",
			args: args{
				dirname: "./testdata/extension",
			},
			assert: func(buf *bytes.Buffer) {
				require.NotZero(t, buf.Len())
				r, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
				require.NoError(t, err)
				require.NotNil(t, r)
				count := 0
				for i := 0; i < len(r.File); i++ {
					file := r.File[i]
					data, err := readZip(file)
					require.NoError(t, err)
					require.NotZero(t, data)
					got := fileExists(filepath.Join("./testdata/extension", file.Name))
					require.True(t, got)
					count++
				}
				require.Equal(t, 3, count)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := bytes.NewBuffer(nil)
			if err := Zip(dst, tt.args.dirname); !errors.Is(err, tt.wantErr) {
				t.Errorf("Zip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.assert != nil {
				tt.assert(dst)
			}
		})
	}
}

func TestZipTo(t *testing.T) {
	basePath, err := os.MkdirTemp("", "testzip")
	require.NoError(t, err)
	defer os.RemoveAll(basePath)

	type args struct {
		filename string
		dirname  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should return error when failed to create new file",
			args: args{
				filename: "",
				dirname:  "",
			},
			wantErr: true,
		},
		{
			name: "should return error when directory name does not exist",
			args: args{
				filename: filepath.Join(basePath, "my.zip"),
				dirname:  "",
			},
			wantErr: true,
		},
		{
			name: "should not return error when all params are valid",
			args: args{
				filename: filepath.Join(basePath, "my.zip"),
				dirname:  "./testdata/extension",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ZipTo(tt.args.filename, tt.args.dirname); (err != nil) != tt.wantErr {
				t.Errorf("ZipTo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
