package mcp

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func Test_handler_unpackHandler(t *testing.T) {
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
		params  unpackParams
		wantErr bool
	}{
		{
			name: "should return error when extension filepath is invalid",
			handler: func() (*handler, serviceParams) {
				h := &handler{opts: &Options{WorkDir: "/"}}
				return h, serviceParams{}
			},
			params: unpackParams{
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
			params: unpackParams{
				Filepath: "some_extension.crx",
			},
			wantErr: true,
		},
		{
			name: "should unpack extension with output dir",
			handler: func() (*handler, serviceParams) {
				svc := NewMockcrx3service(ctrl)

				sp := serviceParams{
					expectExtensionFilename: toPath("./testdata/workspace/extension.crx"),
					expectOutputDir:         toPath("./testdata/workspace/somepath/my_extension"),
				}

				svc.EXPECT().UnpackTo(sp.expectExtensionFilename, sp.expectOutputDir).Return(nil)
				h := &handler{svc: svc, opts: &Options{WorkDir: "testdata/workspace"}}

				return h, sp
			},
			params: unpackParams{
				Filepath:  "./extension.crx",
				OutputDir: "./somepath/my-extension",
			},
			expect: func(t *testing.T, sp serviceParams, res *sdkmcp.CallToolResult) {
				// structured content
				assert.Equal(t, res.StructuredContent.(unpackResult).Filepath, sp.expectOutputDir)
				// content
				assertText(t, res, "Successfully unpacked extension")
			},
		},
		{
			name: "should unpack extension without output dir",
			handler: func() (*handler, serviceParams) {
				svc := NewMockcrx3service(ctrl)

				sp := serviceParams{
					expectExtensionFilename: toPath("./testdata/workspace/extension.crx"),
					expectOutputDir:         toPath("./testdata/workspace/unpacked_extension"),
				}

				svc.EXPECT().UnpackTo(sp.expectExtensionFilename, sp.expectOutputDir).Return(nil)
				h := &handler{svc: svc, opts: &Options{WorkDir: "testdata/workspace"}}

				return h, sp
			},
			params: unpackParams{
				Filepath: "./extension.crx",
			},
			expect: func(t *testing.T, sp serviceParams, res *sdkmcp.CallToolResult) {
				// structured content
				assert.Equal(t, res.StructuredContent.(unpackResult).Filepath, sp.expectOutputDir)
				// content
				assertText(t, res, "Successfully unpacked extension")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, sp := tt.handler()
			got, _, gotErr := h.unpackHandler(context.Background(), nil, tt.params)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("unpackHandler() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("unpackHandler() succeeded unexpectedly")
			}
			tt.expect(t, sp, got)
		})
	}
}

func assertText(t *testing.T, res *sdkmcp.CallToolResult, expected string) {
	assert.Len(t, res.Content, 1)
	data, err := res.Content[0].MarshalJSON()
	assert.NoError(t, err)
	var text = struct {
		Text string `json:"text"`
	}{}
	assert.NoError(t, json.Unmarshal(data, &text))
	assert.Contains(t, text.Text, expected)
}
