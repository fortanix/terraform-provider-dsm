// **********
// Terraform Provider - SDKMS: resource: secret
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.3.7
//       - Date:      27/11/2020
// **********

package dsm

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// [-] Define Security Object
func resourceSecret() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateSecret,
		ReadContext:   resourceReadSecret,
		UpdateContext: resourceUpdateSecret,
		DeleteContext: resourceDeleteSecret,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"obj_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pub_key": {
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
			"custom_metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"key_ops": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allowed_key_justifications_policy": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"allowed_missing_justifications": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"value": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expiry_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rotate": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"rotate_from": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"copied_to": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"copied_from": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replacement": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replaced": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// [C]: Create Security Object
func resourceCreateSecret(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	endpoint := "crypto/v1/keys"
	operation := "PUT"

	plugin_object := map[string]interface{}{
		"operation":   "create",
		"name":        d.Get("name").(string),
		"group_id":    d.Get("group_id").(string),
		"description": d.Get("description").(string),
	}

	if rfcdate := d.Get("expiry_date").(string); len(rfcdate) > 0 {
		layoutRFC := "2006-01-02T15:04:05Z"
		layoutDSM := "20060102T150405Z"
		ddate, newerr := time.Parse(layoutRFC, rfcdate)
		if newerr != nil {
			return diag.FromErr(newerr)
		}
		plugin_object["deactivation_date"] = ddate.Format(layoutDSM)
	}

	if d.Get("rotate").(bool) {
		plugin_object["operation"] = "rotate"
		plugin_object["name"] = d.Get("rotate_from").(string)
		endpoint = "crypto/v1/keys/rekey"
		operation = "POST"
	}

	if err := d.Get("value").(string); len(err) > 0 {
		plugin_object["value"] = d.Get("value").(string)
		plugin_object["obj_type"] = "SECRET"
	} else {
		reqfpi, err := m.(*api_client).FindPluginId("Terraform Plugin")
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: GET sys/v1/plugins: %v", err),
			})
			return diags
		}
		endpoint = fmt.Sprintf("sys/v1/plugins/%s", string(reqfpi))
		operation = "POST"
	}

	allowed_key_justifications_policy, ok := d.GetOk("allowed_key_justifications_policy")
	allowed_missing_justifications, ok2 := d.GetOkExists("allowed_missing_justifications")

	if ok && ok2 {
		if allowed_key_justifications_policy != nil && allowed_missing_justifications != nil {
			plugin_object["google_access_reason_policy"] = map[string]interface{}{
				"allow":                allowed_key_justifications_policy,
				"allow_missing_reason": allowed_missing_justifications,
			}
		}
	} else if ok {
		if allowed_key_justifications_policy != nil {
			plugin_object["google_access_reason_policy"] = map[string]interface{}{
				"allow": allowed_key_justifications_policy,
			}
		}
	} else if ok2 {
		if allowed_missing_justifications != nil {
			plugin_object["google_access_reason_policy"] = map[string]interface{}{
				"allow_missing_reason": allowed_missing_justifications,
			}
		}
	}

	req, err := m.(*api_client).APICallBody(operation, endpoint, plugin_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST sys/v1/plugins: %v", err),
		})
		return diags
	}

	d.SetId(req["kid"].(string))
	return resourceReadSecret(ctx, d, m)
}

