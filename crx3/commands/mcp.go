package commands

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mediabuyerbot/go-crx3/mcp"
	"github.com/spf13/cobra"
)

type mcpOpts struct {
	Address       string
	ShowTools     bool
	Logfile       string
	DisabledTools []string
	WorkDir       string
}

func newMCPCmd(version string) *cobra.Command {
	var opts mcpOpts
	cmd := &cobra.Command{
		Use:   "mcp [mcp-flags]",
		Short: "Start the crx3 MCP server in headless mode",
		Long: `Starts the server over stdio or HTTP (via SSE), depending on whether the --listen flag is provided.

If --listen is specified, the server runs over HTTP using Server-Sent Events (SSE) at the given address.
Otherwise, it runs over standard input/output (stdio) for use with local tools.`,
		Example: `$ crx3 mcp --listen=localhost:3000 # starts over http
$ crx3 mcp  # starts over stdio`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// if we're just showing tools, do that and exit
			if opts.ShowTools {
				// TODO:
				fmt.Println("TODO: print tools")
				return nil
			}

			// set up logging if we have a logfile
			var logWriter io.Writer
			if len(opts.Logfile) > 0 {
				f, err := os.Create(opts.Logfile)
				if err != nil {
					return fmt.Errorf("opening logfile: %v", err)
				}
				log.SetOutput(f)
				defer f.Close()
				logWriter = log.Writer()
			}

			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			// TODO: http, sse
			return mcp.ServeStdIO(ctx, mcp.Options{
				Version: version,
				Logger:  logWriter,
				WorkDir: opts.WorkDir,
			})
		},
	}

	cmd.Flags().StringVarP(&opts.Address, "listen", "l", "", "The address on which to run the mcp server")
	cmd.Flags().StringVarP(&opts.Logfile, "logfile", "f", "", "Filename to log to; if unset, logs to stderr")
	cmd.Flags().BoolVarP(&opts.ShowTools, "tools.show", "s", false, "If set, print tools with instruction and exit")
	cmd.Flags().StringSliceVarP(&opts.DisabledTools, "tools.disabled", "d", []string{}, "Comma-separated list of tool names to disable when running the MCP server")
	cmd.Flags().StringVarP(&opts.WorkDir, "workdir", "w", "", "The working directory in which the server will run. Defaults to the current directory")
	return cmd
}
