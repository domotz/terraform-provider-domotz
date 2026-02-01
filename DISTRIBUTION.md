# Terraform Provider for Domotz - Distribution Guide

## Version 1.0.0

This document describes the available binary distributions for the Terraform Provider for Domotz.

## Available Platforms

The provider has been built for the following platforms:

### macOS (Darwin)
- **Intel (x86_64)**: `terraform-provider-domotz_1.0.0_darwin_amd64.tar.gz` (5.7MB)
  - Compatible with macOS on Intel processors
- **Apple Silicon (ARM64)**: `terraform-provider-domotz_1.0.0_darwin_arm64.tar.gz` (5.2MB)
  - Compatible with M1, M2, M3, M4 Macs

### Linux
- **64-bit (x86_64)**: `terraform-provider-domotz_1.0.0_linux_amd64.tar.gz` (5.5MB)
  - Compatible with most modern Linux distributions
- **ARM64**: `terraform-provider-domotz_1.0.0_linux_arm64.tar.gz` (5.0MB)
  - Compatible with ARM64 servers and Raspberry Pi 4/5
- **ARMv7**: `terraform-provider-domotz_1.0.0_linux_arm.tar.gz` (5.2MB)
  - Compatible with Raspberry Pi 2/3 and other ARMv7 devices

### Windows
- **64-bit**: `terraform-provider-domotz_1.0.0_windows_amd64.zip` (5.7MB)
  - Compatible with modern Windows systems (Windows 10/11, Server 2016+)

## Installation

### Manual Installation

1. Download the appropriate archive for your platform from the `dist/` directory
2. Extract the archive:
   ```bash
   # Linux/macOS
   tar -xzf terraform-provider-domotz_1.0.0_<platform>.tar.gz

   # Windows (PowerShell)
   Expand-Archive terraform-provider-domotz_1.0.0_windows_amd64.zip
   ```

3. Move the binary to your Terraform plugins directory:
   ```bash
   # Linux/macOS
   mkdir -p ~/.terraform.d/plugins/registry.terraform.io/domotz/domotz/1.0.0/<platform>
   cp terraform-provider-domotz_v1.0.0 ~/.terraform.d/plugins/registry.terraform.io/domotz/domotz/1.0.0/<platform>/

   # Windows
   mkdir %APPDATA%\terraform.d\plugins\registry.terraform.io\domotz\domotz\1.0.0\windows_amd64
   copy terraform-provider-domotz_v1.0.0.exe %APPDATA%\terraform.d\plugins\registry.terraform.io\domotz\domotz\1.0.0\windows_amd64\
   ```

### Terraform Registry Installation (Future)

Once published to the Terraform Registry, users can install automatically:

```hcl
terraform {
  required_providers {
    domotz = {
      source  = "domotz/domotz"
      version = "~> 1.0"
    }
  }
}
```

## Verification

All binaries have been built with:
- Go compiler optimizations (`-trimpath`, `-s -w`)
- Static linking (`CGO_ENABLED=0`)
- Version information embedded (`-X main.version=1.0.0`)

### SHA256 Checksums

Verify downloaded archives using the provided checksums:

```bash
# Linux/macOS
shasum -a 256 -c terraform-provider-domotz_1.0.0_SHA256SUMS

# Windows (PowerShell)
Get-FileHash terraform-provider-domotz_1.0.0_windows_amd64.zip -Algorithm SHA256
```

See `terraform-provider-domotz_1.0.0_SHA256SUMS` for all checksums.

## Building from Source

To rebuild binaries for all platforms:

```bash
./build-all-platforms.sh
```

This will create fresh binaries in the `dist/` directory.

## GitHub Release

To create a GitHub release with these binaries:

1. Create a new release on GitHub:
   ```bash
   gh release create v1.0.0 \
     --title "Terraform Provider for Domotz v1.0.0" \
     --notes "Initial release - Production-ready provider for Domotz network monitoring"
   ```

2. Upload all distribution files:
   ```bash
   gh release upload v1.0.0 dist/*.tar.gz dist/*.zip dist/*SHA256SUMS
   ```

## Terraform Registry Publishing

To publish to the Terraform Registry (requires Domotz organization access):

1. Ensure the repository is at `github.com/domotz/terraform-provider-domotz`
2. Sign binaries with GPG key
3. Follow Terraform Registry publishing guide: https://www.terraform.io/registry/providers/publishing

## Platform Support Matrix

| Platform | Architecture | Tested | Notes |
|----------|-------------|--------|-------|
| macOS | Intel (x86_64) | ✅ | macOS 10.15+ |
| macOS | Apple Silicon (ARM64) | ✅ | M1/M2/M3/M4 |
| Linux | x86_64 | ✅ | Most distributions |
| Linux | ARM64 | ⚠️ | Raspberry Pi 4/5, ARM servers |
| Linux | ARMv7 | ⚠️ | Raspberry Pi 2/3 |
| Windows | x86_64 | ⚠️ | Windows 10/11, Server 2016+ |

Legend:
- ✅ Tested and verified
- ⚠️ Built but not tested (should work)

## Questions?

For issues or questions about the provider, please file an issue at:
https://github.com/domotz/terraform-provider-domotz/issues
