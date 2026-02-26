#!/bin/bash

# cmdsetgo bootstrap script
# This script builds the binary and runs the internal install command.

set -e

echo "🚀 Starting cmdsetgo installation..."

# 1. Check for Go
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed. Please install Go from https://golang.org/dl/"
    exit 1
fi

# 2. Build the binary
echo "📦 Building cmdsetgo..."
go build -o cmdsetgo ./cmd/cmdsetgo

# 3. Create a bin directory if it doesn't exist (local to repo)
mkdir -p bin
mv cmdsetgo bin/

# 4. Check if cmdsetgo is in PATH
if ! command -v cmdsetgo &> /dev/null; then
    echo "⚠️  cmdsetgo binary was built successfully in ./bin/cmdsetgo"
    echo "💡 To use it globally, add it to your PATH or move it to /usr/local/bin:"
    echo "   sudo mv ./bin/cmdsetgo /usr/local/bin/"
    echo ""
    echo "Running local install for now..."
    ./bin/cmdsetgo install
else
    echo "✅ cmdsetgo is already in your PATH."
    cmdsetgo install
fi

echo ""
echo "✨ Installation complete!"
echo "Please restart your terminal or source your rc file to start recording."
echo "Run 'cmdsetgo status' to verify."
