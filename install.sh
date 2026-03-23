#!/bin/bash
set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Logging functions
log_info() { echo -e "${GREEN}ℹ️  $1${NC}"; }
log_warn() { echo -e "${YELLOW}⚠️  $1${NC}"; }
log_err()  { echo -e "${RED}❌ $1${NC}"; }

# Cleanup function: remove temporary files on exit (success or failure)
cleanup() {
  if [ -n "${WORK_DIR}" ] && [ -d "${WORK_DIR}" ]; then
    rm -rf "${WORK_DIR}"
    log_info "🧹 Cleaned up temporary files"
  fi
}
# Trap EXIT, INT (Ctrl+C), and TERM signals to ensure cleanup always runs
trap cleanup EXIT INT TERM

echo "🚀 Installing crx3... ${WORK_DIR}"

# =============================================================================
# 1. Detect OS and architecture (lowercase, Go-style naming)
# =============================================================================

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Normalize architecture to Go conventions
case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  i386|i686) ARCH="386" ;;
  *) log_err "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Validate supported operating systems
case "$OS" in
  linux|darwin) ;;
  *) log_err "Unsupported OS: $OS"; exit 1 ;;
esac

log_info "Detected: ${OS} ${ARCH}"

# =============================================================================
# 2. Fetch latest version from GitHub API
# =============================================================================

LATEST_URL="https://api.github.com/repos/mmadfox/go-crx3/releases/latest"
VERSION=$(curl -s "$LATEST_URL" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/' | sed 's/^v//')

if [ -z "$VERSION" ]; then
  log_err "Failed to fetch latest version from GitHub API"
  exit 1
fi

log_info "Latest version: v${VERSION}"

# =============================================================================
# 3. Build download URL (match GoReleaser's naming convention)
# =============================================================================

# Note: project_name in .goreleaser.yml is "crx3", but actual archive names
# use the module name "go-crx3". Real files look like:
# go-crx3_1.6.0_darwin_arm64.tar.gz
PROJECT_NAME="go-crx3"
FILE_NAME="${PROJECT_NAME}_${VERSION}_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/mmadfox/go-crx3/releases/download/v${VERSION}/${FILE_NAME}"

log_info "Downloading: ${FILE_NAME}"

# =============================================================================
# 4. Create working directory in CURRENT folder (not /tmp)
# =============================================================================

# Use PID ($$) to ensure uniqueness for parallel runs
# Prefix with "." to hide from normal ls output
WORK_DIR="./.crx3-install-$$"
mkdir -p "${WORK_DIR}"

# =============================================================================
# 5. Download the archive with error handling
# =============================================================================
if ! curl -sLf --fail "$DOWNLOAD_URL" -o "${WORK_DIR}/${FILE_NAME}"; then
  log_err "Failed to download: ${DOWNLOAD_URL}"
  log_warn "Available files in release v${VERSION}:"
  curl -s "https://api.github.com/repos/mmadfox/go-crx3/releases/tags/v${VERSION}" | \
    grep '"name"' | grep '\.tar\.gz' | sed 's/.*"name": "\([^"]*\)".*/  - \1/'
  exit 1
fi

# =============================================================================
# 6. (Optional) Verify SHA256 checksum for integrity
# =============================================================================
CHECKSUM_URL="https://github.com/mmadfox/go-crx3/releases/download/v${VERSION}/checksums.txt"
if curl -sLf --fail "$CHECKSUM_URL" -o "${WORK_DIR}/checksums.txt" 2>/dev/null; then
  EXPECTED=$(grep "$FILE_NAME" "${WORK_DIR}/checksums.txt" | awk '{print $1}')
  if [ -n "$EXPECTED" ]; then
    ACTUAL=$(sha256sum "${WORK_DIR}/${FILE_NAME}" | awk '{print $1}')
    if [ "$EXPECTED" != "$ACTUAL" ]; then
      log_err "Checksum mismatch!"
      log_warn "Expected: $EXPECTED"
      log_warn "Actual:   $ACTUAL"
      exit 1
    fi
    log_info "✅ Checksum verified"
  fi
fi

# =============================================================================
# 7. Extract the archive
# =============================================================================

log_info "Extracting archive..."
if ! tar -xzf "${WORK_DIR}/${FILE_NAME}" -C "${WORK_DIR}"; then
  log_err "Failed to extract archive"
  exit 1
fi

# =============================================================================
# 8. Locate the binary (may be in a subdirectory like crx3_*/crx3)
# =============================================================================

BIN_PATH=$(find "${WORK_DIR}" -name "go-crx3" -type f | head -1)
if [ -z "$BIN_PATH" ]; then
  log_err "Binary 'crx3' not found in archive"
  exit 1
fi

log_info "Found binary: ${BIN_PATH}"

# =============================================================================
# 9. Determine installation directory
# =============================================================================

# Prefer /usr/local/bin if writable, otherwise use ~/.local/bin
# if [ -w "/usr/local/bin" ]; then
#   INSTALL_DIR="/usr/local/bin"
# else
#   INSTALL_DIR="$HOME/.local/bin"
#   mkdir -p "$INSTALL_DIR"
#   log_warn "No write access to /usr/local/bin"
#   log_warn "Installing to ${INSTALL_DIR}"
#   log_warn "Add to PATH: export PATH=\"${INSTALL_DIR}:\$PATH\""
# fi

# =============================================================================
# 10. Install the binary
# =============================================================================

# log_info "Installing to ${INSTALL_DIR}/crx3"
# cp "${BIN_PATH}" "${INSTALL_DIR}/crx3"
# chmod +x "${INSTALL_DIR}/crx3"

# =============================================================================
# 11. Final verification and user feedback
# =============================================================================

# echo ""
# if command -v crx3 >/dev/null 2>&1; then
#   log_info "✅ crx3 installed successfully!"
#   crx3 version
# else
#   log_warn "⚠️  crx3 installed to ${INSTALL_DIR}, but not in PATH"
#   log_info "Run: export PATH=\"${INSTALL_DIR}:\$PATH\""
# fi

# echo ""
# echo "📚 Usage: crx3 --help"

# =============================================================================
# Cleanup happens automatically via trap (no explicit call needed)
# =============================================================================