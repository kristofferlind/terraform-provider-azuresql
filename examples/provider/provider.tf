# connecting to instance with username and password
provider "azuresql" {
  connection_string = "sqlserver://sa:p@ssw0rd@localhost:1433"
}

# connecting to azure instance (EnvironmentCredential -> ManagedIdentityCredential->AzureCLICredential)
provider "azuresql" {
  connection_string = "sqlserver://some-server.database.windows.net:1433?fedauth=ActiveDirectoryDefault"
}

# check https://github.com/microsoft/go-mssqldb#azure-active-directory-authentication for other fedauth options

# You can manage multiple servers by utilizing provider aliases
provider "azuresql" {
  alias             = "local"
  connection_string = "sqlserver://sa:p@ssw0rd@localhost:1433"
}

resource "azuresql_login" "this" {
  provider = azuresql.local
  name     = "name"
  password = "password"
}

# More information: https://www.terraform.io/language/providers/configuration#alias-multiple-provider-configurations
