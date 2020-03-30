package crx3

import "errors"

var (
	ErrUnknownFileExtension  = errors.New("crx3: unknown file extension")
	ErrUnsupportedFileFormat = errors.New("crx3: unsupported file format")
	ErrExtensionNotSpecified = errors.New("crx3: extension id not specified")
	ErrPathNotFound          = errors.New("crx3: filepath not found")
	ErrPrivateKeyNotFound    = errors.New("crx3: private key not found")
)
