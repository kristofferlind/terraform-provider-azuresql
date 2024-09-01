# terraform-provider-azuresql
Provider for managing MSSQL Server users and permissions. Documentation can be found in /docs or the registry.

Using an unofficial terraform provider is pretty scary and it therefore needs to be really easy to review and fork. That leaves us with the following focus points:
- keep it as simple as possible
- keep dependencies to a minimum

## Maturity
This provider currently only has acceptance tests against a local docker based mssql instance and has received a little bit of manual testing against Azure SQL Server and Azure AD. It's my first time creating a terraform provider and I'm still a bit of a Golang newbie. **Do not use this unless you can first test all your changes in an environment that is ok to break**.

## Dependencies
Other than terraform boilerplate/plumbing this provider only utilizes an mssql driver, specifically github.com/microsoft/go-mssqldb, the official driver for MSSQL with Go from Microsoft.

## Development
This plugin is based on Terraform Plugin Framework and the documentation for that should therefore be a good source of information (https://www.terraform.io/plugin/framework).

If you have [direnv](https://direnv.net/) and [nix package manager](https://nixos.org/) installed. You should be able to just start using commands found in Makefile. If you don't want those you can look at shell.nix to see what's needed.

command | description
---|---
`make test` | build, install and execute tests

Check makefile for more options.

## Known issues
- Need to specify resources with implicit dependencies for them to be created and deleted in the correct order
- aad_user is untested (feature is currently public preview)

## License
Based on [terraform-provider-scaffolding-framework](https://github.com/hashicorp/terraform-provider-scaffolding-framework) and copied its license at time of creation.
