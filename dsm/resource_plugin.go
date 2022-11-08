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
/*

When there is an approval policy configured in plugin, then the request should be
redirected to approval_requests API.

Algorithm:

Create Plugin:

1) Read the inputs
2) Read each group and check if there is approval_policy
	* if approval_policy is configured in any given groups configured,then it will redirect to approval_requests API.
		1) ID will set as approval_request_id and plugin_id as req["request_id"]. User should approve from UI currently
		2) Once user is approved, when user tries to do any change or apply it will read all the plugins
		   and gets the plugin that matches the plugin name.
		3) Adds the plugin_id, Id as req["plugin_id"] and approval_request_id as null.
	* if approval_policy is not configured in any given groups configured,then it will create a plugin.

Update Plugin:

1) Read the inputs
2) Read each group and check if there is approval_policy
	* if approval_policy is configured in any given groups configured,then it will redirect to approval_requests API.
		1) approval_request_id will set as request_id. User should approve from UI currently
		2) Once user is approved, when user tries to do any change or apply it will apply the new changes in state
		   and changes approval_request_id as null
	* if approval_policy is not configured in any given groups configured,then it will create a plugin.


*/
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
			"approval_request_id" : {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
var plugin_endpoint = "sys/v1/plugins"
var approval_endpoint = "sys/v1/approval_requests"

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
	// Checks if any group has approval policy
	isapprovalPolicy := isApprovalPolicy(d.Get("groups").([]interface{}), m)
	// If approval policy exists then it redirects to approval_request API
	if (isapprovalPolicy) {
		return approvalRequestCall(plugin, d, m, "POST")
	}
	req, err := m.(*api_client).APICallBody("POST", plugin_endpoint, plugin)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST %s: %v", plugin_endpoint, err),
		})
		return diags
	}
	d.SetId(req["plugin_id"].(string))
	return resourceReadPlugin(ctx, d, m)
}

// Read
func resourceReadPlugin(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// This will be the case when approval_request API is triggered during create
	if pid := d.Get("plugin_id").(string); len(pid) == 0 && len(d.Get("approval_request_id").(string)) > 0{
		// Checks whether approval_request_id is approved or not
		req, statuscode, _ := m.(*api_client).APICall("GET", fmt.Sprintf(approval_endpoint + "/%s", d.Id()))
		if statuscode == 200 {
			if req["status"] == "APPROVED" {
				req, _ := m.(*api_client).APICallList("GET", fmt.Sprintf(plugin_endpoint))
				for _, data := range req {
					if data.(map[string]interface{})["name"].(string) == d.Get("name").(string) {
						// If plugin available the changes ID as plugin_id
						// and approval_request_id as ""
						d.SetId(data.(map[string]interface{})["plugin_id"].(string))
						d.Set("approval_request_id", "")
						break
					}
				}
			} else if req["status"] == "PENDING" {
				 diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Plugin " + d.Get("name").(string) + " is not yet approved or denied from the required users." +
							  "request_id is: " + d.Id(),
					Detail:   fmt.Sprintf("[W]: API: GET %s: %s", plugin_endpoint, d.Get("name").(string)),
				 })
				return diags
			} else if req["status"] == "DENIED" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Plugin " + d.Get("name").(string) + " is denied from the required user." +
							  "To create a new plugin change the name or delete it from state and recreate it.",
					Detail:   fmt.Sprintf("[W]: API: GET %s: %s", plugin_endpoint, d.Get("name").(string)),
				})
				return diags
			}
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to load the request_id " + d.Id(),
				Detail:   fmt.Sprintf("[E]: API: GET %s: %s", approval_endpoint, d.Id()),
			})
			return diags
		}
	}
	// This will be executed during update
	if approval_rq_id := d.Get("approval_request_id").(string); len(approval_rq_id) > 0{
		req, statuscode, _ := m.(*api_client).APICall("GET", fmt.Sprintf(approval_endpoint + "/%s", approval_rq_id))
		if statuscode == 200 {
			// When it is approved or denied it will make approval_request_id as null
			// And reads the plugin
			if req["status"] != "PENDING" {
				 d.Set("approval_request_id", "")
			} else {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Plugin " + d.Get("name").(string) + " is not yet approved or denied from the required users." +
							  "request_id is: " + approval_rq_id,
					Detail:   fmt.Sprintf("[W]: API: GET %s: %s", plugin_endpoint, d.Get("name").(string)),
				})
				return diags
			}
		} else {
			d.Set("approval_request_id", "")
		}
	}
	// reads the plugin
	req, _, err := m.(*api_client).APICall("GET", fmt.Sprintf(plugin_endpoint + "/%s", d.Id()))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: GET %s: %v", plugin_endpoint, err),
		})
		return diags
	}
	if err := d.Set("name", req["name"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", req["description"].(string)); err != nil {
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
	plugin["name"] = d.Get("name").(string)
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
	isapprovalPolicy := isApprovalPolicy(d.Get("groups").([]interface{}), m)
	if (isapprovalPolicy) {
		return approvalRequestCall(plugin, d, m, "PATCH")
	}
	_, err := m.(*api_client).APICallBody("PATCH", fmt.Sprintf(plugin_endpoint + "/%s", d.Id()), plugin)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: PATCH %s: %v", plugin_endpoint, err),
		})
		return diags
	}
	return resourceReadPlugin(ctx, d, m)
}

// Delete
func resourceDeletePlugin(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if len(d.Get("approval_request_id").(string)) > 0 {
		m.(*api_client).APICall("POST", fmt.Sprintf(approval_endpoint + "/%s/deny", d.Id()))
	}
	if len(d.Get("plugin_id").(string)) > 0 {
		_, statuscode, err := m.(*api_client).APICall("DELETE", fmt.Sprintf(plugin_endpoint + "/%s", d.Id()))
		if (err != nil) && (statuscode != 404) {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: DELETE %s: %v", plugin_endpoint, err),
			})
			return diags
		}
	}
	d.SetId("")
	return nil
}

// Approval request call
func approvalRequestCall(body interface{}, d *schema.ResourceData, m interface{}, method string) diag.Diagnostics{
	var diags diag.Diagnostics

	operation := plugin_endpoint
	if method == "PATCH" {
		operation = plugin_endpoint + "/" + d.Id()
	}
	approval_request_body := map[string]interface{}{
		"operation": operation,
		"body":      body,
		"method":    method,
	}
		req, err := m.(*api_client).APICallBody("POST", approval_endpoint, approval_request_body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST %s %s", approval_endpoint, err),
		})
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "This plugin creation requires approval request. Please get the approval from required users in UI.",
			Detail:   fmt.Sprintf("[W]: API: POST %s, request_id for the plugin: %s", plugin_endpoint, req["request_id"]),
		})
		// sets ID as request_id
		if method == "POST" {
			d.SetId(req["request_id"].(string))
		}
		d.Set("approval_request_id", req["request_id"])
	}
	return diags
}

// Read each group and check if there is an approval policy
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
