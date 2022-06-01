provider "names" {
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

resource "names_name" "sql" {
  name = "one"
  inputs = {
    "region" = "ne"
  }
}

resource "names_name" "web" {
  name   = "two"
  inputs = {}
}

resource "names_name" "app" {
  name = "three"
  inputs = {
    "stage" = "test"
  }
}

output "name_web" {
  value = names_name.web.result
}

output "name_sql" {
  value = names_name.sql.result
}

output "name_app" {
  value = names_name.app.result
}
