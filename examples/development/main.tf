terraform {
  required_version = ">= 0.14"
  required_providers {
    designation = {
      source = "datadrivers/designation"
    }
  }
}
provider "designation" {}

resource "designation_convention" "this" {
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

resource "designation_name" "sql" {
  convention = designation_convention.this.convention

  name = "one"
  inputs = {
    "region" = "ne"
  }
}

resource "designation_name" "web" {
  convention = designation_convention.this.convention

  name   = "two"
  inputs = {}
}

resource "designation_name" "app" {
  convention = designation_convention.this.convention

  name = "three"
  inputs = {
    "stage" = "test"
  }
}

output "name_web" {
  value = designation_name.web.result
}

output "name_sql" {
  value = designation_name.sql.result
}

output "name_app" {
  value = designation_name.app.result
}
