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
func resourcePlugin() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreatePlugin,
		ReadContext:   resourceReadPlugin,
		UpdateContext: resourceUpdatePlugin,
		DeleteContext: resourceDeletePlugin,
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
				Optional: true,
				Default: "",
			},
			"plugin_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default: "standard",
			},
			"default_group": {
				Type:     schema.TypeString,
				Required: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"language": {
				Type:     schema.TypeString,
				Optional: true,
				Default: "LUA",
			},
			"code": {
				 Type:     schema.TypeString,
				 Required: true,
			 },
			 "enabled": {
				 Type:     schema.TypeBool,
				 Optional: true,
				 Default: true,
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
			State: schema.ImportStatePassthrough,
		},
	}
}
var endpoint = "sys/v1/plugins"

 // Create
func resourceCreatePlugin(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics{
	var diags diag.Diagnostics

	plugin := map[string]interface{}{
		"name":          d.Get("name").(string),
		"default_group": d.Get("default_group").(string),
		"description":   d.Get("description").(string),
		"groups":        d.Get("groups").([]interface{}),
	}
	source := make(map[string]string)
	source["language"] = d.Get("language").(string)
	source["code"] = d.Get("code").(string)
	plugin["source"] = source
	if err := d.Get("groups").([]interface{}); len(err) > 0 {
		plugin["add_groups"] = err
	}
	if err := d.Get("plugin_type").(string); len(err) > 0 {
		plugin["plugin_type"] = d.Get("plugin_type")
	}
	if err := d.Get("enabled").(bool); err {
		plugin["enabled"] = d.Get("enabled")
	}
	isapprovalPolicy := isApprovalPolicy(d.Get("groups").([]interface{}), m)
	if (isapprovalPolicy) {
		approval_request_body := map[string]interface{}{
			"operation": endpoint,
			"body":      plugin,
			"method":    "POST",
		}
		req, err := m.(*api_client).APICallBody("POST", "sys/v1/approval_requests", approval_request_body)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: POST %s: sys/v1/approval_requests", err),
			})
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "This plugin creation requires approval request. Please get the approval from required users in UI.",
				Detail:   fmt.Sprintf("[E]: API: POST %s: sys/v1/plugins", req["request_id"]),
			})
		}
		return diags
	} else {
		req, err := m.(*api_client).APICallBody("POST", endpoint, plugin)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: POST %s: %v", endpoint, err),
			})
			return diags
		}
		d.SetId(req["plugin_id"].(string))
		return resourceReadPlugin(ctx, d, m)
	}
}

// Read
func resourceReadPlugin(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, _, err := m.(*api_client).APICall("GET", fmt.Sprintf(endpoint + "/%s", d.Id()))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: GET sys/v1/apps: %v", err),
		})
		return diags
	}
	if err := d.Set("name", req["name"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("plugin_id", req["plugin_id"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("default_group", req["default_group"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enabled", req["enabled"].(bool)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("acct_id", req["acct_id"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("groups", req["groups"].([]interface{})); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("creator", req["creator"]); err != nil {
		return diag.FromErr(err)
	}
	if _, ok := req["source"]; ok {
		if source := req["source"].(map[string]interface{}); len(source) > 0 {
			if err := d.Set("language", source["language"].(string)); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("code", source["code"].(string)); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return diags
}

// Update
func resourceUpdatePlugin(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	plugin := make(map[string]interface{})
	if d.HasChange("name") {
		plugin["name"] = d.Get("name").(string)
	}
	if d.HasChange("description") {
		plugin["description"] = d.Get("description").(string)
	}
	if d.HasChange("enabled") {
		plugin["enabled"] = d.Get("enabled")
	}
	if d.HasChange("code") {
		source := make(map[string]string)
		source["language"] = d.Get("language").(string)
		source["code"] = d.Get("code").(string)
		plugin["source"] = source
	}
	if d.HasChange("default_group") {
		plugin["default_group"] = d.Get("default_group").(string)
	}
	if d.HasChange("groups") {
		old_groups, new_groups := d.GetChange("groups")
		// compute_add_and_del_arrays function is in common.go
		add_group_ids, del_group_ids := compute_add_and_del_arrays(old_groups, new_groups)
		if len(del_group_ids) > 0 {
			plugin["del_groups"] = del_group_ids
		}
		if len(add_group_ids) > 0 {
			plugin["add_groups"] = add_group_ids
		}
	}
	if d.HasChange("plugin_type") {
		plugin["plugin_type"] = d.Get("plugin_type")
	}
	_, err := m.(*api_client).APICallBody("PATCH", fmt.Sprintf(endpoint + "/%s", d.Id()), plugin)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: PATCH %s: %v", endpoint, err),
		})
		return diags
	}
	return resourceReadPlugin(ctx, d, m)
}

// Delete
func resourceDeletePlugin(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	_, statuscode, err := m.(*api_client).APICall("DELETE", fmt.Sprintf(endpoint + "/%s", d.Id()))
	if (err != nil) && (statuscode != 404) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: DELETE sys/v1/apps: %v", err),
		})
		return diags
	}
	d.SetId("")
	return nil
}

// Read Group
func isApprovalPolicy(group_ids ([]interface{}), m interface{}) bool {
	group_ids_arr := make([]string, len(group_ids))
	for i, v := range group_ids {
		group_ids_arr[i] = v.(string)
	}
	for _, group_id := range group_ids_arr {
		req, statuscode, _ := m.(*api_client).APICall("GET", fmt.Sprintf("sys/v1/groups/%s", group_id))
		if (statuscode == 200) {
			if _, ok := req["approval_policy"]; ok {
				return true
			}
		}
	}
	return false
}
