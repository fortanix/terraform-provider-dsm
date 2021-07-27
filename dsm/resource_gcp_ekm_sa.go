// **********
// Terraform Provider - DSM: resource: gcp_ekm_sa
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.2.0
//       - Date:      27/07/2021
//       - Changelog:
//                  - Initial release to support resource for app with GCP EKM SA
// **********

package dsm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// [-] Define App
func resourceGcpEkmSa() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateGcpEkmSa,
		ReadContext:   resourceReadGcpEkmSa,
		UpdateContext: resourceUpdateGcpEkmSa,
		DeleteContext: resourceDeleteGcpEkmSa,
		Schema: map[string]*schema.Schema{
			// service account name = app name
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_group": {
				Type:     schema.TypeString,
				Required: true,
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
				Default:  "",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// [C]: Create GcpEkmSa
func resourceCreateGcpEkmSa(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	app_object := map[string]interface{}{
		"name":          d.Get("name").(string),
		"default_group": d.Get("default_group").(string),
		"add_groups": map[string]interface{}{
			d.Get("default_group").(string): []string{"SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "EXPORT", "MANAGE", "AGREEKEY", "AUDIT"},
		},
		"app_type":    "default",
		"description": d.Get("description").(string),
	}

	// add GCP specifics

	gcp_perm := map[string]interface{}{
		"allow":                []string{"CUSTOMER_INITIATED_SUPPORT", "CUSTOMER_INITIATED_ACCESS", "GOOGLE_INITIATED_SERVICE", "GOOGLE_INITIATED_REVIEW", "GOOGLE_INITIATED_SYSTEM_OPERATION", "THIRD_PARTY_DATA_REQUEST", "REASON_UNSPECIFIED", "REASON_NOT_EXPECTED", "MODIFIED_CUSTOMER_INITIATED_ACCESS"},
		"allow_missing_reason": true,
	}

	app_object["credential"] = map[string]interface{}{
		"googleserviceaccount": gcp_perm,
	}

	req, err := m.(*api_client).APICallBody("POST", "sys/v1/apps", app_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST sys/v1/apps: %s", err),
		})
		return diags
	}

	d.SetId(req["app_id"].(string))
	return resourceReadApp(ctx, d, m)
}

// [R]: Read GcpEkmSa
func resourceReadGcpEkmSa(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, err := m.(*api_client).APICall("GET", fmt.Sprintf("sys/v1/apps/%s", d.Id()))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: GET sys/v1/apps: %s", err),
		})
		return diags
	}

	if err := d.Set("name", req["name"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("app_id", req["app_id"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("default_group", req["default_group"].(string)); err != nil {
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

// [U]: Update GcpEkmSa
func resourceUpdateGcpEkmSa(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

// [D]: Delete GcpEkmSa
func resourceDeleteGcpEkmSa(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	_, err := m.(*api_client).APICall("DELETE", fmt.Sprintf("sys/v1/apps/%s", d.Id()))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: DELETE sys/v1/apps: %s", err),
		})
		return diags
	}

	d.SetId("")
	return nil
}
