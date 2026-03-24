Unpacks a downloaded Chrome extension (.crx file) into a directory structure for inspection or modification.

⚠️ PREREQUISITE: You MUST call `crx3_scan {}` FIRST to retrieve and cache the absolute root path.
Only after that, call `crx3_unpack` using filepaths that are RELATIVE to the workspace root (e.g., "./extensions/app.crx").

<usage>
Use this tool when the user wants to extract the contents of a .crx file. The tool unpacks the extension into a directory, making source files (manifest.json, scripts, assets) accessible for review or editing.
</usage>

<example>
- "Unpack the extension I just downloaded"
- "Extract the contents of ./extensions/abc123.crx"
- "Unpack uBlock Origin to ./source/ folder"
- "I need to inspect the manifest.json of the downloaded extension"
</example>