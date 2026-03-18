package mcp

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/mediabuyerbot/go-crx3"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	//go:embed search.md
	searchDescription string
	searchTitle       = "Search for Chrome extensions"
)

func makeSearchDescription(disabledMarkdown bool) string {
	tmpl, err := template.New("search").Parse(searchDescription)
	if err != nil {
		panic(err)
	}
	params := struct {
		DisabledMarkdown bool
	}{
		DisabledMarkdown: disabledMarkdown,
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, params); err != nil {
		panic(err)
	}
	return strings.TrimSpace(buf.String())
}

type searchParams struct {
	Name string `json:"name" jsonschema:"required,the search query to find Chrome extensions by name"`
}

func (h *handler) searchHandler(ctx context.Context, _ *sdkmcp.CallToolRequest, params searchParams) (*sdkmcp.CallToolResult, any, error) {
	if len(params.Name) == 0 {
		return nil, nil, fmt.Errorf("empty query")
	}

	results, err := crx3.SearchExtensionByName(ctx, params.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to search extension %q: %w", params.Name, err)
	}

	if len(results) == 0 {
		return nil, nil, fmt.Errorf("no extension found %s", params.Name)
	}

	resp := &sdkmcp.CallToolResult{
		StructuredContent: searchOutput{
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
