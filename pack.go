package crx3

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"io"
	"os"
	"path"
	"strings"

	"github.com/mediabuyerbot/go-crx3/pb"

	"google.golang.org/protobuf/proto"
)

const (
	crxExt = ".crx"
	zipExt = ".zip"
	pemExt = ".pem"
)

// Pack packs a zip file or unzipped directory into a crx extension.
// It takes the source 'src' (zip file or directory), target 'dst' CRX file path,
// and a private key 'pk' (optional). If 'pk' is nil, it generates a new private key.
// It creates a CRX extension from the source and writes it to the destination.
func Pack(src string, dst string, pk *rsa.PrivateKey) (err error) {
	var (
		publicKey      []byte
		signedData     []byte
		signature      []byte
		header         []byte
		hasDst         = len(dst) > 0
		isDefaultPk    bool
		isNotCrxSuffix = path.Ext(dst) != crxExt
	)

	if len(src) == 0 || len(dst) == 0 {
		return ErrPathNotFound
	}

	if hasDst && isNotCrxSuffix {
		return ErrUnknownFileExtension
	}

	zipData, err := readZipFile(src)
	if err != nil {
		return err
	}

	// make default private key
	if pk == nil {
		pk, err = NewPrivateKey()
		if err != nil {
			return err
		}
		isDefaultPk = true
	}

	if publicKey, err = makePublicKey(pk); err != nil {
		return err
	}
	if signedData, err = makeSignedData(publicKey); err != nil {
		return err
	}
	if signature, err = makeSign(zipData, signedData, pk); err != nil {
		return err
	}
	if header, err = makeHeader(publicKey, signature, signedData); err != nil {
		return err
	}
	if _, err := zipData.Seek(0, 0); err != nil {
		return err
	}

	if !hasDst {
		crxFilename := strings.TrimSuffix(src, zipExt)
		crxFilename = crxFilename + crxExt
		dst = crxFilename
	}
	if err := writeToCRX(dst, zipData, header); err != nil {
		return err
	}
	if isDefaultPk {
		if err := saveDefaultPrivateKey(dst, pk); err != nil {
			return err
		}
	}
	return nil
}

func readZipFile(filename string) (data io.ReadSeeker, err error) {
	var zipData bytes.Buffer

	switch {
	case isDir(filename):
		if err := Zip(&zipData, filename); err != nil {
			return nil, err
		}
	case isZip(filename):
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		if _, err := io.Copy(&zipData, file); err != nil {
			return nil, err
		}
	default:
		return nil, ErrUnknownFileExtension
	}

	return bytes.NewReader(zipData.Bytes()), nil
}

func writeToCRX(filename string, zipFile io.ReadSeeker, header []byte) error {
	crx, err := os.Create(filename)
	if err != nil {
		return err
	}
	if _, err = crx.Write([]byte("Cr24")); err != nil {
		return err
	}
	if err := binary.Write(crx, binary.LittleEndian, uint32(3)); err != nil {
		return err
	}
	if err := binary.Write(crx, binary.LittleEndian, uint32(len(header))); err != nil {
		return err
	}
	if _, err := crx.Write(header); err != nil {
		return err
	}
	if _, err := io.Copy(crx, zipFile); err != nil {
		return err
	}
	return nil
}

func makeCRXID(publicKey []byte) []byte {
	hash := sha256.New()
	hash.Write(publicKey)
	return hash.Sum(nil)[0:16]
}

func makePublicKey(pk *rsa.PrivateKey) ([]byte, error) {
	return x509.MarshalPKIXPublicKey(&pk.PublicKey)
}

func makeSignedData(publicKey []byte) ([]byte, error) {
	signedData := &pb.SignedData{
		CrxId: makeCRXID(publicKey),
	}
	return proto.Marshal(signedData)
}

func makeSign(r io.Reader, signedData []byte, pk *rsa.PrivateKey) ([]byte, error) {
	sign := sha256.New()
	sign.Write([]byte("CRX3 SignedData\x00"))
	if err := binary.Write(sign, binary.LittleEndian, uint32(len(signedData))); err != nil {
		return nil, err
	}
	sign.Write(signedData)
	if _, err := io.Copy(sign, r); err != nil {
		return nil, err
	}
	return rsa.SignPKCS1v15(rand.Reader, pk, crypto.SHA256, sign.Sum(nil))
}

func makeHeader(pubKey, signature, signedData []byte) ([]byte, error) {
	header := &pb.CrxFileHeader{
		Sha256WithRsa: []*pb.AsymmetricKeyProof{
			{
				PublicKey: pubKey,
				Signature: signature,
			},
		},
		SignedHeaderData: signedData,
	}
	return proto.Marshal(header)
}

func saveDefaultPrivateKey(filename string, pk *rsa.PrivateKey) error {
	pemFilename := strings.TrimSuffix(filename, zipExt)
	pemFilename = pemFilename + pemExt
	return SavePrivateKey(pemFilename, pk)
}
