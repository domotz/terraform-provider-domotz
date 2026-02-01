data "domotz_agent" "main" {
  id = 12345
}

output "agent_info" {
  value = {
    name   = data.domotz_agent.main.display_name
    status = data.domotz_agent.main.status
    team   = data.domotz_agent.main.team_name
  }
}
