package mcp

import (
	"context"
	_ "embed"
	"io"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

//go:embed instruction_crx3.md
var instruction string

type Options struct {
	Version string
	Logger  io.Writer
	WorkDir string
}

type handler struct {
	mcpServer *sdkmcp.Server
	opts      *Options
}

func ServeHTTP(ctx context.Context, addr string) error {
	return nil
}

func ServeStdIO(
	ctx context.Context,
	opts Options,
) error {
	var mcpTransport sdkmcp.Transport
	if opts.Logger != nil {
		mcpTransport = &sdkmcp.LoggingTransport{
			Transport: &sdkmcp.StdioTransport{},
			Writer:    opts.Logger,
		}
	} else {
		mcpTransport = &sdkmcp.StdioTransport{}
	}

	srvOpts := &sdkmcp.ServerOptions{
		Instructions: instruction,
	}
	mcpServer := sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    "crx3",
		Version: opts.Version,
	}, srvOpts)

	h := &handler{opts: &opts, mcpServer: mcpServer}
	h.makeBase64Tool()

	return mcpServer.Run(ctx, mcpTransport)
}
