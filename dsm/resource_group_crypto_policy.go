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
		UpdateContext: resourceCreateGroupCryptoPolicy,
		DeleteContext: resourceDeleteGroupCryptoPolicy,
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
				Computed: true,
			},
			"approval_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cryptographic_policy": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// [C & U]: Create and Update Group Crypto Policy
func resourceCreateGroupCryptoPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	cryptographic_policy := json.RawMessage(d.Get("cryptographic_policy").(string))

	dataSourceGroupRead(ctx, d, m)

	group_crypto_policy_object := make(map[string]interface{})
	group_id := d.Get("group_id").(string)
	if debug_output {
		tflog.Warn(ctx, fmt.Sprintf("Group id: ->%s<-", group_id))
	}
	operation := "PATCH"
	url := fmt.Sprintf("sys/v1/groups/%s", group_id)

	if _, ok := d.GetOk("approval_policy"); ok {
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
	return resourceReadGroupCryptoPolicy(ctx, d, m)
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
	if _, ok := req["approval_policy"]; ok {
		if err := d.Set("approval_policy", fmt.Sprintf("%v", req["approval_policy"])); err != nil {
			return diag.FromErr(err)
		}
	}
	if _, ok := req["cryptographic_policy"]; ok {
		if err := d.Set("cryptographic_policy", fmt.Sprintf("%v", req["cryptographic_policy"])); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if debug_output {
			tflog.Warn(ctx, "[R]: Expected cryptographic policy but found none. Operation might be pending if a quorum policy has been set.")
		}
	}

	return diags
}

// [D]: Delete Group Crypto Policy
func resourceDeleteGroupCryptoPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	dataSourceGroupRead(ctx, d, m)

	group_crypto_policy_object := make(map[string]interface{})
	group_id := d.Get("group_id").(string)
	cryptographic_policy := "remove"
	operation := "PATCH"
	url := fmt.Sprintf("sys/v1/groups/%s", group_id)

	if approval_policy, ok := d.GetOk("approval_policy"); ok {
		if debug_output {
			tflog.Warn(ctx, fmt.Sprintf("[D]: Approval policy is present: %s", approval_policy))
		}
		group_crypto_policy_object["method"] = "PATCH"
		group_crypto_policy_object["operation"] = url
		group_crypto_policy_object["body"] = map[string]interface{}{"cryptographic_policy": cryptographic_policy}
		operation = "POST"
		url = "sys/v1/approval_requests"
	} else {
		if debug_output {
			tflog.Warn(ctx, fmt.Sprintf("[D]: Approval policy is not set: %s", approval_policy))
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
