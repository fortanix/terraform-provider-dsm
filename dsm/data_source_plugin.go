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
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"plugin_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plugin_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"language": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"acct_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator": {
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
