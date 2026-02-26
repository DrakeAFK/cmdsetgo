#!/bin/bash

# cmdsetgo uninstall script
# This script runs the internal uninstall command to remove shell hooks and aliases.

set -e

echo "🗑️ Starting cmdsetgo uninstallation..."

# 1. Build the binary to ensure we have the latest uninstall logic
if command -v go &> /dev/null; then
    echo "📦 Rebuilding cmdsetgo for uninstall..."
    go build -o cmdsetgo ./cmd/cmdsetgo
    mkdir -p bin
    mv cmdsetgo bin/
fi

# 2. Look for the binary
BIN_PATH="./bin/cmdsetgo"
if [[ ! -f "$BIN_PATH" ]]; then
    # Try the root directory if bin/ doesn't exist
    BIN_PATH="./cmdsetgo"
fi

# 2. Run uninstall logic
if [[ -f "$BIN_PATH" ]]; then
    echo "📦 Running uninstall from $BIN_PATH..."
    "$BIN_PATH" uninstall
elif command -v cmdsetgo &> /dev/null; then
    echo "📦 Running uninstall from global PATH..."
    cmdsetgo uninstall
else
    echo "❌ Error: Could not find cmdsetgo binary. Please run this script from the cmdsetgo repo root."
    exit 1
fi

echo ""
echo "📢 IMPORTANT: Complete cleanup requires a fresh terminal"
echo "──────────────────────────────────────────────────────"
echo "1. The shell hook and aliases have been REMOVED from your configuration."
echo "2. Existing terminal tabs will KEEP the hook active until they are closed."
echo "3. Please open a NEW terminal tab to apply these changes."
echo "   (Sourcing your config file is NOT enough to deactivate aliases/functions)"
echo "──────────────────────────────────────────────────────"

# 3. Cleanup local binary if it was built in bin/
if [[ -d "./bin" ]]; then
    echo "🧹 Cleaning up local bin/ directory..."
    rm -rf "./bin"
fi

echo ""
echo "✨ Uninstallation of shell components complete!"
echo "If you want to remove the source code, you can now safely delete the 'cmdsetgo' directory."
