// Code generated by apic. DO NOT EDIT.

package dsm

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"
)

var (
	resourceAzureSobject_createConfig = `resource "dsm_group" "example_group" {
		name = "example_group"
	}

	resource "dsm_sobject" "example_sobject" {		
		name     = "example_sobject"
		group_id = "${dsm_group.example_group.group_id}"
		key_size = 256
		key_ops = [
			"ENCRYPT",
			"DECRYPT",
	  	        "WRAPKEY",
	 	        "UNWRAPKEY",
			"DERIVEKEY",
			"MACGENERATE",
			"MACVERIFY",
			"APPMANAGEABLE",
			"EXPORT"
		]
		obj_type = "AES"
	}
	
	
	resource "dsm_azure_group" "example_azure_group" {
		name = "example_azure_group"
  		description = "Azure Group Test"
  		tenant_id = "%s"
  		secret_key = "%s"
		subscription_id = "%s"
		client_id = "%s"
		url = "%s"
  	}
	resource "dsm_azure_sobject" "example_azure_soject" {
		name     = "example_azure_soject"
		group_id = "${dsm_azure_group.example_azure_group.group_id}"
		key = {
		  kid = "${dsm_sobject.example_sobject.kid}"
		}
	}`
)

func TestAccResourceAzureSobject(t *testing.T) {
	var azure_tenant_id = os.Getenv("AZURE_TENANT_ID")
	var azure_secret_key = os.Getenv("AZURE_SECRET_KEY")
	var azure_subscription_id = os.Getenv("AZURE_SUBSCRIPTION_ID")
	var azure_client_id = os.Getenv("AZURE_CLIENT_ID")
	var azure_url = os.Getenv("AZURE_URL")

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheckAzure(t) },
		CheckDestroy: testAccCheckDestroyAzureSobject,
		Steps: []resource.TestStep{
			{
				Config:             fmt.Sprintf(resourceAzureSobject_createConfig, azure_tenant_id, azure_secret_key, azure_subscription_id, azure_client_id, azure_url),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckDestroyAzureSobject(s *terraform.State) (err error) {
	return err
}