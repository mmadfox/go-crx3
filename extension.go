package crx3

import (
	"crypto/rsa"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Extension represents an extension for google chrome.
type Extension string

// String returns a string representation.
func (e Extension) String() string {
	return string(e)
}

// ID returns an extension id.
func (e Extension) ID() (string, error) {
	return ID(e.String())
}

// IsEmpty checks if the extension is empty.
func (e Extension) IsEmpty() bool {
	return len(e.String()) == 0
}

// IsDir reports whether extension describes a directory.
func (e Extension) IsDir() bool {
	return isDir(e.String())
}

// IsZip reports whether extension describes a zip-archive.
func (e Extension) IsZip() bool {
	return isZip(e.String())
}

// IsCRX3 reports whether extension describes a crc file.
func (e Extension) IsCRX3() bool {
	return isCRC(e.String())
}

// Zip creates a *.zip archive and adds all the files to it.
func (e Extension) Zip() error {
	if e.IsEmpty() {
		return fmt.Errorf("%w: %s", ErrPathNotFound, e)
	}

	filename := strings.TrimRight(e.String(), "/") + zipExt
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return Zip(file, e.String())
}

// Unzip extracts all files from the archive.
func (e Extension) Unzip() error {
	if e.IsEmpty() {
		return fmt.Errorf("%w: %s", ErrPathNotFound, e)
	}

	file, err := os.Open(e.String())
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	unpacked := strings.TrimSuffix(e.String(), zipExt)
	if dirExists(unpacked) {
		index := 1
		for {
			if index >= 100 {
				break
			}
			unpacked = unpacked + "(" + strconv.Itoa(index) + ")"
			if !dirExists(unpacked) {
				break
			}
			index++
		}
	}

	return Unzip(file, stat.Size(), unpacked)
}

// Base64 encodes an extension file to a base64 string.
func (e Extension) Base64() ([]byte, error) {
	if e.IsEmpty() {
		return nil, fmt.Errorf("%w: %s", ErrPathNotFound, e)
	}
	return Base64(e.String())
}

// Unpack unpacks the CRX3 extension into a directory.
func (e Extension) Unpack() error {
	if e.IsEmpty() {
		return fmt.Errorf("%w: %s", ErrPathNotFound, e)
	}
	return Unpack(e.String())
}

// PackTo packs zip file or an unpacked directory into a CRX3 file.
func (e Extension) PackTo(dst string, pk *rsa.PrivateKey) error {
	if e.IsEmpty() {
		return ErrPathNotFound
	}
	return Pack(e.String(), dst, pk)
}

// Pack packs zip file or an unpacked directory into a CRX3 file.
func (e Extension) Pack(pk *rsa.PrivateKey) error {
	if e.IsEmpty() {
		return ErrPathNotFound
	}
	dst := strings.TrimRight(e.String(), "/") + crxExt
	return Pack(e.String(), dst, pk)
}
