// **********
// Terraform Provider - SDKMS: resource: secret
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.1.5
//       - Date:      27/11/2020
// **********

package dsm

import (
	"context"
	"encoding/base64"
	"fmt"

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
			"value": {
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
func resourceCreateSecret(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	security_object := map[string]interface{}{
		"name":        d.Get("name").(string),
		"obj_type":    "SECRET",
		"value":       base64.StdEncoding.EncodeToString([]byte(d.Get("value").(string))),
		"group_id":    d.Get("group_id").(string),
		"description": d.Get("description").(string),
	}

	if err := d.Get("key_ops").([]interface{}); len(err) > 0 {
		security_object["key_ops"] = d.Get("key_ops")
	}

	req, err := m.(*api_client).APICallBody("PUT", "crypto/v1/keys", security_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: PUT crypto/v1/keys: %s", err),
		})
		return diags
	}

	d.SetId(req["kid"].(string))
	return resourceReadSecret(ctx, d, m)
}

// [R]: Read Security Object
func resourceReadSecret(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, err := m.(*api_client).APICall("GET", fmt.Sprintf("crypto/v1/keys/%s", d.Id()))
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
	if err := d.Set("obj_type", req["obj_type"].(string)); err != nil {
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

	return diags
}

// [U]: Update Security Object
func resourceUpdateSecret(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

// [D]: Delete Security Object
func resourceDeleteSecret(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	_, err := m.(*api_client).APICall("DELETE", fmt.Sprintf("crypto/v1/keys/%s", d.Id()))
	if err != nil {
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
