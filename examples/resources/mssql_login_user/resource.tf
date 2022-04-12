resource "mssql_login" "this" {
  name     = "name"
  password = "password"
}

resource "mssql_login_user" "this" {
  name     = mssql_login.this.name
  database = "example_db"
}
