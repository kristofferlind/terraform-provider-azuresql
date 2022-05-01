resource "azuresql_login" "this" {
  name     = "name"
  password = "password"
}

resource "azuresql_login_role" "this" {
  name = azuresql_login.this.name
  role = "securityadmin"
}
