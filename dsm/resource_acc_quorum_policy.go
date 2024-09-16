package dsm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// [-] Define Account Quorum Policy
func resourceAccountQuorumPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateAccountQuorumPolicy,
		ReadContext:   resourceReadAccountQuorumPolicy,
		UpdateContext: resourceUpdateAccountQuorumPolicy,
		DeleteContext: resourceDeleteAccountQuorumPolicy,
		Description: "Adds cryptographic policy to a Fortanix DSM account. Quorum approval policy adds an extra level of protection to sensitive account operations.",
		Schema: map[string]*schema.Schema{
			"acct_id": {
			    Description: "The Fortanix DSM account object id.",
				Type:     schema.TypeString,
				Required: true,
			},
			"approval_policy": {
			    Description: "The Fortanix DSM account object quorum approval policy definition as a JSON string.",
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// [C]: Create Account Quorum Policy
func resourceCreateAccountQuorumPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	policy_object := map[string]interface{}{
		"acct_id":         d.Get("acct_id").(string),
		"approval_policy": json.RawMessage(d.Get("approval_policy").(string)),
	}

	req, err := m.(*api_client).APICallBody("PATCH", fmt.Sprintf("sys/v1/accounts/%s", policy_object["acct_id"]), policy_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST sys/v1/accounts: %v", err),
		})
		return diags
	}

	d.SetId(req["acct_id"].(string))
	return resourceReadAccountQuorumPolicy(ctx, d, m)

}

// [R]: Read Account Quorum Policy
func resourceReadAccountQuorumPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, statuscode, err := m.(*api_client).APICall("GET", fmt.Sprintf("sys/v1/accounts/%s", d.Id()))
	if statuscode == 404 {
		d.SetId("")
	} else {
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: GET sys/v1/accounts: %v", err),
			})
			return diags
		}

		if err := d.Set("acct_id", req["acct_id"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if _, ok := req["approval_policy"]; ok {
			if err := d.Set("approval_policy", fmt.Sprintf("%v", req["approval_policy"])); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return diags
}

// [U]: Update Account Quorum Policy
func resourceUpdateAccountQuorumPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

// [D]: Delete Account Quorum Policy
func resourceDeleteAccountQuorumPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	accountApprovalPolicyRead(ctx, d, m)

	account_quorum_policy_object := make(map[string]interface{})
	acct_id := d.Get("acct_id").(string)
	operation := "PATCH"
	url := "/sys/v1/approval_requests"

	account_quorum_policy_object["method"] = operation
	account_quorum_policy_object["operation"] = fmt.Sprintf("/sys/v1/accounts/%s", acct_id)
	account_quorum_policy_object["body"] = map[string]interface{}{
		"acct_id": acct_id,
		"approval_policy": map[string]interface{}{
			"policy": map[string]interface{}{},
		},
	}
	operation = "POST"

	_, err := m.(*api_client).APICallBody(operation, url, account_quorum_policy_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[D]: API Call: %s %s: %v", operation, url, err),
		})
		return diags
	}

	d.SetId("")
	return diags
}
