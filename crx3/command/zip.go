package command

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
		Use:   "zip [infile]",
		Short: "Creates a *.zip archive and adds all the files to it.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("infile is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			infile := args[0]
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

	cmd.Flags().StringVarP(&opts.Outfile, "outfile", "o", "", "save to the file")

	return cmd
}
