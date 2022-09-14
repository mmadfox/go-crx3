package commands

import (
	"errors"
	"os"
	"path"
	"strings"

	crx3 "github.com/mediabuyerbot/go-crx3"
	"github.com/spf13/cobra"
)

type downloadOpts struct {
	Outfile string
	Unpack  bool
}

func (o downloadOpts) HasNotOutfile() bool {
	return len(o.Outfile) == 0
}

func newDownloadCmd() *cobra.Command {
	var opts downloadOpts
	cmd := &cobra.Command{
		Use:   "download [id or url]",
		Short: "Download the chrome extension from the web store",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("extension is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			extensionID := args[0]
			if strings.HasPrefix(extensionID, "http") {
				extensionID = extractExtensionID(extensionID)
			}
			if opts.HasNotOutfile() {
				pwd, err := os.Getwd()
				if err != nil {
					return err
				}
				opts.Outfile = path.Join(pwd, "extension.crx")
			}

			if !strings.HasSuffix(opts.Outfile, ".crx") {
				opts.Outfile = opts.Outfile + ".crx"
			}

			if err := crx3.DownloadFromWebStore(extensionID, opts.Outfile); err != nil {
				return err
			}
			if opts.Unpack {
				if err := crx3.Unpack(opts.Outfile); err != nil {
					return err
				}
				if err := os.Remove(opts.Outfile); err != nil {
					return err
				}
			}
			return
		},
	}

	cmd.Flags().StringVarP(&opts.Outfile, "outfile", "o", "", "save to file")
	cmd.Flags().BoolVarP(&opts.Unpack, "unpack", "u", true, "unpack")

	return cmd
}

func extractExtensionID(u string) string {
	urlParts := strings.Split(u, "/")
	if len(urlParts) == 0 {
		return ""
	}
	return urlParts[len(urlParts)-1]
}
