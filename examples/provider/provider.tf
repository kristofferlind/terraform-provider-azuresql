provider "azuresql" {
  # will use azuread authentication if fedauth parameter is included (fedauth=ActiveDirectoryDefault for example)
  # fedauth options: https://github.com/denisenkom/go-mssqldb#azure-active-directory-authentication
  # with ActiveDirectoryDefault setting credentials are tried in this order:
  # EnvironmentCredential -> ManagedIdentityCredential->AzureCLICredential
  # connection_string = "sqlserver://some-server.database.windows.net:1433?fedauth=ActiveDirectoryDefault"
  connection_string = "sqlserver://sa:p@ssw0rd@localhost:1433"
}
