// **********
// Terraform Provider - DSM: data source: secret
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.1.5
//       - Date:      05/01/2021
// **********

package dsm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSecret() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecretRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
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
		},
	}
}

func dataSourceSecretRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, err := m.(*api_client).APICallList("GET", "crypto/v1/keys")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: GET crypto/v1/keys: %s", err),
		})
		return diags
	}
	//	Detail:   fmt.Sprintf("%s", req[0].(map[string]interface{})["group_id"]),

	for _, data := range req {
		// FIX-fyoo: maybe set the secret here too
		if data.(map[string]interface{})["name"].(string) == d.Get("name").(string) {
			if err := d.Set("name", data.(map[string]interface{})["name"].(string)); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("kid", data.(map[string]interface{})["kid"].(string)); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("acct_id", data.(map[string]interface{})["acct_id"].(string)); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("creator", data.(map[string]interface{})["creator"]); err != nil {
				return diag.FromErr(err)
			}
			if _, ok := data.(map[string]interface{})["description"]; ok {
				if err := d.Set("description", data.(map[string]interface{})["description"].(string)); err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	d.SetId(d.Get("kid").(string))
	return nil
}
