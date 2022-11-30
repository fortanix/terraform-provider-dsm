package dsm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// [-] Define Group
func resourceExistingGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateExistingGroup,
		ReadContext:   resourceReadExistingGroup,
		UpdateContext: resourceUpdateExistingGroup,
		DeleteContext: resourceDeleteExistingGroup,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"approval_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hmg": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// [C]: Create Group - Not applicable for managing existing groups
func resourceCreateExistingGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

// [U]: Update Group
func resourceUpdateExistingGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var approval_policy_new json.RawMessage = nil
	var hmg_new json.RawMessage
	hmg_object := make(map[string]interface{})
	description_new := ""
	name_new := ""
	hmg_present := false

	if d.HasChange("approval_policy") {
		if approval_policy, ok := d.GetOk("approval_policy"); ok {
			approval_policy_new = json.RawMessage(approval_policy.(string))
			d.Set("approval_policy", nil)
		}
	}

	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			description_new = description.(string)
		}
	}

	if d.HasChange("name") {
		old, new := d.GetChange("name")
		name_new = new.(string)
		d.Set("name", old)
	}

	if hmg, ok := d.GetOk("hmg"); ok {
		hmg_new = json.RawMessage(hmg.(string))
		d.Set("hmg", nil)
	}

	read_diags := resourceReadGroup(ctx, d, m)
	if read_diags != nil {
		return read_diags
	}

	if hmg, ok := d.GetOk("hmg"); ok {
		hmg_id := substr(hmg.(string), 4, 36)
		if debug_output {
			tflog.Warn(ctx, fmt.Sprintf("HMG id: %s", hmg_id))
		}
		hmg_object[hmg_id] = hmg_new
		hmg_present = true
	}

	group_object := make(map[string]interface{})
	group_id := d.Get("group_id").(string)
	if group_id == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM Provider: Group] Unable to find group name",
			Detail:   fmt.Sprintf("[U]: Group not found: %s", name_new),
		})
		return diags
	}
	operation := "PATCH"
	url := fmt.Sprintf("sys/v1/groups/%s", group_id)

	if debug_output {
		tflog.Warn(ctx, fmt.Sprintf("Update operation group id: %s", group_id))
		tflog.Warn(ctx, fmt.Sprintf("New Approval Policy: %s", approval_policy_new))
		tflog.Warn(ctx, fmt.Sprintf("Old Approval Policy: %s", json.RawMessage(d.Get("approval_policy").(string))))
		tflog.Warn(ctx, fmt.Sprintf("New HMG object: %s", hmg_object))
		tflog.Warn(ctx, fmt.Sprintf("Old HMG: %s", d.Get("hmg").(string)))
	}

	if _, ok := d.GetOk("approval_policy"); ok {
		if debug_output {
			tflog.Warn(ctx, "[U]: Approval policy is present.")
		}
		body_object := make(map[string]interface{})
		group_object["method"] = "PATCH"
		group_object["operation"] = url
		if approval_policy_new == nil {
			body_object["approval_policy"] = make(map[string]interface{})
		} else {
			body_object["approval_policy"] = approval_policy_new
		}
		if description_new != "" {
			body_object["description"] = description_new
		}
		if name_new != "" {
			body_object["name"] = name_new
		}
		if hmg_present {
			body_object["mod_hmg"] = hmg_object
		}
		group_object["body"] = body_object
		operation = "POST"
		url = "sys/v1/approval_requests"
	} else {
		if debug_output {
			tflog.Warn(ctx, "[U]: Approval policy is not set.")
		}
		group_object["group_id"] = group_id
		if approval_policy_new == nil {
			group_object["approval_policy"] = make(map[string]interface{})
		} else {
			group_object["approval_policy"] = approval_policy_new
		}
		if description_new != "" {
			group_object["description"] = description_new
		}
		if name_new != "" {
			group_object["name"] = name_new
		}
		if hmg_present {
			group_object["mod_hmg"] = hmg_object
		}
	}

	if debug_output {
		jj, _ := json.Marshal(group_object)
		tflog.Warn(ctx, fmt.Sprintf("Update Group Object: %s", jj))
	}

	resp, err := m.(*api_client).APICallBody(operation, url, group_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[U]: API Call: %s %s: %v", operation, url, err),
		})
		return diags
	}

	if debug_output {
		resp_json, _ := json.Marshal(resp)
		tflog.Warn(ctx, fmt.Sprintf("[U]: API response for group update operation: %s", resp_json))
	}

	d.SetId(group_id)
	return resourceReadGroup(ctx, d, m)
}

// [R]: Read Group
func resourceReadExistingGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		}

		if err := d.Set("name", req["name"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("group_id", req["group_id"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("acct_id", req["acct_id"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("creator", req["creator"]); err != nil {
			return diag.FromErr(err)
		}
		if _, ok := req["description"]; ok {
			if err := d.Set("description", req["description"].(string)); err != nil {
				return diag.FromErr(err)
			}
		}
		if _, ok := req["approval_policy"]; ok {
			if err := d.Set("approval_policy", fmt.Sprintf("%v", req["approval_policy"])); err != nil {
				return diag.FromErr(err)
			}
		}
		if _, ok := req["hmg"]; ok {
			if err := d.Set("hmg", fmt.Sprintf("%v", req["hmg"])); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return diags
}

// [D]: Delete Group - Not much helpful for managing existing groups
func resourceDeleteExistingGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
