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
  name       = "one"
  convention = designation_convention.this.convention
  inputs = {
    "region" = "ne"
  }
}

resource "designation_name" "web" {
  name       = "two"
  convention = designation_convention.this.convention
  inputs     = {}
}

resource "designation_name" "app" {
  name       = "three"
  convention = designation_convention.this.convention
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
