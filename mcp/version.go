package mcp

import (
	"context"
	_ "embed"
	"fmt"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	//go:embed version.md
	versionDescription string
	versionTitle       = "CRX3 Current Version"
)

type versinoResult struct {
	Version string `json:"version"`
}

func (h *handler) versionHandler(ctx context.Context, _ *sdkmcp.CallToolRequest, _ any) (*sdkmcp.CallToolResult, any, error) {
	resp := &sdkmcp.CallToolResult{
		StructuredContent: versinoResult{
			Version: h.opts.Version,
		},
	}

	// render markdown if not disabled
	if !h.opts.DisabledMarkdown {
		vertex := fmt.Sprintf("CRX3 Version: %s", h.opts.Version)
		resp.Content = []sdkmcp.Content{
			&sdkmcp.TextContent{Text: vertex},
		}
	}

	return resp, nil, nil
}
