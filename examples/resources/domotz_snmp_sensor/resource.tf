resource "domotz_snmp_sensor" "cpu_usage" {
  agent_id  = 12345
  device_id = domotz_device.web_server.id

  name     = "CPU Usage"
  oid      = "1.3.6.1.4.1.2021.11.9.0"  # UCD-SNMP-MIB::ssCpuIdle
  category = "OTHER"
}

resource "domotz_snmp_sensor" "memory_usage" {
  agent_id  = 12345
  device_id = domotz_device.web_server.id

  name     = "Memory Usage"
  oid      = "1.3.6.1.4.1.2021.4.6.0"  # UCD-SNMP-MIB::memTotalReal
  category = "OTHER"
}
