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

func Test_handler_downloadHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pwd, _ := os.Getwd()
	toPath := func(other string) string {
		return filepath.Join(pwd, other)
	}

	type serviceParams struct {
		expectExtensionID string
		expectFilepath    string
	}

	tests := []struct {
		name    string
		handler func() (*handler, serviceParams)
		expect  func(*testing.T, serviceParams, *sdkmcp.CallToolResult)
		params  downloadParams
		wantErr bool
	}{
		{
			name: "should return error when extension ID or URL is empty",
			handler: func() (*handler, serviceParams) {
				h := &handler{opts: &Options{WorkDir: "/"}}
				return h, serviceParams{}
			},
			params: downloadParams{
				Target: "",
			},
			wantErr: true,
		},
		{
			name: "should return error when name contains dots",
			handler: func() (*handler, serviceParams) {
				h := &handler{opts: &Options{WorkDir: "/"}}
				return h, serviceParams{}
			},
			params: downloadParams{
				Target: "aaaaaaabbbbbbbbccccccccdddddddddddddddd",
				Name:   "extension.crx",
			},
			wantErr: true,
		},
		{
			name: "should return error when extension URL is invalid",
			handler: func() (*handler, serviceParams) {
				h := &handler{opts: &Options{WorkDir: "/"}}
				return h, serviceParams{}
			},
			params: downloadParams{
				Target: "https://example.com/invalid",
			},
			wantErr: true,
		},
		{
			name: "should return error when extension ID from URL is invalid",
			handler: func() (*handler, serviceParams) {
				h := &handler{opts: &Options{WorkDir: "/"}}
				return h, serviceParams{}
			},
			params: downloadParams{
				Target: "https://chrome.google.com/webstore/detail/invalid-url/aaaaaaabbbbbbbb",
			},
			wantErr: true,
		},
		{
			name: "should download extension with custom path and name",
			handler: func() (*handler, serviceParams) {
				svc := NewMockcrx3service(ctrl)

				sp := serviceParams{
					expectExtensionID: "kpkcennohgffjdgaelocingbmkjnpjgc",
					expectFilepath:    toPath("./testdata/workspace/extensions/extension_kpkcennohgffjdgaelocingbmkjnpjgc.crx"),
				}

				svc.EXPECT().DownloadFromWebStore(sp.expectExtensionID, sp.expectFilepath).Return(nil)
				h := &handler{svc: svc, opts: &Options{WorkDir: "testdata/workspace"}}

				return h, sp
			},
			params: downloadParams{
				Target: "kpkcennohgffjdgaelocingbmkjnpjgc",
				Path:   "./extensions",
				Name:   "extension",
			},
			expect: func(t *testing.T, sp serviceParams, res *sdkmcp.CallToolResult) {
				assert.Equal(t, res.StructuredContent.(downloadResult).Filepath, sp.expectFilepath)
				assert.True(t, res.StructuredContent.(downloadResult).Success)
				assertText(t, res, "Successfully downloaded extension")
			},
		},
		{
			name: "should download extension without custom path and name",
			handler: func() (*handler, serviceParams) {
				svc := NewMockcrx3service(ctrl)

				sp := serviceParams{
					expectExtensionID: "kpkcennohgffjdgaelocingbmkjnpjgc",
					expectFilepath:    toPath("./testdata/workspace/kpkcennohgffjdgaelocingbmkjnpjgc.crx"),
				}

				svc.EXPECT().DownloadFromWebStore(sp.expectExtensionID, sp.expectFilepath).Return(nil)
				h := &handler{svc: svc, opts: &Options{WorkDir: "testdata/workspace"}}

				return h, sp
			},
			params: downloadParams{
				Target: "kpkcennohgffjdgaelocingbmkjnpjgc",
			},
			expect: func(t *testing.T, sp serviceParams, res *sdkmcp.CallToolResult) {
				assert.Equal(t, res.StructuredContent.(downloadResult).Filepath, sp.expectFilepath)
				assert.True(t, res.StructuredContent.(downloadResult).Success)
				assertText(t, res, "Successfully downloaded extension")
			},
		},
		{
			name: "should download extension from Chrome Web Store URL",
			handler: func() (*handler, serviceParams) {
				svc := NewMockcrx3service(ctrl)

				sp := serviceParams{
					expectExtensionID: "kpkcennohgffjdgaelocingbmkjnpjgc",
					expectFilepath:    toPath("./testdata/workspace/extensions/extension_kpkcennohgffjdgaelocingbmkjnpjgc.crx"),
				}

				svc.EXPECT().DownloadFromWebStore(sp.expectExtensionID, sp.expectFilepath).Return(nil)
				h := &handler{svc: svc, opts: &Options{WorkDir: "testdata/workspace"}}

				return h, sp
			},
			params: downloadParams{
				Target: "https://chromewebstore.google.com/detail/example-extension/kpkcennohgffjdgaelocingbmkjnpjgc",
				Path:   "./extensions",
				Name:   "extension",
			},
			expect: func(t *testing.T, sp serviceParams, res *sdkmcp.CallToolResult) {
				assert.Equal(t, res.StructuredContent.(downloadResult).Filepath, sp.expectFilepath)
				assert.True(t, res.StructuredContent.(downloadResult).Success)
				assertText(t, res, "Successfully downloaded extension")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, sp := tt.handler()
			got, _, gotErr := h.downloadHandler(context.Background(), nil, tt.params)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("downloadHandler() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("downloadHandler() succeeded unexpectedly")
			}
			tt.expect(t, sp, got)
		})
	}
}
