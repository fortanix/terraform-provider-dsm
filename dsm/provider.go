// **********
// Terraform Provider - DSM: provider
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.5.7
//       - Date:      27/07/2021
// **********

package dsm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const debug_output = true

// [-] Define Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("DSM_ENDPOINT", ""),
				Required:    true,
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
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DSM_USERNAME", ""),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("DSM_PASSWORD", ""),
			},
			"api_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Default:   "",
			},
			"acct_id": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("DSM_ACCT_ID", ""),
				Required:    true,
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
			"azure_region": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "us-east",
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  600,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"dsm_sobject":             resourceSobject(),
			"dsm_aws_sobject":         resourceAWSSobject(),
			"dsm_aws_group":           resourceAWSGroup(),
			"dsm_azure_sobject":       resourceAzureSobject(),
			"dsm_azure_group":         resourceAzureGroup(),
			"dsm_secret":              resourceSecret(),
			"dsm_group":               resourceGroup(),
			"dsm_group_user_role":     resourceGroupUserRole(),
			"dsm_group_crypto_policy": resourceGroupCryptoPolicy(),
			"dsm_app":                 resourceApp(),
			"dsm_csr":                 resourceCsr(),
			"dsm_gcp_ekm_sa":          resourceGcpEkmSa(),
			"dsm_acc_quorum_policy":   resourceAccountQuorumPolicy(),
			"dsm_acc_crypto_policy":   resourceAccountCryptoPolicy(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"dsm_aws_group":    dataSourceAWSGroup(),
			"dsm_azure_group":  dataSourceAzureGroup(),
			"dsm_secret":       dataSourceSecret(),
			"dsm_group":        dataSourceGroup(),
			"dsm_user":         dataSourceUser(),
			"dsm_role":         dataSourceRole(),
			"dsm_version":      dataSourceVersion(),
			"dsm_app":          dataSourceApp(),
			"dsm_sobject":      dataSourceSobject(),
			"dsm_sobject_info": dataSourceSobjectInfo(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create new API client
	newclient, err := NewAPIClient(d.Get("endpoint").(string), d.Get("port").(int), d.Get("username").(string), d.Get("password").(string), d.Get("api_key").(string), d.Get("acct_id").(string), d.Get("aws_profile").(string), d.Get("aws_region").(string), d.Get("azure_region").(string), d.Get("insecure").(bool), d.Get("timeout").(int))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK]: Unable to configure DSM provider",
			Detail:   fmt.Sprintf("[E]: SDK: Terraform: %s", err),
		})
		return nil, diags
	}

	return newclient, nil
}
