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

func Test_handler_unzipHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pwd, _ := os.Getwd()
	toPath := func(other string) string {
		return filepath.Join(pwd, other)
	}

	type serviceParams struct {
		expectExtensionFilename string
		expectOutputDir         string
	}

	tests := []struct {
		name    string
		handler func() (*handler, serviceParams)
		expect  func(*testing.T, serviceParams, *sdkmcp.CallToolResult)
		params  unzipParams
		wantErr bool
	}{
		{
			name: "should return error when extension filepath is invalid",
			handler: func() (*handler, serviceParams) {
				h := &handler{opts: &Options{WorkDir: "/"}}
				return h, serviceParams{}
			},
			params: unzipParams{
				Filepath: ".../...",
			},
			wantErr: true,
		},
		{
			name: "should return error when extension not found",
			handler: func() (*handler, serviceParams) {
				h := &handler{opts: &Options{WorkDir: "/testdata/workspace"}}
				return h, serviceParams{}
			},
			params: unzipParams{
				Filepath: "some_extension.zip",
			},
			wantErr: true,
		},
		{
			name: "should unzip extension with output dir",
			handler: func() (*handler, serviceParams) {
				svc := NewMockcrx3service(ctrl)

				sp := serviceParams{
					expectExtensionFilename: toPath("./testdata/workspace/extension.zip"),
					expectOutputDir:         toPath("./testdata/workspace/somepath/my-extension"),
				}

				svc.EXPECT().UnzipTo(sp.expectExtensionFilename, sp.expectOutputDir).Return(nil)
				h := &handler{svc: svc, opts: &Options{WorkDir: "testdata/workspace"}}

				return h, sp
			},
			params: unzipParams{
				Filepath:  "./extension.zip",
				OutputDir: "./somepath/my-extension",
			},
			expect: func(t *testing.T, sp serviceParams, res *sdkmcp.CallToolResult) {
				// structured content
				assert.Equal(t, res.StructuredContent.(unzipResult).Filepath, sp.expectOutputDir)
				// content
				assertText(t, res, "Successfully unzipped extension")
			},
		},
		{
			name: "should unzip extension without output dir",
			handler: func() (*handler, serviceParams) {
				svc := NewMockcrx3service(ctrl)

				sp := serviceParams{
					expectExtensionFilename: toPath("./testdata/workspace/extension.zip"),
					expectOutputDir:         toPath("./testdata/workspace/unzipped_extension"),
				}

				svc.EXPECT().UnzipTo(sp.expectExtensionFilename, sp.expectOutputDir).Return(nil)
				h := &handler{svc: svc, opts: &Options{WorkDir: "testdata/workspace"}}

				return h, sp
			},
			params: unzipParams{
				Filepath: "./extension.zip",
			},
			expect: func(t *testing.T, sp serviceParams, res *sdkmcp.CallToolResult) {
				// structured content
				assert.Equal(t, res.StructuredContent.(unzipResult).Filepath, sp.expectOutputDir)
				// content
				assertText(t, res, "Successfully unzipped extension")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, sp := tt.handler()
			got, _, gotErr := h.unzipHandler(context.Background(), nil, tt.params)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("unzipHandler() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("unzipHandler() succeeded unexpectedly")
			}
			tt.expect(t, sp, got)
		})
	}
}
