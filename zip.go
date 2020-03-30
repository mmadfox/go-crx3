package crx3

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// Zip creates a *.zip archive and adds all the files to it.
func Zip(w io.Writer, unpacked string) error {
	if !isDir(unpacked) {
		return ErrPathNotFound
	}
	wz := zip.NewWriter(w)
	defer wz.Close()
	return filepath.Walk(unpacked, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relpath, err := filepath.Rel(unpacked, path)
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
