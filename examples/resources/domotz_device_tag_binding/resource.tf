resource "domotz_device_tag_binding" "web_prod" {
  agent_id  = 12345
  device_id = domotz_device.web_server.id
  tag_id    = domotz_custom_tag.production.id
}
