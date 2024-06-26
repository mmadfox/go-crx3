package crx3

import (
	"bytes"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/mediabuyerbot/go-crx3/pb"

	"google.golang.org/protobuf/proto"
)

const symbols = "abcdefghijklmnopqrstuvwxyz"

// ID returns the extension ID extracted from a CRX (Chrome Extension) file specified by 'filename'.
// It checks if the file is in the CRX format, reads its header and signed data,
// and then converts the CRX ID into a string representation
func ID(filename string) (id string, err error) {
	if !isCRX(filename) {
		return id, ErrUnsupportedFileFormat
	}

	crx, err := os.ReadFile(filename)
	if err != nil {
		return id, err
	}

	var (
		headerSize = binary.LittleEndian.Uint32(crx[8:12])
		metaSize   = uint32(12)
		v          = crx[metaSize : headerSize+metaSize]
		header     pb.CrxFileHeader
		signedData pb.SignedData
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

// IDFromPubKey generates the Chrome Extension ID from a public key.
// It handles PEM formatting, base64 decoding, and SHA-256 hashing to produce the ID.
// Returns the ID or an error if the key processing fails.
func IDFromPubKey(pubKey []byte) (string, error) {
	if len(pubKey) < 64 {
		return "", fmt.Errorf("crx3/id: public key is empty")
	}

	if bytes.Contains(pubKey[:64], []byte("PUBLIC KEY")) {
		pubKey = formatPemKey(pubKey)
	}

	dbuf := make([]byte, base64.StdEncoding.DecodedLen(len(pubKey)))
	n, err := base64.StdEncoding.Decode(dbuf, pubKey)
	if err != nil {
		return "", fmt.Errorf("crx3/id: failed to decode public key: %w", err)
	}

	pubKeyParsed, err := x509.ParsePKIXPublicKey(dbuf[:n])
	if err != nil {
		return "", fmt.Errorf("crx3/id: failed to parse public key: %w", err)
	}
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKeyParsed)
	if err != nil {
		return "", fmt.Errorf("crx3/id: failed to marshal public key: %w", err)
	}
	hash := sha256.New()
	hash.Write(pubKeyBytes)
	digest := hash.Sum(nil)

	idx := strIDx()
	sid := fmt.Sprintf("%x", digest[:16])
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
