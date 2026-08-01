package interfaces_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccInterfacesAssignmentResource exercises the full create/update/import/
// delete lifecycle by assigning a spare device to a new optional interface.
// Skipped unless OPNSENSE_TEST_ASSIGNMENT_SPARE_DEVICE is set: the CI VM is
// single-NIC with that NIC already carrying wan traffic, so there is no
// device it's safe to assign and unassign here.
func TestAccInterfacesAssignmentResource(t *testing.T) {
	device := acctest.AssignmentSpareDevicePreCheck(t)
	if device == "" {
		t.Skip("OPNSENSE_TEST_ASSIGNMENT_SPARE_DEVICE must be set to a spare, non-live device to run this test")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAssignmentResourceConfig(device, "Assignment test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_interfaces_assignment.test", "device", device),
					resource.TestCheckResourceAttr("opnsense_interfaces_assignment.test", "description", "Assignment test"),
					resource.TestCheckResourceAttr("opnsense_interfaces_assignment.test", "lock", "false"),
					resource.TestCheckResourceAttrSet("opnsense_interfaces_assignment.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_interfaces_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccAssignmentResourceConfig(device, "Updated assignment test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_interfaces_assignment.test", "device", device),
					resource.TestCheckResourceAttr("opnsense_interfaces_assignment.test", "description", "Updated assignment test"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccAssignmentResourceConfig(device, description string) string {
	return fmt.Sprintf(`
resource "opnsense_interfaces_assignment" "test" {
  device      = %[1]q
  description = %[2]q
}
`, device, description)
}
