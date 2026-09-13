#!/bin/bash
set -e

# DeepSec Installer
# Installs DeepSec CLI tool - All-in-one Cybersecurity Scanner
# Usage: curl -sSL https://pixelcity.dev/deepsec/install.sh | sh

DEEPSEC_VERSION="${DEEPSEC_VERSION:-v1.0.0}"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
BASE_URL="https://pixelcity.dev/deepsec/releases"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

info()    { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[OK]${NC} $1"; }
warn()    { echo -e "${YELLOW}[WARN]${NC} $1"; }
error()   { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

echo ""
echo -e "${CYAN}  ____                           ____  ______  __ __ ${NC}"
echo -e "${CYAN} |  _ \\  ___  ___ _   _ _ __ ___|  _ \\|  _ \\ \\ V / |${NC}"
echo -e "${CYAN} | | | |/ _ \\/ __| | | | '__/ __| | | | | | | | |  |${NC}"
echo -e "${CYAN} | |_| |  __/\\__ \\ |_| | | | (__| |_| | |_| | | |  |${NC}"
echo -e "${CYAN} |____/ \\___||___/\\__,_|_|  \\___|____/|____/ |_|__|${NC}"
echo -e "${CYAN}                                     v${DEEPSEC_VERSION}${NC}"
echo ""
echo -e "  All-in-one Cybersecurity CLI Tool for Developers & Enterprise"
echo ""

# Check dependencies
check_deps() {
    info "Checking dependencies..."
    
    for cmd in curl uname; do
        if ! command -v "$cmd" &>/dev/null; then
            error "$cmd is required but not installed."
        fi
    done
    
    success "Dependencies OK"
}

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case $ARCH in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        armv7l|armhf)
            ARCH="arm"
            ;;
        *)
            error "Unsupported architecture: $ARCH"
            ;;
    esac
    
    case $OS in
        linux)
            PLATFORM="linux"
            ;;
        darwin)
            PLATFORM="darwin"
            ;;
        mingw*|msys*|cygwin*)
            PLATFORM="windows"
            ;;
        *)
            error "Unsupported OS: $OS"
            ;;
    esac
    
    BINARY_NAME="deepsec"
    if [ "$PLATFORM" = "windows" ]; then
        BINARY_NAME="deepsec.exe"
    fi
    
    DETECTED="${PLATFORM}-${ARCH}"
    info "Detected platform: ${DETECTED}"
}

# Download and install
install_binary() {
    info "Downloading DeepSec ${DEEPSEC_VERSION}..."
    
    DOWNLOAD_URL="${BASE_URL}/${DEEPSEC_VERSION}/deepsec-${DETECTED}"
    if [ "$PLATFORM" = "windows" ]; then
        DOWNLOAD_URL="${DOWNLOAD_URL}.exe"
    fi
    
    # Create install directory
    mkdir -p "$INSTALL_DIR"
    
    TEMP_FILE=$(mktemp)
    
    info "URL: ${DOWNLOAD_URL}"
    
    HTTP_CODE=$(curl -w "%{http_code}" -sL -o "$TEMP_FILE" "$DOWNLOAD_URL" 2>/dev/null)
    
    if [ "$HTTP_CODE" != "200" ]; then
        rm -f "$TEMP_FILE"
        error "Failed to download DeepSec. HTTP status: $HTTP_CODE"
    fi
    
    # Verify it's a valid binary
    file "$TEMP_FILE" | grep -qi "executable\|ELF\|Mach-O" || {
        # Check if it's HTML (404 page)
        if head -c 100 "$TEMP_FILE" | grep -qi "<!DOCTYPE\|<html"; then
            rm -f "$TEMP_FILE"
            error "Downloaded file is not a binary. The release may not exist for your platform."
        fi
    }
    
    chmod +x "$TEMP_FILE"
    mv "$TEMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"
    
    success "Installed to ${INSTALL_DIR}/${BINARY_NAME}"
}

# Update PATH if needed
update_path() {
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        SHELL_RC=""
        if [ -f "$HOME/.bashrc" ]; then
            SHELL_RC="$HOME/.bashrc"
        elif [ -f "$HOME/.zshrc" ]; then
            SHELL_RC="$HOME/.zshrc"
        fi
        
        if [ -n "$SHELL_RC" ]; then
            echo "export PATH=\"\$HOME/.local/bin:\$PATH\"" >> "$SHELL_RC"
            export PATH="$HOME/.local/bin:$PATH"
            warn "Added $INSTALL_DIR to PATH in $SHELL_RC"
            warn "Run 'source $SHELL_RC' or restart your terminal"
        else
            warn "Add $INSTALL_DIR to your PATH manually"
        fi
    fi
}

# Verify installation
verify_install() {
    if command -v deepsec &>/dev/null; then
        DEEPSEC_PATH=$(command -v deepsec)
        success "DeepSec installed successfully!"
        echo ""
        echo -e "  ${GREEN}Location:${NC} $DEEPSEC_PATH"
        echo -e "  ${GREEN}Version:${NC}  $(deepsec --version 2>/dev/null || echo 'installed')"
        echo ""
        echo -e "  ${CYAN}Quick Start:${NC}"
        echo -e "    deepsec scan .              ${CYAN}# Scan current directory${NC}"
        echo -e "    deepsec scan . --scanner sast,secrets  ${CYAN}# Specific scanners${NC}"
        echo -e "    deepsec init                ${CYAN}# Create config file${NC}"
        echo -e "    deepsec --help              ${CYAN}# Show all commands${NC}"
        echo ""
        echo -e "  ${YELLOW}Docs: https://github.com/Pixelcity-dev/Deepsec${NC}"
        echo ""
    else
        warn "Binary installed but not in PATH. Try: ${INSTALL_DIR}/deepsec --help"
    fi
}

# Main
main() {
    check_deps
    detect_platform
    install_binary
    update_path
    verify_install
}

main "$@"
