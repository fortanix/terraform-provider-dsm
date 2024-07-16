// **********
// Terraform Provider - DSM: data source: secret
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.3.7
//       - Date:      28/07/2021
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
		Description: "Returns the Fortanix DSM secret object from the cluster as a Data Source.",
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The secret security object name in Fortanix DSM.",
				Type:     schema.TypeString,
				Required: true,
			},
			"kid": {
				Description: "The unique ID of the secret from Fortanix DSM.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"pub_key": {
				Description: "Public key from DSM (If applicable).",
				Type:     schema.TypeString,
				Computed: true,
			},
			"acct_id": {
				Description: "The account ID from Fortanix DSM.",
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
			"description": {
				Description: "The Fortanix DSM security object description.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"export": {
				Description: "Exports the secret based on the value shown. The value is either True/False.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"value": {
				Description: "The (sensitive) value of the secret shown if exported in base64 format.",
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceSecretRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	//req, err := m.(*api_client).APICallList("GET", "crypto/v1/keys")
	//if err != nil {
	//	diags = append(diags, diag.Diagnostic{
	//		Severity: diag.Error,
	//		Summary:  "[DSM SDK] Unable to call DSM provider API client",
	//		Detail:   fmt.Sprintf("[E]: API: GET crypto/v1/keys: %s", err),
	//	})
	//	return diags
	//}
	//	Detail:   fmt.Sprintf("%s", req[0].(map[string]interface{})["group_id"]),

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
	if d.Get("export").(bool) {
		if err := d.Set("value", req["value"].(string)); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(d.Get("kid").(string))
	return nil
}
