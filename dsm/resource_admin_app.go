// **********
// Terraform Provider - DSM: resource: app
// **********
//       - Author:    ravi.gopal at fortanix dot com
//       - Version:   0.5.33
//       - Date:      03/09/24
// **********

package dsm

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// [-] Define App
func resourceAdminApp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateAdminApp,
		ReadContext:   resourceReadAdminApp,
		UpdateContext: resourceUpdateAdminApp,
		DeleteContext: resourceDeleteAdminApp,
		Description: "Creates a new DSM Admin app. The returned resource object contains the UUID of the app for further references. Using dsm_admin_app following " +
		"apps can be created:\n" +
		"   * `secret`\n" +
		"   * `awsxks`\n" +
		"   * `awsiam`\n" +
		"   * `certificate`\n" +
		"   * `trustedca`\n",
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
				Description: "The Fortanix DSM App credentials.\n\n" +
				"   * When the authentication method is secret. An API key will be mapped here.\n" +
				"   * When the authentication method is awsxks. AWSXKS access and secret keys will be mapped here.\n" +
				"   * When the authentication method is certificate/trustedca. A certificate will be mapped here.\n",
				Type:      schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
				},
			},
			"new_credential": {
				Description: "Set this if you want to rotate/regenerate the API key(secret). The values can be set as true/false.\n\n" +
				"**Note**: This should set only when the type of authentication_method is secret.",
				Type:     schema.TypeBool,
				Optional: true,
			},
			"authentication_method": {
				Description: "The Authentication type of an app.\n" +
				"   * `type`:  Following authentication types are supported.\n" +
				"       * api_key, awsxks, awsiam, certificate and trustedca.\n" +
				"   * `certificate`: Certificate value, this should be configured when the type is certificate.\n" +
				"   * `ca_certificate`: CA certificate value, this should be configured when the type is trustedca.\n" +
				"   One of the following parameters should be given when the type is trustedca.\n" +
				"   * `ip_address`:  IP address value for trusted ca.\n" +
				"   * `dns_name`:  DNS name for trusted ca.\n" +
				"   **Note**: For more details refer the above examples.",
				Type:      schema.TypeMap,
				Required:  true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// [C]: Create App
func resourceCreateAdminApp(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	app_object := map[string]interface{}{
		"name":          d.Get("name").(string),
		"app_type":    "default",
		"description": d.Get("description").(string),
		"role": "admin",
	}
	if am := d.Get("authentication_method").(map[string]interface{}); len(am) > 0 {
		formCredential(d, app_object, am)
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
	return resourceReadAdminApp(ctx, d, m)
}

// [R]: Read App
func resourceReadAdminApp(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	req, _, err = m.(*api_client).APICall("GET", fmt.Sprintf("sys/v1/apps/%s/credential", d.Id()))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: GET sys/v1/apps/-/credential: %v", err),
		})
		return diags
	}

	credential := make(map[string]interface{})
	req_credential := req["credential"].(map[string]interface{})
	if val, ok := req_credential["awsxks"].(map[string]interface{}); ok {
		credential["access_key_id"] =  base64.StdEncoding.EncodeToString([]byte(val["access_key_id"].(string)))
		credential["secret_key"] = base64.StdEncoding.EncodeToString([]byte(val["secret_key"].(string)))
		credential["path_prefix"] = "/crypto/v1/apps/" + d.Id() + "/aws"
	} else if val, ok := req_credential["certificate"].(map[string]interface{}); ok {
		credential["certificate"] = val["certificate"].(string)
	} else if val, ok := req_credential["secret"].(string); ok {
		credential["secret"] = base64.StdEncoding.EncodeToString([]byte(d.Id()+":"+ val))
	}
	if err := d.Set("credential", credential); err != nil {
			return diag.FromErr(err)
	}
	if err := d.Set("new_credential", false); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

// [U]: Update App
func resourceUpdateAdminApp(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	error_summary := "[DSM SDK] Unable to call DSM provider API client"
	if d.Get("new_credential").(bool) {
		reset_secret := map[string]interface{}{
			"credential_migration_period": nil,
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

	app_object := make(map[string]interface{})
	app_object["description"] = d.Get("description")
	app_object["name"] = d.Get("name")
	if d.HasChange("description") {
		app_object["description"] = d.Get("description")
	}
	if d.HasChange("authentication_method") {
		am := d.Get("authentication_method").(map[string]interface{})
		formCredential(d, app_object, am)
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
	return resourceReadAdminApp(ctx, d, m)
}

// [D]: Delete App
func resourceDeleteAdminApp(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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