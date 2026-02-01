# Complete Domotz Provider Example
# This example demonstrates common patterns validated in production

terraform {
  required_providers {
    domotz = {
      source  = "domotz/domotz"
      version = "~> 1.0"
    }
  }
}

# Configure the Domotz Provider
provider "domotz" {
  # API key can be set via DOMOTZ_API_KEY environment variable
  # or specified directly (not recommended for production)
  # api_key = "your-api-key-here"

  # Optional: Override API endpoint (defaults to eu-west-1-cell-1)
  # base_url = "https://api-eu-west-1-cell-1.domotz.com/public-api/v1"
}

# ============================================
# DATA SOURCES
# ============================================

# Query collector information
data "domotz_agent" "primary" {
  id = 200891
}

# List all devices managed by the collector
data "domotz_devices" "all" {
  agent_id = data.domotz_agent.primary.id
}

# Get details for a specific device
data "domotz_device" "critical_switch" {
  agent_id = data.domotz_agent.primary.id
  id       = 12792047
}

# Get device metrics/variables
data "domotz_device_variables" "switch_metrics" {
  agent_id  = data.domotz_agent.primary.id
  device_id = 12792047
}

# ============================================
# FILTERING PATTERNS
# ============================================

# Filter devices by vendor (auto-discovered)
locals {
  # Network equipment by vendor
  ubiquiti_devices = {
    for device in data.domotz_devices.all.devices :
    device.id => device
    if device.vendor == "Ubiquiti Inc"
  }

  apple_devices = {
    for device in data.domotz_devices.all.devices :
    device.id => device
    if device.vendor == "Apple"
  }

  # Filter by importance level
  critical_devices = {
    for device in data.domotz_devices.all.devices :
    device.id => device
    if device.importance == "VITAL"
  }

  # Filter by protocol
  ip_devices = {
    for device in data.domotz_devices.all.devices :
    device.id => device
    if device.protocol == "IP"
  }
}

# ============================================
# TAG MANAGEMENT
# ============================================

# Create custom tags
resource "domotz_custom_tag" "network_equipment" {
  name   = "Network Equipment"
  colour = "blue"
}

resource "domotz_custom_tag" "critical_infrastructure" {
  name   = "Critical Infrastructure"
  colour = "red"
}

resource "domotz_custom_tag" "apple_endpoints" {
  name   = "Apple Endpoints"
  colour = "gray"
}

# Apply tags to filtered devices
resource "domotz_device_tag_binding" "ubiquiti_network_tags" {
  for_each = local.ubiquiti_devices

  agent_id  = data.domotz_agent.primary.id
  device_id = each.value.id
  tag_id    = domotz_custom_tag.network_equipment.id
}

resource "domotz_device_tag_binding" "critical_tags" {
  for_each = local.critical_devices

  agent_id  = data.domotz_agent.primary.id
  device_id = each.value.id
  tag_id    = domotz_custom_tag.critical_infrastructure.id
}

resource "domotz_device_tag_binding" "apple_tags" {
  for_each = local.apple_devices

  agent_id  = data.domotz_agent.primary.id
  device_id = each.value.id
  tag_id    = domotz_custom_tag.apple_endpoints.id
}

# ============================================
# MONITORING - SNMP SENSORS
# ============================================

# Add SNMP sensor to monitor system description
resource "domotz_snmp_sensor" "switch_sysDescr" {
  agent_id   = data.domotz_agent.primary.id
  device_id  = 12792047
  name       = "System Description"
  oid        = "1.3.6.1.2.1.1.1.0"  # sysDescr
  category   = "OTHER"
  value_type = "STRING"
}

# Monitor system uptime
resource "domotz_snmp_sensor" "switch_uptime" {
  agent_id   = data.domotz_agent.primary.id
  device_id  = 12792047
  name       = "System Uptime"
  oid        = "1.3.6.1.2.1.1.3.0"  # sysUpTime
  category   = "OTHER"
  value_type = "NUMERIC"
}

# ============================================
# MONITORING - TCP SENSORS
# ============================================

# Monitor web server ports (ensure port is not already monitored)
# resource "domotz_tcp_sensor" "web_https" {
#   agent_id  = data.domotz_agent.primary.id
#   device_id = 12792047
#   name      = "HTTPS Port"
#   port      = 443
#   category  = "OTHER"
# }

# ============================================
# OUTPUTS
# ============================================

output "agent_summary" {
  description = "Collector information summary"
  value = {
    id            = data.domotz_agent.primary.id
    name          = data.domotz_agent.primary.display_name
    status        = data.domotz_agent.primary.status
    total_devices = length(data.domotz_devices.all.devices)
  }
}

output "device_categorization" {
  description = "Device counts by category"
  value = {
    ubiquiti_count = length(local.ubiquiti_devices)
    apple_count    = length(local.apple_devices)
    critical_count = length(local.critical_devices)
    ip_devices     = length(local.ip_devices)
  }
}

output "critical_switch" {
  description = "Details of the critical switch"
  value = {
    id           = data.domotz_device.critical_switch.id
    name         = data.domotz_device.critical_switch.display_name
    ip           = data.domotz_device.critical_switch.ip_addresses
    metrics      = length(data.domotz_device_variables.switch_metrics.variables)
  }
}

output "tags_created" {
  description = "Custom tags created"
  value = {
    network_equipment       = domotz_custom_tag.network_equipment.id
    critical_infrastructure = domotz_custom_tag.critical_infrastructure.id
    apple_endpoints         = domotz_custom_tag.apple_endpoints.id
  }
}

output "monitoring_sensors" {
  description = "Monitoring sensors created"
  value = {
    snmp = {
      system_description = domotz_snmp_sensor.switch_sysDescr.id
      system_uptime      = domotz_snmp_sensor.switch_uptime.id
    }
  }
}
