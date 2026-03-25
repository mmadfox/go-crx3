#!/bin/bash
set -e

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() { echo -e "${GREEN}ℹ️  $1${NC}"; }
log_warn() { echo -e "${YELLOW}⚠️  $1${NC}"; }
log_err()  { echo -e "${RED}❌ $1${NC}"; }

# Cleanup temp dir
cleanup() {
  [[ -n "${WORK_DIR}" && -d "${WORK_DIR}" ]] && rm -rf "${WORK_DIR}"
}
trap cleanup EXIT INT TERM

log_info "🚀 Installing crx3..."

# =============================================================================
# 1. Detect OS/Arch in GoReleaser format (Capitalized)
# =============================================================================
OS_RAW=$(uname -s)
ARCH_RAW=$(uname -m)

# OS: Linux, Darwin (capitalized)
case "$OS_RAW" in
  Linux)  OS="Linux" ;;
  Darwin) OS="Darwin" ;;
  *) log_err "Unsupported OS: $OS_RAW"; exit 1 ;;
esac

# Arch: Amd64, Arm64, 386 (capitalized first letter)
case "$ARCH_RAW" in
  x86_64)          ARCH="Amd64" ;;
  aarch64|arm64)   ARCH="Arm64" ;;
  i386|i686)       ARCH="386" ;;
  *) log_err "Unsupported architecture: $ARCH_RAW"; exit 1 ;;
esac

log_info "Detected: ${OS} ${ARCH}"

# =============================================================================
# 2. Fetch latest version (or use hardcoded)
# =============================================================================
VERSION="${CRX3_VERSION:-}"  # Allow override via env var
if [ -z "$VERSION" ]; then
  LATEST_URL="https://api.github.com/repos/mmadfox/go-crx3/releases/latest"
  VERSION=$(curl -s "$LATEST_URL" | grep '"tag_name"' | sed -E 's/.*"v?([^"]+)".*/\1/')
fi

if [ -z "$VERSION" ]; then
  log_err "Failed to determine version"
  exit 1
fi
log_info "Version: v${VERSION}"

# =============================================================================
# 3. Build GoReleaser-compatible URL
#    Format: crx3_{VER}_{OS}_{ARCH}.tar.gz
#    OS: Linux/Darwin | ARCH: Amd64/Arm64/386
# =============================================================================
PROJECT="crx3"
FILE_NAME="${PROJECT}_${VERSION}_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/mmadfox/go-crx3/releases/download/v${VERSION}/${FILE_NAME}"

log_info "Downloading: ${FILE_NAME}"

# =============================================================================
# 4. Prepare temp directory
# =============================================================================
WORK_DIR="./.crx3-install-$$"
mkdir -p "${WORK_DIR}"

# =============================================================================
# 5. Download with validation
# =============================================================================
if ! curl -sLf --fail "$DOWNLOAD_URL" -o "${WORK_DIR}/${FILE_NAME}"; then
  log_err "Download failed: ${DOWNLOAD_URL}"
  log_warn "Available assets:"
  curl -s "https://api.github.com/repos/mmadfox/go-crx3/releases/tags/v${VERSION}" | \
    grep '"browser_download_url"' | sed 's/.*: "\([^"]*\)".*/  • \1/'
  exit 1
fi

# Validate it's actually a gzip archive (catch 404 HTML pages)
if ! file "${WORK_DIR}/${FILE_NAME}" | grep -q "gzip compressed"; then
  log_err "Downloaded file is not a valid archive!"
  head -c 300 "${WORK_DIR}/${FILE_NAME}" | cat -v
  exit 1
fi

# =============================================================================
# 6. Optional: Verify checksum
# =============================================================================
CHECKSUM_URL="https://github.com/mmadfox/go-crx3/releases/download/v${VERSION}/checksums.txt"
if curl -sLf --fail "$CHECKSUM_URL" -o "${WORK_DIR}/checksums.txt" 2>/dev/null; then
  EXPECTED=$(grep "$FILE_NAME" "${WORK_DIR}/checksums.txt" 2>/dev/null | awk '{print $1}')
  if [ -n "$EXPECTED" ]; then
    ACTUAL=$(sha256sum "${WORK_DIR}/${FILE_NAME}" | awk '{print $1}')
    if [ "$EXPECTED" != "$ACTUAL" ]; then
      log_err "❌ Checksum mismatch!"
      exit 1
    fi
    log_info "✅ Checksum verified"
  fi
fi

# =============================================================================
# 7. Extract & locate binary
# =============================================================================
log_info "Extracting..."
tar -xzf "${WORK_DIR}/${FILE_NAME}" -C "${WORK_DIR}"

BIN_PATH=$(find "${WORK_DIR}" -name "crx3" -type f -executable | head -1)
if [ -z "$BIN_PATH" ]; then
  log_err "Binary 'crx3' not found in archive"
  exit 1
fi

# =============================================================================
# 8. Install to system path
# =============================================================================
if [ -w "/usr/local/bin" ]; then
  INSTALL_DIR="/usr/local/bin"
else
  INSTALL_DIR="$HOME/.local/bin"
  mkdir -p "$INSTALL_DIR"
  log_warn "Installing to user path: ${INSTALL_DIR}"
  log_warn "Add to PATH: export PATH=\"${INSTALL_DIR}:\$PATH\""
fi

log_info "Installing to ${INSTALL_DIR}/crx3"
cp "$BIN_PATH" "${INSTALL_DIR}/crx3"
chmod +x "${INSTALL_DIR}/crx3"

# =============================================================================
# 9. Verify & finish
# =============================================================================
echo ""
if command -v crx3 >/dev/null 2>&1; then
  log_info "✅ crx3 installed successfully!"
  crx3 version
else
  log_warn "⚠️  Binary installed but not in PATH"
  log_info "Run: export PATH=\"${INSTALL_DIR}:\$PATH\""
fi
echo "📚 Usage: crx3 --help"