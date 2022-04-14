resource "azuresql_login" "this" {
  name     = "name"
  password = "password"
}

resource "azuresql_login_user" "this" {
  name     = azuresql_login.this.name
  database = "example_db"
}
