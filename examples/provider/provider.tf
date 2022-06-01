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
