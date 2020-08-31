package wavefront

import (
	"fmt"
	"testing"

	"github.com/WavefrontHQ/go-wavefront-management-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccWavefrontRole_BasicRole(t *testing.T) {
	var record wavefront.Role
	resourceName := "wavefront_role.role"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontRoleBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontRoleExists(resourceName, &record),
					testAccCheckWavefrontRoleAttributes(&record),
					resource.TestCheckResourceAttr(resourceName, "name", "Test Role"),
				),
			},
		},
	})
}

func TestAccWavefrontRole_AdvancedRole(t *testing.T) {
	var record wavefront.Role
	resourceName := "wavefront_role.role"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontRoleAdvanced(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontRoleExists(resourceName, &record),
					testAccCheckWavefrontRoleAttributes(&record),
					resource.TestCheckResourceAttr(resourceName, "name", "Test Role"),
					resource.TestCheckResourceAttr(resourceName, "description", "Test Role Description"),
					resource.TestCheckResourceAttr(resourceName, "permissions.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "assignees.#", "1"),
				),
			},
			{
				Config: testAccCheckWavefrontRoleAdvancedChanged(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontRoleExists(resourceName, &record),
					testAccCheckWavefrontRoleAttributes(&record),
					resource.TestCheckResourceAttr(resourceName, "name", "Test Role"),
					resource.TestCheckResourceAttr(resourceName, "description", "Test Role Description"),
					resource.TestCheckResourceAttr(resourceName, "permissions.#", "2"),
					resource.TestCheckNoResourceAttr(resourceName, "assignees"),
				),
			},
		},
	})
}

func testAccCheckWavefrontRoleBasic() string {
	return `
resource "wavefront_role" "role" {
  name = "Test Role"
}
`
}

func testAccCheckWavefrontRoleAdvanced() string {
	return `
resource "wavefront_user_group" "user_group" {
  name        = "User Group"
  description = "User Group Description"
}


resource "wavefront_role" "role" {
  name        = "Test Role"
  description = "Test Role Description"
  permissions = [
    "derived_metrics_management",
    "agent_management",
    "alerts_management",
  ]
  assignees = [wavefront_user_group.user_group.id, ]
}
`
}

func testAccCheckWavefrontRoleAdvancedChanged() string {
	return `
resource "wavefront_user_group" "user_group" {
  name        = "User Group"
  description = "User Group Description"
}


resource "wavefront_role" "role" {
  name        = "Test Role"
  description = "Test Role Description"
  permissions = [
    "derived_metrics_management",
    "alerts_management",
  ]
}
`
}

func testAccCheckWavefrontRoleExists(resourceName string, role *wavefront.Role) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		roles := testAccProvider.Meta().(*wavefrontClient).client.Roles()
		if r, ok := state.RootModule().Resources[resourceName]; ok {
			tmp := &wavefront.Role{ID: r.Primary.ID}
			err := roles.Get(tmp)
			if err != nil {
				return err
			}
			*role = *tmp
			return nil
		}
		return fmt.Errorf("not found, %s", resourceName)
	}
}

func testAccCheckWavefrontRoleAttributes(role *wavefront.Role) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		if role.Name != "Test Role" {
			return fmt.Errorf("expected Test Role, got %s", role.Name)
		}

		if role.Description != "" && role.Description != "Test Role Description" {
			return fmt.Errorf("expected Test Role Description, got %s", role.Description)
		}

		return nil
	}
}

func testAccCheckWavefrontRoleDestroy(s *terraform.State) error {

	roles := testAccProvider.Meta().(*wavefrontClient).client.Roles()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "wavefront_role" {
			continue
		}

		results, err := roles.Find(
			[]*wavefront.SearchCondition{
				{
					Key:            "id",
					Value:          rs.Primary.ID,
					MatchingMethod: "EXACT",
				},
			})
		if err != nil {
			return fmt.Errorf("error finding Wavefront Role. %s", err)
		}
		if len(results) > 0 {
			return fmt.Errorf("role still exists")
		}
	}

	return nil
}
