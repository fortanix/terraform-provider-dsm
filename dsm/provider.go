// **********
// Terraform Provider - DSM: provider
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.2.0
//       - Date:      27/07/2021
// **********

package dsm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// [-] Define Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"acct_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"dsm_sobject":     resourceSobject(),
			"dsm_aws_sobject": resourceAWSSobject(),
			//"sdkms_aws_group":   resourceAWSGroup(),
			"dsm_secret":     resourceSecret(),
			"dsm_group":      resourceGroup(),
			"dsm_app":        resourceApp(),
			"dsm_gcp_ekm_sa": resourceGcpEkmSa(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"dsm_aws_group": dataSourceAWSGroup(),
			"dsm_secret":    dataSourceSecret(),
			"dsm_group":     dataSourceGroup(),
			"dsm_version":   dataSourceVersion(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	newclient, err := NewAPIClient(d.Get("endpoint").(string), d.Get("port").(int), d.Get("username").(string), d.Get("password").(string), d.Get("acct_id").(string))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK]: Unable to configure DSM provider",
			Detail:   "[E]: API: Failed to create client",
		})
		return nil, diags
	}
	return newclient, nil
}
