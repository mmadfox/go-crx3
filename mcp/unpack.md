# crx3_unpack

Unpacks a downloaded Chrome extension (.crx file) into a directory structure for inspection or modification.

<usage>
Use this tool when the user wants to extract the contents of a .crx file. The tool unpacks the extension into a directory, making source files (manifest.json, scripts, assets) accessible for review or editing.
</usage>

<params>
Input:
- filepath (string, required): Path to the downloaded .crx file (relative to workspace root or absolute).
  - Example: `./extensions/aapbdbdomjkkjkaonfhkkikfgjllcleb.crx`
- name (string, optional): Custom name for the unpacked extension directory. If not provided, uses extension ID or filename.
- outputDir (string, optional): Directory to save unpacked contents (relative to workspace root). If not provided, creates directory in workspace root.
</params>

<result>
Output:
{{ if not .DisabledMarkdown }}
- Confirmation message with:
  - Unpack status (success/failure)
  - Path to the unpacked extension directory
{{ end }}
StructuredOutput:
```json
# Example (Success):
{
   "filepath": "/workspace/unpacked/aapbdbdomjkkjkaonfhkkikfgjllcleb/"
}

# Example (Failure):
{
   "filepath": ""
}
</result>

<use_cases>
Example use cases:
- "Unpack the extension I just downloaded"
- "Extract the contents of ./extensions/abc123.crx"
- "Unpack uBlock Origin to ./source/ folder"
- "I need to inspect the manifest.json of the downloaded extension"
</use_cases>