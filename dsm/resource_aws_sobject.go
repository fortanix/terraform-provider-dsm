// **********
// Terraform Provider - DSM: resource: aws security object
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.5.8
//       - Date:      27/11/2020
// **********

package dsm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// [-] Structs to define Terraform AWS Security Object
type TFAWSSobjectExternal struct {
	Key_arn           string
	Key_id            string
	Key_state         string
	Key_aliases       string
	Key_deletion_date string
}

// [-] Define AWS Security Object in Terraform
func resourceAWSSobject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateAWSSobject,
		ReadContext:   resourceReadAWSSobject,
		UpdateContext: resourceUpdateAWSSobject,
		DeleteContext: resourceDeleteAWSSobject,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dsm_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"copied_to": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"copied_from": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replacement": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replaced": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kid": {
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
			"custom_metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"aws_tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"external": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"obj_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"key_ops": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pending_window_in_days": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  7,
			},
			"expiry_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rotate": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"DSM", "ALL"}, true),
			},
			"rotate_from": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// [C]: Create AWS Security Object
func resourceCreateAWSSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if rotate := d.Get("rotate").(string); len(rotate) > 0 {
		if rotate_from := d.Get("rotate_from").(string); len(rotate_from) <= 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   "[E]: API: GET crypto/v1/keys/copy: 'rotate_from' missing",
			})
			return diags
		}
	}

	security_object := map[string]interface{}{
		"name":        d.Get("name").(string),
		"group_id":    d.Get("group_id").(string),
		"key":         d.Get("key"),
		"description": d.Get("description").(string),
	}

	if rfcdate := d.Get("expiry_date").(string); len(rfcdate) > 0 {
		layoutRFC := "2006-01-02T15:04:05Z"
		layoutDSM := "20060102T150405Z"
		ddate, newerr := time.Parse(layoutRFC, rfcdate)
		if newerr != nil {
			return diag.FromErr(newerr)
		}
		security_object["deactivation_date"] = ddate.Format(layoutDSM)
	}

	if err := d.Get("custom_metadata").(map[string]interface{}); len(err) > 0 {
		security_object["custom_metadata"] = d.Get("custom_metadata")
	}

	// FYOO: Get tags
	if err := d.Get("aws_tags").(map[string]interface{}); len(err) > 0 {
		if _, cmExists := security_object["custom_metadata"]; !cmExists {
			security_object["custom_metadata"] = make(map[string]interface{})
		}
		for aws_tags_k := range d.Get("aws_tags").(map[string]interface{}) {
			security_object["custom_metadata"].(map[string]interface{})[(fmt.Sprintf("aws-tag-%s", aws_tags_k))] = d.Get("aws_tags").(map[string]interface{})[aws_tags_k]
		}
	}

	if err := d.Get("rotate").(string); len(err) > 0 {
		security_object["name"] = d.Get("rotate_from").(string)
	}

	req, err := m.(*api_client).APICallBody("POST", "crypto/v1/keys/copy", security_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST crypto/v1/keys/copy: %s", err),
		})
		return diags
	}

	d.SetId(req["kid"].(string))
	return resourceReadAWSSobject(ctx, d, m)
}

