package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/mediabuyerbot/go-crx3"
	"github.com/spf13/cobra"
)

func newScanCmd() *cobra.Command {
	var (
		nameFilters string
		rootPath    string
		maxDepth    int
		maxLimit    int
	)

	cmd := &cobra.Command{
		Use:   "scan [path]",
		Short: "Scan a directory for Chrome extensions",
		Long:  `Scan recursively scans a directory for Chrome extensions in CRX, ZIP, or unpacked directory formats`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := rootPath
			if len(args) > 0 {
				path = args[0]
			}
			if path == "" {
				path = "."
			}
			if strings.HasPrefix(path, "~/") || path == "~" {
				home, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to get home directory: %w", err)
				}
				path = strings.Replace(path, "~", home, 1)
			}

			var opts []crx3.ScanOption
			if nameFilters != "" {
				for _, filter := range strings.Split(nameFilters, ",") {
					filter = strings.TrimSpace(filter)
					if filter != "" {
						opts = append(opts, crx3.WithNameFilter(filter))
					}
				}
			}
			if maxDepth >= 0 {
				opts = append(opts, crx3.WithMaxDepth(maxDepth))
			}
			if maxLimit > 0 {
				opts = append(opts, crx3.WithMaxResults(maxLimit))
			}

			var results []*crx3.ExtensionInfo
			for info, err := range crx3.Scan(path, opts...) {
				if err != nil {
					return fmt.Errorf("scan error: %w", err)
				}
				if info != nil {
					results = append(results, info)
				}
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

	cmd.Flags().StringVar(&nameFilters, "filter", "", "Filter extensions by (partial) names, comma-separated")
	cmd.Flags().IntVar(&maxDepth, "depth", 5, "Maximum directory depth to scan (0 = only root, -1 = unlimited)")
	cmd.Flags().IntVar(&maxLimit, "limit", 15, "Maximum number of extensions to find (0 = unlimited)")

	return cmd
}
