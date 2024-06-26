package commands

import (
	"errors"
	"fmt"

	crx3 "github.com/mediabuyerbot/go-crx3"
	"github.com/spf13/cobra"
)

type uppackOpts struct {
	Outfile string
}

func newUnpackCmd() *cobra.Command {
	var opts uppackOpts
	cmd := &cobra.Command{
		Use:   "unpack [extension.crx] [flags]",
		Short: "Unpack chrome extension into current directory or specified directory",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("extension is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			infile, err := toPath(args[0])
			if err != nil {
				return fmt.Errorf("invalid infile path: %w", err)
			}
			outfile, err := toPath(opts.Outfile)
			if err != nil {
				return fmt.Errorf("invalid outfile path: %w", err)
			}
			if len(outfile) > 0 {
				return crx3.UnpackTo(infile, outfile)
			}
			return crx3.Unpack(infile)
		},
	}

	cmd.Flags().StringVarP(&opts.Outfile, "outfile", "o", "", "save to file")

	return cmd
}
