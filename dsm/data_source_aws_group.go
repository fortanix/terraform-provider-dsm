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
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAWSGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAWSGroupRead,
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
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"access_key": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "na",
			},
			"secret_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Default:   "na",
				Sensitive: true,
			},
			"scan": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAWSGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	for _, data := range req {
		prefix_name := fmt.Sprintf("%s-aws-%s", d.Get("name").(string), m.(*api_client).aws_region)
		if data.(map[string]interface{})["name"].(string) == prefix_name {
			jsonbody, err := json.Marshal(data)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "[DSM SDK] Unable to parse DSM provider API client output",
					Detail:   fmt.Sprintf("[E]: API: GET crypto/v1/groups: %s", err),
				})
				return diags
			}

			awsgroup := AWSGroup{}
			if err := json.Unmarshal(jsonbody, &awsgroup); err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "[DSM SDK] Unable to parse DSM provider API client output",
					Detail:   fmt.Sprintf("[E]: API: GET crypto/v1/groups: %s", err),
				})
				return diags
			}
			// FYOO: AWSGroup must conform to this JSON struct - if this crashes, then we have DSM issues
			name := strings.Split(awsgroup.Name, fmt.Sprintf("-aws-%s", m.(*api_client).aws_region))
			d.Set("name", name[0])
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
			if _, ok := data.(map[string]interface{})["description"]; ok {
				d.Set("description", data.(map[string]interface{})["description"].(string))
			}
		}
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
