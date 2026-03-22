# AI Integration Guide for go-crx3-mcp

> 🤖 **Model Context Protocol (MCP) integration for CRX3 Chrome extension tooling**

This document explains how to configure `crx3 mcp` with AI assistants like **Opencode**, **Crush**, **Claude Desktop**, and other MCP-compatible clients.

## 📋 Table of Contents

1. [Overview](#overview)
2. [Available MCP Tools](#available-mcp-tools)
3. [Prerequisites](#prerequisites)
4. [Setup Instructions](#setup-instructions)
   - [Opencode](#opencode-agent)
   - [Crush](#crush-agent)
   - [Claude Desktop](#claude-desktop)
5. [Configuration Options](#configuration-options)
6. [Usage Examples](#usage-examples)

---

## 🌟 Overview

`crx3 mcp` exposes CRX3 extension management capabilities through the **Model Context Protocol (MCP)**, enabling AI assistants to:

- 📦 Pack/unpack Chrome extensions
- ⬇️ Download extensions from the Chrome Web Store
- 🔑 Generate extension IDs from public keys
- 🔍 Analyze extension structure and metadata
- 🗜️ Handle ZIP/archive operations

The MCP server runs in two modes:
| Mode | Command | Use Case |
|------|---------|----------|
| **stdio** | `crx3 mcp` | Local AI clients (Opencode, Crush) |
| **HTTP/SSE** | `crx3 mcp --listen=:3000` | Remote clients, custom integrations **(WIP)** |

---

## 🛠️ Available MCP Tools

Tool                | Purpose
--------------------|--------------------------------------------------
`crx3_search`         | Search Chrome Web Store by name/keywords using DuckDuckGo
`crx3_download`       | Download .crx extension by ID or URL
`crx3_workspace`     | Get absolute path to workspace root
`crx3_unpack`         | Extract .crx file contents to directory
`crx3_pack`           | Pack directory/zip into signed .crx
`crx3_scan`           | List/filter downloaded extensions in workspace
`crx3_unzip`          | Extract .zip archive contents
`crx3_zip`            | Create .zip archive from directory
`crx3_base64`         | Encode file to Base64 string
`crx3_getid`          | Extract Chrome Extension ID from .crx or directory
`crx3_version`        | Show CRX3 tool version

> 💡 Use `crx3 mcp --tools.show` to see the full tool schema in JSON format.

---

## ✅ Prerequisites
1. **Install `crx3 mcp`**:
   ```bash
   go install github.com/mediabuyerbot/go-crx3-mcp@latest
   ```

2. **Verify installation**:
   ```bash
   crx3 mcp --help
   crx3 mcp --tools.show | jq .tools  # requires jq
   ```

3. **Ensure your AI client supports MCP**:
    - Opencode: 1.1.36+ with MCP enabled
    - Crush: v0.41.0+ with MCP enabled
    - Cursor: v0.40+ with MCP enabled
    - Claude Desktop: MCP plugin installed
    - Windsurf: MCP support in settings

---

## ⚙️ Setup Instructions
### ▶️ Opencode
Opencode uses a configuration file (typically .opencode/config.json or opencode.json) to define MCP servers.

1. Locate or create the config file:
Project-local: .opencode/config.json
Global: ~/.opencode/config.json
2. Add the crx3 MCP server:
```json
{
  "mcp": {
    "crx3": {
      "type": "stdio",
      "command": "crx3",
      "args": [
        "mcp"
      ],
      "env": {},
      "tools": {
        "autoApprove": [
          "crx3_pack",
          "crx3_unpack",
          "crx3_download"
        ]
      }
    }
  }
}
```
3. Optional: Disable tools for safety:
```json
{
  "mcp": {
    "crx3": {
      "type": "local",
      "command": "crx3",
      "args": [
        "mcp",
        "--tools.disabled=crx3_download",
        "--workdir=/safe/path"
      ],
      "env": {}
    }
  }
}
```
4. Restart Opencode or reload the workspace.
5. Verify connection:
    - Open the Opencode chat panel
    - Type: @crx3 list available tools
    - You should see the list of CRX3 tools

### ▶️ Crush
Crush supports MCP via its ~/.crush/config.yaml configuration file (or project-local .crush/config.yaml).

1. Edit the Crush config:
```json
 "mcp": {
    "crx3": {
      "type": "stdio",
      "command": "crx3",
      "args": ["mcp", "--workdir=some/dir"]
    }
  }
```
2. Project-specific override (optional):
```json
 "mcp": {
    "crx3": {
      "type": "stdio",
      "command": "crx3",
      "args": ["mcp", "--workdir=some/dir"]
    }
  }
```
3. Restart Crush or run crush reload if supported.
4. @crx3 list available tools

### ▶️ Cursor IDE
1. Open Cursor Settings → **MCP Servers**
2. Add a new server configuration:

```json
{
  "mcpServers": {
    "crx3": {
      "command": "crx3",
      "args": ["mcp", "--workdir=/path/to/extensions"],
      "env": {},
      "disabled": false,
      "autoApprove": ["crx3_pack", "crx3_unpack", "crx3_download"],
      "description": "CRX3 Chrome extension management tools"
    }
  }
}
```

3. Restart Cursor. You should see `crx3` tools available in the AI chat.

5. Test in chat:
```
Can you pack my extension from ./my-extension into ./output.crx3?
```

### ▶️ Claude Desktop

1. Locate Claude Desktop config:
   - macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - Windows: `%APPDATA%\Claude\claude_desktop_config.json`
   - Linux: `~/.config/Claude/claude_desktop_config.json`

2. Add the MCP server entry:

```json
{
  "mcpServers": {
    "crx3": {
      "command": "crx3",
      "args": ["mcp", "--workdir=/path/to/extensions"],
      "env": {
      }
    }
  }
}
```

3. Restart Claude Desktop.

4. Test in chat:
```
Can you pack my extension from ./my-extension into ./output.crx3?
```

---

## ⚙️ Configuration Options

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--listen` | `-l` | string | `""` | Run server over HTTP/SSE at given address (e.g., `:3000`). If unset, uses stdio. |
| `--logfile` | `-f` | string | `""` | Path to log file. If unset, logs to stderr. |
| `--tools.show` | `-s` | bool | `false` | Print available tools + instruction JSON and exit. |
| `--tools.disabled` | `-d` | string[] | `[]` | Comma-separated list of tool names to disable. |
| `--tools.disabledMarkdownOutput` | `-m` | bool | `false` | Return only JSON (no Markdown) in tool responses. |
| `--workdir` | `-w` | string | `.` | Working directory for all file operations. |

### Example: Disable sensitive tools in shared environments

```bash
crx3 mcp --tools.disabled=crx3_download,crx3_unpack --workdir=/safe/dir
```

### Example: Structured output only (for automated clients)

```bash
crx3 mcp --tools.disabledMarkdownOutput --logfile=/var/log/crx3-mcp.json
```

---

## 💬 Usage Examples

### 🤖 Natural Language Prompts for AI
**Download and unpack**:
```
@crx3 Download extension ID "nmmhkkegccagdldgiimedpiccmgmieda" and unpack it to ./extensions/google-pay
```

### 🧪 Testing with `--tools.show`

```bash
crx3 mcp --tools.show | jq .
```

Output:
```json
{
  "instruction": "You are a CRX3 extension management assistant. Use the provided tools to pack, unpack, analyze, and manage Chrome extensions.",
  "tools": [
    {
      "name": "pack_extension",
      "description": "Pack a directory into a CRX3 package",
      "inputSchema": { ... }
    },
    ...
  ]
}
```

---