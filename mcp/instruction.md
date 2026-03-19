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

### `crx3_scan`
- **When to use:**
  - User references a downloaded extension but filepath is unknown from context
  - User wants to browse, filter, or manage their extension library
  - Need to locate a .crx file before unpacking or re-downloading
- **Input handling:**
  - `limit`: Use to restrict results (e.g., `limit: 5` for recent items). Omit or use `0` for all.
  - `filter`: Array of keywords for name-based filtering. Case-insensitive, partial match, OR logic.
- **Result:** List of `ExtensionInfo` objects with `name`, `path`, `type`, `size`, `modified`.


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
