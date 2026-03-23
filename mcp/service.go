package mcp

import (
	"context"
	"crypto/rsa"
	"iter"

	"github.com/mediabuyerbot/go-crx3"
)

type crx3service interface {
	UnpackTo(filename string, dirname string) error
	PackTo(source string, dest string, pk *rsa.PrivateKey) error
	SearchExtensionByName(ctx context.Context, name string) ([]crx3.SearchResult, error)
	Scan(rootPath string, opts ...crx3.ScanOption) iter.Seq2[*crx3.ExtensionInfo, error]
	DownloadFromWebStore(extensionID string, filename string) error
	GetID(filename string) (string, error)
	Base64(filename string) ([]byte, error)
	UnzipTo(filename string, dirname string) error
	ZipTo(source string, dest string) error
	Version() string
}

type impl struct {
	version string
}

func (impl) UnpackTo(filename string, dirname string) error {
	return crx3.UnpackTo(filename, dirname, crx3.UnpackDisableSubdir())
}

func (impl) PackTo(source string, dest string, pk *rsa.PrivateKey) error {
	return crx3.Pack(source, dest, pk)
}

func (impl) SearchExtensionByName(ctx context.Context, name string) ([]crx3.SearchResult, error) {
	return crx3.SearchExtensionByName(ctx, name)
}

func (impl) Scan(rootPath string, opts ...crx3.ScanOption) iter.Seq2[*crx3.ExtensionInfo, error] {
	return crx3.Scan(rootPath, opts...)
}

func (impl) DownloadFromWebStore(extensionID string, filename string) error {
	return crx3.DownloadFromWebStore(extensionID, filename)
}

func (impl) GetID(filename string) (string, error) {
	return crx3.Extension(filename).ID()
}

func (impl) Base64(filename string) ([]byte, error) {
	return crx3.Extension(filename).Base64()
}

func (impl) UnzipTo(filename string, dirname string) error {
	return crx3.UnzipTo(dirname, filename)
}

func (impl) ZipTo(source string, dest string) error {
	return crx3.ZipTo(dest, source)
}

func (s impl) Version() string {
	return s.version
}
