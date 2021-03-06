// **********
// Terraform Provider - SDKMS: resource: security object
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.5.1
//       - Date:      27/11/2020
// **********

package dsm

import (
	"context"
	"fmt"
	"time"

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
				Required: true,
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
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

// createSO: Create Security Object
func createSO(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	endpoint := "crypto/v1/keys"

	security_object := map[string]interface{}{
		"name":        d.Get("name").(string),
		"obj_type":    d.Get("obj_type").(string),
		"key_size":    d.Get("key_size").(int),
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
		security_object["deactivation_date"] = ddate.Format(layoutDSM)
	}

	if err := d.Get("key_ops").([]interface{}); len(err) > 0 {
		security_object["key_ops"] = d.Get("key_ops")
	}

	if err := d.Get("fpe_radix"); err != 0 {
		security_object["fpe"] = map[string]interface{}{
			"radix": d.Get("fpe_radix").(int),
		}
	}

	if err := d.Get("rotate").(string); len(err) > 0 {
		security_object["name"] = d.Get("rotate_from").(string)
		endpoint = "crypto/v1/keys/rekey"
	}

	req, err := m.(*api_client).APICallBody("POST", endpoint, security_object)
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
		if err := d.Set("key_size", int(req["key_size"].(float64))); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("kid", req["kid"].(string)); err != nil {
			return diag.FromErr(err)
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

	// already has been replaced so "rotate" and "rotate_from" does not apply
	_, replacement := d.GetOkExists("replacement")
	_, replaced := d.GetOkExists("replaced")
	if replacement || replaced {
		d.Set("rotate", "")
		d.Set("rotate_from", "")
	}

	if d.HasChange("key_ops") {
		security_object := map[string]interface{}{
			"kid": d.Get("kid").(string),
		}
		security_object["key_ops"] = d.Get("key_ops")

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
					Summary:  "[DSM SDK] Unable to call DSM provider API client",
					Detail:   fmt.Sprintf("[E]: API: PATCH crypto/v1/keys: Sync issue from State and DSM"),
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
