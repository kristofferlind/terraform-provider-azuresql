resource "azuresql_user" "this" {
  name     = "name"
  password = "password"
  database = "example_db"
}

resource "azuresql_user_role" "this" {
  name     = azuresql_user.this.name
  database = azuresql_user.this.database
  role     = "db_reader"
}
