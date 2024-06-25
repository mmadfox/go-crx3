package commands

import (
	"errors"
	"os"

	crx3 "github.com/mediabuyerbot/go-crx3"
	"github.com/spf13/cobra"
)

type zipOpts struct {
	Outfile string
}

func (o zipOpts) HasOutfile() bool {
	return len(o.Outfile) > 0
}

func newZipCmd() *cobra.Command {
	var opts zipOpts
	cmd := &cobra.Command{
		Use:   "zip [filepath]",
		Short: "Add unpacked extension to archive",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("filepath is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			infile, err := toPath(args[0])
			if err != nil {
				return err
			}
			if !opts.HasOutfile() {
				opts.Outfile = infile + ".zip"
			}
			zipFile, err := os.Create(opts.Outfile)
			if err != nil {
				return err
			}
			defer zipFile.Close()

			return crx3.Zip(zipFile, infile)
		},
	}

	cmd.Flags().StringVarP(&opts.Outfile, "outfile", "o", "", "save to file")

	return cmd
}
