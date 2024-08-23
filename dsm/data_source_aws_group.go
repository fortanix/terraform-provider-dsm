// **********
// Terraform Provider - DSM: data source: aws kms group
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.5.0
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

func dataSourceAWSGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAWSGroupRead,
		Description: "Returns the Fortanix DSM AWS KMS mapped group object from the cluster as a Data Source for AWS KMS.",
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The AWS KMS group object name in Fortanix DSM.",
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Description: "The AWS KMS group object ID from Fortanix DSM.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"acct_id": {
				Description: "The Account ID from Fortanix DSM.",
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
			"region": {
				Description: "The AWS region mapped to the group from which keys are imported.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Description: "The AWS KMS group object description.",
				Type:     schema.TypeString,
				Computed: true,
				Default:  "",
			},
			"access_key": {
				Description: "The Access Key ID used to communicate with AWS KMS.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "na",
			},
			"secret_key": {
				Description: "AWS KMS Secret key.",
				Type:      schema.TypeString,
				Optional:  true,
				Default:   "na",
				Sensitive: true,
			},
			"scan": {
				Description: "Syncs keys from AWS KMS to the AWS KMS group in DSM. Value is either true/false.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAWSGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var group_data map[string]interface{}
	original_name := d.Get("name").(string)
	modified_name := fmt.Sprintf("%s-aws-%s", original_name, m.(*api_client).aws_region)

	req, err := m.(*api_client).APICallList("GET", "sys/v1/groups")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: GET sys/v1/groups: %v", err),
		})
		return diags
	}

	// Shashi: First, check for the group name as provided. 
	// If not found, fallback to checking the group name in the "name-aws-region" format used prior to v0.5.33.
	for _, data := range req {
		group_name := data.(map[string]interface{})["name"].(string)
		if group_name == original_name {
			group_data = data.(map[string]interface{})
			break
		}
		if group_name == modified_name && group_data == nil {
			group_data = data.(map[string]interface{})
		}
	}

	if group_data == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Group not found.",
			Detail:   fmt.Sprintf("[E]: No group found with name: %s or %s", original_name, modified_name),
		})
		return diags
	}

	jsonbody, parsing_error := json.Marshal(group_data)
	if parsing_error != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to parse DSM provider API client output",
			Detail:   fmt.Sprintf("[E]: API: GET crypto/v1/groups: %s", parsing_error),
		})
		return diags
	}

	awsgroup := AWSGroup{}
	if err := json.Unmarshal(jsonbody, &awsgroup); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to parse DSM provider API client output",
			Detail:   fmt.Sprintf("[E]: API: GET sys/v1/groups: %s", err),
		})
		return diags
	}
	// FYOO: AWSGroup must conform to this JSON struct - if this crashes, then we have DSM issues
	name := awsgroup.Name
	d.Set("name", name)
	d.Set("group_id", awsgroup.Group_id)
	d.Set("acct_id", awsgroup.Acct_id)

	var creatorInt map[string]interface{}
	creatorRec, _ := json.Marshal(awsgroup.Creator)
	json.Unmarshal(creatorRec, &creatorInt)
	d.Set("creator", creatorInt)
	d.Set("region", m.(*api_client).aws_region)
	// FYOO: there is only one HMG per AWSGroup
	for _, value := range awsgroup.Hmg {
		d.Set("access_key", value.Access_key)
	}
	// FYOO: remove sensitive information
	d.Set("secret_key", "")
	// FYOO: if description is blank, DSM does not return
	if _, ok := group_data["description"]; ok {
		d.Set("description", group_data["description"].(string))
	}

	d.SetId(d.Get("group_id").(string))

	// If Scan is set, then move to scanning for data source
	if d.Get("scan").(bool) {
		check_hmg_req := map[string]interface{}{}
		// Scan the AWS Group first before
		_, err := m.(*api_client).APICallBody("POST", fmt.Sprintf("sys/v1/groups/%s/hmg/check", d.Get("group_id").(string)), check_hmg_req)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: POST sys/v1/groups/-/hmg/check: %v", err),
			})
			return diags
		}

		_, err = m.(*api_client).APICallBody("POST", fmt.Sprintf("sys/v1/groups/%s/hmg/scan", d.Get("group_id").(string)), check_hmg_req)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: POST sys/v1/groups/-/hmg/scan: %v", err),
			})
			return diags
		}
	}
	return nil
}
