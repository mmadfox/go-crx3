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

func Test_handler_base64Handler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pwd, _ := os.Getwd()
	toPath := func(other string) string {
		return filepath.Join(pwd, other)
	}

	type serviceParams struct {
		expectExtensionFilepath string
	}

	tests := []struct {
		name    string
		handler func() (*handler, serviceParams)
		expect  func(*testing.T, serviceParams, *sdkmcp.CallToolResult)
		params  base64Params
		wantErr bool
	}{
		{
			name: "should return error when filepath is invalid",
			handler: func() (*handler, serviceParams) {
				h := &handler{opts: &Options{WorkDir: "/"}}
				return h, serviceParams{}
			},
			params: base64Params{
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
			params: base64Params{
				Filepath: "some_extension.crx",
			},
			wantErr: true,
		},
		{
			name: "should encode extension to base64",
			handler: func() (*handler, serviceParams) {
				svc := NewMockcrx3service(ctrl)

				sp := serviceParams{
					expectExtensionFilepath: toPath("./testdata/workspace/extension.crx"),
				}

				expectedData := "base64_encoded_data"

				svc.EXPECT().Base64(sp.expectExtensionFilepath).Return([]byte(expectedData), nil)
				h := &handler{svc: svc, opts: &Options{WorkDir: "testdata/workspace"}}

				return h, sp
			},
			params: base64Params{
				Filepath: "./extension.crx",
			},
			expect: func(t *testing.T, sp serviceParams, res *sdkmcp.CallToolResult) {
				// structured content
				assert.Equal(t, res.StructuredContent.(base64Result).Data, "base64_encoded_data")
				// content
				assertText(t, res, "Successfully encoded extension")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, sp := tt.handler()
			got, _, gotErr := h.base64Handler(context.Background(), nil, tt.params)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("base64Handler() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("base64Handler() succeeded unexpectedly")
			}
			tt.expect(t, sp, got)
		})
	}
}
