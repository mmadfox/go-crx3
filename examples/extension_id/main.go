package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mediabuyerbot/go-crx3"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	filename := filepath.Join(pwd, "examples", "extension_id", "extension.crx")
	ext := crx3.Extension(filename)

	id, err := ext.ID()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Extension ID: %s\n", id)
}
