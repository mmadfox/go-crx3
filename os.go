package crx3

import (
	"encoding/binary"
	"fmt"
	"net/http"
	"os"
)

func isDir(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func isCRC(filename string) bool {
	return isCRX(filename)
}

func isCRX(filename string) bool {
	size := 12
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()
	buf := make([]byte, size)
	if _, err := file.Read(buf); err != nil {
		return false
	}
	if len(buf) < size {
		return false
	}
	if string(buf[0:4]) != "Cr24" {
		return false
	}
	if binary.LittleEndian.Uint32(buf[4:8]) != 3 {
		return false
	}
	return true
}

func openCrxFile(filename string) (*os.File, error) {
	if err := crxFileExists(filename); err != nil {
		return nil, err
	}
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return fd, nil
}

func crxFileExists(filename string) error {
	if !isCRX(filename) {
		return fmt.Errorf("%v got %s", ErrUnknownFileExtension, filename)
	}
	return nil
}

func isZip(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()
	buf := make([]byte, 512)
	if _, err := file.Read(buf); err != nil {
		return false
	}
	fileType := http.DetectContentType(buf)
	switch fileType {
	case "application/x-gzip", "application/zip":
		return true
	}
	return false
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) || err != nil {
		return false
	}
	return !info.IsDir()
}

func dirExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) || err != nil {
		return false
	}
	return info.IsDir()
}
