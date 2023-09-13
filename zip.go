package crx3

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ZipTo creates a ZIP archive with the specified
// filename and adds all files from the given directory to it.
func ZipTo(filename string, dirname string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return Zip(file, dirname)
}

// Zip creates a *.zip archive and adds all files
// from the specified directory to it.
func Zip(dst io.Writer, dirname string) error {
	if !isDir(dirname) {
		return fmt.Errorf("%w: %s", ErrPathNotFound, dirname)
	}
	wz := zip.NewWriter(dst)
	defer wz.Close()
	return filepath.Walk(dirname,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			relpath, err := filepath.Rel(dirname, path)
			if err != nil {
				return err
			}
			return writeToZip(wz, path, relpath)
		})
}

func writeToZip(w *zip.Writer, filename string, metaname string) error {
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fd.Close()

	info, err := fd.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = metaname
	header.Method = zip.Deflate
	writer, err := w.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fd)
	return err
}
