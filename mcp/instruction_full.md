# CRX3 MCP Server — Complete Tool Instructions

===============================================================================
SESSION INITIALIZATION PROTOCOL — ALWAYS START WITH WORKSPACE
===============================================================================
✅ IMMEDIATELY upon session start (first user message), BEFORE answering any file-related questions:
   1. Call `crx3_workspace {}` automatically — no user prompt needed
   2. Cache the returned `absoluteRootPath` in session memory
   3. Confirm to user (optional but recommended): 
      "Workspace initialized at: {absoluteRootPath}"

✅ This cached workspace root is used for:
   • Resolving all relative paths: `"./file.crx"` → `{cached_root}/file.crx`
   • Constructing absolute paths for user display
   • Validating that tool inputs are workspace-relative
   • Debugging "file not found" errors

✅ If `crx3_workspace` fails or returns unexpected result:
   → Halt file operations
   → Inform user: "⚠️ Unable to determine workspace root. Please check server configuration."
   → Do NOT proceed with unpack/download/pack until workspace is confirmed

[SESSION LIFECYCLE — WORKSPACE AWARENESS]
| Phase              | Action                                                                 |
|--------------------|------------------------------------------------------------------------|
| Session Start      | Auto-call `crx3_workspace {}` → cache `absoluteRootPath`              |
| During Session     | Use cached root for all path resolution; never re-call unless error   |
| Error Recovery     | If file not found: re-call `crx3_workspace` to verify root unchanged  |
| Session End        | Clear cached workspace path (optional, depends on implementation)     |

[WHY THIS MATTERS — COMMON FAILURE SCENARIOS PREVENTED]
❌ Without initialization:
   • User: "Unpack ./ext.crx"
   • LLM: Doesn't know workspace root → guesses or assumes current dir → tool error
   • Result: Retry loop, confused user, wasted tokens


===============================================================================
MANDATORY WORKSPACE CHECK BEFORE FILE OPERATIONS [NEW RULE]
===============================================================================
✅ BEFORE calling any of the following tools, you MUST ensure crx3_workspace has been called in the current session:
• crx3_pack
• crx3_unpack
• crx3_getid
• crx3_base64
• crx3_zip
• crx3_unzip
✅ Implementation rule:
Check if absoluteRootPath is cached in session memory
If NOT cached → call crx3_workspace {} first → cache the result → proceed with target tool
If cached → proceed directly with target tool (no redundant call needed)
✅ This ensures:
• All file operations use the correct workspace root
• Path resolution errors are prevented
• Session state remains consistent even if workspace config changes
✅ Exception: crx3_search, crx3_download, crx3_scan, crx3_version do NOT require prior workspace call (but may still benefit from it for path construction)

===============================================================================
WORKSPACE ISOLATION
===============================================================================
✅ All CRX3 tools operate ONLY within the configured workspace directory.
✅ All file paths passed to tool parameters MUST be RELATIVE to the workspace root.
✅ NEVER pass absolute paths to any CRX3 tool:
   ❌ "/home/user/project/ext.crx"
   ❌ "C:\Users\name\extensions\app.crx"
   ✅ "./extensions/app.crx"
   ✅ "./unpacked/my-extension/"

[PATH FORMAT FOR TOOL INPUTS]
✅ Always use forward slashes: `./dir/file.ext`
✅ Start paths with `./` or `../` (relative to workspace root)
✅ Do NOT use environment variables, tilde `~`, or OS-specific absolute paths

[PATH FORMAT FOR USER DISPLAY]
✅ For user communication, you MAY and SHOULD show absolute paths
✅ Use the `crx3_workspace` tool to retrieve the absolute workspace root
✅ Response format to user: 
   "File saved at: {ABSOLUTE_PATH_FROM_WORKSPACE_TOOL}/{relative_path}"

