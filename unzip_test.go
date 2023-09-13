package crx3

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnzipTo(t *testing.T) {
	basePath, err := os.MkdirTemp("", "unziptest")
	require.NoError(t, err)
	defer os.RemoveAll(basePath)

	type args struct {
		basepath string
		filename string
	}
	tests := []struct {
		name    string
		args    args
		assert  func(args args)
		wantErr bool
	}{
		{
			name: "should return error when zip file does not exists",
			args: args{
				basepath: basePath,
				filename: "/some/file/does/not/exists.zip",
			},
			wantErr: true,
		},
		{
			name: "should return error when basepath does not exists",
			args: args{
				basepath: "/some/base/path",
				filename: "./testdata/bobbyMol.zip",
			},
			wantErr: true,
		},
		{
			name: "should not return error when all params are valid",
			args: args{
				basepath: basePath,
				filename: "./testdata/bobbyMol.zip",
			},
			assert: func(args args) {
				expected := filepath.Join(args.basepath, "bobbyMol")
				require.True(t, isDir(expected))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UnzipTo(tt.args.basepath, tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("UnzipTo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assert != nil {
				tt.assert(tt.args)
			}
		})
	}
}
