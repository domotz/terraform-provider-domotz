data "domotz_device_variables" "metrics" {
  agent_id  = 12345
  device_id = 67890
}

output "device_variables" {
  value = {
    for var in data.domotz_device_variables.metrics.variables :
    var.label => {
      value  = var.value
      unit   = var.unit
      metric = var.metric
    }
  }
}

# Filter for specific metrics
output "cpu_metrics" {
  value = [
    for var in data.domotz_device_variables.metrics.variables :
    var if can(regex("(?i)cpu", var.label))
  ]
}
