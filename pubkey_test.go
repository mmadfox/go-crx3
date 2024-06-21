package crx3

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"io"
	"os"
	"testing"
)

func TestPubkeyFrom(t *testing.T) {
	pk, err := NewPrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	pkpem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(pk),
	})
	zipExtension, err := os.ReadFile("./testdata/bobbyMol.zip")
	if err != nil {
		t.Fatal(err)
	}
	crxExtension := new(bytes.Buffer)
	if err := PackZipToCRX(bytes.NewReader(zipExtension), crxExtension, pk); err != nil {
		t.Fatal(err)
	}

	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should return public key from crx extension",
			args: args{r: crxExtension},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := PubkeyFrom(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("PubkeyFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if ok := verifyKeyPair(pkpem, got); !ok {
				t.Errorf("PubkeyFrom() => public key and private key are not equal")
			}
		})
	}
}

func verifyKeyPair(private, public []byte) bool {
	block, _ := pem.Decode(private)
	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	pubBlock, _ := pem.Decode(public)
	pubKey, err2 := x509.ParsePKCS1PublicKey(pubBlock.Bytes)
	if err2 != nil {
		return false
	}
	return key.PublicKey.Equal(pubKey)
}
