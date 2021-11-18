// Code generated by apic. DO NOT EDIT.

package dsm

import (
	"fmt"
	"os"
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	resourceGcpEKM_createConfig = `resource "dsm_group" "example_group" {
		name = "example_group"
	}
  
	resource "dsm_gcp_ekm_sa" "example_gcp_ekm_sa" {
		name = "%s"
    		default_group = "${dsm_group.example_group.group_id}"
	}`
)

func TestAccResourceGcpEKM(t *testing.T) {
	var google_service_account = os.Getenv("GOOGLE_SERVICE_ACCOUNT")

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheckGcp(t) },
		CheckDestroy: testAccCheckDestroyGcpEKM,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(resourceAwsSobject_createConfig, google_service_account),
			},
		},
	})
}

func testAccCheckDestroyGcpEKM(s *terraform.State) (err error) {
	return err
}
