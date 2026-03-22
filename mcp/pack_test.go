package mcp

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func Test_handler_packHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pwd, _ := os.Getwd()
	toPath := func(other string) string {
		return filepath.Join(pwd, other)
	}

	type serviceParams struct {
		expectSourcePath string
		expectOutputDir  string
		expectPrivateKey string
		expectOutputFile string
		expectName       string
	}

	tests := []struct {
		name    string
		handler func() (*handler, serviceParams)
		expect  func(*testing.T, serviceParams, *sdkmcp.CallToolResult)
		params  packParams
		wantErr bool
	}{
		{
			name: "should return error when source path is empty",
			handler: func() (*handler, serviceParams) {
				h := &handler{opts: &Options{WorkDir: "/"}}
				return h, serviceParams{}
			},
			params: packParams{
				Source: "",
			},
			wantErr: true,
		},
		{
			name: "should return error when source not found",
			handler: func() (*handler, serviceParams) {
				svc := NewMockcrx3service(ctrl)
				h := &handler{svc: svc, opts: &Options{WorkDir: "./"}}
				return h, serviceParams{}
			},
			params: packParams{
				Source: "nonexistent_extension",
			},
			wantErr: true,
		},
		{
			name: "should pack extension with output dir and name",
			handler: func() (*handler, serviceParams) {
				svc := NewMockcrx3service(ctrl)

				sourcePath := toPath("/testdata/workspace/unpacked_extension")
				outputDir := toPath("/testdata/workspace/packed")
				outputFile := filepath.Join(outputDir, "my_extension.crx")

				svc.EXPECT().PackTo(sourcePath, outputFile, nil).Return(nil)
				svc.EXPECT().GetID(outputFile).Return("mocked_extension_id", nil)

				sp := serviceParams{
					expectSourcePath: sourcePath,
					expectOutputDir:  outputDir,
					expectOutputFile: outputFile,
					expectName:       "my_extension",
				}

				h := &handler{svc: svc, opts: &Options{WorkDir: "testdata/workspace"}}

				return h, sp
			},
			params: packParams{
				Source:    "./unpacked_extension",
				OutputDir: "./packed",
				Name:      "my-extension",
			},
			expect: func(t *testing.T, sp serviceParams, res *sdkmcp.CallToolResult) {
				// structured content
				result, ok := res.StructuredContent.(packResult)
				assert.True(t, ok)
				assert.Equal(t, result.Filepath, sp.expectOutputFile)
				assert.Equal(t, result.ExtensionID, "mocked_extension_id")

				// content
				assertText(t, res, "Successfully packed extension")
			},
		},
		{
			name: "should pack extension without output dir and use source base name",
			handler: func() (*handler, serviceParams) {
				svc := NewMockcrx3service(ctrl)

				sourcePath := toPath("/testdata/workspace/unpacked_extension")
				outputDir := toPath("/testdata/workspace/extension.crx")

				svc.EXPECT().PackTo(sourcePath, outputDir, nil).Return(nil)
				svc.EXPECT().GetID(outputDir).Return("mocked_extension_id", nil)

				sp := serviceParams{
					expectSourcePath: sourcePath,
					expectOutputDir:  outputDir,
					expectOutputFile: outputDir,
					expectName:       "unpacked_extension",
				}

				h := &handler{svc: svc, opts: &Options{WorkDir: "testdata/workspace"}}

				return h, sp
			},
			params: packParams{
				Source: "./unpacked_extension",
			},
			expect: func(t *testing.T, sp serviceParams, res *sdkmcp.CallToolResult) {
				// structured content
				result, ok := res.StructuredContent.(packResult)
				assert.True(t, ok)
				assert.Equal(t, result.Filepath, sp.expectOutputFile)
				assert.Equal(t, result.ExtensionID, "mocked_extension_id")
				// content
				assertText(t, res, "Successfully packed extension")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, sp := tt.handler()
			got, _, gotErr := h.packHandler(context.Background(), nil, tt.params)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("packHandler() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("packHandler() succeeded unexpectedly")
			}
			tt.expect(t, sp, got)
		})
	}
}
