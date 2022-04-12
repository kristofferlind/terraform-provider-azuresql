resource "mssql_user" "this" {
  name     = "name"
  password = "password"
  database = "example_db"
}
