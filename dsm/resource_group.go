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
func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateGroup,
		ReadContext:   resourceReadGroup,
		UpdateContext: resourceUpdateGroup,
		DeleteContext: resourceDeleteGroup,
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
			"hmg_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_undo_policy_window_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// FIX ME
// Terraform internal state keeping does not have a native struct for JSON input and what they recommend is to
// keep it as a string. So, what we are keeping locally and what comes as a result is always different. Plus the
// response may have additional fields or the ordering may change. Since change detection is done by Terraform itself
// there is only one way out for us and that is to ignore API responses for JSON fields if the request is successful.
// Since what we get from the yaml is always kept in the internal state without getting overwritten with data from the
// API, change detection works. This is a very pathetic workaround and should be fixed as soon as Terraform improves its SDK.

// [C]: Create Group
func resourceCreateGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	group_object := make(map[string]interface{})

	group_object["name"] = d.Get("name").(string)

	if description, ok := d.GetOk("description"); ok {
		group_object["description"] = description.(string)
	}

	if approval_policy, ok := d.GetOk("approval_policy"); ok {
		group_object["approval_policy"] = json.RawMessage(approval_policy.(string))
	}

	if hmg, ok := d.GetOk("hmg"); ok {
		var hmg_object []json.RawMessage
		hmg_object = append(hmg_object, json.RawMessage(hmg.(string)))
		group_object["add_hmg"] = hmg_object
	}

	if debug_output {
		jj, _ := json.Marshal(group_object)
		tflog.Warn(ctx, fmt.Sprintf("Create Group Object: %s", jj))
	}
	set_key_undo_policy(d, group_object)

	resp, err := m.(*api_client).APICallBody("POST", "sys/v1/groups", group_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST sys/v1/groups: %v", err),
		})
		return diags
	}
	if debug_output {
		resp_json, _ := json.Marshal(resp)
		tflog.Warn(ctx, fmt.Sprintf("[U]: API response for group create operation: %s", resp_json))
	}
	d.SetId(resp["group_id"].(string))
	d.Set("group_id", resp["group_id"].(string))
	set_hmg_id(d, resp)
	return diags
}

// [U]: Update Group
func resourceUpdateGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	hmg_object := make(map[string]interface{})
	hmg_present := false

	if d.HasChange("description") || d.HasChange("name") || d.HasChange("approval_policy") || d.HasChange("hmg") ||
	                                d.HasChange("key_undo_policy_window_time") {
		if debug_output {
			tflog.Warn(ctx, "Group object has changed, calling API")
		}
		if hmg, ok := d.GetOk("hmg"); ok {
			hmg_id := substr(hmg.(string), 4, 36)
			if debug_output {
				tflog.Warn(ctx, fmt.Sprintf("HMG id: %s", hmg_id))
			}
			if d.Get("hmg_id") == nil {
			    resourceReadGroup(ctx, d, m)
			}
			hmg_object[d.Get("hmg_id").(string)] = json.RawMessage(hmg.(string))
			hmg_present = true
		}
		group_object := make(map[string]interface{})
		group_id := d.Get("group_id").(string)
		operation := "PATCH"
		url := fmt.Sprintf("sys/v1/groups/%s", group_id)
		if isSetApprovalPolicy(d, m) {
			if debug_output {
				tflog.Warn(ctx, "[U]: Approval policy is present.")
			}
			body_object := make(map[string]interface{})
			group_object["method"] = "PATCH"
			group_object["operation"] = url
			if approval_policy, ok := d.GetOk("approval_policy"); ok {
				if approval_policy == nil {
					body_object["approval_policy"] = make(map[string]interface{})
				} else {
					body_object["approval_policy"] = json.RawMessage(approval_policy.(string))
				}
			} else {
				// Handle delete quorum policy.
				// when approval_policy is removed from terraform body, the above 'if' condition fails.
				// Hence user, will not be able to remove the quorum policy.
				if (d.HasChange("approval_policy")) {
					body_object["approval_policy"] = make(map[string]interface{})
				}
			}
			if description, ok := d.GetOk("description"); ok {
				if description != "" {
					body_object["description"] = description
				}
			}
			set_key_undo_policy(d, body_object)
			if name, ok := d.GetOk("name"); ok {
				if name != "" {
					body_object["name"] = name
				} else {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "[Name check]",
						Detail:   "Group name can not be null",
					})
					return diags
				}
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

			if approval_policy, ok := d.GetOk("approval_policy"); ok {
				if approval_policy == nil {
					group_object["approval_policy"] = make(map[string]interface{})
				} else {
					group_object["approval_policy"] = json.RawMessage(approval_policy.(string))
				}
			}
			set_key_undo_policy(d, group_object)
			if description, ok := d.GetOk("description"); ok {
				if description != "" {
					group_object["description"] = description
				}
			}

			if name, ok := d.GetOk("name"); ok {
				if name != "" {
					group_object["name"] = name
				} else {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "[Name check]",
						Detail:   "Group name can not be null",
					})
					return diags
				}
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
	}

	return diags
}

