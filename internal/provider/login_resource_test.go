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
				Config: testAccLoginResourceConfig("testuser", "p@ssword1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresql_login.test", "name", "testuser"),
					resource.TestCheckResourceAttr("azuresql_login.test", "password", "p@ssword1"),
				),
			},
			// Update and Read testing
			{
				Config: testAccLoginResourceConfig("testuser", "p@ssword2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresql_login.test", "password", "p@ssword2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccLoginResourceConfig(username string, password string) string {
	config := fmt.Sprintf(`
provider "azuresql" {
  connection_string = "sqlserver://sa:p@ssw0rd@localhost:1433"
}

resource "azuresql_login" "test" {
  name = %[1]q
  password = %[2]q
}
`, username, password)
	return config
}