[VALIDATION CHECKLIST BEFORE TOOL CALL]
□ Does the path start with `./` or `../`?
□ Does the path avoid forbidden characters: `:`, `|`, `<`, `>`, `*`, `?`, `\` (except `/`)?
□ Is the path NOT absolute (not starting with `/` or `C:\`)?
□ If the path contains Cyrillic/special characters, am I prepared for automatic sanitization?

[ERROR PREVENTION]
❌ If path is absolute → DO NOT call the tool; correct the path first
❌ If file not found → first call `crx3_scan` or `crx3_workspace` for diagnostics
❌ If user requests an absolute path for tool input → explain the rule, offer `crx3_workspace` for conversion

===============================================================================
CRITICAL RULE TO PREVENT PATH RESOLUTION ERRORS AND REDUNDANT TOOL CALLS
===============================================================================
[WORKSPACE AS ABSOLUTE ROOT — NON-NEGOTIABLE]
✅ The workspace root (returned by `crx3_workspace`) is the ONLY reference point for all file operations.
✅ NEVER assume "current directory", working directory, or caller context — they are irrelevant to CRX3 tools.
✅ Every relative path like `"./extensions/file.crx"` is resolved as `{workspace_root}/extensions/file.crx`.
✅ If you are unsure where a file is, ALWAYS start with `crx3_workspace` or `crx3_scan` — never guess.

[FILE NOT FOUND — DIAGNOSTIC PROTOCOL]
If a tool returns "file not found", "path does not exist", or similar error:

1️⃣ STOP — do NOT retry with the same path
2️⃣ CALL `crx3_workspace {}` to confirm the absolute root path
3️⃣ CALL `crx3_scan` with appropriate filters to locate the file:
   • By name: `{"filter": ["filename", "extensionId"]}`
   • Broad scan: `{}` to list all .crx files
4️⃣ Use the `filepath` returned by `crx3_scan` directly — it is already workspace-relative and valid
5️⃣ If `crx3_scan` returns no results:
   → Inform user the file is not in workspace
   → Offer to download it via `crx3_search` + `crx3_download`
   → Or ask user to place it inside workspace and rescan

===============================================================================
TOOLS OVERVIEW
===============================================================================

Tool                | Purpose
--------------------|--------------------------------------------------
crx3_search         | Search Chrome Web Store by name/keywords
crx3_download       | Download .crx extension by ID or URL
crx3_workspace      | Get absolute path to workspace root
crx3_unpack         | Extract .crx file contents to directory
crx3_pack           | Pack directory/zip into signed .crx
crx3_scan           | List/filter downloaded extensions in workspace
crx3_unzip          | Extract .zip archive contents
crx3_zip            | Create .zip archive from directory
crx3_base64         | Encode file to Base64 string
crx3_getid          | Extract Chrome Extension ID from .crx or directory
crx3_version        | Show CRX3 tool version

===============================================================================
crx3_search
===============================================================================

WHEN TO USE:
- User provides extension name/keywords but no exact ID
- Need to discover extension ID before download
- Browsing available extensions by category or feature

PARAMETERS:
- query (required): Search query (extension name, keywords, or partial match)
- limit (optional): Maximum number of results to return (default: 10, 0 = all)

CRITICAL RULES:
- Results may include unofficial/malicious extensions — verify source before download
- If multiple results match, present options to user for selection
- Search is case-insensitive and supports partial keyword matching

EXAMPLES:
{"query": "password manager", "limit": 5}
{"query": "adblock"}

===============================================================================
crx3_download
===============================================================================

WHEN TO USE:
- User wants to download a specific extension
- After crx3_search to download selected extension
- Re-download extension for fresh copy or update

PARAMETERS:
- extensionId OR url (required, one of): 32-char Chrome extension ID or full Chrome Web Store URL
- path (optional): Workspace-relative path to save .crx file (default: ./extensions/{extensionId}.crx)

CRITICAL RULES:
- Use workspace-relative paths only: "./downloads/ext.crx"
- NEVER use absolute paths: "/home/user/...", "C:\Users\..."
- Tool auto-extracts extension ID from URLs
- Tool creates parent directories automatically
- After download, cache filepath for subsequent unpack/inspect operations

EXAMPLES:
{"extensionId": "nkbihfbeogaeaoehlefnkodbefgpgknn"}
{"url": "https://chrome.google.com/webstore/detail/nkbihfbeogaeaoehlefnkodbefgpgknn", "path": "./releases/metamask-latest.crx"}

===============================================================================
crx3_workspace
===============================================================================

WHEN TO USE:
- User asks "where are files saved?"
- Need to construct absolute paths for external tools
- Debugging path-related issues

PARAMETERS: None

CRITICAL RULES:
- Use this tool to translate relative paths to absolute for user communication
- Do NOT pass the returned absolute path to other CRX3 tools — they expect relative paths
- Workspace root is configured at server startup and is read-only via this tool

EXAMPLES:
{}

===============================================================================
crx3_unpack
===============================================================================

WHEN TO USE:
- User wants to inspect extension source code
- Before modifying extension files
- After download to verify contents

PARAMETERS:
- filepath (required): Path to .crx file (workspace-relative) call `crx3_workspace`
- outputDir (optional): Target directory for extracted contents (workspace-relative) call `crx3_workspace`

PATH NAMING RULES:
- Allowed: letters (a-z, A-Z, Cyrillic), numbers, hyphens (-), underscores (_), forward slashes (/)
- Forbidden: * ? : | < > \ and leading dots
- Tool auto-sanitizes invalid characters to underscore (_)
- Always use forward slashes (/), relative to workspace root

CRITICAL RULES:
- All paths must be workspace-relative — never absolute call `crx3_workspace`
- Tool auto-creates outputDir if it doesn't exist
- If outputDir name contains Russian/special chars, tool sanitizes — inform user
- After unpack, cache outputDir for future pack/modify operations
- Do NOT manually create the output directory before unpacking.
- Do NOT pass the `outputDir` parameter unless the user explicitly requests a specific location.
- If `outputDir` is needed: Specify ONLY the parent folder (e.g., "./audit/"), do NOT include the extension name or ID.
- Paths: Relative only (must start with "./"), use forward slashes ("/"). Absolute paths are forbidden.

Why this matters:
❌ Error: outputDir: "./unpacked/ext-id/" → Tool creates the directory + nests extracted contents inside another folder → Double nesting / incorrect structure.
✅ Correct: outputDir omitted → Tool automatically creates "./unpacked/{extensionId}/" with proper structure.
// ✅ Allowed (parent folder only)
{"filepath": "./ext/metamask.crx", "outputDir": "./audit/"}
// ❌ FORBIDDEN (causes double nesting)
{"filepath": "./ext.crx", "outputDir": "./unpacked/ext-id/"}

EXAMPLES:
{"filepath": "./extensions/abc123.crx"}
{"filepath": "./ext/metamask.crx", "outputDir": "./audit/metamask-v10/"}
{"filepath": "./ext/test.crx", "outputDir": "./расширения/тест/"}

===============================================================================
crx3_pack
===============================================================================

WHEN TO USE:
- User modified unpacked extension and wants to rebuild .crx
- Create distributable .crx from source directory
- Repack extension with new signing key

PARAMETERS:
- source (required): Path to source directory or .zip file (workspace-relative)
- outputDir (optional): Directory for output .crx (workspace-relative, default: ./packed/)
- name (optional): Custom filename for output .crx (without path)
- privateKey (optional): Path to existing .pem private key for signing (workspace-relative)

CRITICAL RULES:
- source must be existing directory or .zip file — validate before calling
- Tool auto-detects if source is directory or .zip
- All paths workspace-relative — never absolute
- KEY MANAGEMENT:
  * If privateKey omitted → new key generated → inform user + provide .pem path
  * Same key + manifest = same extension ID
  * New key = new extension ID (breaks update chain)
- Cache extensionID and privateKey path after pack for future updates
- If output file exists → suggest alternative name or request overwrite confirmation

EXAMPLES:
{"source": "./unpacked/abc123/"}
{"source": "./modified/react-devtools/", "outputDir": "./release/", "name": "react-devtools-v2.crx", "privateKey": "./keys/react.pem"}
{"source": "./source/my-extension.zip", "outputDir": "./packed/", "name": "my-extension.crx"}

===============================================================================
crx3_scan
===============================================================================

WHEN TO USE:
- User asks "what extensions do I have downloaded?"
- Need to locate .crx filepath before unpack/pack operations
- Browse or filter extension library by name/keyword

PARAMETERS:
- limit (optional): Max results (0 = unlimited, default: 0)
- filter (optional): Array of keywords for name filtering (case-insensitive, OR logic)
- sortBy (optional): Sort by "name", "date", or "size" (default: "date")

CRITICAL RULES:
- Results include only .crx files found in workspace (recursive scan)
- Use filepath from results as input to crx3_unpack/crx3_getid
- If no results: suggest downloading extension or verify workspace via crx3_workspace
- Large libraries: recommend using filter/limit for performance

EXAMPLES:
{}
{"filter": ["adblock", "privacy"], "limit": 5, "sortBy": "name"}

===============================================================================
crx3_unzip
===============================================================================

WHEN TO USE:
- User has .zip with extension source to inspect/modify
- Extract backup archives or downloaded source packages
- Prepare files before packing into .crx

PARAMETERS:
- filepath (required): Path to .zip file (workspace-relative)
- outputDir (optional): Target directory for extracted contents (workspace-relative, default: ./extracted/{zip-name}/)

CRITICAL RULES:
- filepath must point to existing .zip file
- All paths workspace-relative — never absolute
- If outputDir exists, tool merges contents (no automatic cleanup)
- Preserve directory structure and file permissions where possible

EXAMPLES:
{"filepath": "./source/my-extension.zip"}
{"filepath": "./backup/config.zip", "outputDir": "./restored/"}

===============================================================================
crx3_zip
===============================================================================

WHEN TO USE:
- Prepare extension source for distribution
- Create backup archive of modified extension
- Compress files before transmission or storage

PARAMETERS:
- source (required): Path to source directory or file (workspace-relative)
- outputDir (optional): Directory for output .zip (workspace-relative, default: ./archives/)
- name (optional): Custom filename for output .zip (without path)

CRITICAL RULES:
- source must exist and be accessible
- Recursive inclusion: all subdirectories and files included by default
- All paths workspace-relative; invalid filename chars sanitized
- Does not follow symbolic links by default

EXAMPLES:
{"source": "./my-extension/"}
{"source": "./modified/react-devtools/", "outputDir": "./releases/", "name": "react-devtools-v2-source.zip"}

===============================================================================
crx3_base64
===============================================================================

WHEN TO USE:
- Embed binary file in JSON config or HTML data URL
- Transmit file content as text (e.g., API payload)
- Prepare file for systems that require text-only input

PARAMETERS:
- filepath (required): Path to file to encode (workspace-relative)

CRITICAL RULES:
- Works with any file type (.crx, .zip, .json, .js, etc.)
- Large files (>1MB) produce very long strings — warn user
- Base64 increases size by ~33% — consider for bandwidth-sensitive use cases
- All paths workspace-relative

EXAMPLES:
{"filepath": "./unpacked/abc123/manifest.json"}
{"filepath": "./packed/my-extension.crx"}

===============================================================================
crx3_getid
===============================================================================

WHEN TO USE:
- Verify extension identity before installation
- Check if repacked extension preserves original ID
- Debug extension loading issues related to ID mismatch

PARAMETERS:
- filepath (required): Path to .crx file or unpacked extension directory (workspace-relative)

ID GENERATION LOGIC:
1. If manifest.json contains "key" field → extract public key → compute ID
2. Else if .crx header contains public key → compute ID from header
3. Else → error: unsigned/invalid extension

CRITICAL RULES:
- Same public key + manifest = same ID (critical for update chain)
- Repacking with different key → new ID → breaks auto-update
- If extraction fails: extension may be unsigned, corrupted, or modified
- All paths workspace-relative

EXAMPLES:
{"filepath": "./packed/metamask.crx"}
{"filepath": "./unpacked/react-devtools/"}

===============================================================================
crx3_version
===============================================================================

WHEN TO USE:
- Debugging: verify tool version matches documentation
- Support: include version in bug reports
- Audit: confirm deployment version

PARAMETERS: None

CRITICAL RULES:
- Output is read-only informational
- Version format follows semantic versioning (MAJOR.MINOR.PATCH)

EXAMPLES:
{}

===============================================================================
TYPICAL WORKFLOWS
===============================================================================

WORKFLOW: Download → Inspect → Modify → Repack

1. Search & Download:
   crx3_search {"query": "my extension", "limit": 5}
   → Select extensionId
   crx3_download {"extensionId": "abc123..."}

2. Unpack for inspection:
   crx3_unpack {"filepath": "./extensions/abc123.crx"}
   → Returns outputDir: "./unpacked/abc123/"

3. [User modifies files in ./unpacked/abc123/]

4. Repack with same key (preserve ID):
   crx3_pack {
     "source": "./unpacked/abc123/",
     "privateKey": "./packed/abc123.pem"
   }
   → Returns new .crx with same extensionID

5. Verify ID matches:
   crx3_getid {"filepath": "./packed/abc123.crx"}
   → Confirm ID unchanged

---

WORKFLOW: Backup & Archive

1. Scan for extensions:
   crx3_scan {"filter": ["important"]}

2. For each critical extension:
   a. crx3_getid {"filepath": "./ext/example.crx"} → cache ID
   b. crx3_unpack {"filepath": "./ext/example.crx", "outputDir": "./backup/example-src/"}
   c. crx3_zip {"source": "./backup/example-src/", "name": "example-source-backup.zip"}
   d. crx3_base64 {"filepath": "./ext/example.crx"} → store encoded copy in config

---

WORKFLOW: Development Cycle

1. Start from source directory: ./my-extension/
2. Create ZIP for portability:
   crx3_zip {"source": "./my-extension/", "name": "my-extension-src.zip"}
3. Pack to .crx for testing:
   crx3_pack {"source": "./my-extension/", "outputDir": "./build/"}
4. Get ID for manifest/permissions:
   crx3_getid {"filepath": "./build/my-extension.crx"}
5. [Test in Chrome]
6. Iterate: modify source → repack → test

===============================================================================
CRITICAL BEST PRACTICES
===============================================================================

PATH SAFETY (MOST IMPORTANT):

ALWAYS:
- Use forward slashes: "./dir/file.crx"
- Keep paths relative to workspace root
- Use crx3_workspace to get absolute path for user communication only

NEVER:
- Pass absolute paths to CRX3 tools: "/home/...", "C:\Users\..."
- Assume filesystem access outside workspace
- Hardcode OS-specific path separators

---

KEY MANAGEMENT FOR EXTENSION ID PRESERVATION:

- Chrome Extension ID = hash(public_key)
- To preserve ID across repacks:
  1. Save the .pem private key from first pack
  2. Reuse same key in future crx3_pack calls via privateKey parameter
  3. Never lose the .pem — losing it = losing ability to update extension with same ID

Workflow example:
  crx3_pack {"source": "./v1/"}
  → Returns privateKey: "./packed/v1.pem"
  # Save v1.pem securely!
  crx3_pack {"source": "./v2/", "privateKey": "./packed/v1.pem"}
  → Produces .crx with SAME extensionID

---

CONTEXT TRACKING CHECKLIST:

Operation      | Cache These Values                    | Used By Next Step
---------------|---------------------------------------|------------------
crx3_search    | extensionId, name                     | crx3_download
crx3_download  | filepath, extensionId                 | crx3_unpack, crx3_getid
crx3_unpack    | outputDir, sourceCrx                  | manual edit, crx3_pack
crx3_pack      | filepath (new .crx), privateKey, ID   | crx3_getid, distribution
crx3_getid     | extensionID                           | verification, manifest

---

ERROR RECOVERY GUIDE:

Error Scenario              | Suggested Action
----------------------------|------------------------------------------------
"File not found"            | Call crx3_scan to discover valid files
"Invalid manifest"          | Verify source has manifest.json with required fields
"ID extraction failed"      | Extension may be unsigned — check source integrity
"Path not allowed"          | Ensure path is workspace-relative, use forward slashes
"Key generation failed"     | Check workspace write permissions
Multiple search results     | Present options to user; ask for explicit selection

---
