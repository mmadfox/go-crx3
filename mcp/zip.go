package mcp

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	//go:embed zip.md
	zipDescription string
	zipTitle       = "Zip an unpacked Chrome extension"
)

type zipParams struct {
	Source    string `json:"source" jsonschema:"required, path to the source directory"`
	OutputDir string `json:"outputDir,omitempty" jsonschema:"optional, path to the output directory"`
	Name      string `json:"name,omitempty" jsonschema:"optional, name of the output .zip file"`
}

type zipResult struct {
	Filepath string `json:"filepath" jsonschema:"required, path to the zipped .zip file"`
}

func (h *handler) zipHandler(ctx context.Context, _ *sdkmcp.CallToolRequest, params zipParams) (*sdkmcp.CallToolResult, any, error) {
	if len(params.Source) == 0 {
		return nil, nil, fmt.Errorf("source path is required")
	}

	if filepath.IsAbs(params.Source) {
		rel, err := filepath.Rel(h.opts.WorkDir, params.Source)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get relative path: %w", err)
		}
		params.Source = rel
	}

	if len(params.OutputDir) > 0 && filepath.IsAbs(params.OutputDir) {
		rel, err := filepath.Rel(h.opts.WorkDir, params.OutputDir)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get relative path: %w", err)
		}
		params.OutputDir = rel
	}

	sourcePath, err := h.opts.joinPath(params.Source)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to join source path: %w", err)
	}

	if !isDir(sourcePath) {
		return nil, nil, fmt.Errorf("source is not a directory %q", params.Source)
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

	baseName := filepath.Base(params.Source)
	name := strings.TrimSpace(params.Name)
	if len(name) == 0 {
		name = baseName
	}
	name = sanitizeFilename(name)

	outputFile := filepath.Join(targetDir, fmt.Sprintf("%s.zip", name))

	if err := h.svc.ZipTo(sourcePath, outputFile); err != nil {
		return nil, nil, fmt.Errorf("failed to zip %q to %q: %w", params.Source, outputFile, err)
	}

	resp := &sdkmcp.CallToolResult{
		StructuredContent: zipResult{
			Filepath: outputFile,
		},
	}

	// render markdown if not disabled
	if !h.opts.DisabledMarkdown {
		resp.Content = []sdkmcp.Content{
			&sdkmcp.TextContent{Text: fmt.Sprintf("Successfully zipped extension to %q", outputFile)},
		}
	}

	return resp, nil, nil
}
