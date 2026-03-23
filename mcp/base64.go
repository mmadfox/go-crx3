package mcp

import (
	"context"
	_ "embed"
	"fmt"
	"path/filepath"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	//go:embed base64.md
	base64Description string
	base64Title       = "Encode an extension to base64"
)

type base64Params struct {
	Filepath string `json:"filepath" jsonschema:"required, path to the extension file"`
}

type base64Result struct {
	Data string `json:"data" jsonschema:"required, base64-encoded extension data"`
}

func (h *handler) base64Handler(ctx context.Context, _ *sdkmcp.CallToolRequest, params base64Params) (*sdkmcp.CallToolResult, any, error) {
	if filepath.IsAbs(params.Filepath) {
		rel, err := filepath.Rel(h.opts.WorkDir, params.Filepath)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get relative path: %w", err)
		}
		params.Filepath = rel
	}

	extensionFilepath, err := h.opts.joinPath(params.Filepath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to join path: %w", err)
	}
	if !isFile(extensionFilepath) {
		return nil, nil, fmt.Errorf("extension not found %q", params.Filepath)
	}

	data, err := h.svc.Base64(extensionFilepath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encode extension to base64: %w", err)
	}

	resp := &sdkmcp.CallToolResult{
		StructuredContent: base64Result{
			Data: string(data),
		},
	}

	// render markdown if not disabled
	if !h.opts.DisabledMarkdown {
		resp.Content = []sdkmcp.Content{
			&sdkmcp.TextContent{Text: fmt.Sprintf("Successfully encoded extension %q to base64", params.Filepath)},
		}
	}

	return resp, nil, nil
}
