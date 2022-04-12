resource "mssql_user" "this" {
  name     = "name"
  password = "password"
  database = "example_db"
}

resource "mssql_user_role" "this" {
  name     = mssql_user.this.name
  database = mssql_user.this.database
  role     = "db_reader"
}
