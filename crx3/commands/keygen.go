package commands

import (
	"fmt"
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
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			crx3.SetDefaultKeySize(sanitizeKeySize(opts.PrivateKeySize))
			pk, err := crx3.NewPrivateKey()
			if err != nil {
				return err
			}
			if len(args) == 0 {
				key := crx3.PrivateKeyToPEM(pk)
				fmt.Print(string(key))
				return nil
			}
			filename := args[0]
			if !strings.HasSuffix(filename, ".pem") {
				filename = filename + ".pem"
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
