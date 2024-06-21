package crx3

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"io"

	"github.com/mediabuyerbot/go-crx3/pb"
	"google.golang.org/protobuf/proto"
)

// PubkeyFrom extracts the RSA public key and its signature from a
// CRX3 file read from the given Reader.
//
// Parameters:
// - r io.Reader: the source from which the CRX3 content is read.
//
// Returns:
// - []byte: PEM-formatted public key.
// - []byte: Signature of the public key.
// - error: Error detailing any issues encountered.
func PubkeyFrom(r io.Reader) ([]byte, []byte, error) {
	if r == nil {
		return nil, nil, ErrInvalidReader
	}
	ext, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, fmt.Errorf("crx3/pubkey: failed to read the public key: %w", err)
	}
	if len(ext) < 12 {
		return nil, nil, fmt.Errorf("crx3/pubkey: invalid extension size: %d", len(ext))
	}
	var (
		headerSize = binary.LittleEndian.Uint32(ext[8:12])
		metaSize   = uint32(12)
		v          = ext[metaSize : headerSize+metaSize]
		header     pb.CrxFileHeader
	)
	if err := proto.Unmarshal(v, &header); err != nil {
		return nil, nil, err
	}
	if len(header.Sha256WithRsa) == 0 {
		return nil, nil, fmt.Errorf("crx3/pubkey: missing sha256 with rsa signature")
	}
	publicKey, err := publicKeyToPEM(header.Sha256WithRsa[0].PublicKey)
	if err != nil {
		return nil, nil, err
	}
	return publicKey, header.Sha256WithRsa[0].Signature, nil
}

func publicKeyToPEM(pubKeyBytes []byte) ([]byte, error) {
	pubKey, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("crx3/pubkey: failed to parse public key: %w", err)
	}
	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("crx3/pubkey: public key is not of type *rsa.PublicKey")
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(rsaPubKey),
	})
	return pemBytes, nil
}
