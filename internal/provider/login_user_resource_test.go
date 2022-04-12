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
				Config: testAccLoginUserResourceConfig("testuser", "test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mssql_login_user.test", "name", "testuser"),
					resource.TestCheckResourceAttr("mssql_login_user.test", "database", "test"),
				),
			},
			// Update and Read testing
			{
				Config: testAccLoginUserResourceConfig("testuser", "test2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mssql_login_user.test", "database", "test2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccLoginUserResourceConfig(username string, database string) string {
	config := fmt.Sprintf(`
provider "mssql" {
  connection_string = "sqlserver://sa:p@ssw0rd@localhost:1433"
}

resource "mssql_login" "test" {
  name = %[1]q
  password = "P@ssw0rd"
}

resource "mssql_login_user" "test" {
  name = mssql_login.test.name
  database = %[2]q
}
`, username, database)
	return config
}
