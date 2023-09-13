package crx3

import (
	"crypto/rsa"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPack(t *testing.T) {
	basePath, err := os.MkdirTemp("", "packtest")
	require.NoError(t, err)
	defer os.RemoveAll(basePath)

	type args struct {
		src string
		dst string
		pk  *rsa.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		assert  func()
		wantErr bool
	}{
		{
			name: "should return error when src path is empty",
			args: args{
				src: "",
				dst: "/path",
			},
			wantErr: true,
		},
		{
			name: "should return error when dst path is empty",
			args: args{
				src: "/path/to",
				dst: "",
			},
			wantErr: true,
		},
		{
			name: "should return error when file does not crx suffix",
			args: args{
				src: "./testdata/extension",
				dst: "somefile.png",
			},
			wantErr: true,
		},
		{
			name: "should return error when src does not exists",
			args: args{
				src: "path/not/exists",
				dst: "gobyMo.crx",
			},
			wantErr: true,
		},
		{
			name: "should not return when src not zipped",
			args: args{
				src: "./testdata/extension",
				dst: filepath.Join(basePath, "my.crx"),
			},
			assert: func() {
				expectedCrx := filepath.Join(basePath, "my.crx")
				expectedPem := filepath.Join(basePath, "my.crx.pem")
				require.True(t, fileExists(expectedCrx))
				require.True(t, fileExists(expectedPem))
			},
		},
		{
			name: "should not return when src zipped",
			args: args{
				src: "./testdata/bobbyMol.zip",
				dst: filepath.Join(basePath, "my.crx"),
			},
			assert: func() {
				expectedCrx := filepath.Join(basePath, "my.crx")
				expectedPem := filepath.Join(basePath, "my.crx.pem")
				require.True(t, fileExists(expectedCrx))
				require.True(t, fileExists(expectedPem))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Pack(tt.args.src, tt.args.dst, tt.args.pk); (err != nil) != tt.wantErr {
				t.Errorf("Pack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
