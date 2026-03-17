package mcp

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/mediabuyerbot/go-crx3"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

//go:embed search.md
var searchDescription string

type searchParams struct {
	Name string `json:"name" jsonschema:"the search query to use for matching chrome extension"`
}

func (h *handler) makeSearchTool() {
	sdkmcp.AddTool(h.mcpServer, &sdkmcp.Tool{
		Name:        "crx3_search_chrome_extension",
		Description: searchDescription,
	}, h.searchHandler)
}

func (h *handler) searchHandler(ctx context.Context, _ *sdkmcp.CallToolRequest, params searchParams) (*sdkmcp.CallToolResult, any, error) {
	if len(params.Name) == 0 {
		return nil, nil, fmt.Errorf("empty query")
	}

	results, err := crx3.SearchExtensionByName(ctx, params.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to search extension: %w", err)
	}

	if len(results) == 0 {
		return nil, nil, fmt.Errorf("no extension found")
	}

	var sb strings.Builder
	sb.WriteString("## Found chrome extensions:\n")
	for _, etension := range results {
		sb.WriteString("Name: ")
		sb.WriteString(etension.Name)
		sb.WriteString("\n")

		sb.WriteString("URL: [")
		sb.WriteString(etension.Name)
		sb.WriteString("](")
		sb.WriteString(etension.URL)
		sb.WriteString(")\n")
		sb.WriteString("\n")

		sb.WriteString("ExtensionID: ")
		sb.WriteString(etension.ExtensionID)
		sb.WriteString("\n")

		sb.WriteString("---")
		sb.WriteString("\n")
	}

	return textResult(sb.String()), nil, nil
}
