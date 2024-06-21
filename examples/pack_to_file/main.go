package main

import (
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

	// 2. Initialize extension
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// 3. Initialize extension
	filename := filepath.Join(pwd, "examples", "pack_to_file", "bobbyMol.zip")
	fmt.Println("Extension file:", filename)

	// 4. Pack zip file or an unpacked directory into a CRX3 file
	ext := crx3.Extension(filename)
	if err := ext.Pack(pk); err != nil {
		panic(err)
	}

	// 5. Write private key to a file
	if err := crx3.SavePrivateKey(filename+".pem", pk); err != nil {
		panic(err)
	}

	fmt.Printf("CRX file has been written to the file: %s\n", filename+".crx")
}
