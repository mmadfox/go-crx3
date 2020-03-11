package command

import (
	"crypto/rsa"
	"errors"

	crx3 "github.com/mediabuyerbot/go-crx3"
	"github.com/spf13/cobra"
)

type packOpts struct {
	PrivateKey string
	Outfile    string
}

func (o packOpts) HasPem() bool {
	return len(o.PrivateKey) > 0
}

func newPackCmd() *cobra.Command {
	var opts packOpts
	cmd := &cobra.Command{
		Use:   "pack [infile]",
		Short: "Packs zip file or an unpacked directory into a CRX3 file.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("infile is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			unpacked := args[0]
			var pk *rsa.PrivateKey
			if opts.HasPem() {
				pk, err = crx3.LoadPrivateKey(opts.PrivateKey)
				if err != nil {
					return err
				}
			}
			return crx3.Pack(unpacked, opts.Outfile, pk)
		},
	}

	cmd.Flags().StringVarP(&opts.PrivateKey, "pem", "p", "", "extension key filename (*.pem)")
	cmd.Flags().StringVarP(&opts.Outfile, "outfile", "o", "", "extension filename (*.crx)")

	return cmd
}
