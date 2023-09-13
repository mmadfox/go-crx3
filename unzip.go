package crx3

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// UnzipTo extracts the contents of a ZIP archive specified by 'filename' to the 'basepath' directory.
// It opens the ZIP file, creates the necessary directory structure, and extracts all files.
func UnzipTo(basepath string, filename string) error {
	if !isDir(basepath) {
		return fmt.Errorf("%w: does not exists %s",
			ErrPathNotFound, basepath)
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	dirname := filepath.Join(basepath,
		strings.TrimSuffix(filepath.Base(filename), zipExt))

	return Unzip(file, stat.Size(), dirname)
}

// Unzip extracts all files and directories from the provided ZIP archive.
// It takes an io.ReaderAt 'r', the size 'size' of the ZIP archive, and the target directory 'unpacked' for extraction.
// It iterates through the archive, creating directories and writing files as necessary.
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
