package mcp

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/mediabuyerbot/go-crx3"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	//go:embed scan.md
	scanDescription string
	scanTitle       = "Scan Chrome extensions in the workspace"
)

type scanParams struct {
	Limit  int      `json:"limit,omitempty" jsonschema:"maximum number of extensions to return. Use 0 or omit for no limit."`
	Filter []string `json:"filter,omitempty" jsonschema:"list of keywords to filter extensions by name. Case-insensitive partial match. Example: ['react', 'adblock'] matches any extension with 'react' or 'adblock' in the name."`
}

type scanResult struct {
	Results []*crx3.ExtensionInfo `json:"results"`
}

func (h *handler) scanHandler(ctx context.Context, _ *sdkmcp.CallToolRequest, params scanParams) (*sdkmcp.CallToolResult, any, error) {
	if params.Limit < 0 {
		return nil, nil, fmt.Errorf("limit must be non-negative")
	}
	if params.Limit == 0 {
		params.Limit = 20
	}
	if params.Limit > 100 {
		return nil, nil, fmt.Errorf(`limit must be less than 100`)
	}

	scanOpts := make([]crx3.ScanOption, 0, 4)
	if len(params.Filter) > 0 {
		for _, filter := range params.Filter {
			scanOpts = append(scanOpts, crx3.WithNameFilter(filter))
		}
	}
	scanOpts = append(scanOpts, crx3.WithMaxResults(params.Limit))
	scanOpts = append(scanOpts, crx3.WithMaxDepth(10))

	var results []*crx3.ExtensionInfo
	for info, err := range crx3.Scan(h.opts.WorkDir, scanOpts...) {
		if err != nil {
			return nil, nil, fmt.Errorf("scan workspace [%q] error: %w", h.opts.WorkDir, err)
		}
		if info != nil {
			results = append(results, info)
		}
	}
	if len(results) == 0 {
		return nil, nil, fmt.Errorf("No extensions found")
	}

	resp := &sdkmcp.CallToolResult{
		StructuredContent: scanResult{
			Results: results,
		},
	}

	// render markdown if not disabled
	if !h.opts.DisabledMarkdown {
		resp.Content = []sdkmcp.Content{
			&sdkmcp.TextContent{Text: makeScanMarkdownTable(results)},
		}
	}

	return resp, nil, nil
}

func makeScanMarkdownTable(extensions []*crx3.ExtensionInfo) string {
	if len(extensions) == 0 {
		return "No extensions found."
	}

	var sb strings.Builder

	sb.WriteString("| Name | Path | Type | Size | Modified |\n")
	sb.WriteString("|------|------|------|------|----------|\n")

	for _, ext := range extensions {
		name := ext.Name
		if name == "" {
			name = "*unknown*"
		}
		var sizeStr string
		if ext.Type == "dir" {
			sizeStr = "-"
		} else {
			sizeStr = fmt.Sprintf("%d", ext.Size)
		}
		name = escapeMarkdown(name)
		path := escapeMarkdown(ext.Path)
		et := escapeMarkdown(ext.Type)
		sizeStr = escapeMarkdown(sizeStr)
		modified := escapeMarkdown(ext.Modified)

		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
			name, path, et, sizeStr, modified))
	}

	return sb.String()
}

func escapeMarkdown(s string) string {
	s = strings.ReplaceAll(s, "|", "&#124;")
	s = strings.ReplaceAll(s, "*", "\\*")
	s = strings.ReplaceAll(s, "_", "\\_")
	return s
}
