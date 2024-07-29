// **********
// Terraform Provider - SDKMS: resource: security object
// **********
//       - Author:    Ravi Gopal at fortanix dot com
//       - Version:   0.5.1
//       - Date:      21/10/2022
// **********

package dsm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// [-] Define Plugin
func dataSourcePlugin() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadPlugin,
		Description: "Returns the Fortanix DSM plugin object from the cluster as a Resource.",
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The Fortanix DSM plugin object name.",
				Type:     schema.TypeString,
				Required: true,
			},
			"plugin_id": {
				Description: "Plugin object ID from Fortanix DSM.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Description: "The Fortanix DSM plugin object description.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_group": {
				Description: "The Fortanix DSM group object ID that is mapped to the plugin by default.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"groups": {
				Description: "List of other Fortanix DSM group object IDs that are mapped to the plugin.",
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"language": {
				Description: "Programming language for plugin code (Default value is `LUA`).",
				Type:     schema.TypeString,
				Computed: true,
			},
			"code": {
				Description: "Plugin code that will be executed in DSM. Code should be in specified programming language.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Description: "Whether the security object is enabled or disabled. The values are true/false.",
				Type:     schema.TypeBool,
				Computed: true,
			},
			"acct_id": {
				Description: "Account ID from Fortanix DSM.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator": {
				Description: "The creator of the security object from Fortanix DSM.\n" +
				"   * `user`: If the plugin object was created by a user, the computed value will be the matching user id.\n" +
				"   * `app`: If the plugin object was created by a app, the computed value will be the matching app id.",
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// Read
func dataSourceReadPlugin(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, err := m.(*api_client).APICallList("GET", "sys/v1/plugins")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: GET sys/v1/plugins: %v", err),
		})
		return diags
	}
	for _, data := range req {
		if data.(map[string]interface{})["name"].(string) == d.Get("name").(string) {
			if err := d.Set("name", data.(map[string]interface{})["name"].(string)); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("description", data.(map[string]interface{})["description"].(string)); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("plugin_id", data.(map[string]interface{})["plugin_id"].(string)); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("default_group", data.(map[string]interface{})["default_group"].(string)); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("enabled", data.(map[string]interface{})["enabled"].(bool)); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("acct_id", data.(map[string]interface{})["acct_id"].(string)); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("groups", data.(map[string]interface{})["groups"].([]interface{})); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("creator", data.(map[string]interface{})["creator"]); err != nil {
				return diag.FromErr(err)
			}
			if _, ok := data.(map[string]interface{})["source"].(map[string]interface{}); ok {
				source := data.(map[string]interface{})["source"].(map[string]interface{})
				if err := d.Set("language", source["language"].(string)); err != nil {
					return diag.FromErr(err)
				}
				if err := d.Set("code", source["code"].(string)); err != nil {
					return diag.FromErr(err)
				}
			}
			break
		}
	}
	d.SetId(d.Get("plugin_id").(string))
	return nil
}
