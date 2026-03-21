package mcp

import (
	"cmp"
	"context"
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/mediabuyerbot/go-crx3"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	//go:embed download.md
	downloadDescription string
	downloadTitle       = "Download an Chrome extension"
)

type downloadParams struct {
	Target string `json:"url" jsonschema:"required, the extensionId or Chrome Web Store URL to download"`
	Path   string `json:"path,omitempty" jsonschema:"optional, path to save the downloaded .crx file"`
	Name   string `json:"name,omitempty" jsonschema:"optional, name of the downloaded .crx file"`
}

type downloadResult struct {
	Success  bool   `json:"success" jsonschema:"required, whether the download was successful"`
	Filepath string `json:"filepath" jsonschema:"required, the filepath to the downloaded .crx file"`
}

func (h *handler) downloadHandler(ctx context.Context, _ *sdkmcp.CallToolRequest, params downloadParams) (*sdkmcp.CallToolResult, any, error) {
	if len(params.Target) == 0 {
		return nil, nil, fmt.Errorf("extension ID or URL is required")
	}

	if len(params.Name) > 0 && strings.Contains(params.Name, ".") {
		return nil, nil, fmt.Errorf("name must not contain dots or file extensions; provide a clean name without .crx, .zip, etc.")
	}

	var (
		extensionID   string
		extensionName string
	)

	if strings.HasPrefix(params.Target, "http") {
		if !crx3.IsValidChromeWebStoreURL(params.Target) {
			return nil, nil, fmt.Errorf("invalid Chrome Web Store URL")
		}
		extensionID = crx3.ExtractExtensionID(params.Target)
		extensionName, _ = crx3.ExtractExtensionNameFromURL(params.Target)
	} else {
		extensionID = params.Target
	}

	if len(extensionID) == 0 {
		return nil, nil, fmt.Errorf("invalid extension ID or URL")
	}
	if !crx3.IsValidExtensionID(extensionID) {
		return nil, nil, fmt.Errorf("invalid extension ID")
	}

	extensionName = cmp.Or(params.Name, extensionName)
	extensionName = makeName(extensionName, extensionID)

	if len(extensionID) == 0 {
		return nil, nil, fmt.Errorf("invalid extension ID or URL")
	}

	extPath, err := h.opts.joinPath(params.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to join path: %w", err)
	}
	if err := os.MkdirAll(extPath, 0755); err != nil {
		return nil, nil, fmt.Errorf("failed to create directory: %w", err)
	}

	outfile := path.Join(extPath, extensionName)
	if err := h.svc.DownloadFromWebStore(extensionID, outfile); err != nil {
		return nil, nil, fmt.Errorf("failed to download extension: %w", err)
	}

	resp := &sdkmcp.CallToolResult{
		StructuredContent: downloadResult{
			Success:  true,
			Filepath: outfile,
		},
	}

	if !h.opts.DisabledMarkdown {
		var sb strings.Builder
		sb.WriteString("Successfully downloaded extension \"")
		sb.WriteString(extensionID)
		sb.WriteString("\" to \"")
		sb.WriteString(outfile)
		sb.WriteString("\"\n")
		resp.Content = []sdkmcp.Content{
			&sdkmcp.TextContent{Text: sb.String()},
		}
	}

	return resp, nil, nil
}

func makeName(name string, extensionID string) string {
	if len(name) == 0 {
		return fmt.Sprintf("%s.crx", extensionID)
	}
	name = strings.ToLower(name)
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, " ", "_")
	return fmt.Sprintf("%s_%s.crx", name, extensionID)
}
