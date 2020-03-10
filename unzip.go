package crx3

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Unzip extracts all files from the archive.
func Unzip(r io.ReaderAt, size int64, unpacked string) error {
	reader, err := zip.NewReader(r, size)
	if err != nil {
		return err
	}
	if _, err := os.Stat(unpacked); os.IsNotExist(err) {
		if err := os.Mkdir(unpacked, os.ModePerm); err != nil {
			return err
		}
	}
	for _, file := range reader.File {
		fpath := filepath.Join(unpacked, file.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(unpacked)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
			continue
		}
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		rc, err := file.Open()
		if err != nil {
			return err
		}
		if _, err = io.Copy(outFile, rc); err != nil {
			return err
		}
		if err = outFile.Close(); err != nil {
			return err
		}
		if err = rc.Close(); err != nil {
			return err
		}
	}
	return nil
}
