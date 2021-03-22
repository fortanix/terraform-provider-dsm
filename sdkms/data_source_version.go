// **********
// Terraform Provider - SDKMS: data source: cluster version
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.1.2
//       - Date:      27/11/2020
// **********

package sdkms

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVersion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVersionRead,
		Schema: map[string]*schema.Schema{
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_mode": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVersionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, err := m.(*api_client).APICall("GET", "sys/v1/version")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to call SDKMS provider API client",
			Detail:   "[E]: API: sys/v1/version",
		})
		return diags
	}

	if err := d.Set("version", req["version"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("api_version", req["api_version"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("server_mode", req["server_mode"]); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return nil
}
