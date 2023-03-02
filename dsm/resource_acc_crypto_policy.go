package dsm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// [-] Define Account Cryptographic Policy
func resourceAccountCryptoPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateAccountCryptoPolicy,
		ReadContext:   resourceReadAccountCryptoPolicy,
		UpdateContext: resourceUpdateAccountCryptoPolicy,
		DeleteContext: resourceDeleteAccountCryptoPolicy,
		Schema: map[string]*schema.Schema{
			"acct_id": {
				Type:     schema.TypeString,
				Required: true,
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

// [C]: Create Account Crypto Policy
func resourceCreateAccountCryptoPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	cryptographic_policy := json.RawMessage(d.Get("cryptographic_policy").(string))

	accountApprovalPolicyRead(ctx, d, m)

	account_crypto_policy_object := make(map[string]interface{})
	acct_id := d.Get("acct_id").(string)
	operation := "PATCH"
	url := fmt.Sprintf("sys/v1/accounts/%s", acct_id)

	if _, ok := d.GetOk("approval_policy"); ok {
		if debug_output {
			tflog.Warn(ctx, "[C & U]: Approval policy is present.")
		}
		account_crypto_policy_object["method"] = "PATCH"
		account_crypto_policy_object["operation"] = url
		account_crypto_policy_object["body"] = map[string]interface{}{"cryptographic_policy": cryptographic_policy}
		operation = "POST"
		url = "sys/v1/approval_requests"
	} else {
		if debug_output {
			tflog.Warn(ctx, "[C & U]: Approval policy is not set.")
		}
		account_crypto_policy_object["acct_id"] = acct_id
		account_crypto_policy_object["cryptographic_policy"] = cryptographic_policy
	}

	resp, derr := m.(*api_client).APICallBody(operation, url, account_crypto_policy_object)
	if derr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[C & U]: API Call: %s %s: %v", operation, url, derr),
		})
		return diags
	}

	if debug_output {
		resp_json, _ := json.Marshal(resp)
		tflog.Warn(ctx, fmt.Sprintf("[C & U]: API response for cryptographic policy create operation: %s", resp_json))
	}

	d.SetId(acct_id)
	return diags
}

// [U]: Update Account Crypto Policy
func resourceUpdateAccountCryptoPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.HasChange("approval_policy") || d.HasChange("cryptographic_policy") {
		return resourceCreateAccountCryptoPolicy(ctx, d, m)
	}

	return nil
}

// [R]: Read Account Crypto Policy
func resourceReadAccountCryptoPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, statuscode, err := m.(*api_client).APICall("GET", fmt.Sprintf("sys/v1/accounts/%s", d.Id()))
	if statuscode == 404 {
		d.SetId("")
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] DSM provider API client returned not found",
			Detail:   fmt.Sprintf("[R]: API Call: GET sys/v1/accounts/%s", d.Id()),
		})
		return diags
	}
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[R]: API Call: GET sys/v1/accounts/%s: %v", d.Id(), err),
		})
		return diags
	}

	if err := d.Set("acct_id", req["acct_id"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if _, ok := req["approval_policy"]; ok {
		if err := d.Set("approval_policy", fmt.Sprintf("%s", req["approval_policy"])); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

// [D]: Delete Account Crypto Policy
func resourceDeleteAccountCryptoPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	accountApprovalPolicyRead(ctx, d, m)

	account_crypto_policy_object := make(map[string]interface{})
	acct_id := d.Get("acct_id").(string)
	cryptographic_policy := "remove"
	operation := "PATCH"
	url := fmt.Sprintf("sys/v1/accounts/%s", acct_id)

	if _, ok := d.GetOk("approval_policy"); ok {
		if debug_output {
			tflog.Warn(ctx, "[D]: Approval policy is present.")
		}
		account_crypto_policy_object["method"] = "PATCH"
		account_crypto_policy_object["operation"] = url
		account_crypto_policy_object["body"] = map[string]interface{}{"cryptographic_policy": cryptographic_policy}
		operation = "POST"
		url = "sys/v1/approval_requests"
	} else {
		if debug_output {
			tflog.Warn(ctx, "[D]: Approval policy is not set.")
		}
		account_crypto_policy_object["acct_id"] = acct_id
		account_crypto_policy_object["cryptographic_policy"] = cryptographic_policy
	}

	resp, err := m.(*api_client).APICallBody(operation, url, account_crypto_policy_object)
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

	d.SetId(acct_id)
	return nil
}

// Get account details
func accountApprovalPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, statuscode, err := m.(*api_client).APICall("GET", fmt.Sprintf("sys/v1/accounts/%s", d.Get("acct_id").(string)))
	if statuscode == 404 {
		d.SetId("")
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] DSM provider API client returned not found",
			Detail:   fmt.Sprintf("[R]: API Call: GET sys/v1/accounts/%s", d.Get("acct_id").(string)),
		})
		return diags
	} else {
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[R]: API Call: GET sys/v1/accounts/%s: %v", d.Get("acct_id").(string), err),
			})
			return diags
		}
		if debug_output {
			tflog.Warn(ctx, fmt.Sprintf("[R]: API read account id: %s", req["acct_id"]))
		}
		if _, ok := req["approval_policy"]; ok {
			if err := d.Set("approval_policy", fmt.Sprintf("%s", req["approval_policy"])); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if debug_output {
		tflog.Warn(ctx, fmt.Sprintf("[R]: API read account approval policy: %s", req["approval_policy"]))
	}
	return diags
}
