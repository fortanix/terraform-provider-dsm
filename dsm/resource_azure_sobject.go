// **********
// Terraform Provider - DSM: resource: azure security object
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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// [-] Define Azure Security Object in Terraform
func resourceAzureSobject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateAzureSobject,
		ReadContext:   resourceReadAzureSobject,
		UpdateContext: resourceUpdateAzureSobject,
		DeleteContext: resourceDeleteAzureSobject,
		Description: "Creates a new security object in Azure key vault. This is a Bring-Your-Own-Key (BYOK) method and copies an existing DSM local security object to Azure KV as a Customer Managed Key (CMK).",
		Schema: map[string]*schema.Schema{
			"name": {
			    Description: "The security object name.",
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
			    Description: "The Azure group ID in Fortanix DSM into which the key will be generated.",
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
			    Description: "A local security object imported to Fortanix DSM(BYOK) and copied to Azure KV.",
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"links": {
			    Description: "Link between local security object and Azure KV security object.",
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"kid": {
			    Description: "The security object ID from Fortanix DSM.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"acct_id": {
			    Description: "The account ID from Fortanix DSM.",
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
			"rotation_policy": {
				Description: "Policy to rotate a Security Object, configure the below parameters.\n" +
				"   * `interval_days`: Rotate the key for every given number of days.\n" +
				"   * `interval_months`: Rotate the key for every given number of months.\n" +
				"   * `effective_at`: Start of the rotation policy time.\n" +
				"   * `deactivate_rotated_key`: Deactivate original key after rotation true/false.\n" +
				"   * **Note:** Either interval_days or interval_months should be given, but not both.",
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
				},
			},
			"custom_metadata": {
			    Description: "Azure CMK level metadata information.\n" +
			    "   * `azure-key-name`: Key name within Azure KV.\n" +
			    "   * **Note:** By default dsm_azure_sobject creates the key as a software protected key. For a hardware protected key use the below parameter.\n" +
			    "   * `azure-key-type`: Type of a key. It can be used in `PREMIUM` key vault. Values are software/hardware.",
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"obj_type": {
			    Description: "The type of security object.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_size": {
			    Description: "The size of the security object.",
				Type:     schema.TypeInt,
				Optional: true,
			},
			"key_ops": {
			    Description: "The security object operations permitted.\n\n" +
				"| obj_type | key_size/curve | key_ops |\n" +
				"| -------- | -------- |-------- |\n" +
				"| `RSA` | 2048, 3072, 4096 | APPMANAGEABLE, SIGN, VERIFY, ENCRYPT, DECRYPT, WRAPKEY, UNWRAPKEY, EXPORT |\n" +
				"| `EC` | NistP256, NistP384, NistP521,SecP256K1 | APPMANAGEABLE, SIGN, VERIFY, AGREEKEY, EXPORT",
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": {
			    Description: "The security object description.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"enabled": {
			    Description: "Whether the security object will be Enabled or Disabled. The values are true/false.",
				Type:     schema.TypeBool,
				Optional: true,
			},
			"state": {
			    Description: "The key states of the Azure KV key. The values are Created, Deleted, Purged.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"expiry_date": {
			    Description: "The security object expiry date in RFC format.",
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// [C]: Create Azure Security Object
func resourceCreateAzureSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

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
	if rotation_policy := d.Get("rotation_policy").(map[string]interface{}); len(rotation_policy) > 0 {
		security_object["rotation_policy"] = sobj_rotation_policy_write(rotation_policy)
	}

	req, err := m.(*api_client).APICallBody("POST", "crypto/v1/keys/copy", security_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST crypto/v1/keys/copy: %v", err),
		})
		return diags
	}

	d.SetId(req["kid"].(string))
	return resourceReadAzureSobject(ctx, d, m)
}

// [R]: Read Azure Security Object
func resourceReadAzureSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, statuscode, err := m.(*api_client).APICall("GET", fmt.Sprintf("crypto/v1/keys/%s", d.Id()))
	if statuscode == 404 {
		d.SetId("")
	} else {
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: GET crypto/v1/keys: %v", err),
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

		azuresobject := AzureSobject{}
		if err := json.Unmarshal(jsonbody, &azuresobject); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to parse DSM provider API client output",
				Detail:   fmt.Sprintf("[E]: API: GET crypto/v1/keys: %s", err),
			})
			return diags
		}

		// Sync DSM and Terraform attributes
		if err := d.Set("name", azuresobject.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("group_id", req["group_id"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("links", req["links"]); err != nil {
			return diag.FromErr(err)
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
		//external := &TFAWSSobjectExternal{
		//	Key_arn:           awssobject.External.Id.Key_arn,
		//	Key_id:            awssobject.External.Id.Key_id,
		//	Key_state:         awssobject.Custom_metadata.Aws_key_state,
		//	Key_aliases:       awssobject.Custom_metadata.Aws_aliases,
		//	Key_deletion_date: awssobject.Custom_metadata.Aws_deletion_date,
		//}
		//var externalInt map[string]interface{}
		//externalRec, _ := json.Marshal(external)
		//json.Unmarshal(externalRec, &externalInt)
		//if err := d.Set("external", externalInt); err != nil {
		//	return diag.FromErr(err)
		//}
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
		if _, ok := req["rotation_policy"]; ok {
			rotation_policy := sobj_rotation_policy_read(req["rotation_policy"].(map[string]interface{}))
			if err := d.Set("rotation_policy", rotation_policy); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}

// [U]: Update Azure Security Object
func resourceUpdateAzureSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

// [D]: Delete Azure Security Object
func resourceDeleteAzureSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// FIXME: Since deleting, might as well remove the alias if exists
	//if _, ok := d.Get("custom_metadata").(map[string]interface{})["aws-aliases"]; ok {
	//	remove_aws_alias := map[string]interface{}{
	//		"kid": d.Id(),
	//	}
	//	remove_aws_alias["custom_metadata"] = map[string]interface{}{
	//		"aws-aliases": "",
	//		"aws-policy":  d.Get("custom_metadata").(map[string]interface{})["aws-policy"],
	//	}
	//	_, err := m.(*api_client).APICallBody("PATCH", fmt.Sprintf("crypto/v1/keys/%s", d.Id()), remove_aws_alias)
	//	if err != nil {
	//		diags = append(diags, diag.Diagnostic{
	//			Severity: diag.Error,
	//			Summary:  "[DSM SDK] Unable to call DSM provider API client",
	//			Detail:   fmt.Sprintf("[E]: API: PATCH crypto/v1/keys: %s", err),
	//		})
	//		return diags
	//	}
	//}

	// FIXME: Need to schedule deletion then delete the key
	delete_object := make(map[string]interface{})
	if d.Get("custom_metadata").(map[string]interface{})["azure-key-state"] != "deleted" {
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
			Detail:   fmt.Sprintf("[E]: API: DELETE crypto/v1/keys: %v", err),
		})
		return diags
	}

	d.SetId("")
	return nil
}
