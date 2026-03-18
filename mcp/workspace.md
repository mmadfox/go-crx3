# crx3_workspace

Returns the path to the workspace directory used by the CRX3 server to store downloaded Chrome extensions.

<usage>
Use this tool when you need to know where extensions are saved on disk. This is useful for:
- Verifying download location
- Debugging file operations
- Referencing the local path of an installed extension
</usage>

<params>
No input parameters required.
</params>

<result>
Returns a string with the absolute path to the workspace root.
Structured output: { "path": "/path/to/workspace" }
</result>