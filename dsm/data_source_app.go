// **********
// Terraform Provider - DSM: data source: app
// **********
//       - Author:    aman.ahuja@fortanix.com
//       - Version:   0.5.16
//       - Date:      28/07/2022
// **********

package dsm

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceApp() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppRead,
		Description: "Returns the Fortanix DSM app object from the cluster as a Data Source.",
		Schema: map[string]*schema.Schema{
			"app_id": {
				Description: "App id value.",
				Type:     schema.TypeString,
				Required: true,
			},
			"credential": {
				Description: "The Fortanix DSM App API key.",
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"new_credential": {
				Description: "Set this if you want to rotate/regenerate the API key. The values can be set as True/False.",
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func dataSourceAppRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId(d.Get("app_id").(string))

	req, _, err := m.(*api_client).APICall("GET", fmt.Sprintf("sys/v1/apps/%s/credential", d.Id()))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: GET sys/v1/apps/-/credential: %v", err),
		})
		return diags
	}

	if err := d.Set("credential", base64.StdEncoding.EncodeToString([]byte(d.Id()+":"+req["credential"].(map[string]interface{})["secret"].(string)))); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("new_credential", false); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
