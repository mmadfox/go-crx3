package commands

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	crx3 "github.com/mediabuyerbot/go-crx3"
	"github.com/spf13/cobra"
)

type pubkeyOpts struct {
	Outfile string
}

func newPubkeyCmd() *cobra.Command {
	var opts pubkeyOpts
	cmd := &cobra.Command{
		Use:   "pubkey [extension]",
		Short: "Retrieves the public key from the extension header or manifest file",
		Long: `
The public key is searched for first in the manifest file, and if not found, the search continues in the extension header.
		`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("file is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			infile, err := toPath(args[0])
			if err != nil {
				return fmt.Errorf("failed to open infile: %w", err)
			}
			outfile, err := toPath(opts.Outfile)
			if err != nil {
				return fmt.Errorf("failed to open outfile: %w", err)
			}
			ext := crx3.Extension(infile)
			pubkey, _, err := ext.PublicKey()
			if err != nil {
				return err
			}
			if len(outfile) == 0 {
				fmt.Print(string(pubkey))
				return nil
			}
			file, err := os.OpenFile(outfile, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			defer file.Close()
			io.Copy(file, bytes.NewReader(pubkey))
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.Outfile, "outfile", "o", "", "save to file")

	return cmd
}
