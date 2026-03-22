package mcp

import (
	"context"
	_ "embed"
	"fmt"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	//go:embed getid.md
	getidDescription string
	getidTitle       = "Get Chrome extension ID"
)

type getidParams struct {
	Filepath string `json:"filepath" jsonschema:"required, path to the downloaded .crx file or unpacked extension directory"`
}

type getidResult struct {
	ID string `json:"id" jsonschema:"required, Chrome extension ID"`
}

func (h *handler) getidHandler(ctx context.Context, _ *sdkmcp.CallToolRequest, params getidParams) (*sdkmcp.CallToolResult, any, error) {
	extensionFilepath, err := h.opts.joinPath(params.Filepath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to join path: %w", err)
	}
	if !isFile(extensionFilepath) && !isDir(extensionFilepath) {
		return nil, nil, fmt.Errorf("extension not found %q", params.Filepath)
	}

	id, err := h.svc.GetID(extensionFilepath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get extension ID from %q: %w", params.Filepath, err)
	}

	resp := &sdkmcp.CallToolResult{
		StructuredContent: getidResult{
			ID: id,
		},
	}

	if !h.opts.DisabledMarkdown {
		resp.Content = []sdkmcp.Content{
			&sdkmcp.TextContent{Text: fmt.Sprintf("Extension ID: %s", id)},
		}
	}

	return resp, nil, nil
}
