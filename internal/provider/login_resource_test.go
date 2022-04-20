package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLoginResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccLoginResourceConfig("testuser", "p@ssword1", "dash-test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresql_login.test", "name", "testuser"),
					resource.TestCheckResourceAttr("azuresql_login.test", "password", "p@ssword1"),
				),
			},
			// Change default database
			{
				Config: testAccLoginResourceConfig("testuser", "p@ssword1", "contained_test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresql_login.test", "default_database", "contained_test"),
				),
			},
			// Change password
			{
				Config: testAccLoginResourceConfig("testuser", "p@ssword2", "contained_test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresql_login.test", "password", "p@ssword2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccLoginResourceConfig(username string, password string, defaultDatabase string) string {
	config := fmt.Sprintf(`
provider "azuresql" {
  connection_string = "sqlserver://sa:p@ssw0rd@localhost:1433"
}

resource "azuresql_login" "test" {
  name = %[1]q
  password = %[2]q
  default_database = %[3]q
}
`, username, password, defaultDatabase)
	return config
}
