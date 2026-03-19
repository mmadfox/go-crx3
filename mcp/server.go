package mcp

import (
	"context"
	_ "embed"
	"io"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	searchToolName    = "crx3_search"
	downloadToolName  = "crx3_download"
	workspaceToolName = "crx3_workspace"
)

type ToolInfo struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func MakeAllTools(opts *Options) []ToolInfo {
	tools := make([]ToolInfo, 0)
	isNotDisabledTool := func(name string) bool {
		for _, t := range opts.DisabledTools {
			if t == name {
				return false
			}
		}
		return true
	}

	if isNotDisabledTool(searchToolName) {
		tools = append(tools, ToolInfo{
			Name:        searchToolName,
			Title:       searchTitle,
			Description: makeSearchDescription(opts.DisabledMarkdown),
		})
	}

	if isNotDisabledTool(workspaceToolName) {
		tools = append(tools, ToolInfo{
			Name:        workspaceToolName,
			Title:       workspaceTitle,
			Description: workspaceDescription,
		})
	}

	if isNotDisabledTool(downloadToolName) {
		tools = append(tools, ToolInfo{
			Name:        downloadToolName,
			Title:       downloadTitle,
			Description: downloadDescription,
		})
	}

	return tools
}

//go:embed instruction.md
var Instruction string

type Options struct {
	Version          string
	Logger           io.Writer
	WorkDir          string
	DisabledTools    []string
	DisabledMarkdown bool
}

type handler struct {
	opts *Options
}

func ServeHTTP(ctx context.Context, addr string) error {
	return nil
}

func ServeStdIO(
	ctx context.Context,
	tools []ToolInfo,
	opts Options,
) error {
	// TODO:
	// var mcpTransport sdkmcp.Transport
	// if opts.Logger != nil {
	// 	mcpTransport = &sdkmcp.LoggingTransport{
	// 		Transport: &sdkmcp.StdioTransport{},
	// 		Writer:    opts.Logger,
	// 	}
	// } else {
	// 	mcpTransport = &sdkmcp.StdioTransport{}
	// }

	srvOpts := &sdkmcp.ServerOptions{
		Instructions: Instruction,
	}
	mcpServer := sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    "crx3",
		Version: opts.Version,
	}, srvOpts)

	h := &handler{opts: &opts}
	makeTools(mcpServer, h, tools)

	return mcpServer.Run(ctx, &sdkmcp.StdioTransport{})
}

func makeTools(mcpServer *sdkmcp.Server, h *handler, tools []ToolInfo) {
	for _, tool := range tools {
		switch tool.Name {
		case workspaceToolName:
			sdkmcp.AddTool(mcpServer, &sdkmcp.Tool{
				Title:       tool.Title,
				Name:        tool.Name,
				Description: tool.Description,
				Annotations: &sdkmcp.ToolAnnotations{
					IdempotentHint: true,
					ReadOnlyHint:   true,
				},
			}, h.workspaceHandler)
		case searchToolName:
			sdkmcp.AddTool(mcpServer, &sdkmcp.Tool{
				Title:       tool.Title,
				Name:        tool.Name,
				Description: tool.Description,
				Annotations: &sdkmcp.ToolAnnotations{
					IdempotentHint: true,
					ReadOnlyHint:   true,
				},
			}, h.searchHandler)
		case downloadToolName:
			sdkmcp.AddTool(mcpServer, &sdkmcp.Tool{
				Title:       tool.Title,
				Name:        tool.Name,
				Description: tool.Description,
			}, h.downloadHandler)
		}
	}
}
