package commands

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
	cmd.AddCommand(newBase64Cmd())
	cmd.AddCommand(newKeygenCmd())
	cmd.AddCommand(newDownloadCmd())
	cmd.AddCommand(newIDCmd())
	cmd.AddCommand(newPubkeyCmd())

	return cmd
}
