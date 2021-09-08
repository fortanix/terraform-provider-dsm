// **********
// Terraform Provider - DSM: resource: group
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.3.6
//       - Date:      27/11/2020
// **********

package dsm

import (
	"context"
	"fmt"

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
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// [C]: Create AWS Group
func resourceCreateAWSGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	group_object := map[string]interface{}{
		"name":           d.Get("name").(string),
		"description":    d.Get("description").(string),
		"hmg_redundancy": "PriorityFailover",
	}

	group_object["add_hmg"] = []map[string]interface{}{
		{
			"url":        fmt.Sprintf("kms.%s.amazonaws.com", m.(*api_client).aws_region),
			"kind":       "AWSKMS",
			"access_key": d.Get("access_key").(string),
			"secret_key": d.Get("secret_key").(string),
			"hsm_order":  0,
			"tls": map[string]interface{}{
				"mode":              "required",
				"validate_hostname": false,
				"ca": map[string]interface{}{
					"ca_set": "global_roots",
				},
			},
		},
	}

	req, err := m.(*api_client).APICallBody("POST", "sys/v1/groups", group_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST sys/v1/groups: %s", err),
		})
		return diags
	}

	d.SetId(req["group_id"].(string))
	return resourceReadGroup(ctx, d, m)
}

// [R]: Read AWS Group
func resourceReadAWSGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, err := m.(*api_client).APICall("GET", fmt.Sprintf("sys/v1/groups/%s", d.Id()))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: GET sys/v1/groups: %s", err),
		})
		return diags
	}

	if err := d.Set("name", req["name"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("group_id", req["group_id"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("acct_id", req["acct_id"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("creator", req["creator"]); err != nil {
		return diag.FromErr(err)
	}
	if _, ok := req["description"]; ok {
		if err := d.Set("description", req["description"].(string)); err != nil {
			return diag.FromErr(err)
		}
	}
	return diags
}

// [U]: Update AWS Group
func resourceUpdateAWSGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

// [D]: Delete AWS Group
func resourceDeleteAWSGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	_, err := m.(*api_client).APICall("DELETE", fmt.Sprintf("sys/v1/groups/%s", d.Id()))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: Group Not Empty: %s", err),
		})
		return diags
	}

	d.SetId("")
	return nil
}
