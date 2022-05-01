package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLoginRoleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccLoginRoleResourceConfig("testuser", "P@ssw0rd", "securityadmin"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresql_login_role.test", "role", "securityadmin"),
				),
			},
			// Update and Read testing
			{
				Config: testAccLoginRoleResourceConfig("testuser", "P@ssw0rd", "serveradmin"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresql_login_role.test", "role", "serveradmin"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccLoginRoleResourceConfig(username string, password string, role string) string {
	config := fmt.Sprintf(`
provider "azuresql" {
  connection_string = "sqlserver://sa:p@ssw0rd@localhost:1433"
}

resource "azuresql_login" "test" {
  name = %[1]q
  password = %[2]q
}

resource "azuresql_login_role" "test" {
  name = azuresql_login.test.name
  role = %[3]q
}
`, username, password, role)
	return config
}
