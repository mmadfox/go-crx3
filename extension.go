package crx3

import (
	"archive/zip"
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Extension represents an extension for google chrome.
type Extension string

// String returns a string representation.
func (e Extension) String() string {
	return string(e)
}

// ID calculates the Chrome Extension ID for the Extension instance.
// It supports directories, ZIP archives, and CRX3 files. If the extension is unpacked,
// contained in a ZIP archive, or is a CRX3 file with a specified key in its manifest,
// the ID is generated from this key. The function returns an error if the extension is empty,
// the file cannot be read, the key is not found, or the file format is unsupported.
func (e Extension) ID() (string, error) {
	if e.IsEmpty() {
		return "", fmt.Errorf("%w: %s", ErrPathNotFound, e)
	}
	switch {
	case e.IsDir():
		manifest := manifestFile(e.String())
		file, err := os.ReadFile(manifest)
		if err != nil {
			return "", fmt.Errorf("crx3: failed to read file %s: %w", manifest, err)
		}
		pubkey := parseKeyFromManifest(file)
		if len(pubkey) == 0 {
			return "", fmt.Errorf("crx3: failed to parse key from manifest file %s", manifest)
		}
		return IDFromPubKey([]byte(pubkey))
	case e.IsZip():
		pubkey, err := parseManifestFromZip(e.String())
		if err != nil {
			return "", err
		}
		if len(pubkey) == 0 {
			return "", fmt.Errorf("crx3: failed to parse key from manifest file %s", e)
		}
		return IDFromPubKey([]byte(pubkey))
	case e.IsCRX3():
		pubkey, err := parseManifestFromZip(e.String())
		if err != nil || len(pubkey) == 0 {
			return ID(e.String())
		}
		if len(pubkey) > 0 {
			return IDFromPubKey([]byte(pubkey))
		}
	}
	return "", fmt.Errorf("%w: %s", ErrUnknownFileExtension, e)
}

// IsEmpty checks if the extension is empty.
func (e Extension) IsEmpty() bool {
	return len(e.String()) == 0
}

// IsDir reports whether extension describes a directory.
func (e Extension) IsDir() bool {
	return isDir(e.String())
}

// IsZip reports whether extension describes a zip-archive.
func (e Extension) IsZip() bool {
	return isZip(e.String())
}

// IsCRX3 reports whether extension describes a crx file.
func (e Extension) IsCRX3() bool {
	return isCRX(e.String())
}

// Zip creates a *.zip archive and adds all the files to it.
func (e Extension) Zip() error {
	if e.IsEmpty() {
		return fmt.Errorf("%w: %s", ErrPathNotFound, e)
	}

	filename := strings.TrimRight(e.String(), "/") + zipExt
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return Zip(file, e.String())
}

// Unzip extracts all files from the archive.
func (e Extension) Unzip() error {
	if e.IsEmpty() {
		return fmt.Errorf("%w: %s", ErrPathNotFound, e)
	}

	file, err := os.Open(e.String())
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	unpacked := strings.TrimSuffix(e.String(), zipExt)
	if dirExists(unpacked) {
		index := 1
		for {
			if index >= 100 {
				break
			}
			unpacked = unpacked + "(" + strconv.Itoa(index) + ")"
			if !dirExists(unpacked) {
				break
			}
			index++
		}
	}

	return Unzip(file, stat.Size(), unpacked)
}

// Base64 encodes an extension file to a base64 string.
func (e Extension) Base64() ([]byte, error) {
	if e.IsEmpty() {
		return nil, fmt.Errorf("%w: %s", ErrPathNotFound, e)
	}
	return Base64(e.String())
}

// Unpack unpacks the CRX3 extension into a directory.
func (e Extension) Unpack() error {
	if e.IsEmpty() {
		return fmt.Errorf("%w: %s", ErrPathNotFound, e)
	}
	return Unpack(e.String())
}

// PackTo packs zip file or an unpacked directory into a CRX3 file.
func (e Extension) PackTo(dst string, pk *rsa.PrivateKey) error {
	if e.IsEmpty() {
		return ErrPathNotFound
	}
	return Pack(e.String(), dst, pk)
}

// Pack packs zip file or an unpacked directory into a CRX3 file.
func (e Extension) Pack(pk *rsa.PrivateKey) error {
	if e.IsEmpty() {
		return ErrPathNotFound
	}
	dst := strings.TrimRight(e.String(), "/") + crxExt
	return Pack(e.String(), dst, pk)
}

// WriteTo packs the contents of the Extension into a CRX file and writes it to the provided io.Writer.
// This method requires a non-nil *rsa.PrivateKey to sign the CRX package. The Extension must not be empty,
// and its associated zip file must be readable and correctly formatted.
//
// Parameters:
//
//	w  - The io.Writer where the CRX file will be written.
//	pk - The RSA private key used for signing the CRX file.
//
// Returns:
//
//	An error if the Extension is empty, if the private key is nil, if there are issues reading the
//	zip file associated with the Extension, or if there is a failure during the packing process.
//	Errors are wrapped with context to provide more details about the failure.
//
// Usage example:
//
//	ext := Extension("path/to/your/extension/folder") // OR zip file
//	pk, err := rsa.GenerateKey(rand.Reader, 4096)
//	if err != nil {
//	    log.Fatalf("Failed to generate private key: %v", err)
//	}
//
//	var buf bytes.Buffer
//	if err := ext.WriteTo(&buf, pk); err != nil {
//	    log.Printf("Failed to write CRX: %v", err)
//	} else {
//	    // Use buf to save CRX to a file or further processing
//	}
func (e Extension) WriteTo(w io.Writer, pk *rsa.PrivateKey) error {
	if e.IsEmpty() {
		return fmt.Errorf("%w: %s", ErrPathNotFound, e)
	}
	if pk == nil {
		return fmt.Errorf("%w: for extension %s", ErrPrivateKeyNotFound, e)
	}
	reader, err := readZipFile(e.String())
	if err != nil {
		return fmt.Errorf("crx3: failed to read zip file: %w", err)
	}
	return PackZipToCRX(reader, w, pk)
}

// PublicKey extracts the public key based on the extension type.
//   - For directories, it reads the manifest file
//     extracting the 'key' field which contains the public key.
//   - For ZIP archives, it first attempts to parse the manifest file within the ZIP to find the public key.
//     If the manifest is not present, or does not contain the key, it then tries to extract the public key
//     directly from the CRX3 header.
//   - For CRX3 files, it directly reads and parses the public key from the header if available.
//
// Returns the public key in byte slice format and an error if the extraction fails or if the extension type is unsupported.
func (e Extension) PublicKey() ([]byte, []byte, error) {
	if e.IsEmpty() {
		return nil, nil, fmt.Errorf("%w: %s", ErrPathNotFound, e)
	}
	switch {
	case e.IsDir():
		manifest := manifestFile(e.String())
		file, err := os.ReadFile(manifest)
		if err != nil {
			return nil, nil, fmt.Errorf("crx3: failed to read file %s: %w", manifest, err)
		}
		pubkey := parseKeyFromManifest(file)
		if len(pubkey) == 0 {
			return nil, nil, fmt.Errorf("crx3: failed to parse key from manifest file %s", manifest)
		}
		return []byte(pubkey), nil, nil
	case e.IsZip():
		pubkey, err := parseManifestFromZip(e.String())
		if err != nil {
			return nil, nil, err
		}
		if len(pubkey) == 0 {
			return nil, nil, fmt.Errorf("crx3: failed to parse key from manifest file %s", e)
		}
		return []byte(pubkey), nil, nil
	case e.IsCRX3():
		pubkey, err := parseManifestFromZip(e.String())
		if err != nil || len(pubkey) == 0 {
			file, err := os.ReadFile(e.String())
			if err != nil {
				return nil, nil, fmt.Errorf("crx3: failed to read file %s: %w", e.String(), err)
			}
			pubk, _, err := PubkeyFrom(bytes.NewReader(file))
			if err != nil {
				return nil, nil, err
			}
			return formatPemKey(pubk), nil, nil
		}
		if len(pubkey) > 0 {
			return []byte(pubkey), nil, nil
		}
	}
	return nil, nil, fmt.Errorf("%w: %s", ErrUnknownFileExtension, e)
}

func manifestFile(path string) string {
	return filepath.Join(path, "manifest.json")
}

func parseManifestFromZip(e string) (string, error) {
	zipReader, err := zip.OpenReader(e)
	if err != nil {
		return "", fmt.Errorf("crx3: failed to open zip reader: %w", err)
	}
	defer zipReader.Close()
	for _, file := range zipReader.File {
		baseFilename := filepath.Base(file.Name)
		if baseFilename == "manifest.json" {
			fileReader, err := file.Open()
			if err != nil {
				return "", fmt.Errorf("crx3: failed to open file %s: %w", file.Name, err)
			}
			defer fileReader.Close()
			data, err := io.ReadAll(fileReader)
			if err != nil {
				return "", fmt.Errorf("crx3: failed to read file %s: %w", file.Name, err)
			}
			pubkey := parseKeyFromManifest(data)
			if len(pubkey) == 0 {
				return "", fmt.Errorf("crx3: failed to parse key from manifest file %s", e)
			}
			return pubkey, nil
		}
	}
	return "", fmt.Errorf("crx3: failed to find manifest file %s in zip archive", e)
}

func parseKeyFromManifest(data []byte) string {
	key := struct {
		Key string `json:"key"`
	}{}
	if err := json.Unmarshal(data, &key); err != nil {
		return ""
	}
	return key.Key
}

func formatPemKey(key []byte) []byte {
	str := string(key)
	str = strings.TrimSpace(str)
	str = strings.Replace(str, "-----BEGIN RSA PUBLIC KEY-----", "", 1)
	str = strings.Replace(str, "-----END RSA PUBLIC KEY-----", "", 1)
	str = strings.Replace(str, "-----BEGIN PUBLIC KEY-----", "", 1)
	str = strings.Replace(str, "-----END PUBLIC KEY-----", "", 1)
	str = strings.TrimSpace(str)
	str = strings.ReplaceAll(str, "\n", "")
	return []byte(str)
}
