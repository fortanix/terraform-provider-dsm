// **********
// Terraform Provider - DSM: data source: group
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.3.7
//       - Date:      05/01/2021
// **********

package dsm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
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
			"approval_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cryptographic_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hmg": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, err := m.(*api_client).APICallList("GET", "sys/v1/groups")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: GET sys/v1/groups: %v", err),
		})
		return diags
	}

	group_id := ""
	for _, data := range req {
		if data.(map[string]interface{})["name"].(string) == d.Get("name").(string) {
			group_id = data.(map[string]interface{})["group_id"].(string)
			if err := d.Set("name", data.(map[string]interface{})["name"].(string)); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("group_id", group_id); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("acct_id", data.(map[string]interface{})["acct_id"].(string)); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("creator", data.(map[string]interface{})["creator"]); err != nil {
				return diag.FromErr(err)
			}
			if description, ok := data.(map[string]interface{})["description"]; ok {
				if err := d.Set("description", description.(string)); err != nil {
					return diag.FromErr(err)
				}
			}
			if approval_policy, ok := data.(map[string]interface{})["approval_policy"]; ok {
				jsonString, err := json.Marshal(approval_policy)
				if err == nil {
					d.Set("approval_policy", jsonString)
				} else {
					return diag.FromErr(err)
				}
			}
			if cryptographic_policy, ok := data.(map[string]interface{})["cryptographic_policy"]; ok {
				jsonString, err := json.Marshal(cryptographic_policy)
				if err == nil {
					d.Set("cryptographic_policy", jsonString)
				} else {
					return diag.FromErr(err)
				}
			}
			if hmg, ok := data.(map[string]interface{})["hmg"]; ok {
				jsonString, err := json.Marshal(hmg)
				if err == nil {
					d.Set("hmg", jsonString)
				} else {
					return diag.FromErr(err)
				}
			}
		}
	}

	d.SetId(group_id)
	return nil
}
