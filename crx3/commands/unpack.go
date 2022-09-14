package commands

import (
	"errors"

	crx3 "github.com/mediabuyerbot/go-crx3"
	"github.com/spf13/cobra"
)

func newUnpackCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unpack [extension.crx]",
		Short: "Unpack chrome extension into current directory",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("extension is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			infile := args[0]
			return crx3.Unpack(infile)
		},
	}
	return cmd
}
