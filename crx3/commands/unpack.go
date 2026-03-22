package commands

import (
	"errors"
	"fmt"

	crx3 "github.com/mediabuyerbot/go-crx3"
	"github.com/spf13/cobra"
)

func newUnpackCmd() *cobra.Command {
	var opts = struct {
		Outfile       string
		DisableSubDir bool
	}{}
	cmd := &cobra.Command{
		Use:   "unpack [extension.crx] [flags]",
		Short: "Unpack a Chrome extension (.crx) to a directory",
		Long: `Unpack a Chrome extension (.crx file) to a target directory.
By default, creates a subdirectory named after the extension (e.g., 'myext.crx' → './output/myext/').
Use --disable-subdir to extract contents directly into the target directory instead of a subdirectory.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("extension is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			infile, err := toPath(args[0])
			if err != nil {
				return fmt.Errorf("invalid extension filepath: %w", err)
			}
			outfile, err := toPath(opts.Outfile)
			if err != nil {
				return fmt.Errorf("invalid output directory path: %w", err)
			}
			if len(outfile) > 0 {
				var unpackOpts []crx3.UnpackOption
				if opts.DisableSubDir {
					unpackOpts = append(unpackOpts, crx3.UnpackDisableSubdir())
				}
				return crx3.UnpackTo(infile, outfile, unpackOpts...)
			}
			return crx3.Unpack(infile)
		},
	}

	cmd.Flags().StringVarP(&opts.Outfile, "outfile", "o", "", "output directory for unpacked files (default is current directory)")
	cmd.Flags().BoolVarP(&opts.DisableSubDir, "disable-subdir", "s", false, "extract contents directly into the output directory, without creating a subdirectory named after the extension")

	return cmd
}
