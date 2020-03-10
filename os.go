package crx3

import (
	"encoding/binary"
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
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
