data "designation_convention" "this" {
  definition = "(region)-(stage)-(name)-(random)"
  variables = [
    {
      name = "name"
    },
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
      max_length = 4
    },
  ]
}

resource "random_string" "random" {
  length           = 16
  special          = true
  override_special = "/@£$"
}

output "name1" {
  value = provider::designation::generate_name(data.designation_convention.this.convention, {
    "name"   = "one"
    "region" = "ne"
    "random" = random_string.random.result
  })
}
