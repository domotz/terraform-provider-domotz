#!/bin/bash

# Build Terraform Provider for Domotz for all platforms
# Creates binaries for distribution across different operating systems and architectures

set -e  # Exit on error

VERSION="1.0.0"
BINARY_NAME="terraform-provider-domotz"
BUILD_DIR="dist"

echo "=========================================="
echo "Building Terraform Provider for Domotz"
echo "Version: $VERSION"
echo "=========================================="
echo ""

# Clean previous builds
echo "🧹 Cleaning previous builds..."
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"
echo "✅ Clean complete"
echo ""

# Ensure dependencies are up to date
echo "📦 Downloading dependencies..."
go mod download
go mod tidy
echo "✅ Dependencies ready"
echo ""

# Build function
build_binary() {
    local os=$1
    local arch=$2
    local arm_version=$3

    local output_dir="$BUILD_DIR/${BINARY_NAME}_${VERSION}_${os}_${arch}"
    local binary_name="${BINARY_NAME}_v${VERSION}"

    # Add .exe extension for Windows
    if [ "$os" = "windows" ]; then
        binary_name="${binary_name}.exe"
    fi

    mkdir -p "$output_dir"

    echo "🔨 Building for $os/$arch${arm_version:+ (ARM v$arm_version)}..."

    # Set environment variables for cross-compilation
    export GOOS=$os
    export GOARCH=$arch
    export CGO_ENABLED=0

    if [ -n "$arm_version" ]; then
        export GOARM=$arm_version
    fi

    # Build the binary
    go build \
        -trimpath \
        -ldflags="-s -w -X main.version=${VERSION}" \
        -o "$output_dir/$binary_name" \
        .

    # Copy documentation
    cp LICENSE "$output_dir/"
    cp README.md "$output_dir/"

    # Create archive
    echo "📦 Creating archive..."
    cd "$BUILD_DIR"
    if [ "$os" = "windows" ]; then
        zip -q -r "${BINARY_NAME}_${VERSION}_${os}_${arch}.zip" "$(basename $output_dir)"
    else
        tar -czf "${BINARY_NAME}_${VERSION}_${os}_${arch}.tar.gz" "$(basename $output_dir)"
    fi
    cd - > /dev/null

    echo "✅ Built $os/$arch${arm_version:+ (ARM v$arm_version)}"
    echo ""
}

# Build for all platforms
echo "Starting multi-platform build..."
echo ""

# macOS (Darwin)
build_binary "darwin" "amd64" ""
build_binary "darwin" "arm64" ""

# Linux
build_binary "linux" "amd64" ""
build_binary "linux" "arm64" ""
build_binary "linux" "arm" "7"

# Windows
build_binary "windows" "amd64" ""

# Generate SHA256 checksums
echo "🔐 Generating checksums..."
cd "$BUILD_DIR"
if command -v shasum &> /dev/null; then
    shasum -a 256 *.zip *.tar.gz > "${BINARY_NAME}_${VERSION}_SHA256SUMS"
elif command -v sha256sum &> /dev/null; then
    sha256sum *.zip *.tar.gz > "${BINARY_NAME}_${VERSION}_SHA256SUMS"
else
    echo "⚠️  Warning: Neither shasum nor sha256sum found, skipping checksums"
fi
cd - > /dev/null
echo "✅ Checksums generated"
echo ""

# Display build summary
echo "=========================================="
echo "✅ Build Complete!"
echo "=========================================="
echo ""
echo "Built binaries for:"
echo "  • macOS (Intel): darwin/amd64"
echo "  • macOS (Apple Silicon): darwin/arm64"
echo "  • Linux (64-bit): linux/amd64"
echo "  • Linux (ARM64): linux/arm64"
echo "  • Linux (ARMv7): linux/arm"
echo "  • Windows (64-bit): windows/amd64"
echo ""
echo "Output directory: $BUILD_DIR"
echo ""
echo "Contents:"
ls -lh "$BUILD_DIR" | grep -E '\.(zip|tar\.gz|SUMS)$' | awk '{print "  • " $9 " (" $5 ")"}'
echo ""
echo "=========================================="
echo "📦 Ready for distribution!"
echo "=========================================="
