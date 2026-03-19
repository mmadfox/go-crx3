package crx3

import (
	"archive/zip"
	"cmp"
	"fmt"
	"io/fs"
	"iter"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	maxExtensionSize = 300 * 1024 * 1024 // 300MB
	tdir             = "dir"
	tcrx             = "crx"
	tzip             = "zip"
	defaultLayout    = "2006-01-02 15:04"
	unknownName      = "unknown"
	defaultMaxDepth  = 3
)

var extensionNameRe = regexp.MustCompile(`^([a-zA-Z0-9_]+_)?([a-p]{32})(?:\.(crx|zip))?$`)

// ExtensionInfo represents metadata about a Chrome extension found during scanning.
// It includes the extension's name, file path, type (crx, zip, or dir),
// size in bytes, and modification time formatted as a string.
type ExtensionInfo struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Type     string `json:"type"`
	Size     int64  `json:"size"`
	Modified string `json:"modified"`
}

func (e *ExtensionInfo) String() string {
	return fmt.Sprintf(
		"ExtensionInfo{Name: %s, Path: %s, Type: %s, Size: %d, Modified: %s}",
		e.Name, e.Path, e.Type, e.Size, e.Modified)
}

// ScanOption is a function that configures the internal scan filter.
// It is used to pass optional arguments to the Scan function.
type ScanOption func(*scanFilter)

// WithNameFilter returns a ScanOption that filters results by partial
// case-insensitive match of the file or directory name.
//
// For example, WithNameFilter("adblock") will match:
//   - "adblock_plus_pammpkd...crx"
//   - "my_adblock_tool.zip"
//   - "Adblock-Extra" (as a directory)
//
// Multiple filters can be combined; they are joined with logical OR.
func WithNameFilter(query string) ScanOption {
	query = strings.ToLower(query)
	return func(f *scanFilter) {
		f.names = append(f.names, query)
	}
}

// WithMaxDepth returns a ScanOption that limits the directory traversal depth.
// For example, WithMaxDepth(2) will scan only root, its subdirs, and their subdirs.
func WithMaxDepth(depth int) ScanOption {
	return func(f *scanFilter) {
		f.maxDepth = depth
	}
}

// WithMaxResults returns a ScanOption that stops scanning after finding n extensions.
func WithMaxResults(n int) ScanOption {
	return func(f *scanFilter) {
		f.maxCount = n
	}
}

type scanFilter struct {
	names     []string
	maxDepth  int
	maxCount  int
	currCount int
}

// Scan walks the directory tree starting at rootPath and collects information about
// Chrome extensions in CRX, ZIP, or unpacked directory formats.
//
// It identifies extensions by:
//   - Filenames matching the pattern: <name>_<32-character-id>.crx|.zip
//   - Directories containing a "manifest.json" file or a "extension-id" file
//   - Standalone .crx and .zip files that contain a "manifest.json"
//
// For each found extension, an ExtensionInfo struct is created with details
// including name, path, type ("crx", "zip", or "dir"), size, and modification time.
//
// Directories without recognized files are skipped. Files larger than 300MB are ignored.
func Scan(rootPath string, opts ...ScanOption) iter.Seq2[*ExtensionInfo, error] {
	filter := new(scanFilter)
	for _, opt := range opts {
		opt(filter)
	}
	return func(yield func(*ExtensionInfo, error) bool) {
		err := filepath.WalkDir(rootPath,
			func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if isHidden(path) {
					if d.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}

				depth := strings.Count(path[len(rootPath):], string(filepath.Separator))
				if depth < 0 {
					depth = 0
				}

				info, err := d.Info()
				if err != nil {
					return err
				}
				if info != nil && info.Size() > maxExtensionSize {
					return nil
				}

				if strings.HasPrefix(info.Name(), ".") {
					return nil
				}

				if filter.maxDepth >= 0 && depth > filter.maxDepth {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}

				lowerName := strings.ToLower(info.Name())
				if len(filter.names) > 0 {
					var match bool
					for _, q := range filter.names {
						if strings.Contains(lowerName, q) {
							match = true
							break
						}
					}
					if !match {
						return nil
					}
				}

				if filter.maxCount > 0 && filter.currCount >= filter.maxCount {
					return filepath.SkipAll
				}

				// format: some_name_kpkcennohgffjdgaelocingbmkjnpjgc.crx|zip
				matches := extensionNameRe.FindStringSubmatch(info.Name())
				if len(matches) == 4 && IsValidExtensionID(matches[2]) {
					ei := &ExtensionInfo{
						Name:     formatName(matches[1], matches[2]),
						Type:     cmp.Or(matches[3], tdir),
						Size:     info.Size(),
						Path:     path,
						Modified: info.ModTime().Format(defaultLayout),
					}
					if !yield(ei, nil) {
						return filepath.SkipAll
					}
					filter.currCount++
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}

				// format: directory without crx|zip extension and without specifying a special name
				if info.IsDir() {
					crx3Filename := filepath.Join(path, extensionID)
					if fileExists(crx3Filename) {
						ei := &ExtensionInfo{
							Name:     formatName(info.Name(), unknownName),
							Type:     tdir,
							Size:     info.Size(),
							Path:     path,
							Modified: info.ModTime().Format(defaultLayout),
						}
						filter.currCount++
						if !yield(ei, nil) {
							return filepath.SkipAll
						}
						return filepath.SkipDir
					}

					// manifest.json
					manifestFilename := filepath.Join(path, "manifest.json")
					if fileExists(manifestFilename) {
						ei := &ExtensionInfo{
							Name:     formatName(info.Name(), unknownName),
							Type:     tdir,
							Size:     info.Size(),
							Path:     path,
							Modified: info.ModTime().Format(defaultLayout),
						}
						filter.currCount++
						if !yield(ei, nil) {
							return filepath.SkipAll
						}
						return filepath.SkipDir
					}
					return nil
				}

				// crx extension
				if strings.HasSuffix(strings.ToLower(d.Name()), crxExt) {
					ei := &ExtensionInfo{
						Name:     formatName(info.Name(), unknownName),
						Path:     path,
						Type:     tcrx,
						Size:     info.Size(),
						Modified: info.ModTime().Format(defaultLayout),
					}
					filter.currCount++
					if !yield(ei, nil) {
						return filepath.SkipAll
					}
					return nil
				}

				// zip extension
				if strings.HasSuffix(strings.ToLower(d.Name()), zipExt) {
					if manifestExists(path) {
						ei := &ExtensionInfo{
							Name:     formatName(info.Name(), unknownName),
							Path:     path,
							Type:     tzip,
							Size:     info.Size(),
							Modified: info.ModTime().Format(defaultLayout),
						}
						filter.currCount++
						if !yield(ei, nil) {
							return filepath.SkipAll
						}
					}
					return nil
				}
				return nil
			})
		if err != nil {
			if !yield(nil, err) {
				return
			}
		}
	}
}

