package crx3

import (
	"bytes"
	"encoding/base64"
	"io"
)

// Base64 encodes an extension file to a base64 string.
// It returns a bytes and an error encountered while encodes, if any.
func Base64(filename string) (b []byte, err error) {
	extension, err := openCrxFile(filename)
	if err != nil {
		return nil, err
	}
	defer extension.Close()

	return encodeExtensionToBase64Str(extension)
}

func encodeExtensionToBase64Str(file io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	if _, err := io.Copy(encoder, file); err != nil {
		return nil, err
	}
	if err := encoder.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
