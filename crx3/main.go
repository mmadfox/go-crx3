package main

import (
	"fmt"
	"os"

	"github.com/mediabuyerbot/go-crx3/crx3/commands"
)

// TODO: add to ci VERSION=$(git describe --tags --always) ... -ldflags "-X main.Version=$VERSION"
var Version = "v1.6.0"

func main() {
	cli := commands.New(Version)
	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
