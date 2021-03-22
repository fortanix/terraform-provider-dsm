// **********
// Terraform Provider - SDKMS: resource: security object
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.1.2
//       - Date:      27/11/2020
// **********

package sdkms

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"kcv": {
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
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// [C]: Create Security Object
func resourceCreateSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	security_object := map[string]interface{}{
		"name":        d.Get("name").(string),
		"obj_type":    d.Get("obj_type").(string),
		"key_size":    d.Get("key_size").(int),
		"group_id":    d.Get("group_id").(string),
		"description": d.Get("description").(string),
	}

	if err := d.Get("key_ops").([]interface{}); len(err) > 0 {
		security_object["key_ops"] = d.Get("key_ops")
	}

	req, err := m.(*api_client).APICallBody("POST", "crypto/v1/keys", security_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to call SDKMS provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST crypto/v1/keys: %s", err),
		})
		return diags
	}

	d.SetId(req["kid"].(string))
	return resourceReadSobject(ctx, d, m)
}

// [R]: Read Security Object
func resourceReadSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, err := m.(*api_client).APICall("GET", fmt.Sprintf("crypto/v1/keys/%s", d.Id()))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to call SDKMS provider API client",
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
	if err := d.Set("kcv", req["kcv"].(string)); err != nil {
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

	return diags
}

// [U]: Update Security Object
func resourceUpdateSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

// [D]: Delete Security Object
func resourceDeleteSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	_, err := m.(*api_client).APICall("DELETE", fmt.Sprintf("crypto/v1/keys/%s", d.Id()))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to call SDKMS provider API client",
			Detail:   fmt.Sprintf("[E]: API: DELETE crypto/v1/keyss: %s", err),
		})
		return diags
	}

	d.SetId("")
	return nil
}
