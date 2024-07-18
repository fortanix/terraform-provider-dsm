package dsm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// [-] Define Group Cryptographic Policy
func resourceGroupCryptoPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateGroupCryptoPolicy,
		ReadContext:   resourceReadGroupCryptoPolicy,
		UpdateContext: resourceUpdateGroupCryptoPolicy,
		DeleteContext: resourceDeleteGroupCryptoPolicy,
		Description: "Adds cryptographic policy to a existing Fortnanix DSM group.",
		Schema: map[string]*schema.Schema{
			"name": {
			    Description: "The Fortanix DSM group object name.",
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
			    Description: "Group object ID from Fortanix DSM.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"acct_id": {
			    Description: "Account ID from Fortanix DSM.",
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
			    Description: "The Fortanix DSM group object description.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"approval_policy": {
			    Description: "The Fortanix DSM group object quorum approval policy definition as a JSON string.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"cryptographic_policy": {
			    Description: "The Fortanix DSM group object cryptographic policy definition as a JSON string",
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// [C]: Create Group Crypto Policy
func resourceCreateGroupCryptoPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	cryptographic_policy := json.RawMessage(d.Get("cryptographic_policy").(string))

	isSetApprovalPolicy, group_id := dataSourceGroupGetData(d, m)

	group_crypto_policy_object := make(map[string]interface{})
	if group_id == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM Provider: Group Crypto Policy] Unable to find group name",
			Detail:   fmt.Sprintf("[C & U]: Group not found: %s", d.Get("name").(string)),
		})
		return diags
	}
	if debug_output {
		tflog.Warn(ctx, fmt.Sprintf("Group id: ->%s<-", group_id))
	}
	operation := "PATCH"
	url := fmt.Sprintf("sys/v1/groups/%s", group_id)

	if isSetApprovalPolicy {
		if debug_output {
			tflog.Warn(ctx, "[C & U]: Approval policy is present.")
		}
		group_crypto_policy_object["method"] = "PATCH"
		group_crypto_policy_object["operation"] = url
		group_crypto_policy_object["body"] = map[string]interface{}{"cryptographic_policy": cryptographic_policy}
		operation = "POST"
		url = "sys/v1/approval_requests"
	} else {
		if debug_output {
			tflog.Warn(ctx, "[C & U]: Approval policy is not set.")
		}
		group_crypto_policy_object["group_id"] = group_id
		group_crypto_policy_object["cryptographic_policy"] = cryptographic_policy
	}

	resp, err := m.(*api_client).APICallBody(operation, url, group_crypto_policy_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[C & U]: API Call: %s %s: %v", operation, url, err),
		})
		return diags
	}

	if debug_output {
		resp_json, _ := json.Marshal(resp)
		tflog.Warn(ctx, fmt.Sprintf("[C & U]: API response for cryptographic policy create operation: %s", resp_json))
	}

	d.SetId(group_id)
	return diags
}

// [U]: Update Group Crypto Policy
func resourceUpdateGroupCryptoPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.HasChange("name") || d.HasChange("cryptographic_policy") {
		return resourceCreateGroupCryptoPolicy(ctx, d, m)
	}

	return nil
}

// [R]: Read Group Crypto Policy
func resourceReadGroupCryptoPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, statuscode, err := m.(*api_client).APICall("GET", fmt.Sprintf("sys/v1/groups/%s", d.Id()))
	if statuscode == 404 {
		d.SetId("")
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] DSM provider API client returned not found",
			Detail:   fmt.Sprintf("[R]: API Call: GET sys/v1/groups/%s", d.Id()),
		})
		return diags
	}
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[R]: API Call: GET sys/v1/groups/%s: %v", d.Id(), err),
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

	return diags
}

// [D]: Delete Group Crypto Policy
func resourceDeleteGroupCryptoPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	isSetApprovalPolicy, group_id := dataSourceGroupGetData(d, m)

	group_crypto_policy_object := make(map[string]interface{})
	cryptographic_policy := "remove"
	operation := "PATCH"
	url := fmt.Sprintf("sys/v1/groups/%s", group_id)

	if isSetApprovalPolicy {
		if debug_output {
			tflog.Warn(ctx, "[D]: Approval policy is present.")
		}
		group_crypto_policy_object["method"] = "PATCH"
		group_crypto_policy_object["operation"] = url
		group_crypto_policy_object["body"] = map[string]interface{}{"cryptographic_policy": cryptographic_policy}
		operation = "POST"
		url = "sys/v1/approval_requests"
	} else {
		if debug_output {
			tflog.Warn(ctx, "[D]: Approval policy is not set.")
		}
		group_crypto_policy_object["group_id"] = group_id
		group_crypto_policy_object["cryptographic_policy"] = cryptographic_policy
	}

	resp, err := m.(*api_client).APICallBody(operation, url, group_crypto_policy_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[D]: API Call: %s %s: %v", operation, url, err),
		})
		return diags
	}

	if debug_output {
		resp_json, _ := json.Marshal(resp)
		tflog.Warn(ctx, fmt.Sprintf("[D]: API response for cryptographic policy delete operation: %s", resp_json))
	}

	d.SetId(group_id)
	return nil
}

func dataSourceGroupGetData(d *schema.ResourceData, m interface{}) (bool, string) {

	req, err := m.(*api_client).APICallList("GET", "sys/v1/groups")
	if err != nil {
		return false, ""
	}

	group_id := ""
	flag := false
	for _, data := range req {
		if data.(map[string]interface{})["name"].(string) == d.Get("name").(string) {
			group_id = data.(map[string]interface{})["group_id"].(string)
			if _, ok := data.(map[string]interface{})["approval_policy"]; ok {
				flag = true
			}
		}
	}

	return flag, group_id
}
