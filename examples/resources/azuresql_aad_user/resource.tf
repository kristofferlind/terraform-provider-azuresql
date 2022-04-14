resource "azuresql_aad_user" "this" {
  name     = "<display name of aad user/group/application>"
  database = "example_db"
}
