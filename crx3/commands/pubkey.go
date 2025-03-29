package commands

import (
	"archive/zip"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type pubkeyOpts struct {
	PrivateKeyPath string // from private key
	ManifestPath   string // from manifest.json
	ExtensionPath  string // from extension/manifest.json
}

func (opts pubkeyOpts) Validate() bool {
	return len(opts.PrivateKeyPath) != 0 ||
		len(opts.ManifestPath) != 0 ||
		len(opts.ExtensionPath) != 0
}

func newPubkeyCmd() *cobra.Command {
	var opts pubkeyOpts
	cmd := &cobra.Command{
		Use:   "pubkey",
		Short: "Extract the public key from a private key, manifest.json, or CRX/ZIP file. (default: from extension)",
		Long: `The 'pubkey' command extracts a public key from different sources:
- From a private key file (PEM format)
- From the 'key' field in manifest.json of a Chrome extension
- From a CRX or ZIP extension file

The extracted public key is output in DER format and Base64-encoded.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !opts.Validate() {
				if len(args) == 0 {
					return errors.New("you need to specify a source to obtain the public key")
				}
				opts.ExtensionPath = args[0]
			}

			var pubkey string
			var err error
			switch {
			case len(opts.PrivateKeyPath) > 0:
				pubkey, err = extractPublicKeyFromPrivateKey(opts.PrivateKeyPath)
			case len(opts.ExtensionPath) > 0:
				pubkey, err = extractKeyFromExtension(opts.ExtensionPath)
			case len(opts.ManifestPath) > 0:
				pubkey, err = extractKeyFromManifest(opts.ManifestPath)
			}
			if err != nil {
				return err
			}

			fmt.Print(pubkey)

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.PrivateKeyPath, "private", "p", "", "extract public key from a private key file (PEM format)")
	cmd.Flags().StringVarP(&opts.ManifestPath, "manifest", "m", "", "extract public key from the 'key' field in manifest.json")
	cmd.Flags().StringVarP(&opts.ExtensionPath, "extension", "e", "", "extract public key from a CRX or ZIP extension file")

	return cmd
}

func extractPublicKeyFromPrivateKey(privKeyPath string) (string, error) {
	data, err := os.ReadFile(privKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read private key file: %w", err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return "", fmt.Errorf("failed to parse PEM")
	}
	var privKey *rsa.PrivateKey
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		privKey = key
	} else {
		parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return "", fmt.Errorf("failed to parse private key (unsupported format): %w", err)
		}
		var ok bool
		privKey, ok = parsedKey.(*rsa.PrivateKey)
		if !ok {
			return "", fmt.Errorf("parsed PKCS#8 but it's not an RSA private key")
		}
	}
	pubDER, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key: %w", err)
	}
	return base64.StdEncoding.EncodeToString(pubDER), nil
}

func extractKeyFromExtension(crxPath string) (string, error) {
	r, err := zip.OpenReader(crxPath)
	if err != nil {
		return "", err
	}
	defer r.Close()
	var manifestData map[string]interface{}
	for _, f := range r.File {
		if f.Name == "manifest.json" {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()
			err = json.NewDecoder(rc).Decode(&manifestData)
			if err != nil {
				return "", err
			}
			key, ok := manifestData["key"].(string)
			if !ok {
				return "", fmt.Errorf("key not found in manifest.json")
			}
			return strings.TrimSpace(key), nil
		}
	}
	return "", fmt.Errorf("manifest.json not found in CRX")
}

func extractKeyFromManifest(manifestPath string) (string, error) {
	file, err := os.Open(manifestPath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	var manifestData map[string]interface{}
	err = json.NewDecoder(file).Decode(&manifestData)
	if err != nil {
		return "", err
	}
	key, ok := manifestData["key"].(string)
	if !ok {
		return "", fmt.Errorf("key not found in manifest.json")
	}
	return strings.TrimSpace(key), nil
}
