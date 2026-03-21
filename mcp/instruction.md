# The CRX3 MCP Server

These instructions describe how to efficiently work with the CRX3 tools set using the MCP server. You can load this file directly into a session where the CRX3 MCP server is connected.

## Tool Usage Logic

### `crx3_search`
- **When to use:** User provided a name/keywords, but no exact Extension ID is available.
- **Result:** List of found extensions with their `extensionId`.

### `crx3_download`
- **When to use:** User wants to download the extension file (.crx).
- **Input handling:**
  1. If user provided a name → First call `crx3_search`, extract the `extensionId` from results.
  2. If user provided `extensionId` or URL → Call `crx3_download` directly.
- **Storage:** All files are saved in the **workspace root directory**. Subdirectories are created relative to the root.

### `crx3_workspace`
- **When to use:** User asks where files are saved or needs the absolute path to the extension storage.
- **Input:** No parameters required.
- **Output:** Absolute path to the workspace root.

### `crx3_unpack`
- **When to use:** User wants to extract/inspect the contents of a .crx file.
- **Input handling:**
  1. If filepath is known from context → Call `crx3_unpack` directly.
  2. If filepath is unknown → First call `crx3_scan` to locate the file.

<params>
- `filepath` (string, required): Path to the .crx file.
  - Use `path` from `crx3_scan` results OR from previous `crx3_download` output.
  - Must be relative to workspace root (e.g., `./extensions/abc123.crx`).

- `outputDir` (string, optional): Target directory for unpacked contents.
  - **Auto-decision:** If omitted, tool creates directory using extension name/ID (e.g., `./unpacked/abc123/`).
  - **User-specified:** If user requests a specific folder, pass it here (e.g., `./my-tools/react-devtools/`).
  - **Naming rules:**
    - ✅ Allowed: letters (a-z, A-Z, Cyrillic), numbers, hyphens (`-`), underscores (`_`)
    - ❌ Forbidden: special chars (`*`, `?`, `:`, `|`, `<`, `>`), leading dots (`.`), backslashes (`\`)
    - Cross-OS safe: tool auto-converts invalid chars to `_`
  - **Path format:** Always use forward slashes (`/`), relative to workspace root.
</params>

<critical_rules>
1. **Path validation:** Never pass absolute paths (`/home/...`, `C:\Users\...`). All paths must be workspace-relative.
2. **Directory creation:** If `outputDir` doesn't exist, tool creates it automatically — no need to check first.
3. **Name sanitization:** If user provides Russian/special chars in `outputDir`, tool sanitizes them. Inform user if name was modified.
4. **Context tracking:** After unpack, cache the returned `filepath` for future reference (e.g., "the unpacked extension").
</critical_rules>

<examples>
```json
// Auto-generated directory (omit outputDir)
{"filepath": "./ext/abc123.crx"}
→ Creates: ./unpacked/abc123/

// Custom directory (user specified)
{"filepath": "./ext/abc123.crx", "outputDir": "./my-tools/react-devtools/"}
→ Creates: ./my-tools/react-devtools/

// Russian directory name (auto-sanitized)
{"filepath": "./ext/abc123.crx", "outputDir": "./расширения/адблок/"}
→ Creates: ./rasshireniya/adblock/ OR ./расширения/адблок/ (depends on OS support)
</example>

### crx3_pack
- **When to use:**: User wants to create/rebuild a .crx file from an unpacked extension directory or zip archive.
- **Input handling:**:
If source is an unpacked directory → Pass sourceDir to pack it.
If source is a zip file → Pass sourceZip to pack it.
If user modified an unpacked extension → Use crx3_pack to rebuild after changes.
Output: Path to the newly created .crx file (saved in workspace).

<params>
- `sourceDir` (string, optional): Path to unpacked extension directory (relative to workspace).
- Example: `./unpacked/abc123/`

sourceZip (string, optional): Path to zip archive with extension source (relative to workspace).
Example: ./source/my-extension.zip
Tool auto-extracts before packing.
outputPath (string, optional): Target path for the .crx file (relative to workspace).
If omitted: auto-generates ./packed/<name-or-id>.crx
Naming rules: Same as `crx3_unpack` — Cyrillic allowed, special chars sanitized.
keyPath (string, optional): Path to existing private key (.pem) for signing.
If omitted: tool generates new key pair and saves .pem alongside .crx.
overwrite (bool, optional): Allow overwriting existing .crx file. Default: false.
</params>

<critical_rules>
1. **Source validation:** Exactly one of `sourceDir` or `sourceZip` must be provided. If both or none → ask user to clarify.
2. **Workspace-relative paths:** All input/output paths must be relative to workspace root. Never use absolute paths.
3. **Key management:**
- If `keyPath` not provided, new key is generated → inform user and provide path to `.pem`.
- Store key path in context: "Use this key for future updates to maintain extension ID".
5. **ID consistency:** Packing with the same key preserves the Extension ID. Packing with a new key → new ID.
6. **Context tracking:** After pack, cache the `.crx` path and associated key path for future `crx3_download`/`crx3_unpack` workflow.
7. **Conflict handling:** If `outputPath` exists and `overwrite=false` → suggest alternative name or ask user.
</critical_rules>

<examples>
```json
// Pack unpacked directory, auto output path
{"sourceDir": "./unpacked/abc123/"}
→ Creates: ./packed/abc123.crx + ./packed/abc123.pem
// Pack with custom output and existing key
{"sourceDir": "./modified/react-devtools/", "outputPath": "./release/react-devtools-v2.crx", "keyPath": "./keys/react.pem"}
→ Creates: ./release/react-devtools-v2.crx (same ID as original)
// Pack from zip archive
{"sourceZip": "./source/my-extension.zip", "outputPath": "./packed/my-extension.crx"}
→ Extracts zip → packs → ./packed/my-extension.crx
// Russian output path (auto-sanitized)
{"sourceDir": "./распаковано/ублок/", "outputPath": "./сборка/ублок-новая-версия.crx"}
→ Creates: ./sborka/ublock-novaya-versiya.crx
</examples>

### `crx3_scan`
- **When to use:**
  - User references a downloaded extension but filepath is unknown from context
  - User wants to browse, filter, or manage their extension library
  - Need to locate a .crx file before unpacking or re-downloading
- **Input handling:**
  - `limit`: Use to restrict results (e.g., `limit: 5` for recent items). Omit or use `0` for all.
  - `filter`: Array of keywords for name-based filtering. Case-insensitive, partial match, OR logic.


<critical_rules>
1. **Parameter mapping:** Pass extensionId or URL to the `url` parameter of `crx3_download`.
   - Accepts: 32-character ID OR full Chrome Web Store URL.
   - The tool handles ID extraction from URLs automatically.

2. **Workspace awareness:** 
   - All downloads are stored relative to the workspace root.
   - If user asks for the file location, call `crx3_workspace` to get the base path.
   - Do NOT use arbitrary absolute paths (e.g., `/tmp/`, `C:/`) unless explicitly supported by the `path` parameter.

3. **Path parameter:** The `path` parameter in `crx3_download` is **relative to workspace root**, not the filesystem root.
   - ✅ Correct: `./extensions/`, `my-downloads/`
   - ❌ Incorrect: `/var/tmp/`, `C:/Users/...`

4. **Disambiguation:** If `crx3_search` returns multiple results, ask the user to specify which extension to download.

5. **Verification:** After download, you can use `crx3_workspace` + returned filepath to confirm the full location to the user.

</critical_rules>
