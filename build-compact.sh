#!/bin/bash

# Compact Build Script for fps-go-brr
# Based on techniques from: https://words.filippo.io/shrink-your-go-binaries-with-this-one-weird-trick/
#
# This script builds a space-optimized version of the fps-go-brr binary using:
# 1. Go linker flags to strip debugging information (-s -w)
# 2. UPX compression (if available)

set -e

BINARY_NAME="fps-go-brr"
COMPACT_BINARY_NAME="fps-go-brr-compact"

echo "Building compact version of $BINARY_NAME..."

# Step 1: Build with linker flags to strip debugging info
echo "Step 1: Building with -ldflags='-s -w' to strip debugging information..."
go build -ldflags="-s -w" -o "$COMPACT_BINARY_NAME" .

# Get initial size
INITIAL_SIZE=$(stat -c%s "$COMPACT_BINARY_NAME" 2>/dev/null || stat -f%z "$COMPACT_BINARY_NAME" 2>/dev/null || echo "unknown")
echo "Binary size after stripping: $INITIAL_SIZE bytes"

# Step 2: Apply UPX compression if available
if command -v upx >/dev/null 2>&1; then
    echo "Step 2: Applying UPX compression with --brute flag..."
    upx --brute "$COMPACT_BINARY_NAME"
    
    # Get final size
    FINAL_SIZE=$(stat -c%s "$COMPACT_BINARY_NAME" 2>/dev/null || stat -f%z "$COMPACT_BINARY_NAME" 2>/dev/null || echo "unknown")
    echo "Final binary size after UPX compression: $FINAL_SIZE bytes"
    
    if [ "$INITIAL_SIZE" != "unknown" ] && [ "$FINAL_SIZE" != "unknown" ]; then
        REDUCTION=$(echo "scale=1; (1 - $FINAL_SIZE / $INITIAL_SIZE) * 100" | bc -l 2>/dev/null || echo "N/A")
        echo "Total size reduction: ~${REDUCTION}%"
    fi
else
    echo "Step 2: UPX not found - skipping compression"
    echo "To install UPX:"
    echo "  - Ubuntu/Debian: sudo apt install upx-ucl"
    echo "  - macOS: brew install upx"
    echo "  - Arch Linux: sudo pacman -S upx"
    echo "  - Or download from: https://upx.github.io/"
fi

echo ""
echo "Compact build complete: $COMPACT_BINARY_NAME"
echo ""
echo "Note: The compact binary may have slightly slower startup time due to UPX decompression."
echo "For production use, consider whether the size savings are worth the startup overhead."

# Make the script executable
chmod +x "$COMPACT_BINARY_NAME"

# Compare with regular build if it exists
if [ -f "$BINARY_NAME" ]; then
    REGULAR_SIZE=$(stat -c%s "$BINARY_NAME" 2>/dev/null || stat -f%z "$BINARY_NAME" 2>/dev/null || echo "unknown")
    if [ "$REGULAR_SIZE" != "unknown" ] && [ "$FINAL_SIZE" != "unknown" ]; then
        TOTAL_REDUCTION=$(echo "scale=1; (1 - $FINAL_SIZE / $REGULAR_SIZE) * 100" | bc -l 2>/dev/null || echo "N/A")
        echo "Comparison with regular build:"
        echo "  Regular build: $REGULAR_SIZE bytes"
        echo "  Compact build: $FINAL_SIZE bytes"
        echo "  Total reduction: ~${TOTAL_REDUCTION}%"
    fi
fi