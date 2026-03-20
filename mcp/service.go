package mcp

import "github.com/mediabuyerbot/go-crx3"

type crx3service interface {
	UnpackTo(filename string, dirname string) error
}

type impl struct{}

func (impl) UnpackTo(filename string, dirname string) error {
	return crx3.UnpackTo(filename, dirname, crx3.UnpackDisableSubdir())
}
