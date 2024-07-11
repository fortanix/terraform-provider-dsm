// **********
// Terraform Provider - DSM: data source: secret
// **********
//       - Author:    sanjeev at fortanix dot com
//       - Version:   0.3.7
//       - Date:     07/07/2022
// **********

package dsm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSobject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSobjectRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
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
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"export": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"value": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"key_ops": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"key_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"obj_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceSobjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	security_object := map[string]interface{}{
		"name": d.Get("name").(string),
	}

	req, err := m.(*api_client).APICallBody("POST", "crypto/v1/keys/export", security_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST crypto/v1/keys/export: %v", err),
		})
		return diags
	}

	if err := d.Set("name", req["name"].(string)); err != nil {
		return diag.FromErr(err)
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
	if _, ok := req["description"]; ok {
		if err := d.Set("description", req["description"].(string)); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("obj_type", req["obj_type"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("key_size", int(req["key_size"].(float64))); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enabled", req["enabled"].(bool)); err != nil {
		return diag.FromErr(err)
	}
	if d.Get("export").(bool) {
		if err := d.Set("value", req["value"].(string)); err != nil {
			return diag.FromErr(err)
		}
	}
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
	d.SetId(d.Get("kid").(string))
	return nil
}
