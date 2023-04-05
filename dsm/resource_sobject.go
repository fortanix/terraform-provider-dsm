package dsm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// [-] Define Security Object
func resourceSobject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateSobject,
		ReadContext:   resourceReadSobject,
		UpdateContext: resourceUpdateSobject,
		DeleteContext: resourceDeleteSobject,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dsm_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"obj_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"kid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"acct_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			//"kcv": {
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
			"rotate": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"DSM", "ALL"}, true),
			},
			"rotate_from": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"creator": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// Unable to define links
			//"links": {
			//	Type:     schema.TypeMap,
			//	Computed: true,
			//	Elem: &schema.Schema{
			//		Type: schema.TypeList,
			//	},
			//},
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
			"pub_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssh_pub_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"fpe_radix": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"key_ops": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"rsa": {
				Type:     schema.TypeString,
				Optional: true,
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
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"expiry_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"elliptic_curve": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// [-]: Custom Functions
// contains: Need to validate whether a string exists in a []string
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// this function takes a string in JSON format and unmarshals it. If the string is not in correct JSON format, it returns nil.
func unmarshalStringToJson(inputString string) (interface{}, error) {
	type mapFormat map[string]interface{}
	var inputMap mapFormat
	if err := json.Unmarshal([]byte(inputString), &inputMap); err != nil {
		return nil, err
	}

	return inputMap, nil
}

// createSO: Create Security Object
func createSO(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	endpoint := "crypto/v1/keys"
	key_size := d.Get("key_size").(int)
	obj_type := d.Get("obj_type").(string)
	elliptic_curve := d.Get("elliptic_curve").(string)
	method := "POST"

	security_object := map[string]interface{}{
		"name":        d.Get("name").(string),
		"obj_type":    obj_type,
		"group_id":    d.Get("group_id").(string),
		"description": d.Get("description").(string),
	}

	if _, ok := d.GetOk("value"); ok {
		security_object["value"] = d.Get("value").(string)
		method = "PUT"
	} else {
		if obj_type == "EC" {
			if key_size > 0 || len(elliptic_curve) == 0 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Detail:   fmt.Sprintf("key_size should not be specified and elliptic_curve should be specified for %s", obj_type),
				})
				return diags
			} else {
				security_object["elliptic_curve"] = elliptic_curve
			}
		} else if key_size == 0 || len(elliptic_curve) > 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Detail:   fmt.Sprintf("key_size should be specified and elliptic_curve should not be specified for %s", obj_type),
			})
			return diags
		} else {
			security_object["key_size"] = key_size
		}
	}

	if rfcdate := d.Get("expiry_date").(string); len(rfcdate) > 0 {
		layoutRFC := "2006-01-02T15:04:05Z"
		layoutDSM := "20060102T150405Z"
		ddate, newerr := time.Parse(layoutRFC, rfcdate)
		if newerr != nil {
			return diag.FromErr(newerr)
		}
		security_object["deactivation_date"] = ddate.Format(layoutDSM)
	}

	if err := d.Get("key_ops").([]interface{}); len(err) > 0 {
		security_object["key_ops"] = d.Get("key_ops")
	}
	if err := d.Get("rsa").(string); len(err) > 0 {
		rsa_obj, er := unmarshalStringToJson(d.Get("rsa").(string))
		if er != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid json string format for the field 'rsa'.",
				Detail:   fmt.Sprintf("[E]: Input: rsa: %s", err),
			})
			return diags
		}
		security_object["rsa"] = rsa_obj
	}
	allowed_key_justifications_policy, ok := d.GetOk("allowed_key_justifications_policy")
	allowed_missing_justifications, ok2 := d.GetOkExists("allowed_missing_justifications")

	if ok && ok2 {
		if allowed_key_justifications_policy != nil && allowed_missing_justifications != nil {
			security_object["google_access_reason_policy"] = map[string]interface{}{
				"allow":                allowed_key_justifications_policy,
				"allow_missing_reason": allowed_missing_justifications,
			}
		}
	} else if ok {
		if allowed_key_justifications_policy != nil {
			security_object["google_access_reason_policy"] = map[string]interface{}{
				"allow": allowed_key_justifications_policy,
			}
		}
	} else if ok2 {
		if allowed_missing_justifications != nil {
			security_object["google_access_reason_policy"] = map[string]interface{}{
				"allow_missing_reason": allowed_missing_justifications,
			}
		}
	}

	if err := d.Get("fpe_radix"); err != 0 {
		security_object["fpe"] = map[string]interface{}{
			"radix": d.Get("fpe_radix").(int),
		}
	}
	if err := d.Get("custom_metadata").(map[string]interface{}); len(err) > 0 {
		security_object["custom_metadata"] = err
	}
	if err := d.Get("rotate").(string); len(err) > 0 {
		security_object["name"] = d.Get("rotate_from").(string)
		endpoint = "crypto/v1/keys/rekey"
	}
	req, err := m.(*api_client).APICallBody(method, endpoint, security_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST %s: %v", endpoint, err),
		})
		return diags
	}

	d.SetId(req["kid"].(string))
	return diags
}

