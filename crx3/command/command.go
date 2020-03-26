package command

import (
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "crx3",
		Short: "Chrome extensions tools",
	}

	cmd.AddCommand(newPackCmd())
	cmd.AddCommand(newUnpackCmd())
	cmd.AddCommand(newZipCmd())
	cmd.AddCommand(newUnzipCmd())
	cmd.AddCommand(newEncodeCmd())
	cmd.AddCommand(newKeygenCmd())
	cmd.AddCommand(newDownloadCmd())
	cmd.AddCommand(newIDCmd())

	return cmd
}
