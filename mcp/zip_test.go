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

func Test_handler_zipHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pwd, _ := os.Getwd()
	toPath := func(other string) string {
		return filepath.Join(pwd, other)
	}

	type serviceParams struct {
		expectSourcePath string
		expectOutputFile string
	}

	tests := []struct {
		name    string
		handler func() (*handler, serviceParams)
		expect  func(*testing.T, serviceParams, *sdkmcp.CallToolResult)
		params  zipParams
		wantErr bool
	}{
		{
			name: "should return error when source path is empty",
			handler: func() (*handler, serviceParams) {
				h := &handler{opts: &Options{WorkDir: "/"}}
				return h, serviceParams{}
			},
			params: zipParams{
				Source: "",
			},
			wantErr: true,
		},
		{
			name: "should return error when source directory does not exist",
			handler: func() (*handler, serviceParams) {
				h := &handler{opts: &Options{WorkDir: "/testdata/workspace"}}
				return h, serviceParams{}
			},
			params: zipParams{
				Source: "nonexistent_directory",
			},
			wantErr: true,
		},
		{
			name: "should zip extension with output directory and custom name",
			handler: func() (*handler, serviceParams) {
				svc := NewMockcrx3service(ctrl)

				sp := serviceParams{
					expectSourcePath: toPath("./testdata/workspace/extensions"),
					expectOutputFile: toPath("./testdata/workspace/zipped_extensions.zip"),
				}

				svc.EXPECT().ZipTo(sp.expectSourcePath, sp.expectOutputFile).Return(nil)
				h := &handler{svc: svc, opts: &Options{WorkDir: "testdata/workspace"}}

				return h, sp
			},
			params: zipParams{
				Source:    "./extensions",
				OutputDir: "./",
				Name:      "zipped_extensions",
			},
			expect: func(t *testing.T, sp serviceParams, res *sdkmcp.CallToolResult) {
				// structured content
				assert.Equal(t, res.StructuredContent.(zipResult).Filepath, sp.expectOutputFile)
				// content
				assertText(t, res, "Successfully zipped extension")
			},
		},
		{
			name: "should zip extension with default name and output to default directory",
			handler: func() (*handler, serviceParams) {
				svc := NewMockcrx3service(ctrl)

				sp := serviceParams{
					expectSourcePath: toPath("./testdata/workspace/extensions"),
					expectOutputFile: toPath("./testdata/workspace/extensions.zip"),
				}

				svc.EXPECT().ZipTo(sp.expectSourcePath, sp.expectOutputFile).Return(nil)
				h := &handler{svc: svc, opts: &Options{WorkDir: "testdata/workspace"}}

				return h, sp
			},
			params: zipParams{
				Source: "./extensions",
			},
			expect: func(t *testing.T, sp serviceParams, res *sdkmcp.CallToolResult) {
				// structured content
				assert.Equal(t, res.StructuredContent.(zipResult).Filepath, sp.expectOutputFile)
				// content
				assertText(t, res, "Successfully zipped extension")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, sp := tt.handler()
			got, _, gotErr := h.zipHandler(context.Background(), nil, tt.params)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("zipHandler() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("zipHandler() succeeded unexpectedly")
			}
			tt.expect(t, sp, got)
		})
	}
}
