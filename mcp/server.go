package mcp

import (
	"context"
	_ "embed"
	"io"

	"github.com/mediabuyerbot/go-crx3"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	searchToolName    = "crx3_search"
	downloadToolName  = "crx3_download"
	workspaceToolName = "crx3_workspace"
)

var allTools = []string{
	searchToolName,
	workspaceToolName,
	downloadToolName,
}

//go:embed instruction_crx3.md
var instruction string

type Options struct {
	Version          string
	Logger           io.Writer
	WorkDir          string
	DisabledTools    []string
	DisabledMarkdown bool
}

type searchOutput struct {
	Results []crx3.SearchResult `json:"results"`
}

type handler struct {
	opts *Options
}

func ServeHTTP(ctx context.Context, addr string) error {
	return nil
}

func ServeStdIO(
	ctx context.Context,
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
		Instructions: instruction,
	}
	mcpServer := sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    "crx3",
		Version: opts.Version,
	}, srvOpts)

	h := &handler{opts: &opts}
	makeTools(mcpServer, h, &opts)

	return mcpServer.Run(ctx, &sdkmcp.StdioTransport{})
}

func textResult(text string) *sdkmcp.CallToolResult {
	return &sdkmcp.CallToolResult{
		Content: []sdkmcp.Content{&sdkmcp.TextContent{Text: text}},
	}
}

func makeTools(mcpServer *sdkmcp.Server, h *handler, opts *Options) {
	isDisabledTools := func(name string) bool {
		for _, t := range opts.DisabledTools {
			if t == name {
				return true
			}
		}
		return false
	}
	for _, toolName := range allTools {
		if isDisabledTools(toolName) {
			continue
		}
		switch toolName {
		case workspaceToolName:
			sdkmcp.AddTool(mcpServer, &sdkmcp.Tool{
				Title:       workspaceTitle,
				Name:        workspaceToolName,
				Description: workspaceDescription,
			}, h.workspaceHandler)
		case searchToolName:
			sdkmcp.AddTool(mcpServer, &sdkmcp.Tool{
				Title:       searchTitle,
				Name:        searchToolName,
				Description: makeSearchDescription(opts.DisabledMarkdown),
			}, h.searchHandler)
		case downloadToolName:
		}
	}
}
