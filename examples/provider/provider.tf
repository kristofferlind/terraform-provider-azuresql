provider "azuresql" {
  # will use azuread authentication if fedauth parameter is included (fedauth=ActiveDirectoryDefault for example)
  # fedauth options: https://github.com/denisenkom/go-mssqldb#azure-active-directory-authentication
  connection_string = "sqlserver://sa:p@ssw0rd@localhost:1433"
}
