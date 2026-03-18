# Download Chrome Extension

Downloads a Chrome extension from the Chrome Web Store by extension ID or URL.

## Parameters

- `idOrUrl`: The extension ID (e.g. "gighmmpiobklfepjocnamgkkbiglidom") or full Chrome Web Store URL (e.g. "https://chromewebstore.google.com/detail/adblock/gighmmpiobklfepjocnamgkkbiglidom")
- `outfile`: (Optional) Path where to save the downloaded .crx file. If not specified, saves to current directory as "extension.crx"
- `unpack`: (Optional) Whether to unpack the extension after downloading. Defaults to true.

## Returns

Returns a message indicating the download location and status. If unpack is true, also indicates the unpack location.

## Example

```json
{
  "idOrUrl": "https://chromewebstore.google.com/detail/adblock/gighmmpiobklfepjocnamgkkbiglidom",
  "outfile": "/tmp/adblock.crx",
  "unpack": true
}
```

This would download the Adblock extension to /tmp/adblock.crx and unpack it to /tmp/adblock/.