package mcp

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

var multipleUnderscoresRe = regexp.MustCompile(`_+`)

var (
	//go:embed unpack.md
	unpackDescription string
	unpackTitle       = "Unpack a Chrome extension"
)

type unpackParams struct {
	Filepath  string `json:"filepath" jsonschema:"required, path to the downloaded .crx file"`
	OutputDir string `json:"outputDir,omitempty" jsonschema:"optional, path to the output directory"`
}

type unpackResult struct {
	Filepath string `json:"filepath" jsonschema:"required, path to the unpacked extension"`
}

func (h *handler) unpackHandler(ctx context.Context, _ *sdkmcp.CallToolRequest, params unpackParams) (*sdkmcp.CallToolResult, any, error) {
	extensionFilepath, err := h.opts.joinPath(params.Filepath)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to join path: %w", err)
	}
	isFileInvalid := !isFile(extensionFilepath) || filepath.Ext(extensionFilepath) != ".crx"
	if isFileInvalid {
		return nil, nil, fmt.Errorf("Extension not found %q", params.Filepath)
	}

	outputDir := params.OutputDir
	if len(outputDir) == 0 {
		baseName := filepath.Base(params.Filepath)
		baseName = strings.TrimSuffix(baseName, ".crx")
		outputDir = strings.Join([]string{"unpacked", sanitizeFilename(baseName)}, "_")
	}

	if filepath.IsAbs(outputDir) {
		return nil, nil, fmt.Errorf("outputDir must be relative to workspace root, got: %s", outputDir)
	}

	targetDir, err := h.opts.joinPath(outputDir)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to join unpacked path: %w", err)
	}

	if err := h.svc.UnpackTo(extensionFilepath, targetDir); err != nil {
		return nil, nil, fmt.Errorf("Failed to unpack %q to %q: %w", params.Filepath, outputDir, err)
	}

	resp := &sdkmcp.CallToolResult{
		StructuredContent: unpackResult{
			Filepath: targetDir,
		},
	}

	// render markdown if not disabled
	if !h.opts.DisabledMarkdown {
		resp.Content = []sdkmcp.Content{
			&sdkmcp.TextContent{Text: fmt.Sprintf("Successfully unpacked extension to %q", targetDir)},
		}
	}

	return resp, nil, nil
}

func sanitizePath(input string) string {
	result := cyr2latin(input)
	forbiddenChars := []string{
		"*", "?", ":", "|", "<", ">", "\\", "\"", "/",
		"[", "]", "(", ")", "{", "}", "!", "@", "#", "$", "%", "^", "&", "=", "+",
	}

	for _, char := range forbiddenChars {
		result = strings.ReplaceAll(result, char, "_")
	}

	result = strings.Trim(result, ". ")
	result = strings.TrimSpace(result)
	result = multipleUnderscoresRe.ReplaceAllString(result, "_")
	if len(result) > 255 {
		result = result[:255]
	}
	if result == "" {
		result = "unnamed_extension"
	}
	return result
}

func sanitizeFilename(input string) string {
	result := sanitizePath(input)
	result = strings.ReplaceAll(result, ".", "_")
	return result
}

func cyr2latin(text string) string {
	cyrillicToLatin := map[rune]string{
		'a': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d",
		'е': "e", 'ё': "yo", 'ж': "zh", 'з': "z", 'и': "i",
		'й': "y", 'к': "k", 'л': "l", 'м': "m", 'н': "n",
		'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t",
		'у': "u", 'ф': "f", 'х': "kh", 'ц': "ts", 'ч': "ch",
		'ш': "sh", 'щ': "sch", 'ъ': "", 'ы': "y", 'ь': "",
		'э': "e", 'ю': "yu", 'я': "ya",
		'А': "A", 'Б': "B", 'В': "V", 'Г': "G", 'Д': "D",
		'Е': "E", 'Ё': "Yo", 'Ж': "Zh", 'З': "Z", 'И': "I",
		'Й': "Y", 'К': "K", 'Л': "L", 'М': "M", 'Н': "N",
		'О': "O", 'П': "P", 'Р': "R", 'С': "S", 'Т': "T",
		'У': "U", 'Ф': "F", 'Х': "Kh", 'Ц': "Ts", 'Ч': "Ch",
		'Ш': "Sh", 'Щ': "Sch", 'Ъ': "", 'Ы': "Y", 'Ь': "",
		'Э': "E", 'Ю': "Yu", 'Я': "Ya",
	}

	var result strings.Builder
	result.Grow(len(text))

	for _, r := range text {
		if translit, ok := cyrillicToLatin[r]; ok {
			result.WriteString(translit)
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		} else {
			result.WriteRune('_')
		}
	}

	return result.String()
}

func isFile(filename string) bool {
	info, err := os.Stat(filename)
	if info == nil {
		return false
	}
	if info.Size() == 0 || os.IsNotExist(err) || err != nil {
		return false
	}
	return !info.IsDir()
}
