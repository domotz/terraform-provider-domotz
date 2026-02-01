data "domotz_device" "existing" {
  agent_id = 12345
  id       = 67890
}

output "device_info" {
  value = {
    name        = data.domotz_device.existing.display_name
    protocol    = data.domotz_device.existing.protocol
    importance  = data.domotz_device.existing.importance
    ip_addresses = data.domotz_device.existing.ip_addresses
  }
}
