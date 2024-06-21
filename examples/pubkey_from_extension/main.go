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

	rootPath := filepath.Join(pwd, "examples", "pubkey_from_extension")

	crx3.SetDefaultKeySize(2048)

	// 2. Initialize extension
	filename := filepath.Join(rootPath, "extension.crx")
	ext := crx3.Extension(filename)
	pubKey, _, err := ext.PublicKey()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Public key: %s\n", pubKey)
}
