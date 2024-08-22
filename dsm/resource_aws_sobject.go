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
		Description: "Creates a new security object in AWS KMS. This is a Bring-Your-Own-Key (BYOK) method and copies an existing DSM local security object to AWS KMS as a Customer Managed Key (CMK).The returned resource object contains the UUID of the security object for further references.\n" +
		"AWS sobject can also rotate and enable schedule deletion. For more examples, refer Guides/dsm_aws_sobject, Guides/rotate_with_AWS_option and rotate_with_DSM_option.\n\n" +
		"**Temporary Credentials**: AWS sobject can also be created using AWS temporary credentials. Please refer the below example for temporary credentials.\n\n" +
		"**Note**: Once schedule deletion is enabled, AWS sobject can't be modified.",
		Schema: map[string]*schema.Schema{
			"name": {
			    Description: "The security object name.",
				Type:     schema.TypeString,
				Required: true,
			},
			"dsm_name": {
			    Description: "The security object name from Fortanix DSM (matches the name provided during creation).",
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_id": {
			    Description: "The security object group assignment.",
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
			    Description: "A local security object created/imported to Fortanix DSM(BYOK) and copied to AWS KMS.",
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"copied_to": {
			    Description: "List of security objects copied by the current security object.",
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"copied_from": {
			    Description: "Security object that is copied from another security object.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"replacement": {
			    Description: "Replacement of a security object that was rotated.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"replaced": {
			    Description: "Replaced by a security object.",
				Type:     schema.TypeString,
				Computed: true,
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
				Description: "The creator of the group from Fortanix DSM.\n" +
				"   * `user`: If the group was created by a user, the computed value will be the matching user id.\n" +
				"   * `app`: If the group was created by a app, the computed value will be the matching app id.",
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
				"   * **Note:** Either interval_days or interval_months should be given, but not both.",
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
				},
			},
			"custom_metadata": {
			    Description: "AWS KMS key level metadata information.\n" +
			    "   * `aws-aliases`: Key name within AWS KMS.\n" +
			    "   * `aws-policy`: JSON format of AWS policy that should be enforced for the key.\n" +
			    "   * **Note:** Any other DSM custom metadata can be configured.",
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"aws_tags": {
			    Description: "Any other user-defined AWS metadata information.\n" +
			    "   * e.g. test-key = test-value \n" +
			    "   * The above key value pair will be added as `aws-tag-test-key = test-value` \n",
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"external": {
			    Description: "AWS CMK level metadata:\n" +
			    "   * `Key_arn`\n" +
			    "   * `Key_id`\n" +
			    "   * `Key_state`\n" +
			    "   * `Key_aliases`\n" +
			    "   * `Key_deletion_date`\n",
				Type:     schema.TypeMap,
				Computed: true,
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
				"| `AES` | 256 | ENCRYPT, DECRYPT, WRAPKEY, UNWRAPKEY, DERIVEKEY, MACGENERATE, MACVERIFY, APPMANAGEABLE, EXPORT |\n" +
				"| `RSA` | 2048, 3072, 4096 | APPMANAGEABLE, SIGN, VERIFY, ENCRYPT, DECRYPT, WRAPKEY, UNWRAPKEY, EXPORT  |\n" +
				"| `EC` | NistP256, NistP384, NistP521,SecP256K1 | APPMANAGEABLE, SIGN, VERIFY, AGREEKEY, EXPORT",
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
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
			    Description: "Whether the security object will be enabled or disabled. The values are true/false.",
				Type:     schema.TypeBool,
				Optional: true,
				Default: true,
			},
			"state": {
			    Description: "The key states of the AWS key. The supported values are PendingDeletion, Enabled, Disabled and PendingImport.",
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pending_window_in_days": {
			    Description: "input the value for “days” after which the AWS key will be deleted.\n" +
			    "   * The default value is 7 days.\n" +
			    "   * The minimum value is 7 days.",
				Type:     schema.TypeInt,
				Optional: true,
				Default:  7,
				ValidateFunc: validation.IntAtLeast(7),
			},
			"expiry_date": {
			    Description: "The security object expiry date in RFC format.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"rotate": {
			    Description: "The security object rotation. Specify the method to use for key rotation:\n" +
			    "   * `DSM`: To rotate from a DSM local key. The key material of new key will be stored in DSM.\n" +
			    "   * `AWS`: To rotate from a AWS key. The key material of new key will be stored in AWS.\n",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"DSM", "AWS"}, true),
			},
			"rotate_from": {
			    Description: "Name of the security object to be rotated.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"schedule_deletion": {
				Description: "Enable schedule_deletion to delete the key in AWS KMS.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// [C]: Create AWS Security Object
func resourceCreateAWSSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	endpoint := "crypto/v1/keys/copy"
	if rotate := d.Get("rotate").(string); len(rotate) > 0 {
		if rotate_from := d.Get("rotate_from").(string); len(rotate_from) <= 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: POST %s: 'rotate_from' missing", endpoint),
			})
			return diags
		}
	}

	security_object := map[string]interface{}{
		"name":        d.Get("name").(string),
		"group_id":    d.Get("group_id").(string),
		"key":         d.Get("key"),
		"description": d.Get("description").(string),
		"enabled":     d.Get("enabled").(bool),
	}
	if err := d.Get("expiry_date").(string); len(err) > 0 {
		sobj_deactivation_date, date_error := parseTimeToDSM(err)
		if date_error != nil {
			return date_error
		}
		security_object["deactivation_date"] = sobj_deactivation_date
	}
	if err := d.Get("key_ops").([]interface{}); len(err) > 0 {
		security_object["key_ops"] = d.Get("key_ops")
	}
	if err := d.Get("custom_metadata").(map[string]interface{}); len(err) > 0 {
		security_object["custom_metadata"] = d.Get("custom_metadata")
	}
	if rotation_policy := d.Get("rotation_policy").(map[string]interface{}); len(rotation_policy) > 0 {
		security_object["rotation_policy"] = sobj_rotation_policy_write(rotation_policy)
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

	if rotate := d.Get("rotate").(string); len(rotate) > 0 {
		security_object["name"] = d.Get("rotate_from").(string)
		if rotate == "AWS" {
			endpoint = "crypto/v1/keys/rekey"
		}
	}

	req, err := m.(*api_client).APICallBody("POST", endpoint, security_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST %s: %v", endpoint, err),
		})
		return diags
	}

	d.SetId(req["kid"].(string))
    /*
	schedule deletion:
	Few customers may create the key and schedule for the deletion immediately for about 30 days.
	So, schedule_deletion is enabled in create functionality also.
	If schedule_deletion fails, it returns a warning and adds the security object to tf state.
	*/
	if d.Get("schedule_deletion").(bool) {
		schedule_deletion := map[string]interface{}{
			"pending_window_in_days": d.Get("pending_window_in_days").(int),
		}
		_, err := m.(*api_client).APICallBody("POST", fmt.Sprintf("crypto/v1/keys/%s/schedule_deletion", d.Id()), schedule_deletion)
		if err != nil {
		    // to update the tf state
			resourceReadAWSSobject(ctx, d, m)
			return scheduleDeletionWarning(d, err)
		}
	}
	return resourceReadAWSSobject(ctx, d, m)
}

