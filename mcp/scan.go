package mcp

import (
	_ "embed"
)

var (
	//go:embed scan.md
	scanDescription string
	scanTitle       = "List Chrome extensions in the workspace"
)
