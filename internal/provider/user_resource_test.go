package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserResourceConfig("testlocaluser", "P@ssw0rd", "contained_test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresql_user.test", "name", "testlocaluser"),
					resource.TestCheckResourceAttr("azuresql_user.test", "database", "contained_test"),
				),
			},
			// Update and Read testing
			{
				Config: testAccUserResourceConfig("testlocaluser2", "P@ssw0rd", "contained_test2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresql_user.test", "name", "testlocaluser2"),
					resource.TestCheckResourceAttr("azuresql_user.test", "database", "contained_test2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccUserResourceConfig(username string, password string, database string) string {
	config := fmt.Sprintf(`
provider "azuresql" {
  connection_string = "sqlserver://sa:p@ssw0rd@localhost:1433"
}

resource "azuresql_user" "test" {
  name = %[1]q
  password = %[2]q
  database = %[3]q
}
`, username, password, database)
	return config
}
