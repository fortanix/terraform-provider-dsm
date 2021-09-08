// **********
// Terraform Provider - DSM: provider
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.3.7
//       - Date:      27/07/2021
// **********

package dsm

import (
	"context"
	"fmt"

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
			"insecure": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
			"aws_profile": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"aws_region": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "us-east-1",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"dsm_sobject":     resourceSobject(),
			"dsm_aws_sobject": resourceAWSSobject(),
			"dsm_aws_group":   resourceAWSGroup(),
			"dsm_secret":      resourceSecret(),
			"dsm_group":       resourceGroup(),
			"dsm_app":         resourceApp(),
			"dsm_gcp_ekm_sa":  resourceGcpEkmSa(),
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

	// Create new API client
	newclient, err := NewAPIClient(d.Get("endpoint").(string), d.Get("port").(int), d.Get("username").(string), d.Get("password").(string), d.Get("acct_id").(string), d.Get("aws_profile").(string), d.Get("aws_region").(string), d.Get("insecure").(bool))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK]: Unable to configure DSM provider",
			Detail:   fmt.Sprintf("[E]: SDK: Terraform: %s", err),
		})
		return nil, diags
	}

	// Check if AWS profile is set and use it within API client
	//if err := d.Get("aws_profile").(string); len(err) > 0 {
	//	err := loadAWSProfileCreds(d.Get("aws_profile").(string), newclient)
	//	if err != nil {
	//		diags = append(diags, diag.Diagnostic{
	//			Severity: diag.Error,
	//			Summary:  "[DSM SDK]: Unable to configure DSM provider",
	//			Detail:   fmt.Sprintf("[E]: SDK: AWS: %s", err),
	//		})
	//		return nil, diags
	//	}
	//}
	return newclient, nil
}
