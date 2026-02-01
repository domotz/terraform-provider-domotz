terraform {
  required_providers {
    domotz = {
      source = "registry.terraform.io/domotz/domotz"
      version = "~> 0.1"
    }
  }
}

provider "domotz" {
  api_key = var.domotz_api_key
  # base_url = "https://api-eu-west-1-cell-1.domotz.com/public-api/v1" # Optional
}

variable "domotz_api_key" {
  type        = string
  description = "Domotz API key"
  sensitive   = true
}
