// **********
// Terraform Provider - DSM: resource: gcp security object
// **********
//       - Author:    shashidhar naraparaju at fortanix dot com
//       - Version:   0.5.29
//       - Date:      05/04/2024
// **********

package dsm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// [-] Define GCP Security Object in Terraform
func resourceGCPSobject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateGCPSobject,
		ReadContext:   resourceReadGCPSobject,
		UpdateContext: resourceUpdateGCPSobject,
		DeleteContext: resourceDeleteGCPSobject,
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
			"rotation_policy": {
				Type:     schema.TypeMap,
				Optional: true,
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
			"external": {
				Type:     schema.TypeMap,
				Computed: true,
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
				Default: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Default: "Active",
			},
			"expiry_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// [C]: Create GCP Security Object
func resourceCreateGCPSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	if enabled, ok := d.Get("enabled").(bool); ok {
		security_object["enabled"] = enabled
	}
	if err := d.Get("key_ops").([]interface{}); len(err) > 0 {
		security_object["key_ops"] = d.Get("key_ops")
	}
	if err := d.Get("custom_metadata").(map[string]interface{}); len(err) > 0 {
		security_object["custom_metadata"] = d.Get("custom_metadata")
	}
	if rotation_policy := d.Get("rotation_policy").(map[string]interface{}); len(rotation_policy) > 0 {
		security_object["rotation_policy"] = sobj_rotation_policy_write(rotation_policy)
	}
	req, err := m.(*api_client).APICallBody("POST", "crypto/v1/keys/copy", security_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST crypto/v1/keys/copy: %v", err),
		})
		return diags
	}
	d.SetId(req["kid"].(string))
	return resourceReadGCPSobject(ctx, d, m)
}

// [R]: Read GCP Security Object
func resourceReadGCPSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		// Convert returned call into gcp Map
		jsonbody, err := json.Marshal(req)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to parse DSM provider API client output",
				Detail:   fmt.Sprintf("[E]: API: GET crypto/v1/keys: %s", err),
			})
			return diags
		}
		gcpsobject := GCPSobject{}
		if err := json.Unmarshal(jsonbody, &gcpsobject); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to parse DSM provider API client output",
				Detail:   fmt.Sprintf("[E]: API: GET crypto/v1/keys: %s", err),
			})
			return diags
		}
		// Sync DSM and Terraform attributes
		if err := d.Set("name", gcpsobject.Name); err != nil {
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
		if _, ok := req["rotation_policy"]; ok {
			rotation_policy := sobj_rotation_policy_read(req["rotation_policy"].(map[string]interface{}))
			if err := d.Set("rotation_policy", rotation_policy); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return nil
}

// [U]: Update GCP Security Object
func resourceUpdateGCPSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	update_gcp_key := make(map[string]interface{})

	if d.HasChange("name") {
		update_gcp_key["name"] = d.Get("name").(string)
	}
	if d.HasChange("description") {
		update_gcp_key["description"] = d.Get("description").(string)
	}
	if d.HasChange("enabled") {
		update_gcp_key["enabled"] = d.Get("enabled").(bool)
	}
	if d.HasChange("key_ops") {
		update_gcp_key["key_ops"] = d.Get("key_ops")
	}
	if d.HasChange("rotation_policy") {
		if err := d.Get("rotation_policy").(map[string]interface{}); len(err) > 0 {
			update_gcp_key["rotation_policy"] = sobj_rotation_policy_write(d.Get("rotation_policy").(map[string]interface{}))
		}
	}
	if len(update_gcp_key) > 0 {
		_, err := m.(*api_client).APICallBody("PATCH", fmt.Sprintf("crypto/v1/keys/%s", d.Id()), update_gcp_key)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: PATCH crypto/v1/keys: %v", err),
			})
			return diags
		}
	}

	return resourceReadGCPSobject(ctx, d, m)
}

// [D]: Delete GCP Security Object
func resourceDeleteGCPSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Blocker: https://fortanix.atlassian.net/browse/ROFR-4819
	// Backend implementation for deleting GCP sobjects is pending
	return nil
}
