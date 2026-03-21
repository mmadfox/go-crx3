package mcp

import (
	"context"
	"testing"

	"github.com/mediabuyerbot/go-crx3"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func Test_handler_scanHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		handler func() *handler
		expect  func(*testing.T, *sdkmcp.CallToolResult)
		params  scanParams
		wantErr bool
	}{
		{
			name: "should return error when limit is negative",
			handler: func() *handler {
				return &handler{opts: &Options{WorkDir: "/"}}
			},
			params: scanParams{
				Limit: -1,
			},
			wantErr: true,
		},
		{
			name: "should return error when limit exceeds 100",
			handler: func() *handler {
				return &handler{opts: &Options{WorkDir: "/"}}
			},
			params: scanParams{
				Limit: 101,
			},
			wantErr: true,
		},
		{
			name: "should return error when scan fails",
			handler: func() *handler {
				svc := NewMockcrx3service(ctrl)
				svc.EXPECT().Scan("testdata/workspace", gomock.Any()).Return(func(yield func(*crx3.ExtensionInfo, error) bool) {
					yield(nil, assert.AnError)
				})
				return &handler{svc: svc, opts: &Options{WorkDir: "testdata/workspace"}}
			},
			params: scanParams{
				Limit: 10,
			},
			wantErr: true,
		},
		{
			name: "should return error when no extensions found",
			handler: func() *handler {
				svc := NewMockcrx3service(ctrl)
				svc.EXPECT().Scan("testdata/workspace", gomock.Any()).Return(func(yield func(*crx3.ExtensionInfo, error) bool) {
				})
				return &handler{svc: svc, opts: &Options{WorkDir: "testdata/workspace"}}
			},
			params: scanParams{
				Limit: 10,
			},
			wantErr: true,
		},
		{
			name: "should scan workspace and return extensions with markdown enabled",
			handler: func() *handler {
				svc := NewMockcrx3service(ctrl)
				svc.EXPECT().Scan("testdata/workspace", gomock.Any()).Return(func(yield func(*crx3.ExtensionInfo, error) bool) {
					yield(&crx3.ExtensionInfo{
						Name:     "Test Extension",
						Path:     "test/path",
						Type:     "crx",
						Size:     12345,
						Modified: "2024-01-01",
					}, nil)
					yield(&crx3.ExtensionInfo{
						Name:     "Another Extension",
						Path:     "another/path",
						Type:     "dir",
						Size:     0,
						Modified: "2024-01-02",
					}, nil)
				})
				return &handler{svc: svc, opts: &Options{WorkDir: "testdata/workspace", DisabledMarkdown: false}}
			},
			params: scanParams{
				Limit: 10,
			},
			expect: func(t *testing.T, res *sdkmcp.CallToolResult) {
				assert.Equal(t, 2, len(res.StructuredContent.(scanResult).Results))
				assertText(t, res, "| Name | Path | Type | Size | Modified |")
				assertText(t, res, "Test Extension")
				assertText(t, res, "Another Extension")
			},
		},
		{
			name: "should scan workspace with filter and return extensions",
			handler: func() *handler {
				svc := NewMockcrx3service(ctrl)
				svc.EXPECT().Scan("testdata/workspace", gomock.Any()).Return(func(yield func(*crx3.ExtensionInfo, error) bool) {
					yield(&crx3.ExtensionInfo{
						Name:     "React Extension",
						Path:     "react/path",
						Type:     "crx",
						Size:     54321,
						Modified: "2024-01-03",
					}, nil)
				})
				return &handler{svc: svc, opts: &Options{WorkDir: "testdata/workspace", DisabledMarkdown: false}}
			},
			params: scanParams{
				Limit:  5,
				Filter: []string{"react"},
			},
			expect: func(t *testing.T, res *sdkmcp.CallToolResult) {
				assert.Equal(t, 1, len(res.StructuredContent.(scanResult).Results))
				assert.Equal(t, "React Extension", res.StructuredContent.(scanResult).Results[0].Name)
			},
		},
		{
			name: "should scan workspace with markdown disabled",
			handler: func() *handler {
				svc := NewMockcrx3service(ctrl)
				svc.EXPECT().Scan("testdata/workspace", gomock.Any()).Return(func(yield func(*crx3.ExtensionInfo, error) bool) {
					yield(&crx3.ExtensionInfo{
						Name:     "Test Extension",
						Path:     "test/path",
						Type:     "crx",
						Size:     12345,
						Modified: "2024-01-01",
					}, nil)
				})
				return &handler{svc: svc, opts: &Options{WorkDir: "testdata/workspace", DisabledMarkdown: true}}
			},
			params: scanParams{
				Limit: 10,
			},
			expect: func(t *testing.T, res *sdkmcp.CallToolResult) {
				assert.Equal(t, 1, len(res.StructuredContent.(scanResult).Results))
				assert.Empty(t, res.Content)
			},
		},
		{
			name: "should use default limit when limit is 0",
			handler: func() *handler {
				svc := NewMockcrx3service(ctrl)
				svc.EXPECT().Scan("testdata/workspace", gomock.Any()).Return(func(yield func(*crx3.ExtensionInfo, error) bool) {
					yield(&crx3.ExtensionInfo{
						Name:     "Test Extension",
						Path:     "test/path",
						Type:     "crx",
						Size:     12345,
						Modified: "2024-01-01",
					}, nil)
				})
				return &handler{svc: svc, opts: &Options{WorkDir: "testdata/workspace", DisabledMarkdown: false}}
			},
			params: scanParams{
				Limit: 0,
			},
			expect: func(t *testing.T, res *sdkmcp.CallToolResult) {
				assert.Equal(t, 1, len(res.StructuredContent.(scanResult).Results))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.handler()
			got, _, gotErr := h.scanHandler(context.Background(), nil, tt.params)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("scanHandler() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("scanHandler() succeeded unexpectedly")
			}
			tt.expect(t, got)
		})
	}
}

