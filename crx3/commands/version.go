package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of the crx3 tools",
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Printf("crx3 tools version %s\n", version)
		},
	}
	return cmd
}
