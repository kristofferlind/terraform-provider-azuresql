package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLoginUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccLoginUserResourceConfig("testuser2", "test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresql_login_user.test", "name", "testuser2"),
					resource.TestCheckResourceAttr("azuresql_login_user.test", "database", "test"),
				),
			},
			// Update and Read testing
			{
				Config: testAccLoginUserResourceConfig("testuser2", "test2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresql_login_user.test", "database", "test2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccLoginUserResourceConfig(username string, database string) string {
	config := fmt.Sprintf(`
provider "azuresql" {
  connection_string = "sqlserver://sa:p@ssw0rd@localhost:1433"
}

resource "azuresql_login" "test" {
  name = %[1]q
  password = "P@ssw0rd"
  default_database = "master"
}

resource "azuresql_login_user" "test" {
  name = azuresql_login.test.name
  database = %[2]q
}
`, username, database)
	return config
}
