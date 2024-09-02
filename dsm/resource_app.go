// **********
// Terraform Provider - DSM: resource: app
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.5.3
//       - Date:      27/11/2020
// **********

package dsm

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// [-] Define App
func resourceApp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateApp,
		ReadContext:   resourceReadApp,
		UpdateContext: resourceUpdateApp,
		DeleteContext: resourceDeleteApp,
		Description: "Creates a new DSM App of type API key.The returned resource object contains the UUID of the app for further references.\n" +
		"This resource can also rotate/regenerate an API key. Default permissions of any group can be modified.",
		Schema: map[string]*schema.Schema{
			"name": {
			    Description: "The Fortanix DSM App name.",
				Type:     schema.TypeString,
				Required: true,
			},
			"app_id": {
			    Description: "The unique ID of the app from Terraform.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_group": {
			    Description: "The Fortanix DSM group object id to be mapped to the app by default.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"other_group": {
			    Description: "The Fortanix DSM group object id the app needs to be assigned to. If you want to delete the existing groups from an app, remove the ids during update.",
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"acct_id": {
			    Description: "The account ID from Fortanix DSM.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator": {
				Description: "The creator of the app from Fortanix DSM.\n" +
				"   * `user`: If the app was created by a user, the computed value will be the matching user id.\n" +
				"   * `app`: If the app was created by a app, the computed value will be the matching app id.",
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": {
			    Description: "The description of the app.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"credential": {
			    Description: "The Fortanix DSM App API key.",
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"new_credential": {
			    Description: "Set this if you want to rotate/regenerate the API key. The values can be set as true/false.",
				Type:     schema.TypeBool,
				Optional: true,
			},
			"other_group_permissions": {
			    Description: "Incase if you want to change the default permissions of a new group that includes default group. Please refer the example.",
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
			},
			"mod_group_permissions": {
			    Description: "To modify the permissions of any existing group that includes default group. Please refer the example.",
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
			},
			"role": {
				Description: "Role of a DSM app. The allowed values are crypto and admin.\n" +
				"   * `crypto`: To perform the crypto operations.\n" +
				"   * `admin`: To perform the admin operations.",
				Type:     schema.TypeString,
				Optional: true,
				Default: "crypto",
				ValidateFunc: validation.StringInSlice([]string{"crypto", "admin"}, true),
			},
			"secret_size": {
				Description: "Size of an API key. Allowed values are 16, 32 and 64. Default value is 64.",
				Type:     schema.TypeInt,
				Optional: true,
				Default: 64,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// [C]: Create App
func resourceCreateApp(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	role := d.Get("role").(string)
	app_object := map[string]interface{}{
		"name":          d.Get("name").(string),
		//"add_groups": map[string]interface{}{
		//	d.Get("default_group").(string): []string{"SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE", "MACVERIFY", "EXPORT", "MANAGE", "AGREEKEY", "AUDIT"},
		//},
		"app_type":    "default",
		"description": d.Get("description").(string),
		"role": role,
		"secret_size": d.Get("secret_size").(int),
	}
	if err := appRoleValidation(d, app_object, role); err != nil {
	    return err
	}
	req, err := m.(*api_client).APICallBody("POST", "sys/v1/apps", app_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST sys/v1/apps: %v", err),
		})
		return diags
	}

	d.SetId(req["app_id"].(string))
	return resourceReadApp(ctx, d, m)
}

// [R]: Read App
func resourceReadApp(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, _, err := m.(*api_client).APICall("GET", fmt.Sprintf("sys/v1/apps/%s", d.Id()))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: GET sys/v1/apps: %v", err),
		})
		return diags
	}

	if err := d.Set("name", req["name"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("app_id", req["app_id"].(string)); err != nil {
		return diag.FromErr(err)
	}
	// Admin app will not have default_group
	if _, ok := req["default_group"]; ok {
		if err := d.Set("default_group", req["default_group"].(string)); err != nil {
			return diag.FromErr(err)
		}
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
	if err := d.Set("role", req["role"].(string)); err != nil {
		return diag.FromErr(err)
	}

	req, _, err = m.(*api_client).APICall("GET", fmt.Sprintf("sys/v1/apps/%s/credential", d.Id()))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: GET sys/v1/apps/-/credential: %v", err),
		})
		return diags
	}

	if err := d.Set("credential", base64.StdEncoding.EncodeToString([]byte(d.Id()+":"+req["credential"].(map[string]interface{})["secret"].(string)))); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("new_credential", false); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

// [U]: Update App
func resourceUpdateApp(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	error_summary := "[DSM SDK] Unable to call DSM provider API client"
	if d.Get("new_credential").(bool) {
		reset_secret := map[string]interface{}{
			"credential_migration_period": nil,
			"secret_size": d.Get("secret_size").(int),
		}

		_, err := m.(*api_client).APICallBody("POST", fmt.Sprintf("sys/v1/apps/%s/reset_secret", d.Id()), reset_secret)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  error_summary,
				Detail:   fmt.Sprintf("[E]: API: GET sys/v1/apps/-/credential: %v", err),
			})
			return diags
		}
	}
	if d.HasChange("role") {
		old_role, new_role := d.GetChange("role")
		return invokeErrorDiagsWithSummary(error_summary, fmt.Sprintf("[E]: API: PATCH sys/v1/apps/%s: role cannot be changed once it is set. Please retain it to old value: %s -> %s", d.Id(), old_role, new_role))
	}
	if !d.Get("new_credential").(bool) && d.HasChange("secret_size") {
		old_ss, new_ss := d.GetChange("secret_size")
		return invokeErrorDiagsWithSummary(error_summary, fmt.Sprintf("[E]: API: PATCH sys/v1/apps/%s: secret_size cannot be changed once it is set. It can be updated only when new_credential is initiated. Please retain it to old value: %s -> %s", d.Id(), old_ss, new_ss))
	}
	//Modified by Ravi Gopal
	app_object := make(map[string]interface{})
	app_object["description"] = d.Get("description")
	app_object["name"] = d.Get("name")
	if d.HasChange("default_group") {
		if default_group := d.Get("default_group").(string); len(default_group) > 0 {
			app_object["default_group"] = d.Get("default_group")
		}
	}
	if d.HasChange("other_group") {
		getChangesInOtherGroups(d, app_object)
	}
	if d.HasChange("description") {
		app_object["description"] = d.Get("description")
	}
	//Modifies the existing groups
	if d.HasChange("mod_group_permissions") {
		err := getChangesInGroupPermissions(d, app_object)
		if err != nil {
			return err
		}
	}
	if len(app_object) > 0 {
		req, err := m.(*api_client).APICallBody("PATCH", fmt.Sprintf("sys/v1/apps/%s", d.Id()), app_object)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  error_summary,
				Detail:   fmt.Sprintf("[E]: API: POST sys/v1/apps: %v", err),
			})

			return diags
		}
		d.SetId(req["app_id"].(string))
	}
	return resourceReadApp(ctx, d, m)
}

// [D]: Delete App
func resourceDeleteApp(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	_, statuscode, err := m.(*api_client).APICall("DELETE", fmt.Sprintf("sys/v1/apps/%s", d.Id()))
	if (err != nil) && (statuscode != 404) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: DELETE sys/v1/apps: %v", err),
		})
		return diags
	}

	d.SetId("")
	return nil
}