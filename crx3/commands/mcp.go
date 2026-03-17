package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

type mcpOpts struct {
	Address       string
	ShowTools     bool
	Logfile       string
	DisabledTools []string
	// TODO: setup workdir
}

func newMCPCmd() *cobra.Command {
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
			if len(opts.Logfile) > 0 {
				f, err := os.Create(opts.Logfile)
				if err != nil {
					return fmt.Errorf("opening logfile: %v", err)
				}
				log.SetOutput(f)
				defer f.Close()
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.Address, "listen", "l", "", "the address on which to run the mcp server")
	cmd.Flags().StringVarP(&opts.Logfile, "logfile", "f", "", "filename to log to; if unset, logs to stderr")
	cmd.Flags().BoolVarP(&opts.ShowTools, "tools.show", "s", false, "if set, print tools with instruction and exit")
	cmd.Flags().StringSliceVarP(&opts.DisabledTools, "tools.disabled", "d", []string{}, "comma-separated list of tool names to disable when running the MCP server")
	return cmd
}
