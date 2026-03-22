package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mediabuyerbot/go-crx3"
	"github.com/spf13/cobra"
)

func newSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search [name]",
		Short: "Search for Chrome extensions by name",
		Long:  `Search for Chrome extensions using DuckDuckGo and extract relevant results.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			name := args[0]

			results, err := crx3.SearchExtensionByName(ctx, name)
			if err != nil {
				return fmt.Errorf("failed to search extension: %w", err)
			}

			if len(results) == 0 {
				fmt.Fprintln(os.Stderr, "No extensions found.")
				return nil
			}

			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(results); err != nil {
				return fmt.Errorf("failed to encode results: %w", err)
			}

			return nil
		},
	}

	return cmd
}