func Test_makeScanMarkdownTable(t *testing.T) {
	t.Run("should return message when no extensions", func(t *testing.T) {
		result := makeScanMarkdownTable(nil)
		assert.Equal(t, "No extensions found.", result)
	})

	t.Run("should render table with extensions", func(t *testing.T) {
		extensions := []*crx3.ExtensionInfo{
			{
				Name:     "Test Extension",
				Path:     "test/path",
				Type:     "crx",
				Size:     12345,
				Modified: "2024-01-01",
			},
			{
				Name:     "Dir Extension",
				Path:     "dir/path",
				Type:     "dir",
				Size:     0,
				Modified: "2024-01-02",
			},
		}
		result := makeScanMarkdownTable(extensions)
		assert.Contains(t, result, "| Name | Path | Type | Size | Modified |")
		assert.Contains(t, result, "Test Extension")
		assert.Contains(t, result, "12345")
		assert.Contains(t, result, "-")
	})

	t.Run("should handle empty name", func(t *testing.T) {
		extensions := []*crx3.ExtensionInfo{
			{
				Name:     "",
				Path:     "test/path",
				Type:     "crx",
				Size:     100,
				Modified: "2024-01-01",
			},
		}
		result := makeScanMarkdownTable(extensions)
		assert.Contains(t, result, "\\*unknown\\*")
	})
}

func Test_escapeMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "should escape pipe character",
			input:    "test|pipe",
			expected: "test&#124;pipe",
		},
		{
			name:     "should escape asterisk",
			input:    "test*asterisk",
			expected: "test\\*asterisk",
		},
		{
			name:     "should escape underscore",
			input:    "test_underscore",
			expected: "test\\_underscore",
		},
		{
			name:     "should escape multiple special chars",
			input:    "test|*_chars",
			expected: "test&#124;\\*\\_chars",
		},
		{
			name:     "should not modify normal string",
			input:    "normal string",
			expected: "normal string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapeMarkdown(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
