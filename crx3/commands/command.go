package commands

import (
	"os/user"
	"path/filepath"
	"strings"

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

func toPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		home := usr.HomeDir
		if path == "~" {
			return home, nil
		}
		return filepath.Join(home, path[2:]), nil
	}
	return path, nil
}

func sanitizeKeySize(size int) int {
	switch size {
	case 2048, 3072, 4096:
		return size
	default:
		return 2048
	}
}
