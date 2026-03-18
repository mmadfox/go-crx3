package mcp

import (
	"context"
	_ "embed"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/mediabuyerbot/go-crx3"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

//go:embed download.md
var downloadDescription string

type downloadParams struct {
	IdOrUrl string `json:"idOrUrl" jsonschema:"the extension ID or Chrome Web Store URL to download"`
	Outfile string `json:"outfile" jsonschema:"path to save the downloaded .crx file, optional"`
	Unpack  bool   `json:"unpack" jsonschema:"whether to unpack the extension after downloading, defaults to true"`
}

func downloadHandler(opts *Options) sdkmcp.ToolHandlerFor[downloadParams, any] {
	return func(ctx context.Context, _ *sdkmcp.CallToolRequest, params downloadParams) (*sdkmcp.CallToolResult, any, error) {
		if len(params.IdOrUrl) == 0 {
			return nil, nil, fmt.Errorf("extension ID or URL is required")
		}

		// Extract extension ID if URL is provided
		extensionID := params.IdOrUrl
		if strings.HasPrefix(extensionID, "http") {
			extensionID = extractExtensionID(extensionID)
		}

		// Set default output file if not specified
		outfile := params.Outfile
		if len(outfile) == 0 {
			pwd := opts.WorkDir
			if len(pwd) == 0 {
				pwd = "."
			}
			outfile = path.Join(pwd, "extension.crx")
		}

		// Ensure .crx extension
		if !strings.HasSuffix(outfile, ".crx") {
			outfile = outfile + ".crx"
		}

		// Download the extension
		if err := crx3.DownloadFromWebStore(extensionID, outfile); err != nil {
			return nil, nil, fmt.Errorf("failed to download extension: %w", err)
		}

		// Prepare result message
		var sb strings.Builder
		sb.WriteString("Successfully downloaded extension \"")
		sb.WriteString(extensionID)
		sb.WriteString("\" to \"")
		sb.WriteString(outfile)
		sb.WriteString("\"\n")

		// Unpack if requested
		unpack := params.Unpack
		if !unpack {
			// Default to true if not specified
			unpack = true
		}

		if unpack {
			if err := crx3.Unpack(outfile); err != nil {
				return nil, nil, fmt.Errorf("failed to unpack extension: %w", err)
			}
			dir := strings.TrimSuffix(outfile, ".crx")
			sb.WriteString("Unpacked extension to \"")
			sb.WriteString(dir)
			sb.WriteString("\"\n")
		}

		return textResult(sb.String()), nil, nil
	}
}

func extractExtensionID(rawUrl string) string {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return ""
	}
	urlParts := strings.Split(u.Path, "/")
	if len(urlParts) == 0 {
		return ""
	}
	return urlParts[len(urlParts)-1]
}
