# Terraform Provider for Domotz

Manage your [Domotz](https://www.domotz.com/) network monitoring infrastructure as code with Terraform.

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)]()
[![Go Version](https://img.shields.io/badge/go-1.21+-blue)]()
[![Terraform](https://img.shields.io/badge/terraform-1.0+-purple)]()

## About Domotz

[Domotz](https://www.domotz.com/) is a network monitoring and management platform designed for MSPs, system integrators, and IT professionals. It provides:

- **Remote Network Monitoring** - Monitor devices, services, and network performance from anywhere
- **Automated Device Discovery** - Automatic detection of network devices with vendor/model identification
- **Custom Monitoring** - SNMP and TCP port monitoring for tailored alerting
- **Multi-Tenant Architecture** - Manage multiple sites and clients from a single platform
- **Integration APIs** - Comprehensive REST API for automation and integration

## What This Provider Enables

The Terraform Provider for Domotz allows you to:

✅ **Automate Device Tagging** - Automatically tag devices based on vendor, model, or any attribute  
✅ **Deploy Monitoring at Scale** - Create SNMP and TCP sensors across your infrastructure  
✅ **Manage Tags as Code** - Version control your device categorization and organization  
✅ **Query Device Inventory** - Integrate Domotz data into your Terraform workflows  
✅ **Standardize Monitoring** - Ensure consistent monitoring configuration across all sites  

### Use Cases

- **MSP Automation**: Automatically tag and monitor new customer devices as they're discovered
- **Infrastructure as Code**: Manage network monitoring alongside your infrastructure definitions
- **Compliance Workflows**: Ensure critical devices are tagged and monitored according to policy
- **Multi-Site Management**: Deploy consistent monitoring configurations across multiple locations
- **Integration Pipelines**: Feed Domotz device data into other Terraform-managed systems

---

## Table of Contents

- [Installation](#installation)
  - [Option 1: Use Pre-Built Binary](#option-1-use-pre-built-binary-recommended)
  - [Option 2: Build from Source](#option-2-build-from-source)
- [Authentication](#authentication)
- [Quick Start](#quick-start)
- [Provider Configuration](#provider-configuration)
- [Data Sources](#data-sources)
- [Resources](#resources)
- [Complete Example](#complete-example)
- [Building the Provider](#building-the-provider)
- [Contributing](#contributing)
- [License](#license)

---

## Installation

### Option 1: Use Pre-Built Binary (Recommended)

1. **Download the binary** from the releases page or use the included binary in this repository:

```bash
# Create the Terraform plugins directory
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/domotz/domotz/1.0.0/darwin_arm64

# Copy the binary (adjust path and architecture as needed)
cp terraform-provider-domotz ~/.terraform.d/plugins/registry.terraform.io/domotz/domotz/1.0.0/darwin_arm64/
```

**Platform-specific paths:**

- **macOS (Apple Silicon)**: `darwin_arm64`
- **macOS (Intel)**: `darwin_amd64`
- **Linux (64-bit)**: `linux_amd64`
- **Linux (ARM)**: `linux_arm64`
- **Windows (64-bit)**: `windows_amd64`

2. **Configure Terraform** to use the local provider:

```hcl
terraform {
  required_providers {
    domotz = {
      source  = "domotz/domotz"
      version = "1.0.0"
    }
  }
}
```

3. **For development**, create a `.terraformrc` in your home directory to override the provider source:

```hcl
provider_installation {
  dev_overrides {
    "domotz/domotz" = "/path/to/terraform-provider-domotz"
  }

  direct {}
}
```

### Option 2: Build from Source

See [Building the Provider](#building-the-provider) section below.

---

## Authentication

The provider requires a Domotz API key. You can obtain one from the [Domotz Portal](https://portal.domotz.com/).

**Method 1: Environment Variable (Recommended)**

```bash
export DOMOTZ_API_KEY="your-api-key-here"
```

**Method 2: Provider Configuration**

```hcl
provider "domotz" {
  api_key = "your-api-key-here"  # Not recommended for production
}
```

⚠️ **Security Note**: Never commit API keys to version control. Use environment variables or a secrets management solution.

---

## Quick Start

Here's a simple example to get you started:

```hcl
# Configure the provider
provider "domotz" {
  # API key from environment variable DOMOTZ_API_KEY
}

# Query your collector
data "domotz_agent" "my_agent" {
  id = 200891
}

# List all devices
data "domotz_devices" "all" {
  agent_id = data.domotz_agent.my_agent.id
}

# Create a custom tag
resource "domotz_custom_tag" "production" {
  name   = "Production"
  colour = "red"
}

# Tag critical devices
resource "domotz_device_tag_binding" "prod_tags" {
  for_each = {
    for device in data.domotz_devices.all.devices :
    device.id => device
    if device.importance == "VITAL"
  }

  agent_id  = data.domotz_agent.my_agent.id
  device_id = each.value.id
  tag_id    = domotz_custom_tag.production.id
}

# Output summary
output "summary" {
  value = {
    agent         = data.domotz_agent.my_agent.display_name
    total_devices = length(data.domotz_devices.all.devices)
    prod_devices  = length(domotz_device_tag_binding.prod_tags)
  }
}
```

Run:
```bash
terraform init
terraform plan
terraform apply
```

---

## Provider Configuration

```hcl
provider "domotz" {
  # API Key - Required
  # Can also be set via DOMOTZ_API_KEY environment variable
  api_key = "your-api-key"

  # API Base URL - Optional
  # - EU Region (default): `https://api-eu-west-1-cell-1.domotz.com/public-api/v1/`
  # - US Region: `https://api-us-east-1-cell-1.domotz.com/public-api/v1/`
  base_url = "https://api-eu-west-1-cell-1.domotz.com/public-api/v1"
}
```

### Configuration Options

| Argument | Required | Default | Description |
|----------|----------|---------|-------------|
| `api_key` | Yes | - | Domotz API key (can use `DOMOTZ_API_KEY` env var) |
| `base_url` | No | `https://api-eu-west-1-cell-1.domotz.com/public-api/v1` | Domotz API endpoint |

---

## Data Sources

Data sources allow you to query information from Domotz.

### domotz_agent

Retrieve details about a specific Domotz collector.

```hcl
data "domotz_agent" "primary" {
  id = 200891
}

output "agent_status" {
  value = {
    name   = data.domotz_agent.primary.display_name
    status = data.domotz_agent.primary.status
    team   = data.domotz_agent.primary.team_name
  }
}
```

**Attributes:**
- `id` (Required) - Collector ID
- `display_name` (Computed) - Collector display name
- `status` (Computed) - Collector status (ONLINE, OFFLINE)
- `team_id` (Computed) - Team/area ID
- `team_name` (Computed) - Team/area name

---

### domotz_devices

List all devices managed by a specific collector.

```hcl
data "domotz_devices" "all" {
  agent_id = 200891
}

# Filter devices by vendor
locals {
  ubiquiti_devices = {
    for device in data.domotz_devices.all.devices :
    device.id => device
    if device.vendor == "Ubiquiti Inc"
  }
}

output "device_summary" {
  value = {
    total    = length(data.domotz_devices.all.devices)
    ubiquiti = length(local.ubiquiti_devices)
  }
}
```

**Attributes:**
- `agent_id` (Required) - Collector ID
- `devices` (Computed) - List of devices with the following attributes:
  - `id` - Device ID
  - `display_name` - Device display name
  - `protocol` - Device protocol (IP, DUMMY, etc.)
  - `ip_addresses` - List of IP addresses
  - `importance` - Device importance level (VITAL, FLOATING)
  - `vendor` - Auto-discovered device vendor (e.g., "Ubiquiti Inc", "Apple")
  - `model` - Auto-discovered device model (e.g., "USL8LPB", "MacBook")
  - `user_data` - User-editable metadata object

**Filtering Examples:**

```hcl
# By vendor
locals {
  apple_devices = {
    for d in data.domotz_devices.all.devices :
    d.id => d if d.vendor == "Apple"
  }

  # By importance
  critical = {
    for d in data.domotz_devices.all.devices :
    d.id => d if d.importance == "VITAL"
  }

  # By protocol
  ip_devices = {
    for d in data.domotz_devices.all.devices :
    d.id => d if d.protocol == "IP"
  }

  # Multiple conditions
  critical_switches = {
    for d in data.domotz_devices.all.devices :
    d.id => d
    if d.importance == "VITAL" && d.vendor == "Ubiquiti Inc"
  }
}
```

---

### domotz_device

Retrieve details of a specific device.

```hcl
data "domotz_device" "core_switch" {
  agent_id = 200891
  id       = 12792047
}

output "switch_info" {
  value = {
    name         = data.domotz_device.core_switch.display_name
    ip_addresses = data.domotz_device.core_switch.ip_addresses
    importance   = data.domotz_device.core_switch.importance
  }
}
```

**Attributes:**
- `agent_id` (Required) - Collector ID
- `id` (Required) - Device ID
- `display_name` (Computed) - Device display name
- `protocol` (Computed) - Device protocol
- `ip_addresses` (Computed) - List of IP addresses
- `importance` (Computed) - Device importance level
- `user_data` (Computed) - Custom metadata object

---

### domotz_device_variables

Retrieve metrics/variables for a specific device.

```hcl
data "domotz_device_variables" "switch_metrics" {
  agent_id  = 200891
  device_id = 12792047
}

output "metrics_count" {
  value = length(data.domotz_device_variables.switch_metrics.variables)
}

output "sample_metrics" {
  value = [
    for v in data.domotz_device_variables.switch_metrics.variables :
    {
      label = v.label
      value = v.value
      unit  = v.unit
    }
  ]
}
```

**Attributes:**
- `agent_id` (Required) - Collector ID
- `device_id` (Required) - Device ID
- `variables` (Computed) - List of variables with:
  - `id` - Variable ID
  - `label` - Variable label
  - `path` - Variable path
  - `value` - Current value
  - `unit` - Unit of measurement
  - `previous_value` - Previous value
  - `metric` - Metric type

---

## Resources

Resources allow you to create and manage Domotz objects.

### domotz_custom_tag

Create and manage custom tags.

```hcl
resource "domotz_custom_tag" "production" {
  name   = "Production"
  colour = "red"
}

resource "domotz_custom_tag" "network_equipment" {
  name   = "Network Equipment"
  colour = "blue"
}
```

**Arguments:**
- `name` (Required) - Tag name
- `colour` (Required) - Tag color (e.g., "red", "blue", "#FF5733")

**Attributes:**
- `id` (Computed) - Tag ID

**Import:**
```bash
terraform import domotz_custom_tag.production 394382
```

---

### domotz_device_tag_binding

Associate custom tags with devices.

```hcl
# Tag a single device
resource "domotz_device_tag_binding" "switch_prod_tag" {
  agent_id  = 200891
  device_id = 12792047
  tag_id    = domotz_custom_tag.production.id
}

# Tag multiple devices using for_each
resource "domotz_device_tag_binding" "all_prod_tags" {
  for_each = local.critical_devices

  agent_id  = 200891
  device_id = each.value.id
  tag_id    = domotz_custom_tag.production.id
}

# Tag devices by vendor
locals {
  ubiquiti_devices = {
    for device in data.domotz_devices.all.devices :
    device.id => device
    if device.vendor == "Ubiquiti Inc"
  }
}

resource "domotz_device_tag_binding" "network_tags" {
  for_each = local.ubiquiti_devices

  agent_id  = 200891
  device_id = each.value.id
  tag_id    = domotz_custom_tag.network_equipment.id
}
```

**Arguments:**
- `agent_id` (Required) - Collector ID
- `device_id` (Required) - Device ID
- `tag_id` (Required) - Tag ID

**Attributes:**
- `id` (Computed) - Binding ID (format: `{agent_id}:{device_id}:{tag_id}`)

**Import:**
```bash
terraform import domotz_device_tag_binding.example 200891:12792047:394382
```

---

### domotz_snmp_sensor

Create SNMP monitoring sensors.

```hcl
# Monitor system description
resource "domotz_snmp_sensor" "system_description" {
  agent_id   = 200891
  device_id  = 12792047
  name       = "System Description"
  oid        = "1.3.6.1.2.1.1.1.0"  # sysDescr
  category   = "OTHER"
  value_type = "STRING"
}

# Monitor system uptime
resource "domotz_snmp_sensor" "system_uptime" {
  agent_id   = 200891
  device_id  = 12792047
  name       = "System Uptime"
  oid        = "1.3.6.1.2.1.1.3.0"  # sysUpTime
  category   = "OTHER"
  value_type = "NUMERIC"
}

# Monitor interface status
resource "domotz_snmp_sensor" "interface_status" {
  agent_id   = 200891
  device_id  = 12792047
  name       = "Port 1 Status"
  oid        = "1.3.6.1.2.1.2.2.1.8.1"  # ifOperStatus
  category   = "OTHER"
  value_type = "NUMERIC"
}
```

**Arguments:**
- `agent_id` (Required, Forces Replacement) - Collector ID
- `device_id` (Required, Forces Replacement) - Device ID
- `name` (Required, Forces Replacement) - Sensor name
- `oid` (Required, Forces Replacement) - SNMP OID to monitor
- `category` (Required, Forces Replacement) - Sensor category (e.g., "OTHER")
- `value_type` (Required, Forces Replacement) - Value type ("STRING" or "NUMERIC")

**Attributes:**
- `id` (Computed) - Sensor ID

**Common SNMP OIDs:**
- `1.3.6.1.2.1.1.1.0` - System Description
- `1.3.6.1.2.1.1.3.0` - System Uptime
- `1.3.6.1.2.1.1.5.0` - System Name
- `1.3.6.1.2.1.25.3.3.1.2` - CPU Usage
- `1.3.6.1.2.1.25.2.3.1.6` - Memory Usage

**Import:**
```bash
terraform import domotz_snmp_sensor.example 200891:12792047:72336
```

---

### domotz_tcp_sensor

Create TCP port monitoring sensors.

```hcl
# Monitor HTTPS port
resource "domotz_tcp_sensor" "web_https" {
  agent_id  = 200891
  device_id = 12792047
  name      = "HTTPS Port"
  port      = 443
  category  = "OTHER"
}

# Monitor SSH port
resource "domotz_tcp_sensor" "ssh" {
  agent_id  = 200891
  device_id = 12792047
  name      = "SSH Access"
  port      = 22
  category  = "OTHER"
}

# Monitor custom application port
resource "domotz_tcp_sensor" "app_api" {
  agent_id  = 200891
  device_id = 12792047
  name      = "Application API"
  port      = 8080
  category  = "OTHER"
}
```

**Arguments:**
- `agent_id` (Required, Forces Replacement) - Collector ID
- `device_id` (Required, Forces Replacement) - Device ID
- `name` (Required, Forces Replacement) - Sensor name
- `port` (Required, Forces Replacement) - TCP port number to monitor
- `category` (Required, Forces Replacement) - Sensor category (e.g., "OTHER")

**Attributes:**
- `id` (Computed) - Sensor ID

**Common Ports:**
- `22` - SSH
- `80` - HTTP
- `443` - HTTPS
- `3306` - MySQL
- `5432` - PostgreSQL
- `6379` - Redis
- `8080` - Alternative HTTP

⚠️ **Note**: Ensure the port is not already monitored on the device to avoid conflicts (409 error).

**Import:**
```bash
terraform import domotz_tcp_sensor.example 200891:12792047:72337
```

---

### domotz_device

Create external IP devices (external hosts).

```hcl
resource "domotz_device" "external_server" {
  agent_id     = 200891
  display_name = "External Web Server"
  ip_addresses = ["203.0.113.10"]
  importance   = "VITAL"

  user_data = {
    name   = "Production Web Server"
    model  = "Virtual Machine"
    vendor = "AWS"
    type   = "100"
  }
}
```

**Arguments:**
- `agent_id` (Required, Forces Replacement) - Collector ID
- `display_name` (Required) - Device display name
- `ip_addresses` (Required) - List of IP addresses
- `importance` (Optional) - Device importance level ("VITAL", "FLOATING")
- `user_data` (Optional) - Custom metadata object

**Attributes:**
- `id` (Computed) - Device ID

---

## Complete Example

Here's a comprehensive example demonstrating common patterns:

```hcl
terraform {
  required_providers {
    domotz = {
      source  = "domotz/domotz"
      version = "~> 1.0"
    }
  }
}

provider "domotz" {
  # API key from DOMOTZ_API_KEY environment variable
}

# ============================================
# DATA SOURCES
# ============================================

data "domotz_agent" "primary" {
  id = 200891
}

data "domotz_devices" "all" {
  agent_id = data.domotz_agent.primary.id
}

# ============================================
# FILTERING
# ============================================

locals {
  # Network equipment by vendor
  ubiquiti_devices = {
    for device in data.domotz_devices.all.devices :
    device.id => device
    if device.vendor == "Ubiquiti Inc"
  }

  # Critical devices
  critical_devices = {
    for device in data.domotz_devices.all.devices :
    device.id => device
    if device.importance == "VITAL"
  }
}

# ============================================
# TAGS
# ============================================

resource "domotz_custom_tag" "network" {
  name   = "Network Equipment"
  colour = "blue"
}

resource "domotz_custom_tag" "critical" {
  name   = "Critical Infrastructure"
  colour = "red"
}

# ============================================
# TAG BINDINGS
# ============================================

resource "domotz_device_tag_binding" "network_tags" {
  for_each = local.ubiquiti_devices

  agent_id  = data.domotz_agent.primary.id
  device_id = each.value.id
  tag_id    = domotz_custom_tag.network.id
}

resource "domotz_device_tag_binding" "critical_tags" {
  for_each = local.critical_devices

  agent_id  = data.domotz_agent.primary.id
  device_id = each.value.id
  tag_id    = domotz_custom_tag.critical.id
}

# ============================================
# MONITORING
# ============================================

# SNMP monitoring for first Ubiquiti device
resource "domotz_snmp_sensor" "uptime" {
  count = length(local.ubiquiti_devices) > 0 ? 1 : 0

  agent_id   = data.domotz_agent.primary.id
  device_id  = tonumber(keys(local.ubiquiti_devices)[0])
  name       = "System Uptime"
  oid        = "1.3.6.1.2.1.1.3.0"
  category   = "OTHER"
  value_type = "NUMERIC"
}

# ============================================
# OUTPUTS
# ============================================

output "summary" {
  value = {
    agent            = data.domotz_agent.primary.display_name
    total_devices    = length(data.domotz_devices.all.devices)
    network_devices  = length(local.ubiquiti_devices)
    critical_devices = length(local.critical_devices)
    tags_created     = [
      domotz_custom_tag.network.id,
      domotz_custom_tag.critical.id
    ]
  }
}
```

See [`examples/complete-example.tf`](./examples/complete-example.tf) for more examples.

---

## Building the Provider

### Prerequisites

- [Go](https://golang.org/doc/install) 1.21 or later
- [Terraform](https://www.terraform.io/downloads.html) 1.0 or later
- Make (optional, for convenience commands)

### Build Instructions

1. **Clone the repository:**

```bash
git clone https://github.com/domotz/terraform-provider-domotz.git
cd terraform-provider-domotz
```

2. **Build the provider:**

```bash
# Using Make (recommended)
make build

# Or using Go directly
go build -o terraform-provider-domotz
```

3. **Install locally for development:**

```bash
# Create plugin directory
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/domotz/domotz/1.0.0/$(go env GOOS)_$(go env GOARCH)

# Copy binary
cp terraform-provider-domotz ~/.terraform.d/plugins/registry.terraform.io/domotz/domotz/1.0.0/$(go env GOOS)_$(go env GOARCH)/
```

4. **Run tests:**

```bash
# Unit tests
make test

# Acceptance tests (requires DOMOTZ_API_KEY)
export DOMOTZ_API_KEY="your-api-key"
make testacc
```

### Development Setup

For active development, use a dev override in `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "domotz/domotz" = "/path/to/terraform-provider-domotz"
  }

  direct {}
}
```

This allows Terraform to use your local binary without installation.

### Makefile Commands

| Command | Description |
|---------|-------------|
| `make build` | Build the provider binary |
| `make test` | Run unit tests |
| `make testacc` | Run acceptance tests |
| `make fmt` | Format Go code |
| `make lint` | Run linters |
| `make clean` | Remove build artifacts |

---

## Troubleshooting

### Authentication Issues

**Error**: `401 Unauthorized`

**Solution**: Verify your API key is correct:
```bash
# Test API key
curl -H "X-Api-Key: $DOMOTZ_API_KEY" \
  https://api-eu-west-1-cell-1.domotz.com/public-api/v1/agent
```

### TCP Sensor Conflict (409)

**Error**: `failed to create TCP sensor: API error (status 409)`

**Solution**: The port is already monitored on this device. Choose a different port or remove the existing sensor first.

### Provider Not Found

**Error**: `provider registry.terraform.io/domotz/domotz was not found`

**Solution**: Ensure the provider is installed in the correct directory:
```bash
ls ~/.terraform.d/plugins/registry.terraform.io/domotz/domotz/1.0.0/
```

---

## API Documentation

- **Domotz API Portal**: https://portal.domotz.com/api/
- **OpenAPI Specification**: https://api-eu-west-1-cell-1.domotz.com/public-api/v1/meta/open-api-definition
- **Base URL**: 
  - EU Region: `https://api-eu-west-1-cell-1.domotz.com/public-api/v1/`
  - US Region: `https://api-us-east-1-cell-1.domotz.com/public-api/v1/`

---

## Contributing

Contributions are welcome! Please follow these guidelines:

1. **Fork the repository** and create a feature branch
2. **Write tests** for new functionality
3. **Follow Go conventions** and run `make fmt`
4. **Update documentation** including examples
5. **Submit a pull request** with a clear description

### Development Workflow

```bash
# 1. Make changes
vim internal/provider/resource_custom_tag.go

# 2. Format code
make fmt

# 3. Run tests
make test

# 4. Build and test locally
make build
cd test
terraform init
terraform plan
terraform apply

# 5. Commit and push
git commit -m "Add feature X"
git push origin feature-branch
```

---

## Support

- **Issues**: [GitHub Issues](https://github.com/domotz/terraform-provider-domotz/issues)
- **Domotz Support**: https://help.domotz.com/
- **Community**: [Domotz Community Forums](https://community.domotz.com/)

---

## License

This Terraform Provider for Domotz is released under the [Mozilla Public License 2.0](./LICENSE).

---

## Acknowledgments

- Built with the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework)
- Powered by the [Domotz Public API](https://portal.domotz.com/api/)
- Validated with comprehensive testing against production environments

---

## Release Notes

### v1.0.0 (2026-02-01)

**Initial Release**

✅ **Features:**
- Complete data source support for agents, devices, and metrics
- Tag management (create, update, delete)
- Device tagging with bulk operations support
- SNMP sensor creation and management
- TCP sensor creation and management
- Comprehensive filtering patterns
- for_each loop support for bulk operations

✅ **Validation:**
- Tested with 93 devices across multiple vendors
- Production-ready for MSP and enterprise use
