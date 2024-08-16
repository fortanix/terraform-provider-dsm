// **********
// Terraform Provider - DSM: data source: azure kms group
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.5.0
//       - Date:      05/01/2021
// **********

package dsm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAzureGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Returns the Fortanix DSM Azure KV mapped group object from the cluster as a Data Source.",
		ReadContext: dataSourceAzureGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The Azure KV group object name in Fortanix DSM.",
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Description: "The Azure KV group object ID from Fortanix DSM.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"acct_id": {
				Description: "The Account ID from Fortanix DSM.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator": {
				Description: "The creator of the group from Fortanix DSM.\n" +
			    "   * `user`: If the group was created by a user, the computed value will be the matching user id.\n" +
			    "   * `app`: If the group was created by a app, the computed value will be the matching app id.",
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": {
				Description: "Description of the Azure KV Fortanix DSM group.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Description: "The URL of the object in an Azure KV that uniquely identifies the object.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_id": {
				Description: "The Azure registered application id (username).",
				Type:     schema.TypeString,
				Computed: true,
			},
			"subscription_id": {
				Description: "The ID of the Azure AD subscription.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Description: "The tenant/directory id of the Azure subscription.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_key": {
				Description: "A secret string that a registered application in Azure uses to prove its identity (application password).",
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"key_vault_type": {
				Description: "The type of key vaults. The default value is `Standard`.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"scan": {
				Description: "Syncs keys from Azure KV to the Azure group in DSM. Value is either true/false.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAzureGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var group_data map[string]interface{}
	original_name := d.Get("name").(string)
	modified_name := fmt.Sprintf("%s-azure-%s", original_name, m.(*api_client).azure_region)

	req, err := m.(*api_client).APICallList("GET", "sys/v1/groups")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: GET sys/v1/groups: %v", err),
		})
		return diags
	}

	// Shashi: First, check for the group name as provided.
	// If not found, fallback to checking the group name in the "name-azure-region" format used prior to v0.5.33.
	for _, data := range req {
		group_name := data.(map[string]interface{})["name"].(string)
		if group_name == original_name {
			group_data = data.(map[string]interface{})
			break
		}
		if group_name == modified_name && group_data == nil {
			group_data = data.(map[string]interface{})
		}
	}

	if group_data == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Group not found.",
			Detail:   fmt.Sprintf("[E]: No group found with name: %s or %s", original_name, modified_name),
		})
		return diags
	}

	jsonbody, parsing_error := json.Marshal(group_data)
	if parsing_error != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to parse DSM provider API client output.",
			Detail:   fmt.Sprintf("[E]: API: GET sys/v1/groups: %s", parsing_error),
		})
		return diags
	}

	azuregroup := AzureGroup{}
	if err := json.Unmarshal(jsonbody, &azuregroup); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to parse DSM provider API client output",
			Detail:   fmt.Sprintf("[E]: API: GET sys/v1/groups: %s", err),
		})
		return diags
	}

	// FYOO: AzureGroup must conform to this JSON struct - if this crashes, then we have DSM issues
	name := strings.Split(azuregroup.Name, fmt.Sprintf("-azure-%s", m.(*api_client).azure_region))
	d.Set("name", name[0])
	d.Set("group_id", azuregroup.Group_id)
	d.Set("acct_id", azuregroup.Acct_id)

	var creatorInt map[string]interface{}
	creatorRec, _ := json.Marshal(azuregroup.Creator)
	json.Unmarshal(creatorRec, &creatorInt)
	d.Set("creator", creatorInt)
	d.Set("region", m.(*api_client).azure_region)

	// FYOO: there is only one HMG per AzureGroup
	for _, value := range azuregroup.Hmg {
		d.Set("subscription_id", value.Subscription_id)
		d.Set("client_id", value.Client_id)
		d.Set("tenant_id", value.Tenant_id)
		d.Set("key_vault_type", value.Key_vault_type)
		d.Set("url", value.Url)
	}
	// FYOO: remove sensitive information
	d.Set("secret_key", "")
	// FYOO: if description is blank, DSM does not return
	if _, ok := group_data["description"]; ok {
		d.Set("description", group_data["description"].(string))
	}

	d.SetId(d.Get("group_id").(string))

	// If Scan is set, then move to scanning for data source
	if d.Get("scan").(bool) {
		check_hmg_req := map[string]interface{}{}
		// Scan the AWS Group first before
		_, err := m.(*api_client).APICallBody("POST", fmt.Sprintf("sys/v1/groups/%s/hmg/check", d.Get("group_id").(string)), check_hmg_req)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: POST sys/v1/groups/-/hmg/check: %v", err),
			})
			return diags
		}

		_, err = m.(*api_client).APICallBody("POST", fmt.Sprintf("sys/v1/groups/%s/hmg/scan", d.Get("group_id").(string)), check_hmg_req)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: POST sys/v1/groups/-/hmg/scan: %v", err),
			})
			return diags
		}
	}
	return nil
}
