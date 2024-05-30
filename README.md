# Terraform Provider Designation (Terraform Plugin Framework)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.17

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

_Note:_ Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

## Local Test

There is a [makefile](./GNUmakefile) to build the provider and place it in repos root dir.

```sh
make
```

To use the local build version you need tell terraform where to look for it via a terraform config override.

Create `dev.tfrc` in your terraform code folder (f.e. in [dev.tfrc](./examples/development/dev.tfrc)):

```hcl
# dev.tfrc
provider_installation {

  # Use /home/developer/tmp/terraform-nexus as an overridden package directory
  # for the datadrivers/nexus provider. This disables the version and checksum
  # verifications for this provider and forces Terraform to look for the
  # nexus provider plugin in the given directory.
  # relative path also works, but no variable or ~ evaluation
  dev_overrides {
    "datadrivers/nexus" = "../../"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

Tell your shell environment to use override file:

```bash
export TF_CLI_CONFIG_FILE=dev.tfrc
```

Now run your terraform commands (`plan` or `apply`), `init` is ***not*** required.

```bash
# start local nexus
make start-services
# run local terraform code
cd examples/local-development
terraform plan
terraform apply
```

## Create Release

```shell script
cz bump --changelog
git push --tags
```
