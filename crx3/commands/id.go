package commands

import (
	"errors"
	"fmt"

	"github.com/mediabuyerbot/go-crx3"
	"github.com/spf13/cobra"
)

func newIDCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "id [infile]",
		Short: "Generate extension id",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("infile is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			infile := args[0]
			id, err := crx3.ID(infile)
			if err != nil {
				return err
			}
			fmt.Println(id)
			return nil
		},
	}
	return cmd
}
