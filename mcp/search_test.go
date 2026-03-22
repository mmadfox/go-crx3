package mcp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mediabuyerbot/go-crx3"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func Test_handler_searchHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		handler func() *handler
		expect  func(*testing.T, *sdkmcp.CallToolResult)
		params  searchParams
		wantErr bool
	}{
		{
			name: "should return error when query is empty",
			handler: func() *handler {
				return &handler{opts: &Options{WorkDir: "/"}}
			},
			params: searchParams{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "should return error when search fails",
			handler: func() *handler {
				svc := NewMockcrx3service(ctrl)
				svc.EXPECT().SearchExtensionByName(context.Background(), "test").Return(nil, assert.AnError)
				return &handler{svc: svc, opts: &Options{WorkDir: "/"}}
			},
			params: searchParams{
				Name: "test",
			},
			wantErr: true,
		},
		{
			name: "should return error when no extensions found",
			handler: func() *handler {
				svc := NewMockcrx3service(ctrl)
				svc.EXPECT().SearchExtensionByName(context.Background(), "test").Return([]crx3.SearchResult{}, nil)
				return &handler{svc: svc, opts: &Options{WorkDir: "/"}}
			},
			params: searchParams{
				Name: "test",
			},
			wantErr: true,
		},
		{
			name: "should return search results with markdown enabled",
			handler: func() *handler {
				svc := NewMockcrx3service(ctrl)
				results := []crx3.SearchResult{
					{
						Name:        "Test Extension",
						URL:         "https://chrome.google.com/webstore/detail/test",
						ExtensionID: "test-extension-id",
					},
					{
						Name:        "Another Extension",
						URL:         "https://chrome.google.com/webstore/detail/another",
						ExtensionID: "another-extension-id",
					},
				}
				svc.EXPECT().SearchExtensionByName(context.Background(), "test").Return(results, nil)
				return &handler{svc: svc, opts: &Options{WorkDir: "/", DisabledMarkdown: false}}
			},
			params: searchParams{
				Name: "test",
			},
			expect: func(t *testing.T, res *sdkmcp.CallToolResult) {
				// structured content
				assert.Equal(t, 2, len(res.StructuredContent.(searchResult).Results))
				// content
				assertText(t, res, "# Found chrome extensions: 2")
				assertText(t, res, "Test Extension")
				assertText(t, res, "Another Extension")
			},
		},
		{
			name: "should return search results with markdown disabled",
			handler: func() *handler {
				svc := NewMockcrx3service(ctrl)
				results := []crx3.SearchResult{
					{
						Name:        "Test Extension",
						URL:         "https://chrome.google.com/webstore/detail/test",
						ExtensionID: "test-extension-id",
					},
				}
				svc.EXPECT().SearchExtensionByName(context.Background(), "test").Return(results, nil)
				return &handler{svc: svc, opts: &Options{WorkDir: "/", DisabledMarkdown: true}}
			},
			params: searchParams{
				Name: "test",
			},
			expect: func(t *testing.T, res *sdkmcp.CallToolResult) {
				// structured content
				assert.Equal(t, 1, len(res.StructuredContent.(searchResult).Results))
				// no content when markdown is disabled
				assert.Empty(t, res.Content)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.handler()
			got, _, gotErr := h.searchHandler(context.Background(), nil, tt.params)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("searchHandler() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("searchHandler() succeeded unexpectedly")
			}
			tt.expect(t, got)
		})
	}
}

func Test_searchResult_MarshalJSON(t *testing.T) {
	result := searchResult{
		Results: []crx3.SearchResult{
			{
				Name:        "Test Extension",
				URL:         "https://chrome.google.com/webstore/detail/test",
				ExtensionID: "test-extension-id",
			},
		},
	}

	data, err := json.Marshal(result)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
}
