package mcp

import (
	"context"
	_ "embed"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	//go:embed unpack.md
	unpackDescription string
	unpackTitle       = "Unpack a Chrome extension"
)

type unpackParams struct {
	Filepath  string `json:"filepath" jsonschema:"required, path to the downloaded .crx file"`
	Name      string `json:"name,omitempty" jsonschema:"optional, name of the unpacked extension"`
	OutputDir string `json:"outputDir,omitempty" jsonschema:"optional, path to the output directory"`
}

type unpackResult struct {
	Filepath string `json:"filepath" jsonschema:"required, path to the unpacked extension"`
}

func (h *handler) unpackHandler(ctx context.Context, _ *sdkmcp.CallToolRequest, params unpackParams) (*sdkmcp.CallToolResult, any, error) {
	return nil, nil, nil
}
