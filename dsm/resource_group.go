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
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// [C]: Create Group
func resourceCreateGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	group_object := make(map[string]interface{})

	group_object["name"] = d.Get("name").(string)
	if _, ok := d.GetOk("description"); ok {
		group_object["description"] = d.Get("description").(string)
	}
	if _, ok := d.GetOk("approval_policy"); ok {
		group_object["approval_policy"] = json.RawMessage(d.Get("approval_policy").(string))
	}
	if _, ok := d.GetOk("hmg"); ok {
		var hmg_object []json.RawMessage
		hmg_object[0] = json.RawMessage(d.Get("hmg").(string))
		group_object["add_hmg"] = hmg_object
	}
	resp, err := m.(*api_client).APICallBody("POST", "sys/v1/groups", group_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST sys/v1/groups: %v", err),
		})
		return diags
	}

	resp_json, _ := json.Marshal(resp)
	tflog.Warn(ctx, fmt.Sprintf("[U]: API response for group create operation: %s", resp_json))

	d.SetId(resp["group_id"].(string))
	return resourceReadGroup(ctx, d, m)
}

// [U]: Update Group
func resourceUpdateGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var approval_policy_new json.RawMessage = nil
	description_new := ""
	name_new := ""

	if approval_policy, ok := d.GetOk("approval_policy"); ok {
		approval_policy_new = json.RawMessage(approval_policy.(string))
		d.Set("approval_policy", nil)
	}

	if description, ok := d.GetOk("description"); ok {
		description_new = description.(string)
	}

	if d.HasChange("name") {
		name_new = d.Get("name").(string)
	}

	dataSourceGroupRead(ctx, d, m)

	group_object := make(map[string]interface{})
	body_object := make(map[string]interface{})
	group_id := d.Get("group_id").(string)
	operation := "PATCH"
	url := fmt.Sprintf("sys/v1/groups/%s", group_id)

	tflog.Warn(ctx, fmt.Sprintf("Update operation group id: ->%s<-", group_id))
	tflog.Warn(ctx, fmt.Sprintf("New Approval Policy: ->%s<-", approval_policy_new))
	tflog.Warn(ctx, fmt.Sprintf("Old Approval Policy: ->%s<-", json.RawMessage(d.Get("approval_policy").(string))))

	if _, ok := d.GetOk("approval_policy"); ok {
		tflog.Warn(ctx, "[U]: Approval policy is present.")
		group_object["method"] = "PATCH"
		group_object["operation"] = url
		if approval_policy_new == nil {
			body_object["approval_policy"] = make(map[string]interface{})
		} else {
			body_object["approval_policy"] = approval_policy_new
		}
		body_object["description"] = description_new
		if name_new != "" {
			body_object["name"] = name_new
		}
		group_object["body"] = body_object
		operation = "POST"
		url = "sys/v1/approval_requests"
	} else {
		tflog.Warn(ctx, "[U]: Approval policy is not set.")
		group_object["group_id"] = group_id
		if approval_policy_new == nil {
			group_object["approval_policy"] = make(map[string]interface{})
		} else {
			group_object["approval_policy"] = approval_policy_new
		}
		group_object["description"] = description_new
		if name_new != "" {
			group_object["name"] = name_new
		}
	}

	jj, _ := json.Marshal(group_object)
	tflog.Warn(ctx, fmt.Sprintf("Group Object: %s", jj))

	resp, err := m.(*api_client).APICallBody(operation, url, group_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[U]: API Call: %s %s: %v", operation, url, err),
		})
		return diags
	}

	resp_json, _ := json.Marshal(resp)
	tflog.Warn(ctx, fmt.Sprintf("[U]: API response for group update operation: %s", resp_json))

	d.SetId(group_id)
	return resourceReadGroup(ctx, d, m)
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