// [R]: Read Security Object
func resourceReadSecret(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	res, statuscode, err := m.(*api_client).APICall("GET", fmt.Sprintf("crypto/v1/keys/%s", d.Id()))
	if statuscode == 404 {
		d.SetId("")
	} else {
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: GET crypto/v1/keys: %v", err),
			})
			return diags
		}

		if err := d.Set("name", res["name"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("group_id", res["group_id"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("obj_type", res["obj_type"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("kid", res["kid"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if _, ok := res["pub_key"]; ok {
			if err := d.Set("pub_key", res["pub_key"].(string)); err != nil {
				return diag.FromErr(err)
			}
		}
		if err := d.Set("acct_id", res["acct_id"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("creator", res["creator"]); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("custom_metadata", res["custom_metadata"]); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("key_ops", res["key_ops"]); err != nil {
			return diag.FromErr(err)
		}
		if _, ok := res["description"]; ok {
			if err := d.Set("description", res["description"].(string)); err != nil {
				return diag.FromErr(err)
			}
		}
		if _, ok := res["google_access_reason_policy"]; ok {
			google_access_reason_policy := res["google_access_reason_policy"].(map[string]interface{})
			if err := d.Set("allowed_key_justifications_policy", google_access_reason_policy["allow"]); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("allowed_missing_justifications", google_access_reason_policy["allow_missing_reason"]); err != nil {
				return diag.FromErr(err)
			}
		}
		if err := d.Set("enabled", res["enabled"].(bool)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("state", res["state"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if rfcdate, ok := res["deactivation_date"].(string); ok {
			// FYOO: once it's set, you can't remove deactivation date
			layoutRFC := "2006-01-02T15:04:05Z"
			layoutDSM := "20060102T150405Z"
			ddate, newerr := time.Parse(layoutDSM, rfcdate)
			if newerr != nil {
				return diag.FromErr(newerr)
			}
			if newerr = d.Set("expiry_date", ddate.Format(layoutRFC)); newerr != nil {
				return diag.FromErr(newerr)
			}
		}
		if _, ok := res["links"]; ok {
			if links := res["links"].(map[string]interface{}); len(links) > 0 {
				if _, copiedToExists := res["links"].(map[string]interface{})["copiedTo"]; copiedToExists {
					if err := d.Set("copied_to", res["links"].(map[string]interface{})["copiedTo"].([]interface{})); err != nil {
						return diag.FromErr(err)
					}
				}
				if _, copiedFromExists := res["links"].(map[string]interface{})["copiedFrom"]; copiedFromExists {
					if err := d.Set("copied_from", res["links"].(map[string]interface{})["copiedFrom"].(string)); err != nil {
						return diag.FromErr(err)
					}
				}
				if _, replacementExists := res["links"].(map[string]interface{})["replacement"]; replacementExists {
					if err := d.Set("replacement", res["links"].(map[string]interface{})["replacement"].(string)); err != nil {
						return diag.FromErr(err)
					}
				}
				if _, replacedExists := res["links"].(map[string]interface{})["replaced"]; replacedExists {
					if err := d.Set("replaced", res["links"].(map[string]interface{})["replaced"].(string)); err != nil {
						return diag.FromErr(err)
					}
				}
			}
		}
	}
	return diags
}

// [U]: Update Security Object
func resourceUpdateSecret(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var has_changed = false
	var diags diag.Diagnostics

	var plugin_object = map[string]interface{}{
		"kid": d.Get("kid").(string),
	}

	if d.HasChanges("allowed_key_justifications_policy", "allowed_missing_justifications") {

		google_access_reason_policy := make(map[string]interface{})

		google_access_reason_policy["allow"] = d.Get("allowed_key_justifications_policy")
		google_access_reason_policy["allow_missing_reason"] = d.Get("allowed_missing_justifications")

		has_changed = true

		plugin_object["google_access_reason_policy"] = google_access_reason_policy
	}

	if has_changed {
		if debug_output {
			tflog.Warn(ctx, "Sobject has changed, calling API")
		}
		_, err := m.(*api_client).APICallBody("PATCH", fmt.Sprintf("crypto/v1/keys/%s", d.Id()), plugin_object)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: PATCH crypto/v1/keys: %v", err),
			})
			return diags
		}
	}

	// return nil
	return resourceReadSecret(ctx, d, m)
}

// [D]: Delete Security Object
func resourceDeleteSecret(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	_, statuscode, err := m.(*api_client).APICall("DELETE", fmt.Sprintf("crypto/v1/keys/%s", d.Id()))
	if (err != nil) && (statuscode != 404) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: DELETE crypto/v1/keys: %v", err),
		})
		return diags
	}

	d.SetId("")
	return nil
}
