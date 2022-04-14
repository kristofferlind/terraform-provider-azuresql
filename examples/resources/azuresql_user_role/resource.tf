resource "azuresql_user" "this" {
  name     = "name"
  password = "password"
  database = "example_db"
}

resource "azuresql_user_role" "this" {
  # it's currently important that these are implicit dependencies
  # for resources to be created/deleted in the correct order
  name     = azuresql_user.this.name
  database = azuresql_user.this.database
  role     = "db_reader"
}
