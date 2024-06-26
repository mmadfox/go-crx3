package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"

	crx3 "github.com/mediabuyerbot/go-crx3"
	"github.com/spf13/cobra"
)

type unzipOpts struct {
	Outfile string
}

func (o unzipOpts) HasNotOutfile() bool {
	return len(o.Outfile) == 0
}

func newUnzipCmd() *cobra.Command {
	var opts unzipOpts
	cmd := &cobra.Command{
		Use:   "unzip [extension.zip]",
		Short: "Extract all files from the archive",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("extension is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			infile, err := toPath(args[0])
			if err != nil {
				return fmt.Errorf("invalid infile: %w", err)
			}
			zipFile, err := os.Open(infile)
			if err != nil {
				return err
			}
			defer zipFile.Close()
			stat, err := zipFile.Stat()
			if err != nil {
				return err
			}
			if opts.HasNotOutfile() {
				opts.Outfile = strings.TrimSuffix(infile, ".zip")
			}
			return crx3.Unzip(zipFile, stat.Size(), opts.Outfile)
		},
	}

	cmd.Flags().StringVarP(&opts.Outfile, "outfile", "o", "", "save to file")

	return cmd
}
