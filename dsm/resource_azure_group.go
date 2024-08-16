// **********
// Terraform Provider - DSM: resource: azure kms group
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.5.0
//       - Date:      27/11/2020
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

// [-] Define Group
func resourceAzureGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateAzureGroup,
		ReadContext:   resourceReadAzureGroup,
		UpdateContext: resourceUpdateAzureGroup,
		DeleteContext: resourceDeleteAzureGroup,
		Description: "Creates a Fortanix DSM group mapped to Azure Key Vault in the cluster as a resource. This group acts as a container for security objects. The returned resource object contains the UUID of the group for further references.\n",
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
				Optional: true,
				Default:  "",
			},
			"url": {
			    Description: "The URL of the object in an Azure KV that uniquely identifies the object.",
				Type:     schema.TypeString,
				Required: true,
			},
			"client_id": {
			    Description: "The Azure registered application id (username).",
				Type:     schema.TypeString,
				Required: true,
			},
			"subscription_id": {
			    Description: "The ID of the Azure AD subscription.",
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
			    Description: "The tenant/directory id of the Azure subscription.",
				Type:     schema.TypeString,
				Required: true,
			},
			"key_vault_type": {
			    Description: "The type of key vault. The default value is `Standard`. Values are Standard/Premium.",
				Type:     schema.TypeString,
				Optional: true,
				Default :"Standard",
			},
			"secret_key": {
			    Description: "A secret string that a registered application in Azure uses to prove its identity (application password).",
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// [C]: Create Azure Group
func resourceCreateAzureGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	group_object := map[string]interface{}{
		// 0.5.0: Azure KMS Group Name to be predefined as <string>-azure-<region>
		"name":           d.Get("name").(string),
		"description":    d.Get("description").(string),
		"hmg_redundancy": "PriorityFailover",
	}

	group_object["add_hmg"] = []map[string]interface{}{
		{
			"url":             d.Get("url").(string),
			"kind":            "AZUREKEYVAULT",
			"client_id":       d.Get("client_id").(string),
			"tenant_id":       d.Get("tenant_id").(string),
			"subscription_id": d.Get("subscription_id").(string),
			"secret_key":      d.Get("secret_key").(string),
			// 0.5.0: FIXME: key_vault_type currently set to Standard only
			"key_vault_type": d.Get("key_vault_type").(string),
			"hsm_order":      0,
			"tls": map[string]interface{}{
				"mode":              "required",
				"validate_hostname": false,
				"ca": map[string]interface{}{
					"ca_set": "global_roots",
				},
			},
		},
	}

	req, err := m.(*api_client).APICallBody("POST", "sys/v1/groups", group_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST sys/v1/groups: %v", err),
		})
		return diags
	}

	d.SetId(req["group_id"].(string))
	return resourceReadAzureGroup(ctx, d, m)
}

// [R]: Read Azure Group
func resourceReadAzureGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, statuscode, err := m.(*api_client).APICall("GET", fmt.Sprintf("sys/v1/groups/%s", d.Id()))
	if statuscode == 404 {
		d.SetId("")
	} else {
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: GET sys/v1/groups: %v", err),
			})
			return diags
		} else {
			// Convert returned call into AWSGroup Map
			jsonbody, err := json.Marshal(req)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "[DSM SDK] Unable to parse DSM provider API client output",
					Detail:   fmt.Sprintf("[E]: API: GET crypto/v1/groups: %s", err),
				})
				return diags
			}

			azuregroup := AzureGroup{}
			if err := json.Unmarshal(jsonbody, &azuregroup); err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "[DSM SDK] Unable to parse DSM provider API client output",
					Detail:   fmt.Sprintf("[E]: API: GET crypto/v1/groups: %s", err),
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
			if _, ok := req["description"]; ok {
				d.Set("description", req["description"].(string))
			}
		}
	}
	return diags
}

// [U]: Update Azure Group
func resourceUpdateAzureGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.Set("secret_key", "")
	return diags
}

// [D]: Delete Azure Group
func resourceDeleteAzureGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	_, statuscode, err := m.(*api_client).APICall("DELETE", fmt.Sprintf("sys/v1/groups/%s", d.Id()))
	if (err != nil) && (statuscode != 404) && (statuscode != 400) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: DELETE sys/v1/groups: %v", err),
		})
		return diags
	} else {
		if statuscode == 400 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Call to DSM provider API client failed",
				Detail:   fmt.Sprintf("[E]: API: DELETE sys/v1/groups: %s", "Group Not Empty"),
			})
			return diags
		}
	}

	d.SetId("")
	return nil
}
