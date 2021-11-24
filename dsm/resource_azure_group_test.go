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
	resourceAzureGroup_createConfig = `resource "dsm_azure_group" "example_azure_group" {
		name = "example_azure_group"
  		description = "Azure Group Test"
  		tenant_id = "%s"
  		secret_key = "%s"
		subscription_id = "%s"
		client_id = "%s"
		url = "%s"
	}`
	resourceAzureGroup_updateConfig = `resource "dsm_azure_group" "example_azure_group" {
  		name = "example_aws_group_updated"
  		description = "AWS Group Test Update"
	}`
)

func TestAccResourceAzureGroup(t *testing.T) {
	var azure_tenant_id = os.Getenv("AZURE_TENANT_ID")
	var azure_secret_key = os.Getenv("AZURE_SECRET_KEY")
	var azure_subscription_id = os.Getenv("AZURE_SUBSCRIPTION_ID")
	var azure_client_id = os.Getenv("AZURE_CLIENT_ID")
	var azure_url = os.Getenv("AZURE_URL")

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheckAzure(t) },
		CheckDestroy: testAccCheckDestroyAzureGroup,
		Steps: []resource.TestStep{
			{
				Config:             fmt.Sprintf(resourceAzureGroup_createConfig, azure_tenant_id, azure_secret_key, azure_subscription_id, azure_client_id, azure_url),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckDestroyAzureGroup(s *terraform.State) (err error) {
	return err
}
