package command

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	crx3 "github.com/mediabuyerbot/go-crx3"
	"github.com/spf13/cobra"
)

type encodeOpts struct {
	Outfile string
}

func (o encodeOpts) HasOutfile() bool {
	return len(o.Outfile) > 0
}

func newEncodeCmd() *cobra.Command {
	var opts encodeOpts
	cmd := &cobra.Command{
		Use:   "encode [infile]",
		Short: "Encodes file to base64 string",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("infile is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			extension := crx3.Extension(args[0])
			b, err := extension.ToBase64()
			if err != nil {
				return err
			}
			if opts.HasOutfile() {
				file, err := os.Create(opts.Outfile)
				if err != nil {
					return err
				}
				defer file.Close()
				if _, err := io.Copy(file, bytes.NewBuffer(b)); err != nil {
					return err
				}
			} else {
				fmt.Println(string(b))
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.Outfile, "outfile", "o", "", "save base64 data to file")

	return cmd
}
