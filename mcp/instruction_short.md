# CRX3 MCP Server — Condensed Instructions

## ⚡ SESSION INIT (MANDATORY)
• On first message: AUTO-CALL `crx3_workspace {}` → cache `absoluteRootPath`
• Use cached root for ALL path resolution; re-call only on errors
• If workspace fails: halt file ops, notify user, wait for resolution

## 🚫 PATH RULES (NON-NEGOTIABLE)
• Tool inputs: ALWAYS workspace-relative, forward slashes, start with `./` or `../`
• NEVER pass absolute paths to tools: `/home/...`, `C:\...` → reject & correct first
• User display: MAY show absolute paths using cached `absoluteRootPath`
• Forbidden chars in paths: `* ? : | < > \` → auto-sanitized to `_`

## 🛡️ PRE-FLIGHT CHECK (Before file ops)
Before calling: `pack|unpack|getid|base64|zip|unzip`
→ Ensure `absoluteRootPath` is cached
→ If not: call `crx3_workspace {}` first, cache, then proceed
• Exceptions (no workspace required): `search|download|scan|version`

## 🔍 FILE NOT FOUND PROTOCOL
1. STOP — don't retry same path
2. Call `crx3_workspace {}` → verify root
3. Call `crx3_scan {filter:[...]} or {}` → locate file
4. Use returned `filepath` directly (already valid relative path)
5. If still not found: offer `search+download` or ask user to place file in workspace

## 🧰 TOOLS REFERENCE (Essentials)

| Tool | Purpose | Key Params | Critical Notes |
|------|---------|------------|---------------|
| `search` | Find ext by name | `query`, `limit` | Verify source before download |
| `download` | Get .crx by ID/URL | `extensionId`|`url`, `path` | path: relative, auto-creates dirs |
| `workspace` | Get absolute root | none | Use for display ONLY, not tool inputs |
| `unpack` | Extract .crx → dir | `filepath`, `outputDir?` | omit `outputDir` → auto: `./unpacked/{id}/`; NEVER include ext name in outputDir |
| `pack` | Dir/zip → signed .crx | `source`, `outputDir?`, `name?`, `privateKey?` | reuse .pem to preserve ID; new key = new ID |
| `scan` | List .crx in workspace | `limit`, `filter`, `sortBy` | use `filepath` from results for other tools |
| `unzip`/`zip` | Archive ops | `filepath`/`source`, `outputDir?` | relative paths only |
| `base64` | Encode file → string | `filepath` | warn if >1MB (+33% size) |
| `getid` | Extract extension ID | `filepath` (.crx or dir) | ID = hash(pubkey); same key+manifest = same ID |
| `version` | Show tool version | none | informational |

## 🔑 KEY MANAGEMENT (ID Preservation)
• Extension ID = hash(public_key from manifest/.pem)
• To preserve ID on repack: save `.pem` from first `pack`, reuse via `privateKey` param
• Losing .pem = losing update chain for that extension ID

## 🔄 TYPICAL WORKFLOWS (Condensed)

### Download → Inspect → Modify → Repack
```
search{"query":"..."} → select ID → download{extensionId:"..."}
→ unpack{filepath:"./ext/id.crx"} → [edit ./unpacked/id/]
→ pack{source:"./unpacked/id/", privateKey:"./packed/id.pem"}
→ getid{filepath:"./packed/new.crx"} // verify ID unchanged
```

### Backup Cycle
```
scan{filter:[...]} → for each: 
  getid → unpack → zip{source:"./backup/src/"} → base64{filepath:"./ext.crx"}
```

### Dev Loop
```
zip{source:"./src/"} → pack{source:"./src/"} → getid → [test] → iterate
```

## ✅ CONTEXT TRACKING (Cache per op)
• search → `extensionId`, `name`
• download → `filepath`, `extensionId`
• unpack → `outputDir`, `sourceCrx`
• pack → `filepath`, `privateKey`, `extensionID`
• getid → `extensionID` (for verification)

## 🚨 ERROR QUICK-GUIDE
• "File not found" → `scan` or `workspace` to diagnose
• "Invalid manifest" → check `manifest.json` exists + required fields
• "ID extraction failed" → likely unsigned/corrupted extension
• "Path not allowed" → ensure relative, forward slashes, no forbidden chars
• Multiple search results → present options, await user selection

## 💡 PRO TIPS
• Always validate paths before tool calls — prevents 90% of errors
• Use `crx3_scan` proactively when user says "my file" without path
• Cache workspace root once; don't re-call unless error recovery
• When in doubt: workspace → scan → operate with returned paths
```