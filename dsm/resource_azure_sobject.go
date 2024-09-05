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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// [-] Structs to define Terraform AWS Security Object
type TFAzureSobjectExternal struct {
	Version           string
	Azure_key_name    string
	Azure_key_state   string
	Azure_backup      string
}

// [-] Define Azure Security Object in Terraform
func resourceAzureSobject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateAzureSobject,
		ReadContext:   resourceReadAzureSobject,
		UpdateContext: resourceUpdateAzureSobject,
		DeleteContext: resourceDeleteAzureSobject,
		Description: "Creates a new security object in Azure key vault. This is a Bring-Your-Own-Key (BYOK) method and copies an existing DSM local security object to Azure KV as a Customer Managed Key (CMK).\n" +
		"Azure sobject can also rotate, enable soft deletion and purge the key. For examples of rotate and soft deletion, refer Guides/dsm_azure_sobject.\n\n" +
		"**Note**: Once soft deletion is enabled, Azure sobject can't be modified.\n\n" +
		"**Deletion of a dsm_azure_sobject:** Unlike dsm_sobject, deletion of a dsm_azure_sobject is not normal.\n\n" +
		"**Steps to delete a dsm_azure_sobject**:\n\n" +
		"   * Enable soft_deletion as shown in the examples of guides/dsm_azure_sobject.\n" +
		"   * Enable purge_deleted_key after soft_deletion as shown in the examples of guides/dsm_azure_sobject.\n" +
		"   * A dsm_azure_sobject can be deleted completely only when its state is `destroyed`.\n" +
		"   * A dsm_azure_sobject comes to destroyed state when the key is deleted from Azure key vault.\n" +
		"   * To know whether it is in a destroyed state or not, sync keys operation should be performed.\n" +
		"   * Currently, sync keys is not supported by terraform. This can be done in UI by going to the group and HSM/KMS. Then click on `SYNC KEYS`.",
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
			    "   * `azure-key-type`: Type of a key. It can be used in `PREMIUM` key vault. Value is hardware.",
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"obj_type": {
			    Description: "The type of security object.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_size": {
			    Description: "The size of the security object.",
				Type:     schema.TypeInt,
				Computed: true,
			},
			"key_ops": {
			    Description: "The security object operations permitted.\n\n" +
				"| obj_type | key_size/curve | key_ops |\n" +
				"| -------- | -------- |-------- |\n" +
				"| `RSA` | 2048, 3072, 4096 | APPMANAGEABLE, SIGN, VERIFY, ENCRYPT, DECRYPT, WRAPKEY, UNWRAPKEY, EXPORT |\n" +
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
			    Description: "Whether the security object will be Enabled or Disabled. The values are true/false.",
				Type:     schema.TypeBool,
				Optional: true,
				Default: true,
			},
			"state": {
			    Description: "The key states of the Azure KV key. The values are Created, Deleted, Purged.",
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"expiry_date": {
			    Description: "The security object expiry date in RFC format.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"rotate": {
				Description: "The security object rotation. Specify the method to use for key rotation:\n" +
				"   * `DSM`: To use the same key material.\n" +
				"   * `AZURE`: To rotate from a AZURE key. The key material of new key will be stored in AZURE.\n",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"DSM", "AZURE"}, true),
			},
			"rotate_from": {
				Description: "Name of the security object to be rotated.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"soft_deletion": {
				Description: "Enable soft key deletion in Azure key vault. Key is not usable for Sign/Verify, Wrap/Unwrap or Encrypt/Decrypt operations once it is deleted. The supported values are true/false.\n" +
				" **Note:**  This should be enabled only after the creation.",
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"purge_deleted_key": {
				Description: "Purge deleted key in Azure key vault.purging the key makes all data encrypted with it unrecoverable unless you later import the same key material from Fortanix DSM into the Azure key." +
				"The DSM source key is not affected by this operation. The supported values are true/false.\n" +
				" **Note:**  This should be enabled only after the creation.",
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"external": {
				Description: "AWS CMK level metadata:\n" +
				"   * `Version`\n" +
				"   * `Azure_key_name`\n" +
				"   * `Azure_key_state`\n" +
				"   * `Azure_backup`\n",
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
	if d.Get("soft_deletion").(bool) || d.Get("purge_deleted_key").(bool){
        return invokeErrorDiagsNoSummary("[E] soft_deletion or purge_deleted_key should be enabled only after creation.")
    }
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
		"enabled": d.Get("enabled").(bool),
	}

	if err := d.Get("expiry_date").(string); len(err) > 0 {
		sobj_deactivation_date, date_error := parseTimeToDSM(err)
		if date_error != nil {
			return date_error
		}
		security_object["deactivation_date"] = sobj_deactivation_date
	}
	if err := d.Get("custom_metadata").(map[string]interface{}); len(err) > 0 {
		security_object["custom_metadata"] = d.Get("custom_metadata")
	}
	if rotation_policy := d.Get("rotation_policy").(map[string]interface{}); len(rotation_policy) > 0 {
		security_object["rotation_policy"] = sobj_rotation_policy_write(rotation_policy)
	}
	if rotate := d.Get("rotate").(string); len(rotate) > 0 {
		security_object["name"] = d.Get("rotate_from").(string)
		if rotate == "AZURE" {
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
	return resourceReadAzureSobject(ctx, d, m)
}

// [R]: Read Azure Security Object
func resourceReadAzureSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		if err := d.Set("dsm_name", azuresobject.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("obj_type", req["obj_type"]); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("key_size", req["key_size"]); err != nil {
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
		if err := d.Set("custom_metadata", d.Get("custom_metadata").(map[string]interface{})); err != nil {
			return diag.FromErr(err)
		}
		if key_ops_read, ok := req["key_ops"]; ok {
			if err := setKeyOpsTfState(d, key_ops_read); err != nil {
				return err
			}
		}
		external := &TFAzureSobjectExternal{
			Version:           azuresobject.External.Id.Version,
			Azure_key_name:    azuresobject.External.Id.Label,
			Azure_key_state:   azuresobject.Custom_metadata.Azure_key_state,
			Azure_backup:      azuresobject.Custom_metadata.Azure_backup,
		}
		var externalInt map[string]interface{}
		externalRec, _ := json.Marshal(external)
		json.Unmarshal(externalRec, &externalInt)
		if err := d.Set("external", externalInt); err != nil {
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
    var diags diag.Diagnostics
	if d.HasChange("key") {
		return undoTFstate("key", d)
	}
	if d.HasChange("soft_deletion") && d.Get("soft_deletion").(bool) {
		soft_deletion := map[string]interface{}{}
		if d.Get("external").(map[string]interface{})["Azure_key_state"] != "deleted" {
			_, err := m.(*api_client).APICallBody("POST", fmt.Sprintf("crypto/v1/keys/%s/schedule_deletion", d.Id()), soft_deletion)
			if err != nil {
				d.Set("soft_deletion", false)
			    return invokeErrorDiagsWithSummary(error_summary, fmt.Sprintf("[E]: API: POST crypto/v1/keys/%s/schedule_deletion, %v", d.Id(), err))
			}
		} else {
			return showWarning("The security object is already scheduled for the deletion.")
		}
		if !d.Get("purge_deleted_key").(bool){
			return resourceReadAzureSobject(ctx, d, m)
		}
	}
	if d.HasChange("purge_deleted_key") && d.Get("purge_deleted_key").(bool) {
		// To get the latest Azure_key_state status
		resourceReadAzureSobject(ctx, d, m)
		if d.Get("external").(map[string]interface{})["Azure_key_state"] == "deleted" {
			err := deleteKeyMateialBYOKSobject(d, m)
			if err != nil {
				d.Set("purge_deleted_key", false)
				return invokeErrorDiagsWithSummary(error_summary, fmt.Sprintf("[E]: API: POST crypto/v1/keys/%s/delete_key_material, %v", d.Id(), err))
			}
		}
		d.Set("purge_deleted_key", false)
		return showWarning("The purge key cannot be done until the soft deletion is done for the key.")
	}

	update_azure_sobject := map[string]interface{}{
		"kid": d.Id(),
	}
	has_change := false
	if d.HasChange("name") {
		update_azure_sobject["name"] = d.Get("name").(string)
		has_change = true
	}
	if d.HasChange("description") {
		update_azure_sobject["description"] = d.Get("description").(string)
		has_change = true
	}
	if d.HasChange("custom_metadata") {
		if err := d.Get("custom_metadata").(map[string]interface{}); len(err) > 0 {
			update_azure_sobject["custom_metadata"] = d.Get("custom_metadata")
			has_change = true
		}
	}
	if d.HasChange("enabled") {
		/*
		When the key is in destroyed state, then enabled will be set to false.
		In this case terraform plan/apply will detect the changes for enabled.
		Then terraform apply fails, in this scenario we should show a warning to the user.
		*/
		resourceReadAzureSobject(ctx, d, m)
		if d.Get("state").(string) == "Destroyed" {
			return showWarning("The security object is in destroyed state. It can be deleted now.")
		}
		update_azure_sobject["enabled"] = d.Get("enabled").(bool)
		has_change = true
	}
	if d.HasChange("key_ops") {
		update_azure_sobject["key_ops"] = d.Get("key_ops")
		has_change = true
	}
	if d.HasChange("expiry_date") {
		sobj_deactivation_date, date_error := parseTimeToDSM(d.Get("expiry_date").(string))
		if date_error != nil {
			return date_error
		}
		update_azure_sobject["deactivation_date"] = sobj_deactivation_date
		has_change = true
	}
	if d.HasChange("rotation_policy") {
		if rotation_policy := d.Get("rotation_policy").(map[string]interface{}); len(rotation_policy) > 0 {
			update_azure_sobject["rotation_policy"] = sobj_rotation_policy_write(rotation_policy)
			has_change = true
		}
	}
	if has_change {
		_, err := m.(*api_client).APICallBody("PATCH", fmt.Sprintf("crypto/v1/keys/%s", d.Id()), update_azure_sobject)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: PATCH crypto/v1/keys: %v", err),
			})
			// sets back to original tf state
			resourceReadAzureSobject(ctx, d, m)
			return diags
		}
	}
	return nil
}

// [D]: Delete Azure Security Object
func resourceDeleteAzureSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	/* Before destroying, tf state should be updated. If the dsm_azure_sobject state is not in destroyed state,
		It will give an error.
	*/
	resourceReadAzureSobject(ctx, d, m)
	return deleteBYOKDestroyedSobject(d, m)
}
