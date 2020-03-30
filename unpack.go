package crx3

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"strings"

	"github.com/mediabuyerbot/go-crx3/pb"

	"github.com/golang/protobuf/proto"
)

// Unpack unpacks chrome extension into some directory.
func Unpack(filename string) error {
	if !isCRX(filename) {
		return ErrUnsupportedFileFormat
	}
	crx, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	var (
		headerSize = binary.LittleEndian.Uint32(crx[8:12])
		metaSize   = uint32(12)
		v          = crx[metaSize : headerSize+metaSize]
		header     pb.CrxFileHeader
		signedData pb.SignedData
	)

	if err := proto.Unmarshal(v, &header); err != nil {
		return err
	}
	if err := proto.Unmarshal(header.SignedHeaderData, &signedData); err != nil {
		return err
	}

	if len(signedData.CrxId) != 16 {
		return ErrUnsupportedFileFormat
	}

	data := crx[len(v)+int(metaSize):]
	reader := bytes.NewReader(data)
	size := int64(len(data))

	unpacked := strings.TrimRight(filename, crxExt)
	return Unzip(reader, size, unpacked)
}