// [R]: Read AWS Security Object
func resourceReadAWSSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, statuscode, err := m.(*api_client).APICall("GET", fmt.Sprintf("crypto/v1/keys/%s", d.Id()))
	if statuscode == 404 {
		d.SetId("")
	} else {
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: GET crypto/v1/keys: %s", err),
			})
			return diags
		}

		// Convert returned call into AWSSobject Map
		jsonbody, err := json.Marshal(req)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to parse DSM provider API client output",
				Detail:   fmt.Sprintf("[E]: API: GET crypto/v1/keys: %s", err),
			})
			return diags
		}

		awssobject := AWSSobject{}
		if err := json.Unmarshal(jsonbody, &awssobject); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to parse DSM provider API client output",
				Detail:   fmt.Sprintf("[E]: API: GET crypto/v1/keys: %s", err),
			})
			return diags
		}

		// Sync DSM and Terraform attributes
		if err := d.Set("dsm_name", awssobject.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("group_id", req["group_id"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if _, ok := req["links"]; ok {
			if links := req["links"].(map[string]interface{}); len(links) > 0 {
				if _, copiedToExists := req["links"].(map[string]interface{})["copiedTo"]; copiedToExists {
					if err := d.Set("copied_to", req["links"].(map[string]interface{})["copiedTo"].([]interface{})); err != nil {
						return diag.FromErr(err)
					}
				}
				if _, copiedFromExists := req["links"].(map[string]interface{})["copiedFrom"]; copiedFromExists {
					if err := d.Set("copied_from", req["links"].(map[string]interface{})["copiedFrom"].(string)); err != nil {
						return diag.FromErr(err)
					}
				}
				if _, replacementExists := req["links"].(map[string]interface{})["replacement"]; replacementExists {
					if err := d.Set("replacement", req["links"].(map[string]interface{})["replacement"].(string)); err != nil {
						return diag.FromErr(err)
					}
				}
				if _, replacedExists := req["links"].(map[string]interface{})["replaced"]; replacedExists {
					if err := d.Set("replaced", req["links"].(map[string]interface{})["replaced"].(string)); err != nil {
						return diag.FromErr(err)
					}
				}
			}
		}
		if err := d.Set("kid", req["kid"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("acct_id", req["acct_id"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("creator", req["creator"]); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("custom_metadata", req["custom_metadata"]); err != nil {
			return diag.FromErr(err)
		}
		external := &TFAWSSobjectExternal{
			Key_arn:           awssobject.External.Id.Key_arn,
			Key_id:            awssobject.External.Id.Key_id,
			Key_state:         awssobject.Custom_metadata.Aws_key_state,
			Key_aliases:       awssobject.Custom_metadata.Aws_aliases,
			Key_deletion_date: awssobject.Custom_metadata.Aws_deletion_date,
		}
		var externalInt map[string]interface{}
		externalRec, _ := json.Marshal(external)
		json.Unmarshal(externalRec, &externalInt)
		if err := d.Set("external", externalInt); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("key_ops", req["key_ops"]); err != nil {
			return diag.FromErr(err)
		}
		if _, ok := req["description"]; ok {
			if err := d.Set("description", req["description"].(string)); err != nil {
				return diag.FromErr(err)
			}
		}
		if err := d.Set("enabled", req["enabled"].(bool)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("state", req["state"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if rfcdate, ok := req["deactivation_date"]; ok {
			// FYOO: once it's set, you can't remove deactivation date
			layoutRFC := "2006-01-02T15:04:05Z"
			layoutDSM := "20060102T150405Z"
			ddate, newerr := time.Parse(layoutDSM, rfcdate.(string))
			if newerr != nil {
				return diag.FromErr(newerr)
			}
			if newerr = d.Set("expiry_date", ddate.Format(layoutRFC)); newerr != nil {
				return diag.FromErr(newerr)
			}
		}
		// FYOO: clear values that are irrelevant
		d.Set("rotate", "")
		d.Set("rotate_from", "")
	}

	return diags
}

// [U]: Update AWS Security Object
func resourceUpdateAWSSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// already has been replaced so "rotate" and "rotate_from" does not apply
	_, replacement := d.GetOkExists("replacement")
	_, replaced := d.GetOkExists("replaced")
	if replacement || replaced {
		d.Set("rotate", "")
		d.Set("rotate_from", "")
	}

	if d.HasChange("custom_metadata") {
		if err := d.Get("custom_metadata").(map[string]interface{}); len(err) > 0 {
			update_aws_metadata := map[string]interface{}{
				"kid": d.Id(),
			}
			old_custom_metadata, _ := d.GetChange("custom_metadata")
			//update_aws_metadata["custom_metadata"] = old_custom_metadata

			// FYOO: Needs work
			update_aws_metadata["custom_metadata"] = make(map[string]interface{})

			if newAlias, ok := d.Get("custom_metadata").(map[string]interface{})["aws-aliases"]; ok {
				if replacement {
					update_aws_metadata["custom_metadata"].(map[string]interface{})["aws-aliases"] = old_custom_metadata.(map[string]interface{})["aws-aliases"]
				} else {
					update_aws_metadata["custom_metadata"].(map[string]interface{})["aws-aliases"] = newAlias.(string)
				}
			}

			if newPolicy, ok := d.Get("custom_metadata").(map[string]interface{})["aws-policy"]; ok {
				update_aws_metadata["custom_metadata"].(map[string]interface{})["aws-policy"] = newPolicy
			} else {
				update_aws_metadata["custom_metadata"].(map[string]interface{})["aws-policy"] = old_custom_metadata.(map[string]interface{})["aws-policy"]
			}

			for k := range d.Get("custom_metadata").(map[string]interface{}) {
				if strings.HasPrefix(k, "aws-tag-") {
					update_aws_metadata["custom_metadata"].(map[string]interface{})[k] = d.Get("custom_metadata").(map[string]interface{})[k]
				}
			}

			// FYOO: Get tags
			if d.HasChange("aws_tags") {
				if err := d.Get("aws_tags").(map[string]interface{}); len(err) > 0 {
					if _, cmExists := update_aws_metadata["custom_metadata"]; !cmExists {
						update_aws_metadata["custom_metadata"] = make(map[string]interface{})
					}
					for aws_tags_k := range d.Get("aws_tags").(map[string]interface{}) {
						update_aws_metadata["custom_metadata"].(map[string]interface{})[(fmt.Sprintf("aws-tag-%s", aws_tags_k))] = d.Get("aws_tags").(map[string]interface{})[aws_tags_k]
					}
				}
			}

			_, err := m.(*api_client).APICallBody("PATCH", fmt.Sprintf("crypto/v1/keys/%s", d.Id()), update_aws_metadata)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "[DSM SDK] Unable to call DSM provider API client",
					Detail:   fmt.Sprintf("[E]: API: PATCH crypto/v1/keys: %s", err),
				})
				return diags
			}
		}
	}

	return resourceReadAWSSobject(ctx, d, m)
}

// [D]: Delete AWS Security Object
func resourceDeleteAWSSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// FIXME: Since deleting, might as well remove the alias if exists
	if _, ok := d.Get("custom_metadata").(map[string]interface{})["aws-aliases"]; ok {
		remove_aws_alias := map[string]interface{}{
			"kid": d.Id(),
		}
		remove_aws_alias["custom_metadata"] = map[string]interface{}{
			"aws-aliases": "",
			"aws-policy":  d.Get("custom_metadata").(map[string]interface{})["aws-policy"],
		}
		_, err := m.(*api_client).APICallBody("PATCH", fmt.Sprintf("crypto/v1/keys/%s", d.Id()), remove_aws_alias)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: PATCH crypto/v1/keys: %s", err),
			})
			return diags
		}
	}

	// FIXME: Need to schedule deletion then delete the key - default is set to 7 days for now (need to specify)
	delete_object := map[string]interface{}{
		"pending_window_in_days": d.Get("pending_window_in_days").(int),
	}
	if d.Get("custom_metadata").(map[string]interface{})["aws-key-state"] != "PendingDeletion" {
		_, err := m.(*api_client).APICallBody("POST", fmt.Sprintf("crypto/v1/keys/%s/schedule_deletion", d.Id()), delete_object)
		if err != nil {
			return err
		}
	}

	_, statuscode, err := m.(*api_client).APICall("DELETE", fmt.Sprintf("crypto/v1/keys/%s", d.Id()))
	if (err != nil) && (statuscode != 404) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: DELETE crypto/v1/keys: %s", err),
		})
		return diags
	}

	d.SetId("")
	return nil
}
