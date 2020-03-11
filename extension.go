package crx3

import (
	"crypto/rsa"
	"os"
	"path"
	"strings"
)

type Extension string

func (e Extension) String() string {
	return string(e)
}

// Zip creates a *.zip archive and adds all the files to it.
func (e Extension) Zip() error {
	filename := path.Join(e.String(), zipExt)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return Zip(file, e.String())
}

// Unzip extracts all files from the archive.
func (e Extension) Unzip() error {
	file, err := os.Open(e.String())
	if err != nil {
		return err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	unpacked := strings.TrimRight(e.String(), zipExt)
	return Unzip(file, stat.Size(), unpacked)
}

// ToBase64 encodes file to base64.
func (e Extension) ToBase64() ([]byte, error) {
	return ToBase64(e.String())
}

// Unpack unpacks CRX3 extension to directory.
func (e Extension) Unpack() error {
	return Unpack(e.String())
}

// PackTo packs zip file or an unpacked directory into a CRX3 file.
func (e Extension) PackTo(dst string, pk *rsa.PrivateKey) error {
	return Pack(e.String(), dst, pk)
}

// Pack packs zip file or an unpacked directory into a CRX3 file.
func (e Extension) Pack(pk *rsa.PrivateKey) error {
	return Pack(e.String(), "", pk)
}
