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
		Description: "Imports a security object of type Secret. The returned resource object contains the UUID of the security object for further references.\n" +
		"A secret value format should be in a base64 format. Secret can also be rotated.",
		Schema: map[string]*schema.Schema{
			"name": {
			    Description: "The Fortanix DSM secret security object name",
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
			    Description: "The Fortanix DSM security object group assignment.",
				Type:     schema.TypeString,
				Required: true,
			},
			"obj_type": {
			    Description: "The security object key type from Fortanix DSM.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"kid": {
				Description: "Security object ID from Fortanix DSM.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"acct_id": {
			    Description: "Account ID from Fortanix DSM.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator": {
				Description: "The creator of the security object from Fortanix DSM.\n" +
				"   * `user`: If the security object was created by a user, the computed value will be the matching user id.\n" +
				"   * `app`: If the security object was created by a app, the computed value will be the matching app id.",
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"custom_metadata": {
			    Description: "The user defined security object attributes added to the key’s metadata.",
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"key_ops": {
				Description: "The security object key permission from Fortanix DSM.\n" +
				"   * Default is to allow all permissions.",
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": {
			    Description: "The Fortanix DSM security object description.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"enabled": {
			    Description: "Whether the security object is Enabled or Disabled. The values are true/false.",
				Type:     schema.TypeBool,
				Optional: true,
			},
			"value": {
			    Description: "The value of the secret security object Base64 encoded.",
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"state": {
				Description: "The state of the secret security object.\n" +
				"   * Allowed states are: None, PreActive, Active, Deactivated, Compromised, Destroyed, Deleted.",
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"expiry_date": {
			    Description: " The security object expiry date in RFC format.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"rotate": {
			    Description: "boolean value true/false to enable/disable rotation.",
				Type:     schema.TypeBool,
				Optional: true,
			},
			"rotate_from": {
			    Description: "Name of the security object to be rotated from.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"copied_to": {
			    Description: "List of security objects copied by the current security object.",
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"copied_from": {
			    Description: "Security object that is copied to the current security object.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"replacement": {
			    Description: "Replacement of a security object.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"replaced": {
			    Description: "Replaced by a security object.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"allowed_key_justifications_policy": {
			    Description: "The security object key justification policies for GCP External Key Manager. The allowed permissions are:\n" +
			    "   * CUSTOMER_INITIATED_SUPPORT\n" +
			    "   * CUSTOMER_INITIATED_ACCESS\n" +
			    "   * GOOGLE_INITIATED_SERVICE\n" +
			    "   * GOOGLE_INITIATED_REVIEW\n" +
			    "   * GOOGLE_INITIATED_SYSTEM_OPERATION\n" +
			    "   * THIRD_PARTY_DATA_REQUEST\n" +
			    "   * REASON_NOT_EXPECTED\n" +
			    "   * REASON_UNSPECIFIED\n" +
			    "   * MODIFIED_CUSTOMER_INITIATED_ACCESS\n" +
			    "   * MODIFIED_GOOGLE_INITIATED_SYSTEM_OPERATION\n" +
			    "   * GOOGLE_RESPONSE_TO_PRODUCTION_ALERT\n",
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"allowed_missing_justifications": {
			    Description: "Boolean value which allows missing justifications even if not provided to the secret. The values are True / False.",
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// [C]: Create Security Object
func resourceCreateSecret(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	endpoint := "crypto/v1/keys"
	operation := "PUT"

	plugin_object := map[string]interface{}{
		"operation":   "create",
		"name":        d.Get("name").(string),
		"group_id":    d.Get("group_id").(string),
		"description": d.Get("description").(string),
	}

	if err := d.Get("expiry_date").(string); len(err) > 0 {
		sobj_deactivation_date, date_error := parseTimeToDSM(err)
		if date_error != nil {
			return date_error
		}
		plugin_object["deactivation_date"] = sobj_deactivation_date
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
			return invokeErrorDiagsWithSummary("[DSM SDK] Unable to call DSM provider API client", fmt.Sprintf("[E]: API: GET sys/v1/plugins: %v", err))
		}
		endpoint = fmt.Sprintf("sys/v1/plugins/%s", string(reqfpi))
		operation = "POST"
	}
	allowed_key_justifications_policy, allow_exists := d.GetOk("allowed_key_justifications_policy")
	allowed_missing_justifications, allow_missing_justifications_exists := d.GetOk("allowed_missing_justifications")
	
	policy_data := map[string]interface{}{}
	
	if allow_exists && allowed_key_justifications_policy != nil {
		policy_data["allow"] = allowed_key_justifications_policy
	}
	
	if allow_missing_justifications_exists && allowed_missing_justifications != nil {
		policy_data["allow_missing_reason"] = allowed_missing_justifications
	}
	
	// Only update if there are entries
	if len(policy_data) > 0 {
		plugin_object["google_access_reason_policy" ] = policy_data
	}

	req, err := m.(*api_client).APICallBody(operation, endpoint, plugin_object)
	if err != nil {
		return invokeErrorDiagsWithSummary("[DSM SDK] Unable to call DSM provider API client", fmt.Sprintf("[E]: API: POST sys/v1/plugins: %v", err))
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
			return invokeErrorDiagsWithSummary("[DSM SDK] Unable to call DSM provider API client", fmt.Sprintf("[E]: API: GET crypto/v1/keys: %v", err))
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
			tf_state_garp, is_tf_state_garp  := d.GetOk("allowed_key_justifications_policy")
			var is_same_garp bool
			if is_tf_state_garp {
				is_same_garp = compTwoArrays(tf_state_garp, google_access_reason_policy["allow"])
			}
			if is_same_garp {
				if err := d.Set("allowed_key_justifications_policy", tf_state_garp); err != nil {
					return diag.FromErr(err)
				}
			} else if err := d.Set("allowed_key_justifications_policy", google_access_reason_policy["allow"]); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("allowed_missing_justifications", google_access_reason_policy["allow_missing_reason"]); err != nil {
				return diag.FromErr(err)
			}
		} else{
		    /*
		        allowed_key_justifications_policy is either Optional or Computed.
		        It is being made as Computed, because when a key is copied, KAJ will also get copied.
		        In this case, it will become a computed value.

		        If allowed_key_justifications_policy is not set, while updating it shows a difference as it will set to null value.
		        Hence, it needs to be set as an empty value.
		    */
		    empty_array := []string{}
		    d.Set("allowed_key_justifications_policy", empty_array)
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
		if err := d.Set("copied_to", res["copied_to"]); err != nil {
			return diag.FromErr(err)
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

	if d.HasChange("name") {
		plugin_object["name"] = d.Get("name").(string)
		has_changed = true
	}

	if d.HasChange("description") {
		plugin_object["description"] = d.Get("description").(string)
		has_changed = true
	}

	if d.HasChange("custom_metadata") {
		plugin_object["custom_metadata"] = d.Get("custom_metadata").(map[string]interface{})
		has_changed = true
	}

	if d.HasChange("enabled") {
		plugin_object["enabled"] = d.Get("enabled").(bool)
		has_changed = true
	}

	if d.HasChanges("allowed_key_justifications_policy", "allowed_missing_justifications") {
		google_access_reason_policy := make(map[string]interface{})

		if allowed_justifications := d.Get("allowed_key_justifications_policy"); allowed_justifications != nil {
			google_access_reason_policy["allow"] = allowed_justifications
		}
		if allow_missing := d.Get("allowed_missing_justifications"); allow_missing != nil {
			google_access_reason_policy["allow_missing_reason"] = allow_missing
		}

		plugin_object["google_access_reason_policy"] = google_access_reason_policy
		has_changed = true
	}

	// Expiry date cannot be modified if it is already set.
	if d.HasChange("expiry_date") {
		old_expiry_date, new_expiry_date := d.GetChange("expiry_date")
		if old_expiry_date == nil || len(old_expiry_date.(string)) == 0 {
			sobj_deactivation_date, date_error := parseTimeToDSM(d.Get("expiry_date").(string))
			if date_error != nil {
				return date_error
			}
			plugin_object["deactivation_date"] = sobj_deactivation_date
		} else {
			d.Set("expiry_date", old_expiry_date)
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Detail:   fmt.Sprintf("[E]: API: PATCH crypto/v1/keys: expiry_date cannot be changed once it is set. Please retain it to old value: %s -> %s", new_expiry_date, old_expiry_date),
			})
			return diags
		}
		has_changed = true
	}

	if has_changed {
		if debug_output {
			tflog.Warn(ctx, "Secret has changed, calling API.")
		}
		_, err := m.(*api_client).APICallBody("PATCH", fmt.Sprintf("crypto/v1/keys/%s", d.Id()), plugin_object)
		if err != nil {
			return invokeErrorDiagsWithSummary("[DSM SDK] Unable to call DSM provider API client", fmt.Sprintf("[E]: API: PATCH crypto/v1/keys: %v", err))
		}
	}

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