// [R]: Read AWS Security Object
func resourceReadAWSSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, statuscode, err := m.(*api_client).APICall("GET", fmt.Sprintf("crypto/v1/keys/%s?show_destroyed=true&show_deleted=true", d.Id()))
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
		empty_array := []string{}
		if _, ok := req["links"]; ok {
			if links := req["links"].(map[string]interface{}); len(links) > 0 {
				if _, copiedToExists := req["links"].(map[string]interface{})["copiedTo"]; copiedToExists {
					if err := d.Set("copied_to", req["links"].(map[string]interface{})["copiedTo"].([]interface{})); err != nil {
						return diag.FromErr(err)
					}
				} else {
				    d.Set("copied_to", empty_array)
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
		tfstate_custom_metadata := d.Get("custom_metadata").(map[string]interface{})
		if len(tfstate_custom_metadata) > 0 {
			if err := d.Set("custom_metadata", tfstate_custom_metadata); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err := d.Set("custom_metadata", req["custom_metadata"]); err != nil {
				return diag.FromErr(err)
			}
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
		if key_ops_read, ok := req["key_ops"]; ok {
		    if err := setKeyOpsTfState(d, key_ops_read); err != nil {
                return err
            }
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
			sobj_deactivation_date, date_error := parseTimeToDSM(rfcdate.(string))
			if date_error != nil {
				return date_error
			}
			if newerr := d.Set("expiry_date", sobj_deactivation_date); newerr != nil {
				return diag.FromErr(newerr)
			}
		}
		if _, ok := req["rotation_policy"]; ok {
			rotation_policy := sobj_rotation_policy_read(req["rotation_policy"].(map[string]interface{}))
			if err := d.Set("rotation_policy", rotation_policy); err != nil {
				return diag.FromErr(err)
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
	if d.HasChange("schedule_deletion") {
		if d.Get("schedule_deletion").(bool) {
			schedule_deletion := map[string]interface{}{
				"pending_window_in_days": d.Get("pending_window_in_days").(int),
			}
			if d.Get("external").(map[string]interface{})["Key_state"] != "PendingDeletion" {
				_, err := m.(*api_client).APICallBody("POST", fmt.Sprintf("crypto/v1/keys/%s/schedule_deletion", d.Id()), schedule_deletion)
				if err != nil {
					return invokeErrorDiagsWithSummary(error_summary, fmt.Sprintf("[E]: API: POST crypto/v1/keys/%s/schedule_deletion, %v", d.Id(), err))
				}
			}
			return resourceReadAWSSobject(ctx, d, m)
		}
	}
	if d.HasChange("key") {
		return undoTFstate("key", d)
	}
	// already has been replaced so "rotate" and "rotate_from" does not apply
	_, replacement := d.GetOk("replacement")
	_, replaced := d.GetOk("replaced")
	if replacement || replaced {
		d.Set("rotate", "")
		d.Set("rotate_from", "")
	}

	update_aws_sobject := map[string]interface{}{
		"kid": d.Id(),
	}
	has_change := false
	if d.HasChange("custom_metadata") {
		if err := d.Get("custom_metadata").(map[string]interface{}); len(err) > 0 {

			old_custom_metadata, _ := d.GetChange("custom_metadata")
			//update_aws_sobject["custom_metadata"] = old_custom_metadata

			// FYOO: Needs work
			update_aws_sobject["custom_metadata"] = make(map[string]interface{})

			if newAlias, ok := d.Get("custom_metadata").(map[string]interface{})["aws-aliases"]; ok {
				if replacement {
					update_aws_sobject["custom_metadata"].(map[string]interface{})["aws-aliases"] = old_custom_metadata.(map[string]interface{})["aws-aliases"]
				} else {
					update_aws_sobject["custom_metadata"].(map[string]interface{})["aws-aliases"] = newAlias.(string)
				}
			}

			if newPolicy, ok := d.Get("custom_metadata").(map[string]interface{})["aws-policy"]; ok {
				update_aws_sobject["custom_metadata"].(map[string]interface{})["aws-policy"] = newPolicy
			} else {
				update_aws_sobject["custom_metadata"].(map[string]interface{})["aws-policy"] = old_custom_metadata.(map[string]interface{})["aws-policy"]
			}

			for k := range d.Get("custom_metadata").(map[string]interface{}) {
				if strings.HasPrefix(k, "aws-tag-") {
					update_aws_sobject["custom_metadata"].(map[string]interface{})[k] = d.Get("custom_metadata").(map[string]interface{})[k]
				}
			}

			// FYOO: Get tags
			if d.HasChange("aws_tags") {
				if err := d.Get("aws_tags").(map[string]interface{}); len(err) > 0 {
					if _, cmExists := update_aws_sobject["custom_metadata"]; !cmExists {
						update_aws_sobject["custom_metadata"] = make(map[string]interface{})
					}
					for aws_tags_k := range d.Get("aws_tags").(map[string]interface{}) {
						update_aws_sobject["custom_metadata"].(map[string]interface{})[(fmt.Sprintf("aws-tag-%s", aws_tags_k))] = d.Get("aws_tags").(map[string]interface{})[aws_tags_k]
					}
				}
			}
			has_change = true
		}
	}
	if d.HasChange("expiry_date") {
		sobj_deactivation_date, date_error := parseTimeToDSM(d.Get("expiry_date").(string))
		if date_error != nil {
			return date_error
		}
		update_aws_sobject["deactivation_date"] = sobj_deactivation_date
		has_change = true
	}
	if d.HasChange("name") {
		update_aws_sobject["name"] = d.Get("name").(string)
		has_change = true
	}
	if d.HasChange("description") {
		update_aws_sobject["description"] = d.Get("description").(string)
		has_change = true
	}
	if d.HasChange("enabled") {
		update_aws_sobject["enabled"] = d.Get("enabled").(bool)
		has_change = true
	}
	if d.HasChange("rotation_policy") {
		if rotation_policy := d.Get("rotation_policy").(map[string]interface{}); len(rotation_policy) > 0 {
			update_aws_sobject["rotation_policy"] = sobj_rotation_policy_write(rotation_policy)
			has_change = true
		}
	}
	if d.HasChange("key_ops") {
		update_aws_sobject["key_ops"] = d.Get("key_ops")
		has_change = true
	}

	if has_change {
		_, err := m.(*api_client).APICallBody("PATCH", fmt.Sprintf("crypto/v1/keys/%s", d.Id()), update_aws_sobject)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: PATCH crypto/v1/keys: %v", err),
			})
			// sets back to original tf state
			resourceReadAWSSobject(ctx, d, m)
			return diags
		}
	}

	return resourceReadAWSSobject(ctx, d, m)
}

// [D]: Delete AWS Security Object
func resourceDeleteAWSSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//commom.go
	return deleteBYOKDestroyedSobject(d, m)
}
