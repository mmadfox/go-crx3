Generates the Chrome Extension ID from a .crx file or unpacked extension directory.

The identifier is generated from the hash of the public key, which is located in the extension header or declared in the key field of the manifest. If the key is specified in the manifest, the public key is taken from there; otherwise, the search continues in the header.

<usage>
Use this tool when the user wants to get the Chrome extension ID. The tool reads the extension file or directory and extracts or generates the unique extension identifier.
</usage>

<example>
- "Get the extension ID for the downloaded extension"
- "What is the ID of the extension in ./extensions/abc123.crx"
- "Get extension ID from unpacked extension in ./unpacked/react-devtools/"
</example>
