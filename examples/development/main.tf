terraform {
  required_version = ">= 0.14"
  required_providers {
    convention = {
      source = "datadrivers/convention"
    }
  }
}

provider "convention" {
  definition = "(region)-(stage)-(name)-(random)"
  variables = [
    {
      name    = "region"
      default = "we"
    },
    {
      name       = "stage"
      max_length = 4
      default    = "dev"
    },
    {
      name       = "random"
      generated  = true
      max_length = 8
    },
  ]
}

resource "convention_name" "sql" {
  name = "one"
  inputs = {
    "region" = "ne"
  }
}

resource "convention_name" "web" {
  name   = "two"
  inputs = {}
}

resource "convention_name" "app" {
  name = "three"
  inputs = {
    "stage" = "test"
  }
}

output "name_web" {
  value = convention_name.web.result
}

output "name_sql" {
  value = convention_name.sql.result
}

output "name_app" {
  value = convention_name.app.result
}
