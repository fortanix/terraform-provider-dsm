// **********
// Terraform Provider - DSM: resource: aws security object
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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// [-] Define AWS Security Object
func resourceAWSSobject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateAWSSobject,
		ReadContext:   resourceReadAWSSobject,
		UpdateContext: resourceUpdateAWSSobject,
		DeleteContext: resourceDeleteAWSSobject,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"links": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"kid": {
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
			"obj_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"key_ops": {
				Type:     schema.TypeList,
				Optional: true,
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
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pending_window_in_days": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  7,
			},
			"expiry_date": {
				Type:	schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// [C]: Create AWS Security Object
func resourceCreateAWSSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	security_object := map[string]interface{}{
		"name":        d.Get("name").(string),
		"group_id":    d.Get("group_id").(string),
		"key":         d.Get("key"),
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

	if err := d.Get("custom_metadata").(map[string]interface{}); len(err) > 0 {
		security_object["custom_metadata"] = d.Get("custom_metadata")
	}

	req, err := m.(*api_client).APICallBody("POST", "crypto/v1/keys/copy", security_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST crypto/v1/keys/copy: %s", err),
		})
		return diags
	}

	d.SetId(req["kid"].(string))
	return resourceReadAWSSobject(ctx, d, m)
}

// [R]: Read AWS Security Object
func resourceReadAWSSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, statuscode, err := m.(*api_client).APICall("GET", fmt.Sprintf("crypto/v1/keys/%s", d.Id()))
	if statuscode == 404 {
		d.SetId("")
	} else {
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: GET crypto/v1/keys: %s", err),
			})
			return diags
		}

		if err := d.Set("name", req["name"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("group_id", req["group_id"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("links", req["links"]); err != nil {
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
		if err := d.Set("custom_metadata", req["custom_metadata"]); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("key_ops", req["key_ops"]); err != nil {
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
	}

	return nil
}

// [U]: Update AWS Security Object
func resourceUpdateAWSSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

// [D]: Delete AWS Security Object
func resourceDeleteAWSSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// FIXME: Since deleting, might as well remove the alias if exists
	if _, ok := d.Get("custom_metadata").(map[string]interface{})["aws-aliases"]; ok {
		remove_aws_alias := map[string]interface{}{
			"kid": d.Id(),
		}
		remove_aws_alias["custom_metadata"] = map[string]interface{}{
			"aws-aliases": "",
			"aws-policy":  d.Get("custom_metadata").(map[string]interface{})["aws-policy"],
		}
		_, err := m.(*api_client).APICallBody("PATCH", fmt.Sprintf("crypto/v1/keys/%s", d.Id()), remove_aws_alias)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: PATCH crypto/v1/keys: %s", err),
			})
			return diags
		}
	}

	// FIXME: Need to schedule deletion then delete the key - default is set to 7 days for now
	delete_object := map[string]interface{}{
		"pending_window_in_days": d.Get("pending_window_in_days").(int),
	}
	if d.Get("custom_metadata").(map[string]interface{})["aws-key-state"] != "PendingDeletion" {
		_, err := m.(*api_client).APICallBody("POST", fmt.Sprintf("crypto/v1/keys/%s/schedule_deletion", d.Id()), delete_object)
		if err != nil {
			return err
		}
	}

	_, statuscode, err := m.(*api_client).APICall("DELETE", fmt.Sprintf("crypto/v1/keys/%s", d.Id()))
	if (err != nil) && (statuscode != 404) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: DELETE crypto/v1/keys: %s", err),
		})
		return diags
	}

	d.SetId("")
	return nil
}
