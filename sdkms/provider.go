// **********
// Terraform Provider - SDKMS: provider
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.1.3
//       - Date:      27/11/2020
// **********

package sdkms

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
			"sdkms_sobject":     resourceSobject(),
			"sdkms_aws_sobject": resourceAWSSobject(),
			//"sdkms_aws_group":   resourceAWSGroup(),
			"sdkms_group":       resourceGroup(),
			"sdkms_app":         resourceApp(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sdkms_aws_group": dataSourceAWSGroup(),
			"sdkms_group":     dataSourceGroup(),
			"sdkms_version":   dataSourceVersion(),
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
			Summary:  "Unable to configure SDKMS provider",
			Detail:   "[E]: SDKMS API Client Failed to Create",
		})
		return nil, diags
	}
	return newclient, nil
}
