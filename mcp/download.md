# crx3_download

Downloads Chrome extensions as .crx files using an Extension ID or Chrome Web Store URL.

<usage>
Use this tool when the user wants to download a Chrome extension package (.crx file). The tool accepts either a direct Extension ID or a Chrome Web Store URL, downloads the extension, and saves it to the specified or default location.
</usage>

<params>
Input:
- url (string, required): The extensionId or Chrome Web Store URL to download.
  - Accepts formats: `aapbdbdomjkkjkaonfhkkikfgjllcleb` or `https://chrome.google.com/webstore/detail/.../aapbdbdomjkkjkaonfhkkikfgjllcleb`
- path (string, optional): Path to save the downloaded .crx file. If not provided, uses default download location.
</params>

<result>
Output:
{{ if not .DisabledMarkdown }}
- A confirmation message indicating:
  - Download status (success/failure)
  - Full filepath to the downloaded .crx file
{{ end }}
StructuredOutput:
```json
# Example (Success):
{
   "success": true,
   "filepath": "/path/to/extension.crx"
}

# Example (Failure):
{
   "success": false,
   "filepath": ""
}
</result>

<use_cases>
Example use cases:
- "Download the extension with ID 'aapbdbdomjkkjkaonfhkkikfgjllcleb'"
- "Save the uBlock Origin extension to ./downloads/"
- "Download the extension from this URL: https://chrome.google.com/webstore/detail/..."
- "Get the .crx file for the extension I just searched for"
</use_cases>