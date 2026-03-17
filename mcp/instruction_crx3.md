# The CRX3 MCP Server

These instructions describe how to efficiently work with the CRX3 tools set using the MCP server. You can load this file directly into a session where the CRX3 MCP server is connected.

## Detecting a CRX3 Workdir

At the start of every session, you MUST use the `crx3_get_workdir` tool to learn about the CRX3 workspace. ONLY if you are in a CRX3 workspace, you MUST run `crx3_list_extensions` immediately afterwards to identify any existing extensions. The rest of these instructions apply whenever that tool indicates that the user is in a CRX3 workspace.