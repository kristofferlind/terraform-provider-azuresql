resource "azuresql_login" "this" {
  name     = "name"
  password = "password"
}

resource "azuresql_login_user" "this" {
  # it's currently important that this is an implicit dependency
  # for resources to be created/deleted in the correct order
  name     = azuresql_login.this.name
  database = "example_db"
}
