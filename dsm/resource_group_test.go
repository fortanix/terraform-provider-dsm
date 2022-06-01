// Code generated by apic. DO NOT EDIT.

package dsm

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

var (
	resourceGroup_createConfig = `resource "dsm_group" "example_group" {
  		name = "example_group"
	}`
	resourceGroup_updateConfig = `resource "dsm_group" "example_group" {
  		name = "example_group_updated"
	}`
)

func TestAccResourceGroup(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckDestroyGroup,
		Steps: []resource.TestStep{
			{
				Config: resourceGroup_createConfig,
			},
		},
	})
}

func testAccCheckDestroyGroup(s *terraform.State) (err error) {
	return err
}