func isHidden(path string) bool {
	parts := strings.Split(path, string(filepath.Separator))
	for _, part := range parts {
		if len(part) > 0 && part[0] == '.' {
			return true
		}
	}
	return false
}

func formatName(name string, extensionID string) string {
	if len(name) == 0 {
		return extensionID
	}
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.TrimSpace(name)
	name = strings.TrimSuffix(name, zipExt)
	name = strings.TrimSuffix(name, crxExt)
	return name
}

func manifestExists(path string) bool {
	r, err := zip.OpenReader(path)
	if err != nil {
		return false
	}
	defer r.Close()
	for _, f := range r.File {
		if f.Name == "manifest.json" {
			return true
		}
	}
	return false
}

// TODO:
// type header struct {
// 	Version          uint32
// 	HeaderLength     uint32
// 	KeyID            []byte
// 	SignedHeaderData []byte
// }

// func (h *header) PublicKeyID() string {
// 	return hex.EncodeToString(h.KeyID)
// }

// func readHeader(r io.Reader) (*header, error) {
// 	var magic [4]byte
// 	if err := binary.Read(r, binary.LittleEndian, &magic); err != nil {
// 		return nil, fmt.Errorf("failed to read magic: %w", err)
// 	}
// 	if string(magic[:]) != "Cr24" {
// 		return nil, ErrUnsupportedFileFormat
// 	}

// 	var version uint32
// 	if err := binary.Read(r, binary.LittleEndian, &version); err != nil {
// 		return nil, fmt.Errorf("failed to read version: %w", err)
// 	}
// 	if version != 3 {
// 		return nil, fmt.Errorf("unsupported CRX version: %d, expected 3", version)
// 	}

// 	var headerSize uint32
// 	if err := binary.Read(r, binary.LittleEndian, &headerSize); err != nil {
// 		return nil, fmt.Errorf("failed to read header size: %w", err)
// 	}

// 	headerBuf := make([]byte, headerSize)
// 	if _, err := io.ReadFull(r, headerBuf); err != nil {
// 		return nil, fmt.Errorf("failed to read header data: %w", err)
// 	}

// 	var crxHeader pb.CrxFileHeader
// 	if err := proto.Unmarshal(headerBuf, &crxHeader); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal CrxFileHeader: %w", err)
// 	}

// 	var signedData pb.SignedData
// 	if err := proto.Unmarshal(crxHeader.SignedHeaderData, &signedData); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal SignedData: %w", err)
// 	}

// 	if len(signedData.CrxId) != 16 {
// 		return nil, ErrUnsupportedFileFormat
// 	}

// 	return &header{
// 		Version:          version,
// 		HeaderLength:     headerSize + 12,
// 		KeyID:            signedData.CrxId,
// 		SignedHeaderData: crxHeader.SignedHeaderData,
// 	}, nil
// }
