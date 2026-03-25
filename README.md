# crx3

[![Coverage Status](https://coveralls.io/repos/github/mmadfox/go-crx3/badge.svg?branch=master)](https://coveralls.io/github/mmadfox/go-crx3?branch=master)
[![Documentation](https://pkg.go.dev/badge/github.com/mediabuyerbot/go-crx3.svg)](https://pkg.go.dev/github.com/mediabuyerbot/go-crx3)
[![Go Report Card](https://goreportcard.com/badge/github.com/mediabuyerbot/go-crx3)](https://goreportcard.com/report/github.com/mediabuyerbot/go-crx3)
![Actions](https://github.com/mmadfox/go-crx3/actions/workflows/cover.yml/badge.svg)

> 🤖 **AI-ready Chrome extension tooling via Model Context Protocol (MCP)**  
> A modern, comprehensive toolset for CRX3 Chrome extension management — pack, unpack, download, analyze, and automate with AI.

👉 **[MCP Configuration Guide](./MCP.md)** — Setup instructions for Cursor, Claude Desktop, Opencode, Crush & more.

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| 📦 **Pack & Unpack** | Create and extract signed CRX3 packages with key management |
| 🗜️ **Zip & Unzip** | Handle extension archives with flexible output paths |
| ⬇️ **Download** | Fetch extensions from Chrome Web Store by ID or URL |
| 🔑 **ID Generation** | Generate extension IDs from public keys (SHA-256) |
| 🔍 **Analyze** | Inspect manifest, permissions, and extension metadata |
| 🤖 **MCP Support** | Native AI integration via Model Context Protocol |

---

## 🤖 MCP Integration

This tool is **AI-compatible** through [Model Context Protocol (MCP)](https://modelcontextprotocol.io), enabling seamless interaction with:

- ✅ **Cursor IDE** — Chat with your extensions, auto-generate packs
- ✅ **Claude Desktop** — Natural language CRX3 operations
- ✅ **Opencode** — Automated extension workflows
- ✅ **Crush** — AI-assisted extension management
- ✅ **Any MCP client** — HTTP/SSE or stdio transport

Automate CRX3 operations through natural language commands powered by AI.

```bash
# Start MCP server (stdio mode for local AI clients)
crx3 mcp

# Or run over HTTP/SSE for remote clients 
crx3 mcp --listen=locahost:3000
crx3 mcp --listen=localhost:3000 --sse 
```

📖 See **[MCP.md](./MCP.md)** for detailed setup instructions.

---

## 🛠️ Command Reference

| Tool | Purpose |
|------|---------|
| `crx3 pack` | Pack directory/zip into signed `.crx` extension |
| `crx3 unpack` | Extract `.crx` file contents to directory |
| `crx3 download` | Download `.crx` extension by ID or Chrome Web Store URL |
| `crx3 search` | Search Chrome Web Store by name/keywords (via DuckDuckGo) |
| `crx3 zip` | Create `.zip` archive from directory |
| `crx3 unzip` | Extract `.zip` archive contents |
| `crx3 base64` | Encode file to Base64 string |
| `crx3 getid` | Extract Chrome Extension ID from `.crx` or directory |
| `crx3 scan` | List/filter downloaded extensions in workspace |
| `crx3 workspace` | Get absolute path to workspace root |
| `crx3 version` | Show CRX3 tool version |
| `crx3 mcp` | Start MCP server for AI integration |

> 💡 All commands support `--help` for detailed usage: `crx3 pack --help`

---

## 📦 Installation

### Via Homebrew (macOS/Linux)
```bash
brew tap mmadfox/tap https://github.com/mmadfox/homebrew-tap
brew install mmadfox/tap/crx3
```

### Via install.sh (Linux/macOS)
```bash
sudo curl -sSfL https://raw.githubusercontent.com/mmadfox/go-crx3/master/install.sh | bash -s
```

### Via Go
```bash
go install github.com/mediabuyerbot/go-crx3/crx3@latest
```

### Via Release Binary
Download pre-built binaries from [Releases](https://github.com/mmadfox/go-crx3/releases).


### Verify Installation
```bash
crx3 version
# Output: crx3 version x.y.z
```

---

## 🚀 Quick Start Examples

### Pack an extension
```bash
# Basic pack (auto-generates key if missing)
crx3 pack ./my-extension -o ./build/extension.crx3

# Pack with existing private key
crx3 pack ./my-extension -p ./keys/private.pem -o ./build/extension.crx3
```

### Unpack and inspect
```bash
# Extract CRX3 to directory
crx3 unpack ./extension.crx3 -o ./extracted 

# View manifest
cat ./extracted/manifest.json

# Without creating a subdirectory 
crx3 unpack ./extension.crx3 -o ./extracted -s
```

### Download from Chrome Web Store
```bash
# By extension ID
crx3 download blipmdconlkpinefehnmjammfjpmpbjk -o ./extensions/lighthouse

# Wihtout auto unpack
crx3 download blipmdconlkpinefehnmjammfjpmpbjk --unpack=false

# By full URL
crx3 download "https://chrome.google.com/webstore/detail/lighthouse/blipmdconlkpinefehnmjammfjpmpbjk"
```

### Generate extension ID
```bash
# From existing CRX3 file
crx3 id ./extension.crx3
# Output: dgmchnekcpklnjppdmmjlgpmpohmpmgp

# From public key
crx3 id -k ./keys/public.pem
```

### Key management
```bash
# Generate new RSA key pair
crx3 keygen ./keys/my-key.pem

# The same command creates both private.pem and extracts public key internally
```

### Archive operations
```bash
# Zip a directory
crx3 zip ./my-extension -o ./archive.zip

# Unzip archive
crx3 unzip ./archive.zip -o ./output
```

### Base64 encoding (for embedding)
```bash
# Encode CRX3 to base64
crx3 base64 ./extension.crx3

# Save to file
crx3 base64 ./extension.crx3 -o ./encoded.txt
```

### Search Chrome Web Store
```bash
# Search by keyword
crx3 search "ad blocker" 
```

---

## 💻 Code Examples (Go API)

### Pack a zip file or directory
```go
import crx3 "github.com/mediabuyerbot/go-crx3"

// Basic pack
if err := crx3.Extension("/path/to/file.zip").Pack(nil); err != nil {
    panic(err)
}

// Pack with private key
pk, err := crx3.LoadPrivateKey("/path/to/key.pem")
if err != nil { panic(err) }
if err := crx3.Extension("/path/to/file.zip").Pack(pk); err != nil {
    panic(err)
}

// Pack to custom output path
if err := crx3.Extension("/path/to/file.zip").PackTo("/path/to/ext.crx", pk); err != nil {
    panic(err)
}
```

### Unpack extension
```go
import crx3 "github.com/mediabuyerbot/go-crx3"

if err := crx3.Extension("/path/to/ext.crx").Unpack(); err != nil {
   panic(err)
}
```

### Download from Web Store
```go
import crx3 "github.com/mediabuyerbot/go-crx3"

extensionID := "blipmdconlkpinefehnmjammfjpmpbjk"
if err := crx3.DownloadFromWebStore(extensionID, "/path/to/ext.crx"); err != nil {
    panic(err)
}
```

### Generate/Load Keys
```go
import crx3 "github.com/mediabuyerbot/go-crx3"

// Generate new key
pk, err := crx3.NewPrivateKey()
if err != nil { panic(err) }

// Save and load
if err := crx3.SavePrivateKey("/path/to/key.pem", pk); err != nil { panic(err) }
pk, err = crx3.LoadPrivateKey("/path/to/key.pem")
```

### Helper Functions
```go
import crx3 "github.com/mediabuyerbot/go-crx3"

// Type detection
crx3.Extension("/path/to/ext.zip").IsZip()   // true
crx3.Extension("/path/to/ext").IsDir()       // true  
crx3.Extension("/path/to/ext.crx").IsCRX3()  // true

// Get extension ID
id, err := crx3.Extension("/path/to/ext.crx").ID()
```

### Base64 Encoding
```go
import crx3 "github.com/mediabuyerbot/go-crx3"

b, err := crx3.Extension("/path/to/ext.crx").Base64()
if err != nil { panic(err) }
fmt.Println(string(b))
```

---

## ⚙️ Advanced Usage

### MCP Server Modes

| Mode | Command | Use Case |
|------|---------|----------|
| **stdio** | `crx3 mcp` | Local AI clients (Cursor, Claude Desktop) |
| **HTTP/SSE** | `crx3 mcp --listen=:3000` | Remote clients, custom integrations |

### Common Flags
```bash
# Global flags
--workdir, -w    # Set working directory (default: .)
--logfile, -f    # Path to log file (default: stderr)

# MCP-specific flags  
--listen, -l     # Run HTTP server at address (e.g., :3000)
--tools.disabled, -d  # Disable specific tools (comma-separated)
--tools.show, -s      # List available tools and exit
```

### Example: Secure MCP Setup
```bash
# Run MCP with restricted tools and sandboxed directory
crx3 mcp \
  --workdir=/safe/extensions \
  --tools.disabled=crx3_download,crx3_unpack \
  --logfile=/var/log/crx3-mcp.log
```

---

## 🧪 Development

```bash
# Generate protobuf code
make proto

# Run tests with coverage
make test/cover

# View coverage report
go tool cover -html=coverage.out
```

---

## 📚 Resources

- 📘 [CRX3 Format Spec](https://developer.chrome.com/docs/extensions/mv3/architecture-overview/)
- 🤖 [Model Context Protocol](https://modelcontextprotocol.io)
- 🐛 [Bug Reports & Issues](https://github.com/mmadfox/go-crx3/issues)

---

## 📄 License

go-crx3 is released under the **Apache 2.0 License**.  
See [LICENSE](https://github.com/mediabuyerbot/go-crx3/blob/master/LICENSE) for details.

---
