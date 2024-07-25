// **********
// Terraform Provider - DSM: data source: cluster version
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.3.7
//       - Date:      27/11/2020
// **********

package dsm

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
		Description: "Returns the Fortanix DSM version of the cluster as a Data Source.",
		Schema: map[string]*schema.Schema{
			"version": {
				Description: "The Fortanix DSM version.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_version": {
				Description: "The Fortanix DSM API version.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_mode": {
				Description: "The Fortanix DSM execution environment.\n" +
			    "   * `SGX`: The Fortanix DSM running in Intel® SGX environment.",
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVersionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, _, err := m.(*api_client).APICall("GET", "sys/v1/version")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
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
