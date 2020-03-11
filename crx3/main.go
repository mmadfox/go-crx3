package main

import (
	"fmt"
	"os"

	"github.com/mediabuyerbot/go-crx3/crx3/command"
)

func main() {
	cli := command.New()
	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
