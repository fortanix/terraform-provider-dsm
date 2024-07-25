package dsm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRoleRead,
		Description: "Returns the Fortanix DSM role object from the cluster as a Data Source.",
		Schema: map[string]*schema.Schema{
			"role_id": {
				Description: "Role object ID from Fortanix DSM.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Description: "Role name in Fortanix DSM.",
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, err := m.(*api_client).APICallList("GET", "sys/v1/roles")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: GET sys/v1/roles: %v", err),
		})
		return diags
	}

	role_id := ""
	for _, data := range req {
		if data.(map[string]interface{})["name"].(string) == d.Get("name").(string) {
			role_id = data.(map[string]interface{})["role_id"].(string)
			if err := d.Set("name", d.Get("name").(string)); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("role_id", role_id); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(role_id)
	return nil

}
