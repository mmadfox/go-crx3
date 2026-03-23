package mcp

import (
	"context"
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	//go:embed unzip.md
	unzipDescription string
	unzipTitle       = "Unzip a Chrome extension"
)

type unzipParams struct {
	Filepath  string `json:"filepath" jsonschema:"required, path to the downloaded .zip file"`
	OutputDir string `json:"outputDir,omitempty" jsonschema:"optional, path to the output directory"`
}

type unzipResult struct {
	Filepath string `json:"filepath" jsonschema:"required, path to the unzipped extension"`
}

func (h *handler) unzipHandler(ctx context.Context, _ *sdkmcp.CallToolRequest, params unzipParams) (*sdkmcp.CallToolResult, any, error) {
	if filepath.IsAbs(params.Filepath) {
		rel, err := filepath.Rel(h.opts.WorkDir, params.Filepath)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get relative path: %w", err)
		}
		params.Filepath = rel
	}

	if len(params.OutputDir) > 0 && filepath.IsAbs(params.OutputDir) {
		rel, err := filepath.Rel(h.opts.WorkDir, params.OutputDir)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get relative path: %w", err)
		}
		params.OutputDir = rel
	}

	extensionFilepath, err := h.opts.joinPath(params.Filepath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to join path: %w", err)
	}
	isFileInvalid := !isFile(extensionFilepath) || filepath.Ext(extensionFilepath) != ".zip"
	if isFileInvalid {
		return nil, nil, fmt.Errorf("extension not found %q", params.Filepath)
	}

	outputDir := params.OutputDir
	if len(outputDir) == 0 {
		baseName := filepath.Base(params.Filepath)
		baseName = strings.TrimSuffix(baseName, ".zip")
		outputDir = strings.Join([]string{"unzipped", sanitizeFilename(baseName)}, "_")
	}

	if filepath.IsAbs(outputDir) {
		return nil, nil, fmt.Errorf("outputDir must be relative to workspace root, got: %s", outputDir)
	}

	targetDir, err := h.opts.joinPath(outputDir)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to join unzipped path: %w", err)
	}

	if err := h.svc.UnzipTo(extensionFilepath, targetDir); err != nil {
		return nil, nil, fmt.Errorf("failed to unzip %q to %q: %w", params.Filepath, outputDir, err)
	}

	resp := &sdkmcp.CallToolResult{
		StructuredContent: unzipResult{
			Filepath: targetDir,
		},
	}

	// render markdown if not disabled
	if !h.opts.DisabledMarkdown {
		resp.Content = []sdkmcp.Content{
			&sdkmcp.TextContent{Text: fmt.Sprintf("Successfully unzipped extension to %q", targetDir)},
		}
	}

	return resp, nil, nil
}
