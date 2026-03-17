# crx3_search_chrome_extension

Searches for Chrome extensions by name using DuckDuckGo and returns structured information about matching extensions.

<usage>
Use this tool when the user wants to find a Chrome browser extension by its name. The tool performs a search and returns the extension's name, URL (from the Chrome Web Store), and its unique Extension ID.
</usage>

<params>
Input:
- name (string): The name or keyword of the Chrome extension to search for.
</params>

<result>
Output:
- A markdown-formatted list of found extensions, including:
  - Name
  - Clickable link to the extension page
  - Extension ID (useful for direct installation or further processing)
</result>

<examples>
Example use cases:
- "Find the Chrome extension called 'React Developer Tools'"
- "Search for ad blockers on Chrome"
- "Get the Extension ID for 'uBlock Origin'"
</examples>