# crx3_scan

Scans the workspace directory and returns detailed information about downloaded Chrome extensions (.crx files) and unpacked extension directories.

<usage>
Use this tool when you need to discover what extensions are available in the workspace. This is essential when:
- The user references a previously downloaded extension without providing a specific path
- You need to locate a file before unpacking or processing it
- The user wants to browse, filter, or manage their extension library
</usage>

<params>
Input:
- limit (int, optional): Maximum number of extensions to return. Use `0` or omit for no limit.
  - Example: `10` to get only the 10 most recent extensions
- filter (array of strings, optional): Keywords to filter extensions by name. Case-insensitive partial match.
  - Example: `["react", "devtools"]` matches extensions with "react" OR "devtools" in the name
  - Example: `["adblock", "privacy"]` matches any extension containing either keyword
</params>

<result>
Output:
{{ if not .DisabledMarkdown }}
- A markdown-formatted table or list of extensions, including:
  - **Name**: Extension name (if detectable) or filename
  - **Path**: Relative path within workspace (use this for `crx3_download` or `crx3_unpack`)
  - **Type**: `crx` (packed) or `directory` (unpacked)
  - **Size**: File size in bytes (for .crx) or directory indicator
  - **Modified**: Last modified timestamp (ISO 8601 or human-readable)
{{ end }}
StructuredOutput:
```json
# Example:
[
    {
       "name": "React Developer Tools",
       "path": "./myext/fmkadmapgofadopljbjfkapdkoieni.crx",
       "type": "crx",
       "size": 2847392,
       "modified": "2026-03-20T14:32:10Z"
    },
    {
       "name": "uBlock Origin",
       "path": "./unpacked/cjpalhdlnbpafiamejdnhcphjbkeiagm/",
       "type": "dir",
       "size": 1200,
       "modified": "2026-03-19T09:15:42Z"
    },
      {
       "name": "uBlock Origin",
       "path": "./cjpalhdlnbpafiamejdnhcphjbkeiagm.zip",
       "type": "zip",
       "size": 245,
       "modified": "2026-03-19T09:15:42Z"
    }
]
</result>

<use_cases>
Example use cases:
- "What extensions are already downloaded?"
- "Find extensions with 'react' in the name"
- "List only the 5 most recent downloads"
- "Show me all unpacked extensions"
- "I need to unpack the ad blocker I downloaded earlier" (when path is unknown)
- "Filter extensions by keywords: ['privacy', 'security']"
</use_cases>