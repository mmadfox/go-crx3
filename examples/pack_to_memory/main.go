package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mediabuyerbot/go-crx3"
)

func main() {
	// 1. Generate a new private key OR load from file
	pk, err := crx3.NewPrivateKey()
	if err != nil {
		panic(err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// 2. Initialize extension
	filename := filepath.Join(pwd, "examples", "pack_to_memory", "bobbyMol.zip")
	fmt.Println("Extension file:", filename)

	ext := crx3.Extension(filename)
	buf := new(bytes.Buffer)

	// 3. Pack zip file or an unpacked directory into a CRX3 file
	if err := ext.WriteTo(buf, pk); err != nil {
		panic(err)
	}

	fmt.Printf("CRX file has been written to the buffer\n")
	fmt.Printf("CRX file size: %d bytes\n", buf.Len())
}
