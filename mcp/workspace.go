package mcp

import (
	"context"
	_ "embed"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	//go:embed workspace.md
	workspaceDescription string
	workspaceTitle       = "Get the workspace directory where Chrome extensions are stored"
)

func (h *handler) workspaceHandler(ctx context.Context, _ *sdkmcp.CallToolRequest, _ any) (*sdkmcp.CallToolResult, any, error) {
	return &sdkmcp.CallToolResult{
		StructuredContent: struct {
			Path string `json:"path"`
		}{
			Path: h.opts.WorkDir,
		},
	}, nil, nil
}
