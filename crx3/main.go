package main

import (
	"fmt"
	"os"

	"github.com/mediabuyerbot/go-crx3/crx3/commands"
)

var Version = "v1.7.0"

func main() {
	cli := commands.New(Version)
	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
