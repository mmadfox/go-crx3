package crx3

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
)

const symbols = "abcdefghijklmnopqrstuvwxyz"

// ID returns an extension id.
func ID(filename string) (id string, err error) {
	if !isCRC(filename) {
		return id, ErrUnsupportedFileFormat
	}

	crx, err := ioutil.ReadFile(filename)
	if err != nil {
		return id, err
	}

	var (
		headerSize = binary.LittleEndian.Uint32(crx[8:12])
		metaSize   = uint32(12)
		v          = crx[metaSize : headerSize+metaSize]
		header     CrxFileHeader
		signedData SignedData
	)

	if err := proto.Unmarshal(v, &header); err != nil {
		return id, err
	}
	if err := proto.Unmarshal(header.SignedHeaderData, &signedData); err != nil {
		return id, err
	}

	idx := strIDx()
	sid := fmt.Sprintf("%x", signedData.CrxId[:16])
	buf := bytes.NewBuffer(nil)
	for _, char := range sid {
		index := idx[char]
		buf.WriteString(string(symbols[index]))
	}
	return buf.String(), nil
}

func strIDx() map[rune]int {
	index := make(map[rune]int)
	src := "0123456789abcdef"
	for i, char := range src {
		index[char] = i
	}
	return index
}
