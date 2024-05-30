data "designation_convention" "this" {
  definition = "(region)-(stage)-(name)-(type)"
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
      name    = "type"
      default = "service"
    },
  ]
}
