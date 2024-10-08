---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "generate_name function - terraform-provider-designation"
subcategory: ""
description: |-
  Generate a name for a convention and with specified parameters
---

# function: generate_name

Generate a name for a convention and with specified parameters

## Example Usage

```terraform
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

output "app_name" {
  value = provider::designation::generate_name(data.designation_convention.this.convention, {
    "name"   = "one"
    "region" = "ne"
    "random" = random_string.random.result
  })
}
```

## Signature

<!-- signature generated by tfplugindocs -->
```text
generate_name(convention string, inputs map of string) string
```

## Arguments

<!-- arguments generated by tfplugindocs -->
1. `convention` (String) The validated convention formated as a json string
1. `inputs` (Map of String) Map of input values for variables in provider defined convention

