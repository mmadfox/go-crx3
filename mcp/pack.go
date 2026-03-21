package mcp

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"crypto/rsa"

	"github.com/mediabuyerbot/go-crx3"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	//go:embed pack.md
	packDescription string
	packTitle       = "Pack an unpacked Chrome extension"
)

type packParams struct {
	Source     string `json:"source" jsonschema:"required, path to the source zip file or unpacked directory"`
	OutputDir  string `json:"outputDir,omitempty" jsonschema:"optional, path to the output directory"`
	Name       string `json:"name,omitempty" jsonschema:"optional, name of the output .crx file"`
	PrivateKey string `json:"privateKey,omitempty" jsonschema:"optional, path to the private key PEM file"`
}

type packResult struct {
	Filepath    string `json:"filepath" jsonschema:"required, path to the packed .crx file"`
	PrivateKey  string `json:"privateKey,omitempty" jsonschema:"optional, path to the generated private key PEM file"`
	ExtensionID string `json:"extensionID,omitempty" jsonschema:"optional, Google Chrome extension ID"`
}

func (h *handler) packHandler(ctx context.Context, _ *sdkmcp.CallToolRequest, params packParams) (*sdkmcp.CallToolResult, any, error) {
	if len(params.Source) == 0 {
		return nil, nil, fmt.Errorf("source path is required")
	}

	sourcePath, err := h.opts.joinPath(params.Source)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to join source path: %w", err)
	}

	if !isFile(sourcePath) && !isDir(sourcePath) {
		return nil, nil, fmt.Errorf("source not found %q", params.Source)
	}

	outputDir := params.OutputDir
	if len(outputDir) == 0 {
		outputDir = "."
	}

	if filepath.IsAbs(outputDir) {
		return nil, nil, fmt.Errorf("outputDir must be relative to workspace root, got: %s", outputDir)
	}

	targetDir, err := h.opts.joinPath(outputDir)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to join output path: %w", err)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return nil, nil, fmt.Errorf("failed to create directory: %w", err)
	}

	var pk *rsa.PrivateKey
	if len(params.PrivateKey) > 0 {
		pemPath, err := h.opts.joinPath(params.PrivateKey)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to join private key path: %w", err)
		}
		pk, err = crx3.LoadPrivateKey(pemPath)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load private key: %w", err)
		}
	}

	sourceExt := path.Ext(sourcePath)
	var baseName string
	if sourceExt == ".zip" || sourceExt == ".crx" {
		baseName = strings.TrimSuffix(path.Base(params.Source), sourceExt)
	} else {
		baseName = path.Base(params.Source)
	}

	name := strings.TrimSpace(params.Name)
	if len(name) == 0 {
		name = baseName
	}
	name = sanitizeFilename(name)

	outputFile := path.Join(targetDir, fmt.Sprintf("%s.crx", name))

	if err := h.svc.PackTo(sourcePath, outputFile, pk); err != nil {
		return nil, nil, fmt.Errorf("failed to pack %q to %q: %w", params.Source, outputFile, err)
	}

	extID, err := crx3.ID(outputFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract extension ID: %w", err)
	}

	resp := &sdkmcp.CallToolResult{
		StructuredContent: packResult{
			Filepath:    outputFile,
			ExtensionID: extID,
		},
	}

	if !h.opts.DisabledMarkdown {
		resp.Content = []sdkmcp.Content{
			&sdkmcp.TextContent{Text: fmt.Sprintf("Successfully packed extension to %q\nExtension ID: `%s`", outputFile, extID)},
		}
	}

	return resp, nil, nil
}

func isDir(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return info.IsDir()
}
