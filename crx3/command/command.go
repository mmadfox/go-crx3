package command

import (
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "crx3",
		Short: "Provides tools for working with the crx3 extension.",
	}

	cmd.AddCommand(newPackCmd())
	cmd.AddCommand(newUnpackCmd())
	cmd.AddCommand(newZipCmd())
	cmd.AddCommand(newUnzipCmd())
	cmd.AddCommand(newEncodeCmd())
	cmd.AddCommand(newKeygenCmd())
	cmd.AddCommand(newDownloadCmd())

	return cmd
}
