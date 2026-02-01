---
page_title: "Domotz Provider"
subcategory: ""
description: |-
  Official Terraform provider for Domotz network monitoring platform.
---

# Domotz Provider

The Domotz provider allows you to manage your Domotz network monitoring infrastructure as code. It provides resources for managing devices, tags, sensors, and more.

## Features

- **Device Management**: Create and manage external IP devices
- **Custom Tags**: Create custom tags and associate them with devices
- **Sensor Configuration**: Configure SNMP and TCP port sensors
- **Data Querying**: Query agents, devices, and device metrics
- **Import Support**: Import existing resources into Terraform state

## Authentication

The provider requires a Domotz API key for authentication. You can obtain an API key from your Domotz account settings.

## Example Usage

```terraform
terraform {
  required_providers {
    domotz = {
      source  = "registry.terraform.io/domotz/domotz"
      version = "~> 0.1"
    }
  }
}

provider "domotz" {
  api_key = var.domotz_api_key
}

# Create an external device
resource "domotz_device" "web_server" {
  agent_id     = 12345
  display_name = "Production Web Server"
  ip_addresses = ["203.0.113.42"]
  importance   = "VITAL"

  user_data {
    name   = "Web Server"
    vendor = "AWS"
    type   = "Server"
  }
}

# Create a custom tag
resource "domotz_custom_tag" "production" {
  name   = "production"
  colour = "#FF5733"
}

# Bind tag to device
resource "domotz_device_tag_binding" "web_prod" {
  agent_id  = 12345
  device_id = domotz_device.web_server.id
  tag_id    = domotz_custom_tag.production.id
}
```

## Schema

### Required

- `api_key` (String, Sensitive) - Domotz API key for authentication. Can also be set via `DOMOTZ_API_KEY` environment variable.

### Optional

- `base_url` (String) - Base URL for the Domotz API. Defaults to `https://api-eu-west-1-cell-1.domotz.com/public-api/v1`. Can also be set via `DOMOTZ_BASE_URL` environment variable.

## Resources

- [domotz_device](resources/device.md) - Manage external IP devices
- [domotz_custom_tag](resources/custom_tag.md) - Manage custom tags
- [domotz_device_tag_binding](resources/device_tag_binding.md) - Bind tags to devices
- [domotz_snmp_sensor](resources/snmp_sensor.md) - Configure SNMP sensors
- [domotz_tcp_sensor](resources/tcp_sensor.md) - Configure TCP port sensors

## Data Sources

- [domotz_agent](data-sources/agent.md) - Query agent details
- [domotz_device](data-sources/device.md) - Query device details
- [domotz_devices](data-sources/devices.md) - List all devices for an agent
- [domotz_device_variables](data-sources/device_variables.md) - Query device metrics/variables

## Import

All resources support import. See individual resource documentation for import syntax.

## Support

For issues and feature requests, please use the [GitHub issue tracker](https://github.com/domotz/terraform-provider-domotz/issues).
