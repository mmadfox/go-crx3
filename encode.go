package crx3

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

var ErrUnknownFileExtension = errors.New("crx3: unknown file extension")

// ToBase64 encodes file to base64.
func ToBase64(filename string) (b []byte, err error) {
	if !isCRC(filename) {
		return nil, ErrUnknownFileExtension
	}
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	if _, err := io.Copy(encoder, bufio.NewReader(fd)); err != nil {
		return nil, err
	}
	if err := encoder.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
