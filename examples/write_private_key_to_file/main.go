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

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	pf := filepath.Join(pwd, "examples", "write_private_key_to_file", "private.pem")
	file, err := os.OpenFile(pf, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 2. Save the private key to a file
	if err := crx3.WritePrivateKey(file, pk); err != nil {
		panic(err)
	}

	fmt.Printf("Private key has been written to the file: %s\n", file.Name())
}
