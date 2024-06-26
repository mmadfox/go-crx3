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
		Short: "Generate id from header extension or manifest file",
		Long: `
The identifier is generated from the hash of the public key, which is located in the extension header or declared in the key field of the manifest.
If the key is specified in the manifest, the public key is taken from there; otherwise, the search continues in the header.
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("infile is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			infile, err := toPath(args[0])
			if err != nil {
				return err
			}
			ext := crx3.Extension(infile)
			id, err := ext.ID()
			if err != nil {
				return err
			}
			fmt.Println(id)
			return nil
		},
	}
	return cmd
}
