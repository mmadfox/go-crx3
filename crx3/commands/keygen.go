package commands

import (
	"fmt"
	"strings"

	crx3 "github.com/mediabuyerbot/go-crx3"
	"github.com/spf13/cobra"
)

const pemExt = ".pem"

type keygenOpts struct {
	PrivateKeySize int
}

func newKeygenCmd() *cobra.Command {
	var opts keygenOpts
	cmd := &cobra.Command{
		Use:   "keygen [file]",
		Short: "Create a new private key",
		Long: `
If no file is specified, the private key is printed to stdout. 
Otherwise, the private key is saved to the specified file. 
If the file does not have a .pem extension, it is added automatically.
Size of the private key can be set with the --size or -s flag. Sizes of 2048, 3072, or 4096 bits are allowed.
		`,
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
			filename, err := toPath(args[0])
			if err != nil {
				return fmt.Errorf("invalid infile path: %w", err)
			}
			if !strings.HasSuffix(filename, pemExt) {
				filename = filename + pemExt
			}
			return crx3.SavePrivateKey(filename, pk)
		},
	}

	cmd.Flags().IntVarP(&opts.PrivateKeySize, "size", "s", 2048, "private key size")

	return cmd
}
