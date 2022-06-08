resource "convention_convention" "this" {
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
  name       = "one"
  convention = convention_convention.this.convention
  inputs = {
    "region" = "ne"
  }
}

resource "convention_name" "web" {
  name       = "two"
  convention = convention_convention.this.convention
  inputs     = {}
}

resource "convention_name" "app" {
  name       = "three"
  convention = convention_convention.this.convention
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
