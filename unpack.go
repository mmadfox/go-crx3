package crx3

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"

	"github.com/mediabuyerbot/go-crx3/pb"

	"google.golang.org/protobuf/proto"
)

// UnpackTo unpacks a CRX (Chrome Extension) file specified by 'filename' to the directory 'dirname'.
// If 'dirname' does not exist, it creates the directory before unpacking.
func UnpackTo(filename string, dirname string) error {
	if !isDir(dirname) {
		if err := os.Mkdir(dirname, 0755); err != nil {
			return err
		}
	}

	return unpack(filename, dirname)
}

// Unpack unpacks a CRX (Chrome Extension) file specified by 'filename' to its original contents.
// It checks if the file is in the CRX format, reads its header and signed data,
// and then extracts and decompresses the original contents.
// The unpacked contents are placed in a directory with the same name as the original file (without the '.crx' extension).
func Unpack(filename string) error {
	return unpack(filename, "")
}

func unpack(filename string, dirname string) error {
	// check if the file is in the CRX format.
	if len(filename) == 0 || !isCRX(filename) {
		return ErrUnsupportedFileFormat
	}

	// read the entire CRX file into memory.
	crx, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// extract header and signed data from the CRX file.
	var (
		headerSize = binary.LittleEndian.Uint32(crx[8:12])
		metaSize   = uint32(12)
		v          = crx[metaSize : headerSize+metaSize]
		header     pb.CrxFileHeader
		signedData pb.SignedData
	)

	// unmarshal the header data.
	if err = proto.Unmarshal(v, &header); err != nil {
		return err
	}

	// unmarshal the header data.
	if err = proto.Unmarshal(header.SignedHeaderData, &signedData); err != nil {
		return err
	}

	// check if the CRX ID has the expected length.
	if len(signedData.CrxId) != 16 {
		return ErrUnsupportedFileFormat
	}

	data := crx[len(v)+int(metaSize):]
	reader := bytes.NewReader(data)
	size := int64(len(data))

	var unpacked string
	if len(dirname) > 0 {
		fn := filepath.Base(filename)
		unpacked = filepath.Join(dirname, strings.TrimSuffix(fn, crxExt))
	} else {
		unpacked = strings.TrimSuffix(filename, crxExt)
	}

	return Unzip(reader, size, unpacked)
}
