data "domotz_devices" "all" {
  agent_id = 12345
}

output "device_count" {
  value = length(data.domotz_devices.all.devices)
}

output "device_names" {
  value = [for device in data.domotz_devices.all.devices : device.display_name]
}
