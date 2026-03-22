package mcp

import (
	"context"
	"testing"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestVersionHandler(t *testing.T) {
	h := &handler{
		opts: &Options{
			Version: "v1.0.0",
		},
	}

	ctx := context.Background()
	req := &sdkmcp.CallToolRequest{}
	params := unzipParams{}

	resp, _, err := h.versionHandler(ctx, req, params)
	if err != nil {
		t.Fatalf("versionHandler returned an error: %v", err)
	}

	if resp == nil {
		t.Fatal("response is nil")
	}

	result, ok := resp.StructuredContent.(versinoResult)
	if !ok {
		t.Fatal("data is not of type versinoRsult")
	}

	if result.Version != h.opts.Version {
		t.Errorf("expected version %s, got %s", h.opts.Version, result.Version)
	}

	if !h.opts.DisabledMarkdown {
		if len(resp.Content) == 0 {
			t.Error("expected content to be rendered when markdown is not disabled")
		} else {
			textContent, ok := resp.Content[0].(*sdkmcp.TextContent)
			if !ok {
				t.Error("content is not of type TextContent")
			} else if textContent.Text != "CRX3 Version: v1.0.0" {
				t.Errorf("expected text content \"CRX3 Version: v1.0.0\", got \"%s\"", textContent.Text)
			}
		}
	}
}

func TestVersionHandlerMarkdownDisabled(t *testing.T) {
	h := &handler{
		opts: &Options{
			Version:          "v1.0.0",
			DisabledMarkdown: true,
		},
	}

	ctx := context.Background()
	req := &sdkmcp.CallToolRequest{}
	params := unzipParams{}

	resp, _, err := h.versionHandler(ctx, req, params)
	if err != nil {
		t.Fatalf("versionHandler returned an error: %v", err)
	}

	if len(resp.Content) > 0 {
		t.Error("expected no content when markdown is disabled")
	}
}
