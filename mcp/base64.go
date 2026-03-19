package mcp

import (
	"context"
	_ "embed"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

//go:embed base64.md
var base64Description string

func (h *handler) base64Handler(ctx context.Context, _ *sdkmcp.CallToolRequest, _ any) (*sdkmcp.CallToolResult, any, error) {
	return nil, nil, nil
}
