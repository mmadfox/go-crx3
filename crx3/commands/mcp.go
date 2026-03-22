package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mediabuyerbot/go-crx3/mcp"
	"github.com/spf13/cobra"
)

type mcpOpts struct {
	Address          string
	ShowTools        bool
	Logfile          string
	DisabledTools    []string
	WorkDir          string
	DisabledMarkdown bool
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

			mcpOpts := &mcp.Options{
				Version:          version,
				Logger:           logWriter,
				WorkDir:          opts.WorkDir,
				DisabledMarkdown: opts.DisabledMarkdown,
				DisabledTools:    opts.DisabledTools,
			}
			allTools := mcp.MakeAllTools(mcpOpts)

			// if we're just showing tools, do that and exit
			if opts.ShowTools {
				encoder := json.NewEncoder(os.Stdout)
				encoder.SetIndent("", "  ")
				encoder.SetEscapeHTML(false)
				if err := encoder.Encode(struct {
					Instruction string         `json:"instruction"`
					Tools       []mcp.ToolInfo `json:"tools"`
				}{
					Instruction: mcp.Instruction,
					Tools:       allTools,
				}); err != nil {
					return err
				}
				return nil
			}

			// validate workdir
			opts.WorkDir, err = validateAndNormalizeWorkdir(opts.WorkDir)
			if err != nil {
				return fmt.Errorf("workdir validatation exit with error: %w", err)
			}

			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			// TODO: http, sse
			return mcp.ServeStdIO(ctx, allTools, mcp.Options{
				Version:          version,
				Logger:           logWriter,
				WorkDir:          opts.WorkDir,
				DisabledMarkdown: opts.DisabledMarkdown,
				DisabledTools:    opts.DisabledTools,
			})
		},
	}

	cmd.Flags().StringVarP(&opts.Address, "listen", "l", "", "The address on which to run the mcp server")
	cmd.Flags().StringVarP(&opts.Logfile, "logfile", "f", "", "Filename to log to; if unset, logs to stderr")
	cmd.Flags().BoolVarP(&opts.ShowTools, "tools.show", "s", false, "If set, print tools with instruction and exit")
	cmd.Flags().StringSliceVarP(&opts.DisabledTools, "tools.disabled", "d", []string{}, "Comma-separated list of tool names to disable when running the MCP server")
	cmd.Flags().BoolVarP(&opts.DisabledMarkdown, "tools.disabledMarkdownOutput", "m", false, "If set, disables human-readable text output (Markdown) in tool responses. Only structured data (JSON) will be returned. Intended for automated clients that consume structured content directly")
	cmd.Flags().StringVarP(&opts.WorkDir, "workdir", "w", "", "The working directory in which the server will run. Defaults to the current directory")

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}

func validateAndNormalizeWorkdir(workdir string) (string, error) {
	if workdir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		defaultCrx3Path := filepath.Join(home, "crx3_extensions")
		if err := os.MkdirAll(defaultCrx3Path, 0755); err != nil {
			return "", fmt.Errorf("failed to create default crx3_extensions directory: %w", err)
		}
		workdir = defaultCrx3Path
	}

	if strings.HasPrefix(workdir, "~/") || workdir == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		workdir = strings.Replace(workdir, "~", home, 1)
	}

	absPath, err := filepath.Abs(workdir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("workdir does not exist: %s", absPath)
		}
		return "", fmt.Errorf("failed to stat workdir: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("workdir is not a directory: %s", absPath)
	}

	testFile := filepath.Join(absPath, ".crx3-test-write-"+strconv.FormatInt(time.Now().UnixNano(), 10)+".tmp")
	testData := []byte("test")

	if err := os.WriteFile(testFile, testData, 0644); err != nil {
		return "", fmt.Errorf("workdir is not writable: %s: %w", absPath, err)
	}
	defer os.Remove(testFile)

	data, err := os.ReadFile(testFile)
	if err != nil {
		return "", fmt.Errorf("workdir is not readable: %s: %w", absPath, err)
	}
	if !bytes.Equal(data, testData) {
		return "", fmt.Errorf("workdir read/write test failed: corrupted data")
	}

	return absPath, nil
}