// [C]: Terraform Func: resourceCreateSobject
func resourceCreateSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if rotate := d.Get("rotate").(string); len(rotate) > 0 {
		if rotate_from := d.Get("rotate_from").(string); len(rotate_from) <= 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   "[E]: API: GET crypto/v1/keys/rekey: 'rotate_from' missing",
			})
			return diags
		}
	}

	if err := createSO(ctx, d, m); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: GET crypto/v1/keys: %v", err),
		})
		return diags
	}

	return resourceReadSobject(ctx, d, m)
}

// [R]: Terraform Func: resourceReadSobject
func resourceReadSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, statuscode, err := m.(*api_client).APICall("GET", fmt.Sprintf("crypto/v1/keys/%s", d.Id()))
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

		if err := d.Set("dsm_name", req["name"].(string)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("group_id", req["group_id"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("obj_type", req["obj_type"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if req["origin"] != "External" {
			if _, ok := req["key_size"]; ok {
				if err := d.Set("key_size", int(req["key_size"].(float64))); err != nil {
					return diag.FromErr(err)
				}
			}
			if _, ok := req["elliptic_curve"]; ok {
				if err := d.Set("elliptic_curve", req["elliptic_curve"].(string)); err != nil {
					return diag.FromErr(err)
				}
			}
		}
		if _, ok := req["google_access_reason_policy"]; ok {
			google_access_reason_policy := req["google_access_reason_policy"].(map[string]interface{})
			if err := d.Set("allowed_key_justifications_policy", google_access_reason_policy["allow"]); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("allowed_missing_justifications", google_access_reason_policy["allow_missing_reason"]); err != nil {
				return diag.FromErr(err)
			}
		}
		if err := d.Set("kid", req["kid"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if _, ok := req["pub_key"]; ok {
			if err := d.Set("pub_key", req["pub_key"].(string)); err != nil {
				return diag.FromErr(err)
			}
		}
		if err := d.Set("acct_id", req["acct_id"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("creator", req["creator"]); err != nil {
			return diag.FromErr(err)
		}
		//if err := d.Set("links", req["links"]); err != nil {
		//	return diag.FromErr(err)
		//}
		// FYOO: Fix this later - some wierd reaction to TypeList/TypeMap within TF
		if err := d.Set("copied_to", req["copied_to"]); err != nil {
			return diag.FromErr(err)
		}

		if _, ok := req["links"]; ok {
			if links := req["links"].(map[string]interface{}); len(links) > 0 {
				if _, copiedToExists := req["links"].(map[string]interface{})["copiedTo"]; copiedToExists {
					if err := d.Set("copied_to", req["links"].(map[string]interface{})["copiedTo"].([]interface{})); err != nil {
						return diag.FromErr(err)
					}
				}
				if _, copiedFromExists := req["links"].(map[string]interface{})["copiedFrom"]; copiedFromExists {
					if err := d.Set("copied_from", req["links"].(map[string]interface{})["copiedFrom"].(string)); err != nil {
						return diag.FromErr(err)
					}
				}
				if _, replacementExists := req["links"].(map[string]interface{})["replacement"]; replacementExists {
					if err := d.Set("replacement", req["links"].(map[string]interface{})["replacement"].(string)); err != nil {
						return diag.FromErr(err)
					}
				}
				if _, replacedExists := req["links"].(map[string]interface{})["replaced"]; replacedExists {
					if err := d.Set("replaced", req["links"].(map[string]interface{})["replaced"].(string)); err != nil {
						return diag.FromErr(err)
					}
				}
			}
		}
		if err := d.Set("custom_metadata", req["custom_metadata"]); err != nil {
			return diag.FromErr(err)
		}
		if err := req["fpe"]; err != nil {
			if err := d.Set("fpe_radix", int(req["fpe"].(map[string]interface{})["radix"].(float64))); err != nil {
				return diag.FromErr(err)
			}
		}
		// FYOO: Fix TypeList sorting error
		key_ops := make([]string, len(req["key_ops"].([]interface{})))
		if err := d.Get("key_ops").([]interface{}); len(err) > 0 {
			if len(d.Get("key_ops").([]interface{})) == len(req["key_ops"].([]interface{})) {
				for idx, key_op := range d.Get("key_ops").([]interface{}) {
					key_ops[idx] = fmt.Sprint(key_op)
				}
			} else {
				req_key_ops := make([]string, len(req["key_ops"].([]interface{})))
				for idx, key_op := range req["key_ops"].([]interface{}) {
					req_key_ops[idx] = fmt.Sprint(key_op)
				}
				final_idx := 0
				for _, key_op := range d.Get("key_ops").([]interface{}) {
					if contains(req_key_ops, fmt.Sprint(key_op)) {
						key_ops[final_idx] = fmt.Sprint(key_op)
						final_idx = final_idx + 1
					}
				}
			}
		} else {
			for idx, key_op := range req["key_ops"].([]interface{}) {
				key_ops[idx] = fmt.Sprint(key_op)
			}
		}
		if err := d.Set("key_ops", key_ops); err != nil {
			return diag.FromErr(err)
		}
		if _, ok := req["description"]; ok {
			if err := d.Set("description", req["description"].(string)); err != nil {
				return diag.FromErr(err)
			}
		}
		if err := d.Set("enabled", req["enabled"].(bool)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("state", req["state"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if rfcdate, ok := req["deactivation_date"]; ok {
			// FYOO: once it's set, you can't remove deactivation date
			layoutRFC := "2006-01-02T15:04:05Z"
			layoutDSM := "20060102T150405Z"
			ddate, newerr := time.Parse(layoutDSM, rfcdate.(string))
			if newerr != nil {
				return diag.FromErr(newerr)
			}
			if newerr = d.Set("expiry_date", ddate.Format(layoutRFC)); newerr != nil {
				return diag.FromErr(newerr)
			}
		}
		if err := req["obj_type"].(string); err == "RSA" {
			openssh_pub_key, err := PublicPEMtoOpenSSH([]byte(req["pub_key"].(string)))
			if err != nil {
				return err
			} else {
				if err := d.Set("ssh_pub_key", openssh_pub_key); err != nil {
					return diag.FromErr(err)
				}
			}
		}

		// FYOO: clear values that are irrelevant
		d.Set("rotate", "")
		d.Set("rotate_from", "")
	}
	return diags
}

// [U]: Terraform Func: resourceUpdateSobject
func resourceUpdateSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var has_changed = false

	// already has been replaced so "rotate" and "rotate_from" does not apply
	_, replacement := d.GetOk("replacement")
	_, replaced := d.GetOk("replaced")
	if replacement || replaced {
		d.Set("rotate", "")
		d.Set("rotate_from", "")
	}

	var security_object = map[string]interface{}{
		"kid": d.Get("kid").(string),
	}
	if d.HasChange("rsa") {
		rsa_obj, err := unmarshalStringToJson(d.Get("rsa").(string))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid json string format for the field 'rsa'.",
				Detail:   fmt.Sprintf("[E]: Input: rsa: %s", err),
			})
			return diags
		}
		security_object["rsa"] = rsa_obj
		has_changed = true
	}
	if d.HasChanges("allowed_key_justifications_policy", "allowed_missing_justifications") {
		google_access_reason_policy := make(map[string]interface{})
		allowed_key_justifications_policy, ok1 := d.GetOk("allowed_key_justifications_policy")
		allowed_missing_justifications, ok2 := d.GetOk("allowed_missing_justifications")

		if !ok1 && !ok2 {
			security_object["google_access_reason_policy"] = "remove"

		} else {
			google_access_reason_policy["allow"] = allowed_key_justifications_policy
			google_access_reason_policy["allow_missing_reason"] = allowed_missing_justifications
			security_object["google_access_reason_policy"] = google_access_reason_policy
		}

		has_changed = true
	}
	if d.HasChange("key_ops") {
		security_object["key_ops"] = d.Get("key_ops")
		has_changed = true
	}
	if d.HasChange("description") {
		security_object["description"] = d.Get("description")
		has_changed = true
	}
	if d.HasChange("name") {
		security_object["name"] = d.Get("name")
		has_changed = true
	}
	if d.HasChange("custom_metadata") {
		security_object["custom_metadata"] = d.Get("custom_metadata").(map[string]interface{})
		has_changed = true
	}
	if d.HasChange("elliptic_curve") {
		old_ec, new_ec := d.GetChange("elliptic_curve")
		d.Set("elliptic_curve", old_ec)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "elliptic_curve cannot modify on update",
			Detail:   fmt.Sprintf("[E]: API: PATCH crypto/v1/keys: elliptic_curve cannot change on update. Please retain it to old value: %s -> %s", old_ec, new_ec),
		})
		return diags
	}

	if has_changed {
		if debug_output {
			tflog.Warn(ctx, "Sobject has changed, calling API")
		}
		req, err := m.(*api_client).APICallBody("PATCH", fmt.Sprintf("crypto/v1/keys/%s", d.Id()), security_object)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: PATCH crypto/v1/keys: %v", err),
			})
			return diags
		}

		key_ops := make([]string, len(req["key_ops"].([]interface{})))
		if err := d.Get("key_ops").([]interface{}); len(err) > 0 {
			if len(d.Get("key_ops").([]interface{})) == len(req["key_ops"].([]interface{})) {
				for idx, key_op := range d.Get("key_ops").([]interface{}) {
					key_ops[idx] = fmt.Sprint(key_op)
				}
			} else {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "[DSM SDK] Error in processing DSM provider API response key_ops",
					Detail:   "[E]: API: PATCH crypto/v1/keys: Sync issue from State and DSM",
				})
				return diags
			}
		} else {
			for idx, key_op := range req["key_ops"].([]interface{}) {
				key_ops[idx] = fmt.Sprint(key_op)
			}
		}
		if err := d.Set("key_ops", key_ops); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceReadSobject(ctx, d, m)
}

// [D]: Terraform Func: resourceDeleteSobject
func resourceDeleteSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
