resource "domotz_device" "web_server" {
  agent_id     = 12345
  display_name = "Production Web Server"

  ip_addresses = ["203.0.113.42"]

  user_data {
    name   = "Web Server"
    model  = "Cloud Instance"
    vendor = "AWS"
    type   = "Server"
  }

  importance = "VITAL"
}
