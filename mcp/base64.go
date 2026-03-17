package mcp

import (
	"context"
	_ "embed"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

//go:embed base64.md
var base64Description string

type base64ToolParams struct {
}

func (h *handler) makeBase64Tool() {
	sdkmcp.AddTool(h.mcpServer, &sdkmcp.Tool{
		Name:        "crx3_base64",
		Description: base64Description,
	}, h.base64Handler)
}

func (h *handler) base64Handler(ctx context.Context, _ *sdkmcp.CallToolRequest, params base64ToolParams) (*sdkmcp.CallToolResult, any, error) {
	return nil, nil, nil
}
