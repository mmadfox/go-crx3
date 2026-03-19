# crx3_workspace

Returns the absolute path to the workspace directory used by the CRX3 server to store downloaded Chrome extensions.

<usage>
Use this tool when the user needs to know where extension files are physically stored on disk. This is useful for:
- Locating downloaded .crx files
- Verifying the storage configuration
- Providing the full path to an extension for external tools or scripts
</usage>

<params>
Input:
- No input parameters required.
</params>

<result>
Output:
{{ if not .DisabledMarkdown }}
- The absolute filesystem path to the workspace root directory.
{{ end }}
StructuredOutput:
```json
# Example:
{
   "path": "/home/user/.crx3/workspace"
}
</result>

<use_cases>
Example use cases:
- "Where are the extensions saved?"
- "What is the workspace path for CRX3?"
- "I need the full path to the downloaded .crx file"
- "Show me the root directory for extension storage"
</use_cases>