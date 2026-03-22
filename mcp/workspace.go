package mcp

import (
	"context"
	_ "embed"
	"path/filepath"
	"strings"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	//go:embed workspace.md
	workspaceDescription string
	workspaceTitle       = "Get the workspace directory where Chrome extensions are stored"
)

func (h *handler) workspaceHandler(ctx context.Context, _ *sdkmcp.CallToolRequest, _ any) (*sdkmcp.CallToolResult, any, error) {
	return &sdkmcp.CallToolResult{
		StructuredContent: struct {
			Path string `json:"path"`
		}{
			Path: h.opts.WorkDir,
		},
	}, nil, nil
}

func (opts *Options) joinPath(otherPath string) (string, error) {
	baseAbs, err := filepath.Abs(opts.WorkDir)
	if err != nil {
		return "", err
	}
	baseAbs = filepath.Clean(baseAbs) + string(filepath.Separator)
	cleanedRel := filepath.Clean(otherPath)
	fullPath := filepath.Join(baseAbs, cleanedRel)
	resolved, err := filepath.Abs(fullPath)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(resolved+string(filepath.Separator), baseAbs) {
		return "", &PathTraversalError{Base: baseAbs, Attempted: resolved}
	}
	return resolved, nil
}

type PathTraversalError struct {
	Base      string
	Attempted string
}

func (e *PathTraversalError) Error() string {
	return "attempted path traversal outside base directory"
}

func (e *PathTraversalError) Is(target error) bool {
	_, ok := target.(*PathTraversalError)
	return ok
}
