package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserRoleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserRoleResourceConfig("testuser", "P@ssw0rd", "contained_test", "db_datareader"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresql_user_role.test_role", "role", "db_datareader"),
				),
			},
			// Update and Read testing
			{
				Config: testAccUserRoleResourceConfig("testuser", "P@ssw0rd", "contained_test", "db_datawriter"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresql_user_role.test_role", "role", "db_datawriter"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccUserRoleResourceConfig(username string, password string, database string, role string) string {
	config := fmt.Sprintf(`
provider "azuresql" {
  connection_string = "sqlserver://sa:p@ssw0rd@localhost:1433"
}

resource "azuresql_user" "test" {
  name = %[1]q
  password = %[2]q
  database = %[3]q
}

resource "azuresql_user_role" "test_role" {
  name = azuresql_user.test.name
  database = azuresql_user.test.database
  role = %[4]q
}
`, username, password, database, role)
	return config
}