// [R]: Read Group
func resourceReadGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, statuscode, err := m.(*api_client).APICall("GET", fmt.Sprintf("sys/v1/groups/%s", d.Id()))
	if statuscode == 404 {
		d.SetId("")
	} else {
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: GET sys/v1/groups/%s: %v", d.Id(), err),
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
		if description, ok := req["description"]; ok {
			if err := d.Set("description", description.(string)); err != nil {
				return diag.FromErr(err)
			}
		}
		set_hmg_id(d, req)
	}
	return diags
}

// [D]: Delete Group
func resourceDeleteGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	dataSourceGroupRead(ctx, d, m)
	group_id := d.Get("group_id").(string)

	_, statuscode, err := m.(*api_client).APICall("DELETE", fmt.Sprintf("sys/v1/groups/%s", group_id))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: DELETE sys/v1/groups/%s %d: %v", group_id, statuscode, err),
		})
		return diags
	}
	if statuscode == 400 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Call to DSM provider API client failed",
			Detail:   fmt.Sprintf("[E]: API: DELETE sys/v1/groups/%s Group is not empty", group_id),
		})
		return diags
	}

	d.SetId("")
	return nil
}

// [R]: Check Approval Policy presence
func isSetApprovalPolicy(d *schema.ResourceData, m interface{}) bool {
	req, statuscode, _ := m.(*api_client).APICall("GET", fmt.Sprintf("sys/v1/groups/%s", d.Id()))
	if statuscode == 200 {
		if _, ok := req["approval_policy"]; ok {
			return true
		}
	}
	return false
}
// handle the key undo policy for both create and update
func set_key_undo_policy(d *schema.ResourceData, obj map[string]interface{}) {
    if key_undo_policy_window_time, ok := d.GetOk("key_undo_policy_window_time"); ok {
        key_history_policy := make(map[string]interface{})
        key_history_policy["undo_time_window"] = key_undo_policy_window_time.(int)
        obj["key_history_policy"] = key_history_policy
    } else if (!d.IsNewResource()) {
		// Handle when undo policy is removed from terraform
		// The below condition checks whether undo policy is defined or not.
		if (d.HasChange("key_undo_policy_window_time")) {
			obj["key_history_policy"] = "remove"
		}
	}
}

func set_hmg_id(d *schema.ResourceData, resp map[string]interface{}) diag.Diagnostics {
    // set the hmg_id to update the cdc/byok attributes
    if _, ok := resp["hmg"]; ok {
        if hmg := resp["hmg"].(map[string]interface{}); len(hmg) > 0 {
            for k := range hmg {
                if err := d.Set("hmg_id", k); err != nil {
                    return diag.FromErr(err)
                }
            }
        }
    }
    return nil
}