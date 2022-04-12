# terraform-provider-mssql
Provider for managing MSSQL Server users and permissions. Documentation can be found in /docs or the registry.

Using an unofficial terraform provider is pretty scary and it therefore needs to be really easy to review and fork. That leaves us with the following focus points:
- keep it as simple as possible
- keep dependencies to a minimum

## Maturity
This provider has currently only been tested against a local docker based mssql instance. It's the first time I even modify code in a terraform provider and I'm still a bit of a newbie when it comes to Golang. **Do not use this unless you can first test all your changes in an environment that is ok to break**.

## Dependencies
Other than terraform boilerplate/plumbing this provider only utilizes an mssql driver, specifically github.com/denisenkom/go-mssqldb, which is also mentioned in Microsoft's documentation for working with mssql using golang.

## Structure
main.go and everything that resides in internal/provider directory is terraform boilerplate and dependencies, internal/manager contains the actual implementation and it's dependencies.

## Development
This plugin is based on Terraform Plugin Framework and the documentation for that should therefore be a good source of information (https://www.terraform.io/plugin/framework).

`nix-shell --pure` open shell with all prerequisites installed
`make test` build, install and execute tests

Check makefile for more options.
