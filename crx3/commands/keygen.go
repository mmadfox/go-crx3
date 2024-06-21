package commands

import (
	"errors"
	"strings"

	crx3 "github.com/mediabuyerbot/go-crx3"
	"github.com/spf13/cobra"
)

type keygenOpts struct {
	PrivateKeySize int
}

func newKeygenCmd() *cobra.Command {
	var opts keygenOpts
	cmd := &cobra.Command{
		Use:   "keygen [file]",
		Short: "Create a new private key",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("infile is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			crx3.SetDefaultKeySize(sanitizeKeySize(opts.PrivateKeySize))
			filename := args[0]
			if !strings.HasSuffix(filename, ".pem") {
				filename = filename + ".pem"
			}
			pk, err := crx3.NewPrivateKey()
			if err != nil {
				return err
			}
			return crx3.SavePrivateKey(filename, pk)
		},
	}

	cmd.Flags().IntVarP(&opts.PrivateKeySize, "size", "s", 4096, "private key size")

	return cmd
}

func sanitizeKeySize(size int) int {
	switch size {
	case 2048, 3072, 4096:
		return size
	default:
		return 4096
	}
}
