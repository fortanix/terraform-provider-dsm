// **********
// Terraform Provider - DSM: resource: group
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.5.0
//       - Date:      27/11/2020
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

// [-] Define Group
func resourceAWSGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateAWSGroup,
		ReadContext:   resourceReadAWSGroup,
		UpdateContext: resourceUpdateAWSGroup,
		DeleteContext: resourceDeleteAWSGroup,
		Description: "Creates a Fortanix DSM group mapped to AWS KMS in the cluster as a resource. This group acts as a container for security objects. The returned resource object contains the UUID of the group for further references.",
		Schema: map[string]*schema.Schema{
			"name": {
			    Description: "The name follows the nomenclature of `<Custom Group Name>-aws-<Region>`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
			    Description: "The unique ID for AWS KMS Mapped group from Fortanix DSM.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"acct_id": {
			    Description: "The Account ID from Fortanix DSM.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator": {
				Description: "The creator of the group from Fortanix DSM.\n" +
				"   * `user`: If the group was created by a user, the computed value will be the matching user id.\n" +
				"   * `app`: If the group was created by a app, the computed value will be the matching app id.",
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
			    Description: "The description of the AWS KMS group.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"access_key": {
			    Description: "The Access Key ID to set for AWS KMS group for programmatic (API) access to AWS Services.",
				Type:     schema.TypeString,
				Required: true,
			},
			"secret_key": {
			    Description: "The Secret Access Key to set for AWS KMS group for programmatic (API) access to AWS Services.",
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// [C]: Create AWS Group
func resourceCreateAWSGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	group_object := map[string]interface{}{
		// 0.4.1: AWS KMS Group Name to be predefined as <string>-aws-<region>
		"name":           fmt.Sprintf("%s-aws-%s", d.Get("name").(string), m.(*api_client).aws_region),
		"description":    d.Get("description").(string),
		"hmg_redundancy": "PriorityFailover",
	}

	group_object["add_hmg"] = []map[string]interface{}{
		{
			"url":       fmt.Sprintf("kms.%s.amazonaws.com", m.(*api_client).aws_region),
			"kind":      "AWSKMS",
			"hsm_order": 0,
			"tls": map[string]interface{}{
				"mode":              "required",
				"validate_hostname": false,
				"ca": map[string]interface{}{
					"ca_set": "global_roots",
				},
			},
		},
	}

	// 0.5.0: parse optionals
	access_key, access_key_exists := d.GetOk("access_key")
	if access_key_exists {
		group_object["add_hmg"].([]map[string]interface{})[0]["access_key"] = access_key.(string)
	}
	secret_key, secret_key_exists := d.GetOk("secret_key")
	if secret_key_exists {
		group_object["add_hmg"].([]map[string]interface{})[0]["secret_key"] = secret_key.(string)
	}

	req, err := m.(*api_client).APICallBody("POST", "sys/v1/groups", group_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST sys/v1/groups: %v", err),
		})
		return diags
	}

	d.SetId(req["group_id"].(string))
	return resourceReadAWSGroup(ctx, d, m)
}

// [R]: Read AWS Group
func resourceReadAWSGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, statuscode, err := m.(*api_client).APICall("GET", fmt.Sprintf("sys/v1/groups/%s", d.Id()))
	if statuscode == 404 {
		d.SetId("")
	} else {
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: GET sys/v1/groups: %v", err),
			})
			return diags
		} else {
			// Convert returned call into AWSGroup Map
			jsonbody, err := json.Marshal(req)
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
			if _, ok := req["description"]; ok {
				d.Set("description", req["description"].(string))
			}
		}
	}
	return diags
}

// [U]: Update AWS Group
func resourceUpdateAWSGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.Set("secret_key", "")
	return diags
}

// [D]: Delete AWS Group
func resourceDeleteAWSGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	_, statuscode, err := m.(*api_client).APICall("DELETE", fmt.Sprintf("sys/v1/groups/%s", d.Id()))
	if (err != nil) && (statuscode != 404) && (statuscode != 400) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: DELETE sys/v1/groups: %v", err),
		})
		return diags
	} else {
		if statuscode == 400 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Call to DSM provider API client failed",
				Detail:   fmt.Sprintf("[E]: API: DELETE sys/v1/groups: %s", "Group Not Empty"),
			})
			return diags
		}
	}

	d.SetId("")
	return nil
}
