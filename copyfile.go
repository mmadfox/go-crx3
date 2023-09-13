package crx3

import (
	"fmt"
	"io"
	"os"
)

// CopyFile copies the contents of the source file 'src' to the destination file 'dst'.
// It returns the number of bytes copied and any encountered error.
// If the source file does not exist, is not a regular file, or other errors occur during the copy,
// it returns an error with a descriptive message.
func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)

	return nBytes, err
}
