package crx3

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/mediabuyerbot/go-crx3/pb"

	"google.golang.org/protobuf/proto"
)

// PackZipToCRX reads a ZIP archive from the provided Reader, signs it using
// the provided RSA private key, and writes the signed CRX file to the provided Writer.
// This function is essential for producing production-ready CRX files that require
// digital signatures to be installed in browsers. The function will return an error
// if any issues occur during the zip reading, signing, or CRX writing processes.
func PackZipToCRX(zip io.ReadSeeker, w io.Writer, pk *rsa.PrivateKey) error {
	if zip == nil || w == nil || pk == nil {
		return fmt.Errorf("crx3/pack: zip or writer or privateKey is nil")
	}
	publicKey, err := makePublicKey(pk)
	if err != nil {
		return fmt.Errorf("crx3/pack: failed to make public key: %w", err)
	}
	signedData, err := makeSignedData(publicKey)
	if err != nil {
		return fmt.Errorf("crx3/pack: failed to make signed data: %w", err)
	}
	signature, err := makeSign(zip, signedData, pk)
	if err != nil {
		return fmt.Errorf("crx3/pack: failed to make signature: %w", err)
	}
	header, err := makeHeader(publicKey, signature, signedData)
	if err != nil {
		return fmt.Errorf("crx3/pack: failed to make header: %w", err)
	}
	if _, err := zip.Seek(0, 0); err != nil {
		return fmt.Errorf("crx3/pack: failed to seek zip: %w", err)
	}
	if err := copyZipToCRX(w, zip, header); err != nil {
		return fmt.Errorf("crx3/pack: failed to copy zip to crx data: %w", err)
	}
	return nil
}

// WritePrivateKey writes the RSA private key to the provided io.Writer in the PEM format.
// The function expects a non-nil *rsa.PrivateKey. If the key is nil, it returns an
// ErrPrivateKeyNotFound error. This function handles the marshalling of the private key
// into PKCS#8 format and then encodes it into PEM format before writing.
//
// Parameters:
//
//	w    : An io.Writer to which the PEM encoded private key will be written.
//	key  : A non-nil *rsa.PrivateKey that will be marshalled and written.
//
// Returns:
//
//	An error if the private key is nil, if there is a marshalling error, or if writing
//	to the io.Writer fails. The error includes a descriptive message to aid in debugging.
//
// Usage example:
//
//	file, err := os.Create("private_key.pem")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer file.Close()
//
//	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	if err := WritePrivateKey(file, privKey); err != nil {
//	    log.Printf("Failed to write private key: %v", err)
//	}
//
// Note:
//
//	This function does not close the io.Writer; the caller is responsible for managing
//	the writer's lifecycle, including opening and closing it.
func WritePrivateKey(w io.Writer, key *rsa.PrivateKey) error {
	if key == nil {
		return ErrPrivateKeyNotFound
	}
	bytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return fmt.Errorf("crx3/pack: failed to marshal private key: %w", err)
	}
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: bytes,
	}
	if _, err := w.Write(pem.EncodeToMemory(block)); err != nil {
		return fmt.Errorf("crx3/pack: failed to write private key: %w", err)
	}
	return nil
}

// ReadZipFile opens and reads the contents of a zip file specified by 'filename'.
// It can handle both direct paths to zip files or directories. If 'filename' is a directory,
// the function zips its contents into a buffer and returns a reader for that buffer.
// If 'filename' is a zip file, it reads the file into a buffer and returns a reader for it.
// The function returns a *bytes.Reader to allow random access reads, which is particularly
// useful for large files. It returns an error if the file cannot be opened, read, or if the
// path does not correspond to a zip file or directory.
func ReadZipFile(filename string) (*bytes.Reader, error) {
	return readZipFile(filename)
}

func readZipFile(filename string) (*bytes.Reader, error) {
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

func copyZipToCRX(crx io.Writer, zipFile io.ReadSeeker, header []byte) error {
	if _, err := crx.Write([]byte("Cr24")); err != nil {
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
