package mcp

import (
	"context"
	_ "embed"
	"fmt"
	"strconv"
	"strings"

	"github.com/mediabuyerbot/go-crx3"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	//go:embed search.md
	searchDescription string
	searchTitle       = "Search for Chrome extensions"
)

type searchParams struct {
	Name string `json:"name" jsonschema:"required,the search query to find Chrome extensions by name"`
}

type searchResult struct {
	Results []crx3.SearchResult `json:"results"`
}

func (h *handler) searchHandler(ctx context.Context, _ *sdkmcp.CallToolRequest, params searchParams) (*sdkmcp.CallToolResult, any, error) {
	if len(params.Name) == 0 {
		return nil, nil, fmt.Errorf("empty query")
	}

	results, err := h.svc.SearchExtensionByName(ctx, params.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to search extension %q: %w", params.Name, err)
	}

	if len(results) == 0 {
		return nil, nil, fmt.Errorf("no extension found %s", params.Name)
	}

	resp := &sdkmcp.CallToolResult{
		StructuredContent: searchResult{
			Results: results,
		},
	}

	// render markdown if not disabled
	if !h.opts.DisabledMarkdown {
		var sb strings.Builder
		sb.WriteString("# Found chrome extensions: ")
		sb.WriteString(strconv.Itoa(len(results)))
		sb.WriteString("\n")

		for i, extension := range results {
			row := i + 1
			sb.WriteString(strconv.Itoa(row))
			sb.WriteString(". ")
			sb.WriteString(extension.Name)
			sb.WriteString("\n")
			sb.WriteString("URL: [")
			sb.WriteString(extension.Name)
			sb.WriteString("](")
			sb.WriteString(extension.URL)
			sb.WriteString(")\n")
			sb.WriteString("ExtensionID: ")
			sb.WriteString(extension.ExtensionID)
			sb.WriteString("\n")
		}

		resp.Content = []sdkmcp.Content{
			&sdkmcp.TextContent{Text: sb.String()},
		}
	}

	return resp, nil, nil
}
