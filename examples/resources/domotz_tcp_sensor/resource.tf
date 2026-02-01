resource "domotz_tcp_sensor" "https_check" {
  agent_id  = 12345
  device_id = domotz_device.web_server.id

  name     = "HTTPS Service"
  port     = 443
  category = "OTHER"
}

resource "domotz_tcp_sensor" "ssh_check" {
  agent_id  = 12345
  device_id = domotz_device.web_server.id

  name     = "SSH Service"
  port     = 22
  category = "OTHER"
}
