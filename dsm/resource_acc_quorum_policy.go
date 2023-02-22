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
		Schema: map[string]*schema.Schema{
			"acct_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"approval_policy": {
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

	jsonString, errjson := json.Marshal(d.Get("approval_policy"))
	if errjson == nil {
		return diag.FromErr(errjson)
	}

	policy_object := map[string]interface{}{
		"acct_id":         d.Get("acct_id").(string),
		"approval_policy": string(jsonString),
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
			jsonString, err := json.Marshal(req["approval_policy"])
			if err == nil {
				d.Set("approval_policy", jsonString)
			} else {
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
	return nil
}
